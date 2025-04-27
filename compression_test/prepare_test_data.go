package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"

	"github.com/andybalholm/brotli"
)

func main() {
	// Create test_data directory in web-test/src
	testDataDir := filepath.Join("web-test", "src", "test_data")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		fmt.Printf("Error creating test data directory: %v\n", err)
		return
	}

	// Import all dictionaries
	entries, err := common.ImportAllDictionaries()
	if err != nil {
		fmt.Printf("Error importing dictionaries: %v\n", err)
		return
	}

	// Group entries by dictionary
	dictEntries := make(map[string][]common.Entry)
	for _, entry := range entries {
		dictName := ""
		switch entry.(type) {
		case *jmdict.Word:
			dictName = "jmdict"
		case *jmnedict.Name:
			dictName = "jmnedict"
		case *kanjidic.Kanji:
			dictName = "kanjidic"
		}
		if dictName != "" {
			dictEntries[dictName] = append(dictEntries[dictName], entry)
		}
	}

	// Process each dictionary
	for dictName, entries := range dictEntries {
		// Convert entries to JSON
		data, err := json.Marshal(entries)
		if err != nil {
			fmt.Printf("Error marshaling %s: %v\n", dictName, err)
			continue
		}

		// Save raw JSON
		rawPath := filepath.Join(testDataDir, fmt.Sprintf("%s.json", dictName))
		if err := os.WriteFile(rawPath, data, 0644); err != nil {
			fmt.Printf("Error writing raw JSON for %s: %v\n", dictName, err)
			continue
		}

		// Compress with gzip
		var gzipBuf bytes.Buffer
		gzipWriter := gzip.NewWriter(&gzipBuf)
		if _, err := gzipWriter.Write(data); err != nil {
			fmt.Printf("Error compressing %s with gzip: %v\n", dictName, err)
			continue
		}
		gzipWriter.Close()

		gzipPath := filepath.Join(testDataDir, fmt.Sprintf("%s.json.gz", dictName))
		if err := os.WriteFile(gzipPath, gzipBuf.Bytes(), 0644); err != nil {
			fmt.Printf("Error writing gzip for %s: %v\n", dictName, err)
			continue
		}

		// Compress with brotli
		var brotliBuf bytes.Buffer
		brotliWriter := brotli.NewWriter(&brotliBuf)
		if _, err := brotliWriter.Write(data); err != nil {
			fmt.Printf("Error compressing %s with brotli: %v\n", dictName, err)
			continue
		}
		brotliWriter.Close()

		brotliPath := filepath.Join(testDataDir, fmt.Sprintf("%s.json.br", dictName))
		if err := os.WriteFile(brotliPath, brotliBuf.Bytes(), 0644); err != nil {
			fmt.Printf("Error writing brotli for %s: %v\n", dictName, err)
			continue
		}

		fmt.Printf("Successfully processed %s\n", dictName)
		fmt.Printf("  Raw size: %d bytes\n", len(data))
		fmt.Printf("  Gzip size: %d bytes\n", gzipBuf.Len())
		fmt.Printf("  Brotli size: %d bytes\n", brotliBuf.Len())
	}
}
