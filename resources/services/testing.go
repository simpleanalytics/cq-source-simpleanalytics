package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func testServer(t *testing.T, filename string) *httptest.Server {
	contents, err := os.ReadFile(filepath.Join("..", "..", "internal", "simpleanalytics", "testdata", filename))
	if err != nil {
		t.Fatalf("unexpected error reading testdata file: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Write(contents)
	}))
	return ts
}
