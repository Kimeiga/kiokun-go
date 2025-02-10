package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler - handles requests to /api/word/{word}
func Handler(w http.ResponseWriter, r *http.Request) {
	// Extract word from path
	word := strings.TrimPrefix(r.URL.Path, "/api/word/")
	if word == "" {
		http.Error(w, "Word parameter is required", http.StatusBadRequest)
		return
	}

	// Construct path to gzipped file
	filename := filepath.Join("dictionary", word+".json.gz")

	// Open the gzipped file
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Word not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		http.Error(w, "Error reading compressed file", http.StatusInternalServerError)
		return
	}
	defer gzReader.Close()

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Copy the uncompressed content to response
	_, err = io.Copy(w, gzReader)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}
}
