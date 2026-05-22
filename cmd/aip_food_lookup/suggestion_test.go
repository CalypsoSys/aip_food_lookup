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

func TestSuggestHandlerSucceedsAndLogsWhenSlackFails(t *testing.T) {
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer slackServer.Close()

	tempDir := t.TempDir()
	store = newFoodStore(tempDir)
	store.errorLogPath = filepath.Join(tempDir, "errors.log")
	store.suggestionSink = newSuggestionSink(appConfig{
		SlackFeedbackWebhookURL: slackServer.URL,
	})

	body := strings.NewReader(`{"inputText":"cassava chips","allowed":true}`)
	request := httptest.NewRequest(http.MethodPost, "/suggest", body)
	response := httptest.NewRecorder()

	suggestHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "suggested_allowed.txt"))
	if err != nil {
		t.Fatalf("expected local suggestion file: %v", err)
	}
	if !strings.Contains(string(content), "cassava chips") {
		t.Fatalf("expected suggestion in local file, got %q", string(content))
	}

	errorLog, err := os.ReadFile(filepath.Join(tempDir, "errors.log"))
	if err != nil {
		t.Fatalf("expected Slack failure in error log: %v", err)
	}
	if !strings.Contains(string(errorLog), "slack suggestion failed") {
		t.Fatalf("expected Slack failure log, got %q", string(errorLog))
	}
}

func TestSuggestHandlerSucceedsAndLogsWhenLocalWriteFailsButSlackWorks(t *testing.T) {
	var slackPayload struct {
		Text string `json:"text"`
	}
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&slackPayload); err != nil {
			t.Fatalf("expected Slack JSON payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer slackServer.Close()

	tempDir := t.TempDir()
	store = newFoodStore(filepath.Join(tempDir, "missing"))
	store.errorLogPath = filepath.Join(tempDir, "errors.log")
	store.suggestionSink = newSuggestionSink(appConfig{
		SlackFeedbackWebhookURL: slackServer.URL,
	})

	body := strings.NewReader(`{"inputText":"cassava chips","allowed":false}`)
	request := httptest.NewRequest(http.MethodPost, "/suggest", body)
	response := httptest.NewRecorder()

	suggestHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(slackPayload.Text, "*AIP Food Lookup suggestion*") ||
		!strings.Contains(slackPayload.Text, "cassava chips") ||
		!strings.Contains(slackPayload.Text, "not allowed") {
		t.Fatalf("unexpected Slack payload: %#v", slackPayload)
	}

	errorLog, err := os.ReadFile(filepath.Join(tempDir, "errors.log"))
	if err != nil {
		t.Fatalf("expected local write failure in error log: %v", err)
	}
	if !strings.Contains(string(errorLog), "suggestion file write failed") {
		t.Fatalf("expected local write failure log, got %q", string(errorLog))
	}
}

func TestBuildSuggestionSlackMessageEscapesSlackSpecialCharacters(t *testing.T) {
	message := buildSuggestionSlackMessage(requestData{
		InputText: "chips <maybe> & dip",
		Allowed:   true,
	})

	if !strings.Contains(message, "*AIP Food Lookup suggestion*") {
		t.Fatalf("expected Slack heading, got %q", message)
	}
	if !strings.Contains(message, "chips &lt;maybe&gt; &amp; dip") {
		t.Fatalf("expected escaped Slack value, got %q", message)
	}
}
