package z

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
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

func TestHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Test", "true")
	z := &Z{r: req}
	val := z.Header("X-Test")
	if val != "true" {
		t.Errorf("Expected header 'true', got '%s'", val)
	}
}

func TestCookie(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test", Value: "123"})
	z := &Z{r: req}
	cookie, err := z.Cookie("test")
	if err != nil {
		t.Fatalf("Cookie failed: %v", err)
	}
	if cookie.Value != "123" {
		t.Errorf("Expected cookie value '123', got '%s'", cookie.Value)
	}
}

func TestFormFile(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.txt"`)
	h.Set("Content-Type", "text/plain")
	part, _ := mw.CreatePart(h)
	part.Write([]byte("this is a test"))
	mw.Close()

	req, _ := http.NewRequest("POST", "/", buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	z := &Z{r: req}
	file, header, err := z.FormFile("file")
	if err != nil {
		t.Fatalf("FormFile failed: %v", err)
	}
	defer file.Close()

	if header.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", header.Filename)
	}

	content, _ := io.ReadAll(file)
	if string(content) != "this is a test" {
		t.Errorf("Expected file content 'this is a test', got '%s'", string(content))
	}
}
