package z

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
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

func TestDefaultLoggingMiddleware(t *testing.T) {
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	slog.SetDefault(logger)

	req := httptest.NewRequest("POST", "/default", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusCreated, "Default Logging!")
	}

	loggingMiddleware := Middlewares.Logging()
	wrappedHandler := loggingMiddleware(handler)
	wrappedHandler(z)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	if rr.Body.String() != "Default Logging!" {
		t.Errorf("Expected body %q, got %q", "Default Logging!", rr.Body.String())
	}
	logStr := logOutput.String()
	if !strings.Contains(logStr, `"method":"POST"`) {
		t.Errorf("Log output should contain method POST")
	}
	if !strings.Contains(logStr, `"path":"/default"`) {
		t.Errorf("Log output should contain path /default")
	}
	if !strings.Contains(logStr, `"status":201`) {
		t.Errorf("Log output should contain status 201")
	}
}

func TestLoggingMiddleware_RequestID(t *testing.T) {
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	slog.SetDefault(logger)

	req := httptest.NewRequest("GET", "/with-request-id", nil)
	req.Header.Set("X-Request-ID", "abc-123")
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "RequestID Test")
	}

	loggingMiddleware := Middlewares.LoggingWithCfg(LoggingConfig{LogResponseBody: true})
	wrappedHandler := loggingMiddleware(handler)
	wrappedHandler(z)

	logStr := logOutput.String()
	if !strings.Contains(logStr, `"request_id":"abc-123"`) {
		t.Errorf("Log output should contain request_id abc-123")
	}
}

func TestLoggingMiddleware_RequestBodyLogged(t *testing.T) {
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	slog.SetDefault(logger)

	bodyContent := "foobar-body"
	req := httptest.NewRequest("POST", "/with-body", strings.NewReader(bodyContent))
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "Body Test")
	}

	loggingMiddleware := Middlewares.LoggingWithCfg(LoggingConfig{LogRequestBody: true})
	wrappedHandler := loggingMiddleware(handler)
	wrappedHandler(z)

	logStr := logOutput.String()
	if !strings.Contains(logStr, `"request_body":"foobar-body"`) {
		t.Errorf("Log output should contain request_body foobar-body")
	}
}

func TestLoggingMiddleware_RequestBodyError(t *testing.T) {
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	slog.SetDefault(logger)

	brokenBody := io.NopCloser(brokenReader{})
	req := httptest.NewRequest("POST", "/error", brokenBody)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "Error Test")
	}

	loggingMiddleware := Middlewares.LoggingWithCfg(LoggingConfig{LogRequestBody: true})
	wrappedHandler := loggingMiddleware(handler)
	wrappedHandler(z)

	logStr := logOutput.String()
	if !strings.Contains(logStr, "Error reading request body") {
		t.Errorf("Expected error log for reading request body")
	}
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestLoggingMiddleware_LogFileOutput(t *testing.T) {
	logFile := "test_log_output.json"
	defer os.Remove(logFile)

	req := httptest.NewRequest("GET", "/logfile", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "LogFile Test")
	}

	loggingMiddleware := Middlewares.LoggingWithCfg(LoggingConfig{
		LogFilePath:     logFile,
		LogResponseBody: true,
	})
	wrappedHandler := loggingMiddleware(handler)
	wrappedHandler(z)

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	logStr := string(data)
	if !strings.Contains(logStr, "LogFile Test") {
		t.Errorf("Log file should contain response body 'LogFile Test', got: %s", logStr)
	}
	if !strings.Contains(logStr, `"method":"GET"`) {
		t.Errorf("Log file should contain method GET")
	}
	if !strings.Contains(logStr, `"path":"/logfile"`) {
		t.Errorf("Log file should contain path /logfile")
	}
	if !strings.Contains(logStr, `"status":200`) {
		t.Errorf("Log file should contain status 200")
	}
}

func TestLoggingMiddleware_LogFileOpenError(t *testing.T) {
	invalidPath := "/invalid_path/test_log_output.json"

	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	slog.SetDefault(logger)

	req := httptest.NewRequest("GET", "/logfileerror", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "Should not log to file")
	}

	loggingMiddleware := Middlewares.LoggingWithCfg(LoggingConfig{
		LogFilePath:     invalidPath,
		LogResponseBody: true,
	})
	wrappedHandler := loggingMiddleware(handler)
	wrappedHandler(z)

	logStr := logOutput.String()
	if !strings.Contains(logStr, "Failed to open log file") {
		t.Errorf("Expected error log for failed log file open, got: %s", logStr)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		panic("test panic")
	}

	recoveryMiddleware := Middlewares.Recovery()
	wrappedHandler := recoveryMiddleware(handler)
	wrappedHandler(z)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	if rr.Body.String() != "Internal Server Error" {
		t.Errorf("Expected body %q, got %q", "Internal Server Error", rr.Body.String())
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "Hello, World!")
	}

	requestIDMiddleware := Middlewares.RequestID()
	wrappedHandler := requestIDMiddleware(handler)
	wrappedHandler(z)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	if req.Header.Get("X-Request-ID") == "" {
		t.Error("Expected X-Request-ID header to be set in the request")
	}

	if rr.Header().Get("X-Request-ID") == "" {
		t.Error("Expected X-Request-ID header to be set in the response")
	}
}

func TestCORSMiddleware(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {}

	corsMiddleware := Middlewares.CORS()
	wrappedHandler := corsMiddleware(handler)
	wrappedHandler(z)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin header to be %q, got %q", "*", rr.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSMiddleware_NonOptions_CallsNext(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	called := false
	handler := func(z *Z) { called = true; z.String(http.StatusOK, "ok") }

	corsMiddleware := Middlewares.CORS()
	wrappedHandler := corsMiddleware(handler)
	wrappedHandler(z)

	if !called {
		t.Fatalf("expected next to be called for non-OPTIONS request")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("expected CORS headers to be set on non-OPTIONS request as well")
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {}

	securityHeadersMiddleware := Middlewares.SecurityHeaders()
	wrappedHandler := securityHeadersMiddleware(handler)
	wrappedHandler(z)

	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options header to be %q, got %q", "nosniff", rr.Header().Get("X-Content-Type-Options"))
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		time.Sleep(100 * time.Millisecond)
		z.String(http.StatusOK, "Hello, World!")
	}

	timeoutMiddleware := Middlewares.TimeoutWithCfg(TimeoutConfig{Timeout: 50 * time.Millisecond})
	wrappedHandler := timeoutMiddleware(handler)
	wrappedHandler(z)

	if rr.Code != http.StatusGatewayTimeout {
		t.Errorf("Expected status code %d, got %d", http.StatusGatewayTimeout, rr.Code)
	}
}

func TestTimeoutMiddleware_Default_NoTimeout(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	z := &Z{rw: rr, r: req}

	handler := func(z *Z) {
		z.String(http.StatusOK, "fast")
	}

	timeoutMiddleware := Middlewares.Timeout()
	wrappedHandler := timeoutMiddleware(handler)
	wrappedHandler(z)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
