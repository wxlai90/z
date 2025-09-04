package z

import (
	"bytes"
	"io"
	"mime/multipart"
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
