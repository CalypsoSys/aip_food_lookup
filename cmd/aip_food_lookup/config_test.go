package main

import "testing"

func TestLoadConfigReadsNestedYamlEnvironment(t *testing.T) {
	t.Setenv("AIP__API__ListenAddress", ":9090")
	t.Setenv("AIP__API__DataFolder", "/app/data")
	t.Setenv("AIP__API__AllowedOrigins__0", "https://hashimojoe.com")
	t.Setenv("AIP__API__AllowedOrigins__1", "https://www.hashimojoe.com")
	t.Setenv("AIP__API__RequireGatewaySecret", "true")
	t.Setenv("AIP__API__GatewaySecretHeaderName", "X-Test-Gateway")
	t.Setenv("AIP__API__GatewaySecret", "test-secret")
	t.Setenv("AIP__API__SlackFeedbackWebhookUrl", "https://hooks.slack.com/services/test")
	t.Setenv("AIP__API__RateLimit__Enabled", "true")
	t.Setenv("AIP__API__RateLimit__SearchPermitLimit", "44")
	t.Setenv("AIP__API__RateLimit__WritePermitLimit", "12")
	t.Setenv("AIP__API__RateLimit__FeedbackPermitLimit", "3")
	t.Setenv("AIP__API__RateLimit__WindowSeconds", "30")

	config := loadConfig()

	if config.ListenAddress != ":9090" {
		t.Fatalf("expected listen address :9090, got %q", config.ListenAddress)
	}
	if config.DataFolder != "/app/data" {
		t.Fatalf("expected data folder /app/data, got %q", config.DataFolder)
	}
	if len(config.AllowedOrigins) != 2 || config.AllowedOrigins[0] != "https://hashimojoe.com" {
		t.Fatalf("unexpected allowed origins: %#v", config.AllowedOrigins)
	}
	if !config.RequireGatewaySecret {
		t.Fatal("expected gateway secret to be required")
	}
	if config.GatewaySecretHeaderName != "X-Test-Gateway" || config.GatewaySecret != "test-secret" {
		t.Fatalf("unexpected gateway settings: %#v", config)
	}
	if config.SlackFeedbackWebhookURL != "https://hooks.slack.com/services/test" {
		t.Fatalf("unexpected Slack webhook: %q", config.SlackFeedbackWebhookURL)
	}
	if !config.RateLimit.Enabled ||
		config.RateLimit.SearchPermitLimit != 44 ||
		config.RateLimit.WritePermitLimit != 12 ||
		config.RateLimit.FeedbackPermitLimit != 3 ||
		config.RateLimit.WindowSeconds != 30 {
		t.Fatalf("unexpected rate limit config: %#v", config.RateLimit)
	}
}
