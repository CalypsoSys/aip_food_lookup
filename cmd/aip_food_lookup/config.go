package main

import (
	"os"
	"strconv"
	"strings"
)

type rateLimitConfig struct {
	Enabled             bool
	SearchPermitLimit   int
	WritePermitLimit    int
	FeedbackPermitLimit int
	WindowSeconds       int
}

type appConfig struct {
	ListenAddress           string
	DataFolder              string
	AccessLogPath           string
	ErrorLogPath            string
	AllowedOrigins          []string
	RequireGatewaySecret    bool
	GatewaySecretHeaderName string
	GatewaySecret           string
	SlackFeedbackWebhookURL string
	FeedbackJSONLPath       string
	RequestBodyLimitBytes   int64
	RateLimit               rateLimitConfig
}

func loadConfig() appConfig {
	return appConfig{
		ListenAddress:           envString(":8080", "AIP__API__ListenAddress", "AIP_LISTEN_ADDRESS"),
		DataFolder:              envString("data", "AIP__API__DataFolder", "AIP_DATA_FOLDER"),
		AccessLogPath:           envString("logs/access.log", "AIP__API__AccessLogPath", "AIP_ACCESS_LOG_PATH"),
		ErrorLogPath:            envString("logs/errors.log", "AIP__API__ErrorLogPath", "AIP_ERROR_LOG_PATH"),
		AllowedOrigins:          envList("AIP__API__AllowedOrigins", "AIP_ALLOWED_ORIGINS"),
		RequireGatewaySecret:    envBool(false, "AIP__API__RequireGatewaySecret", "AIP_REQUIRE_GATEWAY_SECRET"),
		GatewaySecretHeaderName: envString("X-Internal-Api-Key", "AIP__API__GatewaySecretHeaderName", "AIP_GATEWAY_SECRET_HEADER_NAME"),
		GatewaySecret:           envString("", "AIP__API__GatewaySecret", "AIP_GATEWAY_SECRET"),
		SlackFeedbackWebhookURL: envString("", "AIP__API__SlackFeedbackWebhookUrl", "AIP_SLACK_FEEDBACK_WEBHOOK_URL"),
		FeedbackJSONLPath:       envString("", "AIP__API__FeedbackJSONLPath", "AIP_FEEDBACK_JSONL_PATH"),
		RequestBodyLimitBytes:   int64(envInt(32768, "AIP__API__RequestBodyLimitBytes", "AIP_REQUEST_BODY_LIMIT_BYTES")),
		RateLimit: rateLimitConfig{
			Enabled:             envBool(false, "AIP__API__RateLimit__Enabled", "AIP_RATE_LIMIT_ENABLED"),
			SearchPermitLimit:   envInt(300, "AIP__API__RateLimit__SearchPermitLimit", "AIP_RATE_LIMIT_SEARCH_PERMIT_LIMIT"),
			WritePermitLimit:    envInt(60, "AIP__API__RateLimit__WritePermitLimit", "AIP_RATE_LIMIT_WRITE_PERMIT_LIMIT"),
			FeedbackPermitLimit: envInt(10, "AIP__API__RateLimit__FeedbackPermitLimit", "AIP_RATE_LIMIT_FEEDBACK_PERMIT_LIMIT"),
			WindowSeconds:       envInt(60, "AIP__API__RateLimit__WindowSeconds", "AIP_RATE_LIMIT_WINDOW_SECONDS"),
		},
	}
}

func envString(defaultValue string, names ...string) string {
	for _, name := range names {
		value := strings.TrimSpace(os.Getenv(name))
		if value != "" {
			return value
		}
	}
	return defaultValue
}

func envBool(defaultValue bool, names ...string) bool {
	value := envString("", names...)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func envInt(defaultValue int, names ...string) int {
	value := envString("", names...)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func envList(indexedPrefix string, commaName string) []string {
	values := []string{}
	for i := 0; ; i++ {
		value := strings.TrimSpace(os.Getenv(indexedPrefix + "__" + strconv.Itoa(i)))
		if value == "" {
			break
		}
		values = append(values, value)
	}

	if len(values) > 0 {
		return values
	}

	commaValue := strings.TrimSpace(os.Getenv(commaName))
	if commaValue == "" {
		return values
	}

	for _, value := range strings.Split(commaValue, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}
