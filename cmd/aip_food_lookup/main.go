//go:build linux && !appengine && !heroku
// +build linux,!appengine,!heroku

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CalypsoSys/godoublemetaphone/pkg/godoublemetaphone"
)

type apiFood struct {
	allowed                 bool
	name                    string
	primaryShortMetaphone   uint16
	alternateShortMetaphone uint16
	category                string
}

type ResponseData struct {
	PossibleAllowed    []string `json:"possible_allowed"`
	PossibleDisallowed []string `json:"possible_disallowed"`
}

// RequestData represents the structure of the JSON request
type RequestData struct {
	InputText string `json:"inputText"`
	Allowed   bool   `json:"allowed"`
}

const (
	allowedNotAllowedLimit  = 10000
	allowedNotAllowedMinLen = 3
	allowedNotAllowedMaxLen = 50

	encodedPrivateKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBcjdmMlVYQ1lEWEtWa1BjN2wzWEZRb2dyTWV6REVZN2J2M01CL2syZ0YrcUJFbGJCCi9DeVRsdUZCVUk2YzRwVFRBbFVCcVBZQ1VxTTg0QnNPakY5TXJNTlExTTY2SWVnTmEzalJQckEzYjFya09kaEYKWCtQTTVEUHdCc3MzK3FaR1hiMjJka1lYRHhhZWp6UFVPelZyclBldUZuWTZTVTJoS1lUOFQ2clBnN0xMb0drQwowOFk5aXNUa3JDc3RWN3BxQUZ6eE1SWStFZXk4ZnoySy9CL05YOHhGejFGRndReWw2c3g2dWZMaDd3VWhyQ0NOCmdlQ1o2bEZCR1E3OXYzUXVsSmNwNVA4YWV5WU9UR2RGSm9RTGx6SU00dWlsM3hSM0FEakZOTzNQUnduTWJzN3MKckNiYVYyUytjVFNONmxweS8zQXFNSEVkNHZYdFRXVlNxaitaTFFJREFRQUJBb0lCQURQVW5IaExNTjZPbE9Wcwp0NHdtZGVmNUNGeXlqSnRxT3hGT21DRHR2ckl2UHFFdExBejVEUk90SDdubVJ3cGlnWmNuZ2RUWHM3bXlZcXRyCjc2K2lFSmpKQjlldG5xT1BzaDJvUm5ncVBEL0JYSjVmVjU5QUwxaUVwV0Vyb2poeHdVRzNTdEc2UE9UN2RBdWoKYXcrSDQxbml1TnZ4UmFJSG51a0RTL1VuMmd2clBQMitINTM3VmNZTWhZTFVIZmsxN1lBVlZycUFHUzBQUlVhKwpmRGROZUJNUGtkUGI0VEZmR0lTemxldElaVlJwbVM3YnBMSi9naVR4cERvZ3kvMWNocHc0Y3RwQ25LSzN6UEY2Cmd6VFpHS3didDhodU0wc3JoR1JpdXhSbUNia3M5eHBVMUwyTnRSZit4WEZLWWF2K0tPejRoSVk1Nmd3eWU2RUUKNjZMdFNqVUNnWUVBM1hvTlhwV1JrUWo1WC93cEZiL2ZKSzVEc2ZDajdWQnRTaDFKZlUxMytvWmdBcldEYVNncApuZlh3WEtVWktKMTZSNi9WWThVMWJvb0hKejN3ZUN2Rkc0N1IwOTlpWUJ0U2xjdUM2Myt1bkYrUXE3M01TK1I5Cm0vTUc2Z0RzQXgzZFZJVXNObDhzWHB0TDJKUFdwdWpKMVFCS05kRnArSkdvS0hzSEVnbFB0SnNDZ1lFQXl4dnkKamNTczh5YzQ3S1BvdFQ4YWVGaUNwb0ptSEJZZ3R5QnhoL3ZralFuZTFCNTAwTE9pUXRuR3pTUmhYUlZ1OTVjOQpIajErWENUZS9WSURyVHpObVlqcVR6Z0pDekVjdjZRY3ZqTnF5NHR0UWExOHRnYWhNT1RRdG4zdGdmVWxhMkYvCnFERWIrc0NKRXNpNUxHaFlqQ2dhV09Ba0lZcTM0eG9POGQvcjhkY0NnWUVBMkZqSXhKTlFyaC9aRW1WTmNQeU0KS3RXOE5RNy80dXRFeHpoU3VIODdhMU5tYUY4TmJtU1lPc0NyT3FUZ0xhZWZjbldWK3E4REllYmRVLzBTY1NFNAptMUhwTUpHdkZIaThOSzJuUndyajg4YjZtSG1BSHNhbDJQZ08wZmx5a3h6U1B5VVQ2azBRRjU2VitZdDVESFNyCjdGRXJMT1ZUSWtpT3ZuUm5sTHZaeTI4Q2dZQkZ2K0U2QWpLS2hndXNhRlYvK0oyMGVtRFRvYkJETU80bk5VTUgKdWQ4dytCVEhyM1hhUGZZWkV3U01hbFB0VFhFQUliWGhicWk0S0FsVDRSaFdJNjFQYm85WWlSdkI5aW16UGo2SQpxc3VmL3MrVVlHbVZjUTFsNXc0dHZXMFUxZ1QxclZQVGhKbmhNTUZoN0FCN1daSWUvNTZjcXN4OW9FK3A4OGJ5CkZUM0huUUtCZ0FaV3E5WmVoVFhDcWpvTTZibjhqL1g0RHdMUm81ZWp3WThONDVlbk1XNVJ1RlV1WThTUFQrbGIKWDZHNHdIdFpLcFdYNlRaNjQ0Ump2K05XRDJMZ3lqcmJhdkV4eFRmOGNEMEFkdFVveFUzbHBaM0g4RmJEK3dtZwpmZkUxSUxSK2tGM0FiWlV0S242NVk4UUUvaFJzMkUxOWMwV2ZJb3UvaFNXM25IQmxEZlFoCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="
	encodedPublicKey  = "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBcjdmMlVYQ1lEWEtWa1BjN2wzWEYKUW9nck1lekRFWTdidjNNQi9rMmdGK3FCRWxiQi9DeVRsdUZCVUk2YzRwVFRBbFVCcVBZQ1VxTTg0QnNPakY5TQpyTU5RMU02NkllZ05hM2pSUHJBM2IxcmtPZGhGWCtQTTVEUHdCc3MzK3FaR1hiMjJka1lYRHhhZWp6UFVPelZyCnJQZXVGblk2U1UyaEtZVDhUNnJQZzdMTG9Ha0MwOFk5aXNUa3JDc3RWN3BxQUZ6eE1SWStFZXk4ZnoySy9CL04KWDh4RnoxRkZ3UXlsNnN4NnVmTGg3d1VockNDTmdlQ1o2bEZCR1E3OXYzUXVsSmNwNVA4YWV5WU9UR2RGSm9RTApseklNNHVpbDN4UjNBRGpGTk8zUFJ3bk1iczdzckNiYVYyUytjVFNONmxweS8zQXFNSEVkNHZYdFRXVlNxaitaCkxRSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K"
	stupid            = "V29uJ3Qgc2VjdXJlIGFueXRoaW5nLCBidXQgcGxlYXNlIGRvIG5vdCBhYnVzZSwgYXBpIGZvb2QgbG9va3VwIHYxLjAuMA=="
)

var (
	dataFolder            string
	nameFoosMap           map[string]*apiFood
	allowedSuggestions    map[string]bool
	notAllowedSuggestions map[string]bool
)

func decodePrivateKey(encodedPrivateKey string) (*rsa.PrivateKey, error) {
	privPEM, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func decrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

func setCORS(w http.ResponseWriter, r *http.Request) bool {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Nonsense-I-Know")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	nonSense := r.Header.Get("Nonsense-I-Know")

	if nonSense == "" {
		http.Error(w, "nonsense is missing", http.StatusBadRequest)
		return true
	}
	decrtptedNonsense, err := base64.StdEncoding.DecodeString(nonSense)
	if err != nil || len(decrtptedNonsense) == 0 {
		http.Error(w, "cannot decrypt nonsense", http.StatusBadRequest)
		return true
	}
	decrtptedStupid, err := base64.StdEncoding.DecodeString(stupid)
	if err != nil || len(decrtptedStupid) == 0 {
		http.Error(w, "cannot decrypt stupid", http.StatusBadRequest)
		return true
	}
	privateKey, err := decodePrivateKey(encodedPrivateKey)
	if err != nil || len(decrtptedStupid) == 0 {
		http.Error(w, "cannot get private key", http.StatusBadRequest)
		return true
	}
	test, err := decrypt(privateKey, decrtptedNonsense)
	if err != nil || len(test) == 0 {
		http.Error(w, "cannot decrypt with private key", http.StatusBadRequest)
		return true
	}
	if !bytes.Equal(test, decrtptedStupid) {
		http.Error(w, "wrong key", http.StatusBadRequest)
		return true
	}

	return false
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	var response ResponseData
	if key == "longlistjoe" {
		response = ResponseData{
			PossibleAllowed:    []string{"a one", "a two", "a three", "a four", "a five", "a six", "a seven", "a eight", "a nine", "a ten"},
			PossibleDisallowed: []string{"d one", "d two", "d three", "d four", "d five", "d six", "d seven", "d eight", "d nine", "d ten"},
		}
	} else {
		response = match(key)
	}

	// Convert the response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON data to the response
	w.Write(jsonData)
}

func suggestHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	// Parse the JSON request
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	if len(requestData.InputText) > allowedNotAllowedMaxLen {
		http.Error(w, "Suggestion to long", http.StatusBadRequest)
		return
	}

	requestData.InputText = stripNonASCII(requestData.InputText)

	if len(requestData.InputText) < allowedNotAllowedMinLen {
		http.Error(w, "Suggestion to short", http.StatusBadRequest)
		return
	}

	appendToFile(requestData.Allowed, requestData.InputText)
	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	nameFoosMap = map[string]*apiFood{}
	allowedSuggestions = map[string]bool{}
	notAllowedSuggestions = map[string]bool{}

	dataFolder = os.Getenv("AIP_DATA_FOLDER")
	processDirectory(dataFolder)

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/suggest", suggestHandler)
	http.ListenAndServe(":8080", nil)
}

func match(name string) ResponseData {
	name = strings.ToLower(name)
	sdm := godoublemetaphone.NewShortDoubleMetaphone(name)
	primary := sdm.PrimaryShortKey()
	possibleAllowed := []string{}
	possibleDisallowed := []string{}
	soundPossibleAllowed := []string{}
	soundPossibleDisallowed := []string{}
	for k, v := range nameFoosMap {
		if strings.HasPrefix(k, name) {
			if v.allowed {
				possibleAllowed = append(possibleAllowed, v.name)
			} else {
				possibleDisallowed = append(possibleDisallowed, v.name)
			}
		} else if math.Abs(float64(int(primary)-int(v.primaryShortMetaphone))) < 10 {
			if v.allowed {
				soundPossibleAllowed = append(possibleAllowed, v.name)
			} else {
				soundPossibleDisallowed = append(possibleDisallowed, v.name)
			}
		}
	}

	if len(possibleAllowed) == 0 && len(soundPossibleAllowed) > 0 {
		possibleAllowed = soundPossibleAllowed
	}

	if len(possibleDisallowed) == 0 && len(soundPossibleDisallowed) > 0 {
		possibleDisallowed = soundPossibleDisallowed
	}

	return ResponseData{PossibleAllowed: possibleAllowed, PossibleDisallowed: possibleDisallowed}
}

func getParentFolder(path string) (string, error) {
	dir := filepath.Dir(path)
	parent := filepath.Base(dir)
	return parent, nil
}

func getFileNameWithoutExtension(filePath string) (string, error) {
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return name, nil
}

func processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	allowed, _ := getParentFolder(filePath)
	if allowed != "allowed" && allowed != "not_allowed" {
		return errors.New("must be allowed or not allowed")
	}
	category, _ := getFileNameWithoutExtension(filePath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lower := strings.ToLower(line)

			isAllowed := allowed == "allowed"
			if cat, exists := nameFoosMap[lower]; !exists {
				sdm := godoublemetaphone.NewShortDoubleMetaphone(line)

				nameFoosMap[lower] = &apiFood{
					allowed:                 isAllowed,
					name:                    line,
					primaryShortMetaphone:   sdm.PrimaryShortKey(),
					alternateShortMetaphone: sdm.AlternateShortKey(),
					category:                category,
				}
			} else {
				if cat.allowed != isAllowed || cat.category != category {
					fmt.Println(line, cat)
				}
				fmt.Println(line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processDirectory(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".dat" {
			if err := processFile(path); err != nil {
				fmt.Println("Error processing file:", err)
			}
		} else if filepath.Base(path) == "suggested_allowed.txt" {
			loadCurrentSuggested(path, allowedSuggestions)
		} else if filepath.Base(path) == "suggested_not_allowed.txt" {
			loadCurrentSuggested(path, notAllowedSuggestions)
		}

		return nil
	})

	return err
}

func appendToFile(alloweded bool, text string) error {
	var fileName string
	var cache map[string]bool

	text = strings.ToLower(strings.TrimSpace(text))
	if _, exists := nameFoosMap[text]; exists {
		return nil
	}

	if alloweded {
		fileName = "suggested_allowed.txt"
		cache = allowedSuggestions
	} else {
		fileName = "suggested_not_allowed.txt"
		cache = notAllowedSuggestions
	}

	if _, exists := allowedSuggestions[text]; exists {
		return nil
	}
	if _, exists := notAllowedSuggestions[text]; exists {
		return nil
	}

	if len(cache) > allowedNotAllowedLimit {
		return errors.New("allowed/not allowed exceded")
	}

	filePath := path.Join(dataFolder, fileName)
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffered writer to efficiently write to the file
	writer := bufio.NewWriter(file)

	// Write the text followed by a new line
	_, err = fmt.Fprintln(writer, text)
	if err != nil {
		return err
	}

	// Flush the writer to ensure the data is written to the file
	err = writer.Flush()
	if err != nil {
		return err
	}

	cache[text] = true

	return nil
}

func stripNonASCII(input string) string {
	// Use a regular expression to match non-ASCII characters
	regex := regexp.MustCompile("[^[:ascii:]]")
	return regex.ReplaceAllString(input, "")
}

func loadCurrentSuggested(filePath string, cache map[string]bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		cache[line] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
