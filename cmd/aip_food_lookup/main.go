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

var (
	nameFoosMap map[string]*apiFood
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	response := match(key)

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

func main() {
	nameFoosMap = map[string]*apiFood{}

	dataFolder := os.Getenv("AIP_DATA_FOLDER")
	processDirectory(dataFolder)

	match("Sel")

	http.HandleFunc("/", handler)
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
