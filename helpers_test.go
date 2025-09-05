package z

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveUploadedFile_Success(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("testfile", "example.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.WriteString(fileWriter, "Hello, World!")
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		t.Fatalf("Failed to parse multipart form: %v", err)
	}

	z := &Z{
		rw: httptest.NewRecorder(),
		r:  req,
	}

	tmpDir := t.TempDir()
	dstPath := filepath.Join(tmpDir, "saved_example.txt")

	err = z.SaveUploadedFile("testfile", dstPath)
	if err != nil {
		t.Fatalf("SaveUploadedFile failed: %v", err)
	}

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	content, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	expected := "Hello, World!"
	if string(content) != expected {
		t.Errorf("File content mismatch: got %q, want %q", string(content), expected)
	}
}

func TestSaveUploadedFile_KeyNotFound(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)

	z := &Z{
		rw: httptest.NewRecorder(),
		r:  req,
	}

	err := z.SaveUploadedFile("nonexistent", "/tmp/test.txt")
	if err == nil {
		t.Fatal("Expected error for non-existent key, got nil")
	}

	if !strings.Contains(err.Error(), "failed to get form file") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSaveUploadedFile_DirectoryCreation(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("testfile", "example.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.WriteString(fileWriter, "Test content")
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		t.Fatalf("Failed to parse multipart form: %v", err)
	}

	z := &Z{
		rw: httptest.NewRecorder(),
		r:  req,
	}

	tmpDir := t.TempDir()
	dstPath := filepath.Join(tmpDir, "subdir", "anotherdir", "saved.txt")

	err = z.SaveUploadedFile("testfile", dstPath)
	if err != nil {
		t.Fatalf("SaveUploadedFile failed with directory creation: %v", err)
	}

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		t.Fatal("File was not created in nested directory")
	}
}

type fakeResponseWriter struct{}

func (f *fakeResponseWriter) Header() http.Header        { return http.Header{} }
func (f *fakeResponseWriter) Write([]byte) (int, error)  { return 0, nil }
func (f *fakeResponseWriter) WriteHeader(statusCode int) {}

func TestSaveUploadedFile_CreateDestFail(t *testing.T) {
	tmpDir := t.TempDir()
	dstPath := filepath.Join(tmpDir, "as_dir")
	if err := os.MkdirAll(dstPath, 0o755); err != nil {
		t.Fatalf("setup mkdir failed: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "a.txt")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	_, _ = io.WriteString(fileWriter, "x")
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		t.Fatalf("ParseMultipartForm: %v", err)
	}

	z := &Z{rw: &fakeResponseWriter{}, r: req}

	if err := z.SaveUploadedFile("file", dstPath); err == nil || !strings.Contains(err.Error(), "failed to create destination file") {
		t.Fatalf("expected create dest error, got %v", err)
	}
}

func TestSaveUploadedFile_CopyFail(t *testing.T) {
	old := copyFile
	copyFile = func(dst io.Writer, src io.Reader) (int64, error) { return 0, io.ErrClosedPipe }
	t.Cleanup(func() { copyFile = old })

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, _ := writer.CreateFormFile("file", "a.txt")
	_, _ = io.WriteString(fileWriter, "x")
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_ = req.ParseMultipartForm(32 << 20)

	z := &Z{rw: &fakeResponseWriter{}, r: req}
	dst := filepath.Join(t.TempDir(), "out.txt")
	if err := z.SaveUploadedFile("file", dst); err == nil || !strings.Contains(err.Error(), "failed to write file") {
		t.Fatalf("expected copy error, got %v", err)
	}
}

func TestSaveUploadedFile_MkdirFail(t *testing.T) {
	oldMk := mkdirAll
	mkdirAll = func(path string, perm os.FileMode) error { return os.ErrPermission }
	t.Cleanup(func() { mkdirAll = oldMk })

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, _ := writer.CreateFormFile("file", "a.txt")
	_, _ = io.WriteString(fileWriter, "x")
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_ = req.ParseMultipartForm(32 << 20)

	z := &Z{rw: &fakeResponseWriter{}, r: req}
	dst := filepath.Join(t.TempDir(), "sub", "out.txt")
	if err := z.SaveUploadedFile("file", dst); err == nil || !strings.Contains(err.Error(), "failed to create directory") {
		t.Fatalf("expected mkdir error, got %v", err)
	}
}
