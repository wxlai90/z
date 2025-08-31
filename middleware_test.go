package z

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareFunc(t *testing.T) {
	called := false
	handler := func(z *Z) { called = true }
	mw := func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			called = true
			next(z)
		}
	}
	wrapped := mw(handler)
	wrapped(nil)
	if !called {
		t.Error("MiddlewareFunc did not call handler")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	slog.SetDefault(logger)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "Hello, World!")
	}

	loggingMiddleware := Middlewares.LoggingWithCfg(LoggingConfig{LogRequestBody: true, LogResponseBody: true})

	wrappedHandler := loggingMiddleware(handler)

	wrappedHandler(z)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "Hello, World!" {
		t.Errorf("Expected body %q, got %q", "Hello, World!", rr.Body.String())
	}

	logStr := logOutput.String()

	if !strings.Contains(logStr, `"method":"GET"`) {
		t.Errorf("Log output should contain method GET")
	}
	if !strings.Contains(logStr, `"path":"/"`) {
		t.Errorf("Log output should contain path /")
	}
	if !strings.Contains(logStr, `"status":200`) {
		t.Errorf("Log output should contain status 200")
	}
	if !strings.Contains(logStr, `"response_body":"Hello, World!"`) {
		t.Errorf("Log output should contain response body")
	}
}
