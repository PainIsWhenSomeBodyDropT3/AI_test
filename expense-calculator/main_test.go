package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Test the server handles static file requests properly
func TestStaticFileServer(t *testing.T) {
	// Create a temporary directory to hold static files for testing
	tempDir, err := os.MkdirTemp("", "static-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test HTML file
	testHTML := "<html><body>Test File</body></html>"
	testFilePath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(testFilePath, []byte(testHTML), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a file server handler using the temp directory
	fs := http.FileServer(http.Dir(tempDir))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	// Create a test request to access the test file
	req := httptest.NewRequest("GET", "/test.html", nil)
	rec := httptest.NewRecorder()

	// Call the handler with the request and response recorder
	handler.ServeHTTP(rec, req)

	// Check the status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Check the content
	if rec.Body.String() != testHTML {
		t.Errorf("Expected body %q, got %q", testHTML, rec.Body.String())
	}
}

// Test the server returns 404 for non-existent files
func TestStaticFileServerNotFound(t *testing.T) {
	// Create a temporary directory to hold static files for testing
	tempDir, err := os.MkdirTemp("", "static-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file server handler using the temp directory
	fs := http.FileServer(http.Dir(tempDir))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	// Create a test request to access a non-existent file
	req := httptest.NewRequest("GET", "/non-existent.html", nil)
	rec := httptest.NewRecorder()

	// Call the handler with the request and response recorder
	handler.ServeHTTP(rec, req)

	// Check the status code
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
