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
