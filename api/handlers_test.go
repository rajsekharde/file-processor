package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandleUpload(t *testing.T) {
	// Cleanup after test
	defer os.Remove("../file-storage/uploads/testfile.txt")

	// Create a pipe to simulate a multipart file upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("newFile", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("hello file content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	HandleUpload(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %d", res.StatusCode)
	}

	// Verify that the file was actually written to the local storage path
	savedFilePath := "../file-storage/uploads/testfile.txt"
	if _, err := os.Stat(savedFilePath); os.IsNotExist(err) {
		t.Errorf("Expected file to be saved to %s, but it doesn't exist", savedFilePath)
	}
}

func TestHandleDownload(t *testing.T) {
	defer os.Remove("../file-storage/completed/downloadme.txt")

	// Create a dummy file in the completed folder for testing download
	fileName := "downloadme.txt"
	filePath := filepath.Join("../file-storage/completed", fileName)
	err := os.WriteFile(filePath, []byte("download content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	// Setup request using Go 1.22+ PathValue routing emulation
	req := httptest.NewRequest(http.MethodGet, "/download/"+fileName, nil)
	
	// If you are using Go 1.22+ standard ServeMux, use a router to bind path values, 
	// or mock it using path values mechanism if needed. Alternatively, test via an httptest Server:
	mux := http.NewServeMux()
	mux.HandleFunc("GET /download/{filename}", HandleDownload)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %d", res.StatusCode)
	}

	// Check Content-Disposition header is present
	disposition := res.Header.Get("Content-Disposition")
	expectedDisposition := `attachment; filename="downloadme.txt"`
	if disposition != expectedDisposition {
		t.Errorf("Expected Content-Disposition %q, got %q", expectedDisposition, disposition)
	}

	// Check body content
	data, _ := io.ReadAll(res.Body)
	if string(data) != "download content" {
		t.Errorf("Expected file body 'download content', got %q", string(data))
	}
}

func TestHandlePostTask(t *testing.T) {
	// Spin up a fake local worker test server
	workerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "output_file": "out.txt"}`))
	}))
	defer workerServer.Close()

	// Override the global WorkerURL with the test server's dynamic URL
	originalURL := WorkerURL
	WorkerURL = workerServer.URL
	defer func() { WorkerURL = originalURL }() // Restore original URL after test

	// Perform the test request
	jsonBody := []byte(`{"FileName": "input.txt", "OutputFormat": "pdf"}`)
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	HandlePostTask(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", res.StatusCode)
	}
}