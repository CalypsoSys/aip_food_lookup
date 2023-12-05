//go:build linux && !appengine && !heroku
// +build linux,!appengine,!heroku

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

var (
	dataFolder  string
	nameFoosMap map[string]*apiFood
)

func setCORS(w http.ResponseWriter, r *http.Request) bool {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
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

	if len(requestData.InputText) < 3 {
		http.Error(w, "Suggestion to short", http.StatusBadRequest)
	}

	appendToFile(requestData.Allowed, requestData.InputText)
	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	nameFoosMap = map[string]*apiFood{}

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
		}

		return nil
	})

	return err
}

func appendToFile(alloweded bool, text string) error {
	var fileName string
	if alloweded {
		fileName = "suggested_allowed.txt"
	} else {
		fileName = "suggested_not_allowed.txt"
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

	return nil
}
