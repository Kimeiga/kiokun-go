package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"

	"github.com/andybalholm/brotli"
)

// This matches the structure from processor.go
type WordGroup struct {
	WordJapanese []jmdict.Word    `json:"w_j,omitempty"`
	NameJapanese []jmnedict.Name  `json:"n_j,omitempty"`
	CharJapanese []kanjidic.Kanji `json:"c_j,omitempty"`
}

func main() {
	outputDir := "output_test" // Hardcoded to look at the debug output

	fmt.Printf("Analyzing debug dictionary files in %s...\n", outputDir)
	startTime := time.Now()

	// Statistics counters
	var totalFiles, errorFiles int
	var filesWithJMdict, filesWithJMNedict, filesWithKanjidic int
	var multiDictFiles, wordKanjiFiles, wordNameFiles, kanjiNameFiles, wordKanjiNameFiles int

	// List of all combined entries for verification
	var combinedEntries []string
	var wordKanjiEntries []string
	var wordNameEntries []string
	var kanjiNameEntries []string
	var tripleEntries []string

	// Walk directory and process each file
	err := filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Continue with next file
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process .json.br files
		if filepath.Ext(path) != ".br" {
			return nil
		}

		totalFiles++

		// Progress update every 100 files (smaller number for test directory)
		if totalFiles%100 == 0 {
			fmt.Printf("\rProcessed %d files, found %d combined entries...", totalFiles, multiDictFiles)
		}

		// Extract the form/term from the filename
		form := filepath.Base(path)
		form = form[:len(form)-8] // Remove .json.br suffix

		// Open and decompress the file
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("\nWarning: Error opening %s: %v\n", path, err)
			errorFiles++
			return nil // Continue with next file
		}
		defer file.Close()

		brReader := brotli.NewReader(file)

		// Decode the JSON data
		var group WordGroup
		decoder := json.NewDecoder(brReader)
		if err := decoder.Decode(&group); err != nil {
			fmt.Printf("\nWarning: Error decoding JSON for %s: %v\n", path, err)
			errorFiles++
			return nil // Continue with next file
		}

		// Update counters based on content
		hasWord := len(group.WordJapanese) > 0
		hasKanji := len(group.CharJapanese) > 0
		hasName := len(group.NameJapanese) > 0

		if hasWord {
			filesWithJMdict++
		}

		if hasName {
			filesWithJMNedict++
		}

		if hasKanji {
			filesWithKanjidic++
		}

		// Count combinations and save entry form for reporting
		if (hasWord && hasKanji) || (hasWord && hasName) || (hasKanji && hasName) {
			multiDictFiles++
			combinedEntries = append(combinedEntries, form)

			if hasWord && hasKanji && !hasName {
				wordKanjiFiles++
				wordKanjiEntries = append(wordKanjiEntries, form)
			}

			if hasWord && hasName && !hasKanji {
				wordNameFiles++
				wordNameEntries = append(wordNameEntries, form)
			}

			if hasKanji && hasName && !hasWord {
				kanjiNameFiles++
				kanjiNameEntries = append(kanjiNameEntries, form)
			}

			if hasWord && hasKanji && hasName {
				wordKanjiNameFiles++
				tripleEntries = append(tripleEntries, form)
			}
		}

		return nil
	})

	fmt.Printf("\r") // Clear the progress line

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during file scan: %v\n", err)
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Print detailed statistics
	fmt.Printf("\nAnalysis results after processing %d files (with %d errors):\n", totalFiles, errorFiles)
	fmt.Printf("- Files with JMdict (word) entries: %d\n", filesWithJMdict)
	fmt.Printf("- Files with JMNedict (name) entries: %d\n", filesWithJMNedict)
	fmt.Printf("- Files with Kanjidic (kanji) entries: %d\n", filesWithKanjidic)
	fmt.Printf("\nCombined entries:\n")
	fmt.Printf("- Files with data from multiple dictionaries: %d\n", multiDictFiles)
	fmt.Printf("- Files with both word and kanji data: %d\n", wordKanjiFiles)
	fmt.Printf("- Files with both word and name data: %d\n", wordNameFiles)
	fmt.Printf("- Files with both kanji and name data: %d\n", kanjiNameFiles)
	fmt.Printf("- Files with word, kanji, and name data: %d\n", wordKanjiNameFiles)

	// Print all combined entries
	if len(combinedEntries) > 0 {
		fmt.Printf("\nAll combined entries (%d): %v\n", len(combinedEntries), combinedEntries)
	} else {
		fmt.Printf("\nNo combined entries found!\n")
	}

	// Print specific combination types if found
	if len(wordKanjiEntries) > 0 {
		fmt.Printf("\nEntries with word+kanji data (%d): %v\n", len(wordKanjiEntries), wordKanjiEntries)
	}

	if len(wordNameEntries) > 0 {
		fmt.Printf("\nEntries with word+name data (%d): %v\n", len(wordNameEntries), wordNameEntries)
	}

	if len(kanjiNameEntries) > 0 {
		fmt.Printf("\nEntries with kanji+name data (%d): %v\n", len(kanjiNameEntries), kanjiNameEntries)
	}

	if len(tripleEntries) > 0 {
		fmt.Printf("\nEntries with word+kanji+name data (%d): %v\n", len(tripleEntries), tripleEntries)
	}

	fmt.Printf("\nAnalysis completed in %v\n", elapsedTime)
}
