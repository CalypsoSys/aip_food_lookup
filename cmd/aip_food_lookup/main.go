//go:build linux && !appengine && !heroku
// +build linux,!appengine,!heroku

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	allFoods       map[string]string
	allowedFood    map[string]string
	disAllowedFood map[string]string
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Fprint(w, "Hello, this is your Go HTTP service!")
}

func main() {
	allFoods = map[string]string{}
	allowedFood = map[string]string{}
	disAllowedFood = map[string]string{}

	dataFolder := os.Getenv("AIP_DATA_FOLDER")
	processDirectory(dataFolder)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
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
		line := scanner.Text()
		lower := strings.ToLower(line)
		if cat, exists := allFoods[lower]; !exists {
			allFoods[lower] = fmt.Sprintf("%s-%s", allowed, category)
		} else {
			if cat != fmt.Sprintf("%s-%s", allowed, category) {
				fmt.Println(line, cat)
			}
			fmt.Println(line)
		}

		if allowed == "allowed" {
			allowedFood[line] = category
		} else {
			disAllowedFood[line] = category
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
