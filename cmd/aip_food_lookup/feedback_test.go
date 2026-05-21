package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSlackFeedbackSinkFallsBackToJsonlWhenSlackFails(t *testing.T) {
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer slackServer.Close()

	tempDir := t.TempDir()
	sink := newFeedbackSink(appConfig{
		DataFolder:              tempDir,
		ErrorLogPath:            filepath.Join(tempDir, "errors.log"),
		SlackFeedbackWebhookURL: slackServer.URL,
	})

	err := sink.submitFeedback(feedbackRequest{
		Name:    "Joe",
		Email:   "joe@example.com",
		Subject: "Hello",
		Message: "Works?",
		Source:  "android",
	})
	if err != nil {
		t.Fatalf("expected JSONL fallback to hide Slack error, got %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "feedback.jsonl"))
	if err != nil {
		t.Fatalf("expected fallback feedback file: %v", err)
	}
	if !strings.Contains(string(content), `"message":"Works?"`) {
		t.Fatalf("expected feedback message in fallback file, got %q", string(content))
	}

	errorLog, err := os.ReadFile(filepath.Join(tempDir, "errors.log"))
	if err != nil {
		t.Fatalf("expected Slack failure in error log: %v", err)
	}
	if !strings.Contains(string(errorLog), "slack feedback failed") {
		t.Fatalf("expected Slack failure log, got %q", string(errorLog))
	}
}

func TestBuildFeedbackSlackMessageEscapesSlackSpecialCharacters(t *testing.T) {
	message := buildFeedbackSlackMessage(feedbackRequest{
		Name:    "Joe <Admin>",
		Email:   "joe@example.com",
		Subject: "Hello & help",
		Message: "Please inspect > ingredients",
		Source:  "android",
	})

	if !strings.Contains(message, "*AIP Food Lookup feedback*") {
		t.Fatalf("expected Slack heading, got %q", message)
	}
	if !strings.Contains(message, "Joe &lt;Admin&gt;") ||
		!strings.Contains(message, "Hello &amp; help") ||
		!strings.Contains(message, "Please inspect &gt; ingredients") {
		t.Fatalf("expected escaped Slack values, got %q", message)
	}
}
