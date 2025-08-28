package z

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestBindBody(t *testing.T) {
	body := `{"name":"test"}`
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	z := &Z{r: req}
	type payload struct {
		Name string `json:"name"`
	}
	var p payload
	z.r.Body = io.NopCloser(bytes.NewBufferString(body))
	err := z.BindBody(&p)
	if err != nil {
		t.Fatalf("BindBody failed: %v", err)
	}
	if p.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", p.Name)
	}
}

func TestBindBodyNil(t *testing.T) {
	z := &Z{r: &http.Request{Body: nil}}
	err := z.BindBody(&struct{}{})
	if err == nil {
		t.Error("Expected error for nil body")
	}
}

func TestPathValue(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/123", nil)
	z := &Z{r: req}
	z.r.SetPathValue("id", "123")
	id := z.PathValue("id")
	if id != "123" {
		t.Errorf("Expected id '123', got '%s'", id)
	}
}

func TestQuery(t *testing.T) {
	req, _ := http.NewRequest("GET", "/search?q=test", nil)
	z := &Z{r: req}
	q := z.Query("q")
	if q != "test" {
		t.Errorf("Expected query 'test', got '%s'", q)
	}
}
