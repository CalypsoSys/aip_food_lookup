package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGatewaySecretMiddlewareProtectsApiRoutes(t *testing.T) {
	config := appConfig{
		RequireGatewaySecret:    true,
		GatewaySecretHeaderName: "X-Internal-Api-Key",
		GatewaySecret:           "secret",
	}
	handler := gatewaySecretMiddleware(config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/search?key=apple", nil))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}

	request := httptest.NewRequest(http.MethodGet, "/search?key=apple", nil)
	request.Header.Set("X-Internal-Api-Key", "secret")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
}

func TestGatewaySecretMiddlewareProtectsAdminReloadWithoutPublicApiFlag(t *testing.T) {
	config := appConfig{
		GatewaySecretHeaderName: "X-Internal-Api-Key",
		GatewaySecret:           "secret",
	}
	handler := gatewaySecretMiddleware(config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, adminReloadPath, nil))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}

	request := httptest.NewRequest(http.MethodPost, adminReloadPath, nil)
	request.Header.Set("X-Internal-Api-Key", "secret")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
}

func TestGatewaySecretMiddlewareRejectsAdminReloadWithoutConfiguredSecret(t *testing.T) {
	handler := gatewaySecretMiddleware(appConfig{}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, adminReloadPath, nil))

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", response.Code)
	}
}

func TestGatewaySecretMiddlewareSkipsHealth(t *testing.T) {
	config := appConfig{
		RequireGatewaySecret: true,
		GatewaySecret:        "secret",
	}
	handler := gatewaySecretMiddleware(config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
}

func TestCorsMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	config := appConfig{
		AllowedOrigins: []string{"https://hashimojoe.com"},
	}
	handler := corsMiddleware(config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/search", nil)
	request.Header.Set("Origin", "https://hashimojoe.com")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "https://hashimojoe.com" {
		t.Fatalf("unexpected allow origin header: %q", got)
	}
}

func TestRateLimitMiddlewareLimitsByRouteGroup(t *testing.T) {
	apiRateLimiter = &fixedWindowRateLimiter{windows: make(map[string]rateWindow)}
	config := appConfig{
		RateLimit: rateLimitConfig{
			Enabled:             true,
			FeedbackPermitLimit: 1,
			WindowSeconds:       60,
		},
	}
	handler := rateLimitMiddleware(config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader("{}"))
	request.RemoteAddr = "198.51.100.10:12345"
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected first request status 204, got %d", response.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader("{}"))
	request.RemoteAddr = "198.51.100.10:12345"
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request status 429, got %d", response.Code)
	}
}
