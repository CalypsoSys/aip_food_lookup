package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/CalypsoSys/godoublemetaphone/pkg/godoublemetaphone"
)

const (
	allowedNotAllowedLimit  = 10000
	allowedNotAllowedMinLen = 3
	allowedNotAllowedMaxLen = 50
	feedbackMessageMaxLen   = 2000
	feedbackFieldMaxLen     = 200
	adminReloadPath         = "/admin/reload"
)

type responseData struct {
	Allowed    []string `json:"allowed"`
	NotAllowed []string `json:"not_allowed"`
}

type requestData struct {
	InputText string `json:"inputText"`
	Allowed   bool   `json:"allowed"`
}

type adminReloadResponse struct {
	OK                   bool   `json:"ok"`
	AllowedCategories    int    `json:"allowedCategories,omitempty"`
	NotAllowedCategories int    `json:"notAllowedCategories,omitempty"`
	Foods                int    `json:"foods,omitempty"`
	Error                string `json:"error,omitempty"`
}

type feedbackRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Source  string `json:"source"`
}

type apiFood struct {
	allowed                 bool
	name                    string
	primaryShortMetaphone   uint16
	alternateShortMetaphone uint16
	category                string
}

type foodStore struct {
	allowedCategories     []string
	notAllowedCategories  []string
	allowedSuggestions    map[string]bool
	notAllowedSuggestions map[string]bool
	dataFolder            string
	errorLogPath          string
	feedbackSink          feedbackSink
	suggestionSink        suggestionSink
	nameFoods             map[string]*apiFood
}

type feedbackSink interface {
	submitFeedback(feedbackRequest) error
}

type fileFeedbackSink struct {
	dataFolder string
	filePath   string
}

type suggestionSink interface {
	submitSuggestion(requestData) error
}

var (
	store     = newFoodStore("")
	storeLock sync.RWMutex
)

// newFoodStore initializes the in-memory index and suggestion caches.
func newFoodStore(dataFolder string) *foodStore {
	return &foodStore{
		allowedCategories:     []string{},
		notAllowedCategories:  []string{},
		allowedSuggestions:    make(map[string]bool),
		notAllowedSuggestions: make(map[string]bool),
		dataFolder:            dataFolder,
		feedbackSink:          fileFeedbackSink{dataFolder: dataFolder},
		nameFoods:             make(map[string]*apiFood),
	}
}

func main() {
	config := loadConfig()

	store = newFoodStore(config.DataFolder)
	store.errorLogPath = config.ErrorLogPath
	store.feedbackSink = newFeedbackSink(config)
	store.suggestionSink = newSuggestionSink(config)
	if err := store.processDirectory(config.DataFolder); err != nil {
		fmt.Println("error loading data:", err)
	}

	mux := http.NewServeMux()
	registerHandlers(mux)

	server := &http.Server{
		Addr:              config.ListenAddress,
		Handler:           buildHTTPHandler(config, mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err := server.ListenAndServe()
	fmt.Println(err)
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", healthHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/suggest", suggestHandler)
	mux.HandleFunc("/feedback", feedbackHandler)
	mux.HandleFunc("/categories", categoriesHandler)
	mux.HandleFunc("/subcategory", subCategoryHandler)
	mux.HandleFunc(adminReloadPath, adminReloadHandler)
}

// healthHandler gives load balancers and local smoke tests a simple API check.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(w, "AIP Food Lookup API")
}

// feedbackHandler validates app feedback and stores it for later Slack plumbing.
func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request feedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	normalized, err := normalizeFeedback(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currentStore := getStore()
	if err := currentStore.feedbackSink.submitFeedback(normalized); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// searchHandler returns matching allowed and not allowed foods for a query.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		http.Error(w, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	currentStore := getStore()
	response := currentStore.match(key, r.URL.Query().Get("type"))
	commonResponse(w, response)
}

// suggestHandler records user suggestions after basic length and ASCII cleanup.
func suggestHandler(w http.ResponseWriter, r *http.Request) {
	var request requestData
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	request.InputText = stripNonASCII(strings.TrimSpace(request.InputText))
	if len(request.InputText) > allowedNotAllowedMaxLen {
		http.Error(w, "Suggestion too long", http.StatusBadRequest)
		return
	}
	if len(request.InputText) < allowedNotAllowedMinLen {
		http.Error(w, "Suggestion too short", http.StatusBadRequest)
		return
	}
	currentStore := getStore()
	if err := currentStore.submitSuggestion(request.Allowed, request.InputText); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// categoriesHandler returns the available top-level allowed/not allowed groups.
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	currentStore := getStore()
	commonResponse(w, responseData{
		Allowed:    currentStore.allowedCategories,
		NotAllowed: currentStore.notAllowedCategories,
	})
}

// subCategoryHandler returns foods for one allowed/not allowed category group.
func subCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("cat")
	if category == "" {
		http.Error(w, "Category parameter is missing", http.StatusBadRequest)
		return
	}

	subCategory := r.URL.Query().Get("sub")
	if subCategory == "" {
		http.Error(w, "Sub Category parameter is missing", http.StatusBadRequest)
		return
	}

	currentStore := getStore()
	response := currentStore.subCategory(category, subCategory)
	commonResponse(w, response)
}

// adminReloadHandler reloads catalog files without restarting the container.
func adminReloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nextStore, err := reloadFoodStore()
	if err != nil {
		writeErrorLog(getStore().errorLogPath, fmt.Sprintf("catalog reload failed: %v", err))
		writeAdminReloadResponse(w, http.StatusInternalServerError, adminReloadResponse{
			OK:    false,
			Error: "catalog reload failed",
		})
		return
	}

	writeAdminReloadResponse(w, http.StatusOK, adminReloadResponse{
		OK:                   true,
		AllowedCategories:    len(nextStore.allowedCategories),
		NotAllowedCategories: len(nextStore.notAllowedCategories),
		Foods:                len(nextStore.nameFoods),
	})
}

func reloadFoodStore() (*foodStore, error) {
	currentStore := getStore()
	if currentStore == nil {
		return nil, errors.New("food store is not initialized")
	}

	dataFolder := strings.TrimSpace(currentStore.dataFolder)
	if dataFolder == "" {
		dataFolder = "data"
	}

	nextStore := newFoodStore(dataFolder)
	nextStore.errorLogPath = currentStore.errorLogPath
	nextStore.feedbackSink = currentStore.feedbackSink
	nextStore.suggestionSink = currentStore.suggestionSink
	if err := nextStore.processDirectory(dataFolder); err != nil {
		return nil, err
	}

	setStore(nextStore)
	return nextStore, nil
}

func writeAdminReloadResponse(w http.ResponseWriter, statusCode int, response adminReloadResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func getStore() *foodStore {
	storeLock.RLock()
	defer storeLock.RUnlock()
	return store
}

func setStore(nextStore *foodStore) {
	storeLock.Lock()
	defer storeLock.Unlock()
	store = nextStore
}

// commonResponse serializes the API's shared response envelope.
func commonResponse(w http.ResponseWriter, response responseData) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

// match combines prefix matching with Double Metaphone sound matching.
func (s *foodStore) match(name string, typeSearch string) responseData {
	name = strings.ToLower(strings.TrimSpace(name))
	sdm := godoublemetaphone.NewShortDoubleMetaphone(name)

	possibleAllowed := []string{}
	possibleNotAllowed := []string{}
	soundPossibleAllowed := []string{}
	soundPossibleNotAllowed := []string{}

	textSearch := true
	soundSearch := true
	if typeSearch == "searchbytext" {
		soundSearch = false
	} else if typeSearch == "searchbysound" {
		textSearch = false
	}

	for key, food := range s.nameFoods {
		if textSearch && strings.HasPrefix(key, name) {
			if food.allowed {
				possibleAllowed = append(possibleAllowed, food.name)
			} else {
				possibleNotAllowed = append(possibleNotAllowed, food.name)
			}
		}

		if soundSearch && fuzzySoundMatch(name, sdm, food) {
			if food.allowed {
				soundPossibleAllowed = append(soundPossibleAllowed, food.name)
			} else {
				soundPossibleNotAllowed = append(soundPossibleNotAllowed, food.name)
			}
		}
	}

	possibleAllowed = append(possibleAllowed, soundPossibleAllowed...)
	possibleNotAllowed = append(possibleNotAllowed, soundPossibleNotAllowed...)

	return responseData{
		Allowed:    sortedUnique(possibleAllowed),
		NotAllowed: sortedUnique(possibleNotAllowed),
	}
}

func fuzzySoundMatch(query string, queryMetaphone godoublemetaphone.ShortDoubleMetaphone, food *apiFood) bool {
	if spellingDistanceAllowed(query, food.name) {
		return true
	}
	for _, token := range searchableTokens(food.name) {
		if spellingDistanceAllowed(query, token) {
			return true
		}
	}
	if !metaphoneKeysMatch(queryMetaphone, food) {
		return false
	}

	query = strings.ToLower(strings.TrimSpace(query))
	candidate := strings.ToLower(strings.TrimSpace(food.name))
	limit := spellingDistanceLimit(query)
	if len(query) > 4 {
		limit++
	}
	if levenshteinDistance(query, candidate) <= limit {
		return true
	}
	for _, token := range searchableTokens(candidate) {
		if levenshteinDistance(query, token) <= limit {
			return true
		}
	}
	return false
}

func metaphoneKeysMatch(query godoublemetaphone.ShortDoubleMetaphone, food *apiFood) bool {
	queryKeys := validMetaphoneKeys(query.PrimaryShortKey(), query.AlternateShortKey())
	foodKeys := validMetaphoneKeys(food.primaryShortMetaphone, food.alternateShortMetaphone)
	for queryKey := range queryKeys {
		if foodKeys[queryKey] {
			return true
		}
	}
	return false
}

func validMetaphoneKeys(keys ...uint16) map[uint16]bool {
	valid := make(map[uint16]bool, len(keys))
	for _, key := range keys {
		if key == godoublemetaphone.METAPHONE_INVALID_KEY {
			continue
		}
		valid[key] = true
	}
	return valid
}

func spellingDistanceAllowed(query string, candidate string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	candidate = strings.ToLower(strings.TrimSpace(candidate))
	if query == "" || candidate == "" {
		return false
	}
	if query == candidate {
		return true
	}
	if query[0] != candidate[0] {
		return false
	}

	distance := levenshteinDistance(query, candidate)
	return distance <= spellingDistanceLimit(query)
}

func searchableTokens(candidate string) []string {
	candidate = strings.ToLower(candidate)
	fields := strings.FieldsFunc(candidate, func(r rune) bool {
		return r < 'a' || r > 'z'
	})
	tokens := []string{}
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) >= 3 {
			tokens = append(tokens, field)
		}
	}
	return tokens
}

func spellingDistanceLimit(query string) int {
	switch length := len(query); {
	case length <= 4:
		return 1
	case length <= 7:
		return 3
	default:
		return 3
	}
}

func levenshteinDistance(a string, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}

	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			current[j] = minInt(
				previous[j]+1,
				current[j-1]+1,
				previous[j-1]+cost,
			)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func minInt(values ...int) int {
	minimum := values[0]
	for _, value := range values[1:] {
		if value < minimum {
			minimum = value
		}
	}
	return minimum
}

// subCategory filters loaded foods by MAUI-compatible category route values.
func (s *foodStore) subCategory(category string, subCategory string) responseData {
	response := responseData{
		Allowed:    []string{},
		NotAllowed: []string{},
	}

	var allowed bool
	var output *[]string
	if category == "Allowed" {
		allowed = true
		output = &response.Allowed
	} else if category == "Not Allowed" {
		allowed = false
		output = &response.NotAllowed
	} else {
		return response
	}

	for _, food := range s.nameFoods {
		if food.allowed == allowed && food.category == subCategory {
			*output = append(*output, food.name)
		}
	}
	sort.Strings(*output)
	return response
}

// normalizeFeedback applies minimal privacy-conscious validation and cleanup.
func normalizeFeedback(request feedbackRequest) (feedbackRequest, error) {
	request.Name = stripNonASCII(strings.TrimSpace(request.Name))
	request.Email = stripNonASCII(strings.TrimSpace(request.Email))
	request.Subject = stripNonASCII(strings.TrimSpace(request.Subject))
	request.Message = stripNonASCII(strings.TrimSpace(request.Message))
	request.Source = stripNonASCII(strings.TrimSpace(request.Source))

	if request.Message == "" {
		return request, errors.New("Message is required")
	}
	if request.Name == "" && request.Email == "" {
		return request, errors.New("Name or email is required")
	}
	if len(request.Name) > feedbackFieldMaxLen ||
		len(request.Email) > feedbackFieldMaxLen ||
		len(request.Subject) > feedbackFieldMaxLen ||
		len(request.Source) > feedbackFieldMaxLen {
		return request, errors.New("Feedback field too long")
	}
	if len(request.Message) > feedbackMessageMaxLen {
		return request, errors.New("Feedback message too long")
	}
	if request.Subject == "" {
		request.Subject = "App feedback"
	}
	if request.Source == "" {
		request.Source = "mobile"
	}
	return request, nil
}

// sortedUnique removes duplicates and makes endpoint responses stable.
func sortedUnique(list []string) []string {
	unique := make(map[string]bool)
	result := []string{}
	for _, item := range list {
		if _, exists := unique[item]; exists {
			continue
		}
		unique[item] = true
		result = append(result, item)
	}
	sort.Strings(result)
	return result
}

func getParentFolder(p string) string {
	return filepath.Base(filepath.Dir(p))
}

func getFileNameWithoutExtension(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// processFile loads one .dat category file into the in-memory food index.
func (s *foodStore) processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	allowedFolder := getParentFolder(filePath)
	category := getFileNameWithoutExtension(filePath)

	if allowedFolder == "allowed" {
		s.allowedCategories = append(s.allowedCategories, convertPhrase(category))
	} else if allowedFolder == "not_allowed" {
		s.notAllowedCategories = append(s.notAllowedCategories, convertPhrase(category))
	} else {
		return errors.New("must be allowed or not_allowed")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		nameLower := strings.ToLower(line)
		if _, exists := s.nameFoods[nameLower]; exists {
			continue
		}

		sdm := godoublemetaphone.NewShortDoubleMetaphone(line)
		s.nameFoods[nameLower] = &apiFood{
			allowed:                 allowedFolder == "allowed",
			name:                    line,
			primaryShortMetaphone:   sdm.PrimaryShortKey(),
			alternateShortMetaphone: sdm.AlternateShortKey(),
			category:                category,
		}
	}

	return scanner.Err()
}

// processDirectory walks the configured data folder and loads all food data.
func (s *foodStore) processDirectory(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(p) == ".dat" {
			return s.processFile(p)
		}

		if filepath.Base(p) == "suggested_allowed.txt" {
			_ = loadCurrentSuggested(p, s.allowedSuggestions)
		}
		if filepath.Base(p) == "suggested_not_allowed.txt" {
			_ = loadCurrentSuggested(p, s.notAllowedSuggestions)
		}
		return nil
	})
	sort.Strings(s.allowedCategories)
	sort.Strings(s.notAllowedCategories)
	return err
}

// submitFeedback writes one JSON line until a Slack sink is added later.
func (s fileFeedbackSink) submitFeedback(request feedbackRequest) error {
	filePath := s.filePath
	if filePath == "" {
		filePath = path.Join(s.dataFolder, "feedback.jsonl")
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	if _, err = file.Write(append(payload, '\n')); err != nil {
		return err
	}
	return nil
}

// submitSuggestion attempts local storage and Slack notification independently.
func (s *foodStore) submitSuggestion(allowed bool, text string) error {
	request := requestData{
		InputText: text,
		Allowed:   allowed,
	}

	localErr := s.appendSuggestion(allowed, text)
	if localErr != nil {
		writeErrorLog(s.errorLogPath, fmt.Sprintf("suggestion file write failed: %v", localErr))
	}

	if s.suggestionSink == nil {
		return localErr
	}

	slackErr := s.suggestionSink.submitSuggestion(request)
	if slackErr != nil {
		writeErrorLog(s.errorLogPath, fmt.Sprintf("slack suggestion failed: %v", slackErr))
	}

	if localErr == nil || slackErr == nil {
		return nil
	}
	return errors.Join(localErr, slackErr)
}

// appendSuggestion persists a new user suggestion if it is not already known.
func (s *foodStore) appendSuggestion(allowed bool, text string) error {
	key := strings.ToLower(strings.TrimSpace(text))
	if _, exists := s.nameFoods[key]; exists {
		return nil
	}
	if _, exists := s.allowedSuggestions[key]; exists {
		return nil
	}
	if _, exists := s.notAllowedSuggestions[key]; exists {
		return nil
	}

	var fileName string
	var cache map[string]bool
	if allowed {
		fileName = "suggested_allowed.txt"
		cache = s.allowedSuggestions
	} else {
		fileName = "suggested_not_allowed.txt"
		cache = s.notAllowedSuggestions
	}

	if len(cache) > allowedNotAllowedLimit {
		return errors.New("suggestion limit exceeded")
	}

	filePath := path.Join(s.dataFolder, fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err = fmt.Fprintln(writer, key); err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}

	cache[key] = true
	return nil
}

// stripNonASCII keeps suggestion storage simple and compatible with old data.
func stripNonASCII(input string) string {
	regex := regexp.MustCompile("[^[:ascii:]]")
	return regex.ReplaceAllString(input, "")
}

// loadCurrentSuggested hydrates duplicate-detection caches for suggestion files.
func loadCurrentSuggested(filePath string, cache map[string]bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if line != "" {
			cache[line] = true
		}
	}

	return scanner.Err()
}

// convertPhrase maps category filenames to the labels used by the MAUI app.
func convertPhrase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(string(word[0])) + word[1:]
	}
	return strings.Join(words, " and ")
}
