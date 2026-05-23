package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CalypsoSys/godoublemetaphone/pkg/godoublemetaphone"
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

func TestFeedbackHandlerWritesValidFeedback(t *testing.T) {
	tempDir := t.TempDir()
	store = newFoodStore(tempDir)

	body := strings.NewReader(`{"name":"Joe","email":"","subject":"","message":"Great app","source":"mobile"}`)
	request := httptest.NewRequest(http.MethodPost, "/feedback", body)
	response := httptest.NewRecorder()

	feedbackHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "feedback.jsonl"))
	if err != nil {
		t.Fatalf("expected feedback file to be written: %v", err)
	}
	if !strings.Contains(string(content), `"message":"Great app"`) {
		t.Fatalf("expected feedback message in file, got %q", string(content))
	}
	if !strings.Contains(string(content), `"subject":"App feedback"`) {
		t.Fatalf("expected default feedback subject, got %q", string(content))
	}
}

func TestFeedbackHandlerRequiresNameOrEmailAndMessage(t *testing.T) {
	store = newFoodStore(t.TempDir())

	body := strings.NewReader(`{"message":"","source":"mobile"}`)
	request := httptest.NewRequest(http.MethodPost, "/feedback", body)
	response := httptest.NewRecorder()

	feedbackHandler(response, request)

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

func TestSoundSearchPorkDoesNotReturnUnrelatedCatalogItems(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	result := testStore.match("pork", "searchbysound")

	if !contains(result.Allowed, "Pork") {
		t.Fatalf("expected Pork in sound search results, got allowed=%#v notAllowed=%#v", result.Allowed, result.NotAllowed)
	}
	for _, unrelated := range []string{"Berries", "Brain", "Pears", "Perch"} {
		if contains(result.Allowed, unrelated) {
			t.Fatalf("expected %s not to appear in pork sound search results, got %#v", unrelated, result.Allowed)
		}
	}
	for _, unrelated := range []string{"Barley", "Butter"} {
		if contains(result.NotAllowed, unrelated) {
			t.Fatalf("expected %s not to appear in pork sound search results, got %#v", unrelated, result.NotAllowed)
		}
	}
}

func TestSoundSearchHelpsWithPorkMisspelling(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	result := testStore.match("porc", "searchbysound")

	if !contains(result.Allowed, "Pork") {
		t.Fatalf("expected Pork for porc sound search, got allowed=%#v notAllowed=%#v", result.Allowed, result.NotAllowed)
	}
}

func TestSoundSearchHelpsWithCommonCatalogMisspellings(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	tests := []struct {
		name      string
		query     string
		allowed   bool
		want      string
		unrelated string
	}{
		{
			name:      "porc finds pork",
			query:     "porc",
			allowed:   true,
			want:      "Pork",
			unrelated: "Perch",
		},
		{
			name:    "tomatoe finds tomatoes",
			query:   "tomatoe",
			allowed: false,
			want:    "Tomatoes",
		},
		{
			name:    "berys finds berries",
			query:   "berys",
			allowed: true,
			want:    "Berries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testStore.match(tt.query, "searchbysound")
			results := result.NotAllowed
			if tt.allowed {
				results = result.Allowed
			}
			if !contains(results, tt.want) {
				t.Fatalf("expected %s for %q sound search, got allowed=%#v notAllowed=%#v", tt.want, tt.query, result.Allowed, result.NotAllowed)
			}
			if tt.unrelated != "" && contains(result.Allowed, tt.unrelated) {
				t.Fatalf("expected %s not to appear for %q sound search, got %#v", tt.unrelated, tt.query, result.Allowed)
			}
		})
	}
}

func TestSoundSearchCumminFindsCuminPhrasesWithoutPumpkin(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	result := testStore.match("cummin", "searchbysound")

	for _, want := range []string{"Black Cumin", "Cumin Seed"} {
		if !contains(result.NotAllowed, want) {
			t.Fatalf("expected %s for cummin sound search, got allowed=%#v notAllowed=%#v", want, result.Allowed, result.NotAllowed)
		}
	}
	if contains(result.Allowed, "Pumpkin") {
		t.Fatalf("expected Pumpkin not to appear for cummin sound search, got %#v", result.Allowed)
	}
}

func TestCombinedSearchCumminFindsCuminPhrasesWithoutPumpkin(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	result := testStore.match("cummin", "")

	for _, want := range []string{"Black Cumin", "Cumin Seed"} {
		if !contains(result.NotAllowed, want) {
			t.Fatalf("expected %s for cummin combined search, got allowed=%#v notAllowed=%#v", want, result.Allowed, result.NotAllowed)
		}
	}
	if contains(result.Allowed, "Pumpkin") {
		t.Fatalf("expected Pumpkin not to appear for cummin combined search, got %#v", result.Allowed)
	}
}

func TestCombinedSearchPorkKeepsTextMatchAndRejectsUnrelatedSoundMatches(t *testing.T) {
	testStore := newFoodStore("../../data")
	if err := testStore.processDirectory("../../data"); err != nil {
		t.Fatalf("processDirectory returned error: %v", err)
	}

	result := testStore.match("pork", "")

	if !contains(result.Allowed, "Pork") {
		t.Fatalf("expected Pork in combined search results, got allowed=%#v notAllowed=%#v", result.Allowed, result.NotAllowed)
	}
	for _, unrelated := range []string{"Berries", "Brain", "Pears", "Perch"} {
		if contains(result.Allowed, unrelated) {
			t.Fatalf("expected %s not to appear in combined pork results, got %#v", unrelated, result.Allowed)
		}
	}
	for _, unrelated := range []string{"Barley", "Butter"} {
		if contains(result.NotAllowed, unrelated) {
			t.Fatalf("expected %s not to appear in combined pork results, got %#v", unrelated, result.NotAllowed)
		}
	}
}

func TestSpellingDistanceAllowedBalancesTyposAndShortWordNoise(t *testing.T) {
	tests := []struct {
		query     string
		candidate string
		want      bool
	}{
		{query: "porc", candidate: "Pork", want: true},
		{query: "pork", candidate: "Perch", want: false},
		{query: "tomatoe", candidate: "Tomatoes", want: true},
		{query: "berys", candidate: "Berries", want: true},
		{query: "pork", candidate: "Barley", want: false},
		{query: "cummin", candidate: "Cumin", want: true},
		{query: "cummin", candidate: "Pumpkin", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.query+"_"+tt.candidate, func(t *testing.T) {
			if got := spellingDistanceAllowed(tt.query, tt.candidate); got != tt.want {
				t.Fatalf("spellingDistanceAllowed(%q, %q) = %v, want %v", tt.query, tt.candidate, got, tt.want)
			}
		})
	}
}

func TestSoundSearchMatchesAlternateMetaphoneKey(t *testing.T) {
	query := godoublemetaphone.NewShortDoubleMetaphone("apple")
	testStore := newFoodStore("")
	testStore.nameFoods["applf"] = &apiFood{
		allowed:                 true,
		name:                    "Applf",
		primaryShortMetaphone:   query.PrimaryShortKey() + 1,
		alternateShortMetaphone: query.PrimaryShortKey(),
	}

	result := testStore.match("apple", "searchbysound")

	if !contains(result.Allowed, "Applf") {
		t.Fatalf("expected alternate metaphone key match, got %#v", result.Allowed)
	}
}

func TestSoundSearchDoesNotMatchNearbyNumericMetaphoneKey(t *testing.T) {
	query := godoublemetaphone.NewShortDoubleMetaphone("apple")
	testStore := newFoodStore("")
	testStore.nameFoods["near numeric key"] = &apiFood{
		allowed:                 true,
		name:                    "Near Numeric Key",
		primaryShortMetaphone:   query.PrimaryShortKey() + 1,
		alternateShortMetaphone: godoublemetaphone.METAPHONE_INVALID_KEY,
	}

	result := testStore.match("apple", "searchbysound")

	if contains(result.Allowed, "Near Numeric Key") {
		t.Fatalf("expected nearby numeric metaphone key not to match, got %#v", result.Allowed)
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
