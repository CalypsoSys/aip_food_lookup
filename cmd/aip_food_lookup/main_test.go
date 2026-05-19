package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessDirectoryLoadsCategoriesAndSearchData(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	if len(testStore.allowedCategories) == 0 {
		t.Fatal("expected allowed categories to load")
	}
	if len(testStore.notAllowedCategories) == 0 {
		t.Fatal("expected not allowed categories to load")
	}

	result := testStore.match("App", "searchbytext")
	if !contains(result.Allowed, "Apples") {
		t.Fatalf("expected Apples in allowed search results, got %#v", result.Allowed)
	}
}

func TestSuggestHandlerWritesValidSuggestion(t *testing.T) {
	tempDir := t.TempDir()
	store = newFoodStore(tempDir)

	body := strings.NewReader(`{"inputText":"cassava chips","allowed":true}`)
	request := httptest.NewRequest(http.MethodPost, "/suggest", body)
	response := httptest.NewRecorder()

	suggestHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}

	path := filepath.Join(tempDir, "suggested_allowed.txt")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected suggestion file to be written: %v", err)
	}
	if !strings.Contains(string(content), "cassava chips") {
		t.Fatalf("expected suggestion in file, got %q", string(content))
	}
}

func TestSuggestHandlerRejectsShortSuggestion(t *testing.T) {
	store = newFoodStore(t.TempDir())

	body := strings.NewReader(`{"inputText":"ab","allowed":false}`)
	request := httptest.NewRequest(http.MethodPost, "/suggest", body)
	response := httptest.NewRecorder()

	suggestHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestSearchHandlerReturnsJsonResults(t *testing.T) {
	store = newFoodStore("../../data")
	if err := store.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/search?key=App&type=searchbytext", nil)
	response := httptest.NewRecorder()

	searchHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}

	var result responseData
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("expected JSON response: %v", err)
	}
	if !contains(result.Allowed, "Apples") {
		t.Fatalf("expected Apples in allowed response, got %#v", result.Allowed)
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
