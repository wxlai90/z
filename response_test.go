package z

import (
	"net/http/httptest"
	"testing"
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
