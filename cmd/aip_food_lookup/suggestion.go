package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type slackSuggestionSink struct {
	webhookURL string
	client     *http.Client
}

func newSuggestionSink(config appConfig) suggestionSink {
	if strings.TrimSpace(config.SlackFeedbackWebhookURL) == "" {
		return nil
	}

	return slackSuggestionSink{
		webhookURL: config.SlackFeedbackWebhookURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s slackSuggestionSink) submitSuggestion(request requestData) error {
	payload, err := json.Marshal(map[string]any{
		"text":   buildSuggestionSlackMessage(request),
		"mrkdwn": true,
	})
	if err != nil {
		return err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, s.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("slack returned status %d", response.StatusCode)
	}
	return nil
}

func buildSuggestionSlackMessage(request requestData) string {
	status := "not allowed"
	if request.Allowed {
		status = "allowed"
	}

	return fmt.Sprintf(
		"*AIP Food Lookup suggestion*\n*Food:* %s\n*Suggested as:* %s",
		escapeSlackValue(request.InputText),
		escapeSlackValue(status),
	)
}
