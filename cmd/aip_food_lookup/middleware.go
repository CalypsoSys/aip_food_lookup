package main

import (
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var logFileLock sync.Mutex

type statusCaptureWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (w *statusCaptureWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCaptureWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	count, err := w.ResponseWriter.Write(data)
	w.bytes += count
	return count, err
}

func buildHTTPHandler(config appConfig, next http.Handler) http.Handler {
	return recoverMiddleware(config,
		accessLogMiddleware(config,
			corsMiddleware(config,
				rateLimitMiddleware(config,
					bodyLimitMiddleware(config,
						gatewaySecretMiddleware(config, next))))))
}

func corsMiddleware(config appConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowedOrigin(config.AllowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-AIP-Client, X-AIP-App-Version")
		}

		if r.Method == http.MethodOptions {
			if origin != "" && !isAllowedOrigin(config.AllowedOrigins, origin) {
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(allowedOrigins []string, origin string) bool {
	if len(allowedOrigins) == 0 {
		return false
	}

	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || strings.EqualFold(allowedOrigin, origin) {
			return true
		}
	}
	return false
}

func gatewaySecretMiddleware(config appConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions ||
			!config.RequireGatewaySecret ||
			!requiresGatewaySecret(r.URL.Path) ||
			config.GatewaySecret == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerName := config.GatewaySecretHeaderName
		if strings.TrimSpace(headerName) == "" {
			headerName = "X-Internal-Api-Key"
		}

		providedSecret := r.Header.Get(headerName)
		if subtle.ConstantTimeCompare([]byte(providedSecret), []byte(config.GatewaySecret)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func requiresGatewaySecret(requestPath string) bool {
	switch requestPath {
	case "/search", "/suggest", "/feedback", "/categories", "/subcategory":
		return true
	default:
		return false
	}
}

func bodyLimitMiddleware(config appConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.RequestBodyLimitBytes > 0 && requestMayHaveBody(r.Method) {
			r.Body = http.MaxBytesReader(w, r.Body, config.RequestBodyLimitBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func requestMayHaveBody(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}

func recoverMiddleware(config appConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				writeErrorLog(config.ErrorLogPath, fmt.Sprintf("panic path=%s error=%v", sanitizeLogValue(r.URL.RequestURI()), recovered))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func accessLogMiddleware(config appConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		capture := &statusCaptureWriter{ResponseWriter: w}
		next.ServeHTTP(capture, r)

		statusCode := capture.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		line := fmt.Sprintf(
			"%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" %dms",
			sanitizeLogValue(remoteIP(r)),
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			sanitizeLogValue(r.Method),
			sanitizeLogValue(r.URL.RequestURI()),
			sanitizeLogValue(r.Proto),
			statusCode,
			capture.bytes,
			sanitizeLogValue(r.Referer()),
			sanitizeLogValue(r.UserAgent()),
			time.Since(start).Milliseconds(),
		)
		writeLogLine(config.AccessLogPath, line)
	})
}

func remoteIP(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); value != "" {
		parts := strings.Split(value, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}
	return "unknown"
}

func sanitizeLogValue(value string) string {
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func writeErrorLog(path string, message string) {
	writeLogLine(path, fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), sanitizeLogValue(message)))
}

func writeLogLine(path string, line string) {
	if strings.TrimSpace(path) == "" {
		return
	}

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return
	}

	logFileLock.Lock()
	defer logFileLock.Unlock()

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	_, _ = fmt.Fprintln(file, line)
}

type fixedWindowRateLimiter struct {
	mu      sync.Mutex
	windows map[string]rateWindow
}

type rateWindow struct {
	start time.Time
	count int
}

var apiRateLimiter = &fixedWindowRateLimiter{windows: make(map[string]rateWindow)}

func rateLimitMiddleware(config appConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.RateLimit.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		windowSeconds := config.RateLimit.WindowSeconds
		if windowSeconds <= 0 {
			windowSeconds = 60
		}

		group, limit := rateLimitGroup(config, r)
		if limit <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		key := group + ":" + remoteIP(r)
		if !apiRateLimiter.allow(key, limit, time.Duration(windowSeconds)*time.Second) {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", windowSeconds))
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rateLimitGroup(config appConfig, r *http.Request) (string, int) {
	switch r.URL.Path {
	case "/feedback":
		return "feedback", config.RateLimit.FeedbackPermitLimit
	case "/suggest":
		return "write", config.RateLimit.WritePermitLimit
	default:
		return "read", config.RateLimit.SearchPermitLimit
	}
}

func (l *fixedWindowRateLimiter) allow(key string, limit int, window time.Duration) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	current := l.windows[key]
	if current.start.IsZero() || now.Sub(current.start) >= window {
		l.windows[key] = rateWindow{start: now, count: 1}
		return true
	}

	if current.count >= limit {
		return false
	}

	current.count++
	l.windows[key] = current
	return true
}
