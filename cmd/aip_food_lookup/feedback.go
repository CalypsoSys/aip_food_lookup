package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type slackFeedbackSink struct {
	webhookURL string
	client     *http.Client
	fallback   feedbackSink
	errorLog   string
}

func newFeedbackSink(config appConfig) feedbackSink {
	fallback := fileFeedbackSink{
		dataFolder: config.DataFolder,
		filePath:   config.FeedbackJSONLPath,
	}

	if strings.TrimSpace(config.SlackFeedbackWebhookURL) == "" {
		return fallback
	}

	return slackFeedbackSink{
		webhookURL: config.SlackFeedbackWebhookURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		fallback: fallback,
		errorLog: config.ErrorLogPath,
	}
}

func (s slackFeedbackSink) submitFeedback(request feedbackRequest) error {
	if err := s.submitSlack(request); err == nil {
		return nil
	} else {
		writeErrorLog(s.errorLog, fmt.Sprintf("slack feedback failed: %v", err))
		if fallbackErr := s.fallback.submitFeedback(request); fallbackErr != nil {
			return errors.Join(err, fallbackErr)
		}
		return nil
	}
}

func (s slackFeedbackSink) submitSlack(request feedbackRequest) error {
	payload, err := json.Marshal(map[string]any{
		"text":   buildFeedbackSlackMessage(request),
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

func buildFeedbackSlackMessage(request feedbackRequest) string {
	return fmt.Sprintf(
		"*AIP Food Lookup feedback*\n*Name:* %s\n*Email:* %s\n*Subject:* %s\n*Source:* %s\n*Message:*\n%s",
		escapeSlackValue(request.Name),
		escapeSlackValue(request.Email),
		escapeSlackValue(request.Subject),
		escapeSlackValue(request.Source),
		escapeSlackValue(request.Message),
	)
}

func escapeSlackValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "none"
	}

	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(value)
}
