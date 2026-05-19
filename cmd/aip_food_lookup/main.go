package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/CalypsoSys/godoublemetaphone/pkg/godoublemetaphone"
)

const (
	allowedNotAllowedLimit  = 10000
	allowedNotAllowedMinLen = 3
	allowedNotAllowedMaxLen = 50
)

type responseData struct {
	Allowed    []string `json:"allowed"`
	NotAllowed []string `json:"not_allowed"`
}

type requestData struct {
	InputText string `json:"inputText"`
	Allowed   bool   `json:"allowed"`
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
	nameFoods             map[string]*apiFood
}

var store = newFoodStore("")

// newFoodStore initializes the in-memory index and suggestion caches.
func newFoodStore(dataFolder string) *foodStore {
	return &foodStore{
		allowedCategories:     []string{},
		notAllowedCategories:  []string{},
		allowedSuggestions:    make(map[string]bool),
		notAllowedSuggestions: make(map[string]bool),
		dataFolder:            dataFolder,
		nameFoods:             make(map[string]*apiFood),
	}
}

func main() {
	dataFolder := os.Getenv("AIP_DATA_FOLDER")
	if dataFolder == "" {
		dataFolder = "data"
	}

	store = newFoodStore(dataFolder)
	if err := store.processDirectory(dataFolder); err != nil {
		fmt.Println("error loading data:", err)
	}

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/suggest", suggestHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/subcategory", subCategoryHandler)

	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}

// healthHandler gives load balancers and local smoke tests a simple API check.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(w, "AIP Food Lookup API")
}

// setCORS applies permissive CORS headers for the mobile and web clients.
func setCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// searchHandler returns matching allowed and not allowed foods for a query.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		http.Error(w, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	response := store.match(key, r.URL.Query().Get("type"))
	commonResponse(w, response)
}

// suggestHandler records user suggestions after basic length and ASCII cleanup.
func suggestHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

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
	if err := store.appendSuggestion(request.Allowed, request.InputText); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// categoriesHandler returns the available top-level allowed/not allowed groups.
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	commonResponse(w, responseData{
		Allowed:    store.allowedCategories,
		NotAllowed: store.notAllowedCategories,
	})
}

// subCategoryHandler returns foods for one allowed/not allowed category group.
func subCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

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

	response := store.subCategory(category, subCategory)
	commonResponse(w, response)
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
	primary := sdm.PrimaryShortKey()

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

		if soundSearch && math.Abs(float64(int64(primary)-int64(food.primaryShortMetaphone))) < 10 {
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
