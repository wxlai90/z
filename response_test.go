package z

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	z.String(200, "hello")
	if rw.Code != 200 {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}
	if rw.Body.String() != "hello" {
		t.Errorf("Expected body 'hello', got '%s'", rw.Body.String())
	}
}

func TestJSON(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	resp := map[string]string{"msg": "ok"}
	z.JSON(201, resp)
	if rw.Code != 201 {
		t.Errorf("Expected status 201, got %d", rw.Code)
	}
	if ct := rw.Header().Get("content-type"); ct != "application/json" {
		t.Errorf("Expected content-type application/json, got %s", ct)
	}
}

func TestOk(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	z.Ok("ok")
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if rw.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got '%s'", rw.Body.String())
	}
}

func TestOkJSON(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	resp := map[string]string{"msg": "ok"}
	z.OkJSON(resp)
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if ct := rw.Header().Get("content-type"); ct != "application/json" {
		t.Errorf("Expected content-type application/json, got %s", ct)
	}
}

func TestSetHeader(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	z.SetHeader("X-Test", "true")
	if val := rw.Header().Get("X-Test"); val != "true" {
		t.Errorf("Expected header X-Test to be 'true', got '%s'", val)
	}
}

func TestSetCookie(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	cookie := &http.Cookie{
		Name:    "test",
		Value:   "123",
		Expires: time.Now().Add(24 * time.Hour),
	}
	z.SetCookie(cookie)
	h := rw.Header().Get("Set-Cookie")
	if !strings.Contains(h, "test=123") {
		t.Errorf("Expected Set-Cookie header to contain 'test=123', got '%s'", h)
	}
}

func TestError(t *testing.T) {
	rw := httptest.NewRecorder()
	z := &Z{rw: rw}
	err := errors.New("test error")
	z.Error(err, http.StatusInternalServerError)
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "test error") {
		t.Errorf("Expected body to contain 'test error', got '%s'", rw.Body.String())
	}
}
