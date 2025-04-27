package analyzer

import (
	"encoding/json"
	"flag"
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

// WordGroup represents the combined data for a single word/character
// This must match the structure in processor/processor.go
type WordGroup struct {
	WordJapanese []jmdict.Word    `json:"w_j,omitempty"`
	NameJapanese []jmnedict.Name  `json:"n_j,omitempty"`
	CharJapanese []kanjidic.Kanji `json:"c_j,omitempty"`
}

// HasMultipleDictData returns true if this group has data from multiple dictionaries
func (wg *WordGroup) HasMultipleDictData() bool {
	sources := 0
	if len(wg.WordJapanese) > 0 {
		sources++
	}
	if len(wg.NameJapanese) > 0 {
		sources++
	}
	if len(wg.CharJapanese) > 0 {
		sources++
	}
	return sources > 1
}

// RunAnalysis is the main entry point for analyzing dictionary outputs
func RunAnalysis() {
	// Parse command line arguments
	outputDir := flag.String("dir", "output", "Directory containing output files")
	scanAll := flag.Bool("scanall", false, "Scan all files instead of stopping after finding examples")
	flag.Parse()

	fmt.Printf("Analyzing dictionary files in %s, looking for combined entries...\n", *outputDir)
	startTime := time.Now()

	// Statistics counters
	var totalFiles, errorFiles int
	var filesWithJMdict, filesWithJMNedict, filesWithKanjidic int
	var multiDictFiles, wordKanjiFiles, wordNameFiles, kanjiNameFiles, wordKanjiNameFiles int

	// Examples to print
	var wordExample *WordGroup
	var wordKanjiExample *WordGroup
	var wordKanjiNameExample *WordGroup
	var kanjiExample *WordGroup
	var nameExample *WordGroup
	var wordExampleFilename string
	var wordKanjiFilename string
	var wordKanjiNameFilename string
	var kanjiExampleFilename string
	var nameExampleFilename string

	// Check if we have found all example types we're looking for
	examplesComplete := func() bool {
		if *scanAll {
			return false // Continue scanning all files if requested
		}
		return wordExample != nil &&
			kanjiExample != nil &&
			nameExample != nil &&
			wordKanjiExample != nil &&
			wordKanjiNameExample != nil
	}

	// Walk directory and process each file
	err := filepath.WalkDir(*outputDir, func(path string, d fs.DirEntry, err error) error {
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

		// Progress update every 5000 files
		if totalFiles%5000 == 0 {
			fmt.Printf("\rProcessed %d files, found %d combined entries...", totalFiles, multiDictFiles)
		}

		// Stop if we've found all example types
		if examplesComplete() {
			fmt.Printf("\nFound all example types, stopping scan...\n")
			return filepath.SkipAll
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
			if wordExample == nil {
				wordExample = &group
				wordExampleFilename = form
				fmt.Printf("\nFound example word entry: %s\n", form)
			}
		}

		if hasName {
			filesWithJMNedict++
			if nameExample == nil {
				nameExample = &group
				nameExampleFilename = form
				fmt.Printf("\nFound example name entry: %s\n", form)
			}
		}

		if hasKanji {
			filesWithKanjidic++
			if kanjiExample == nil {
				kanjiExample = &group
				kanjiExampleFilename = form
				fmt.Printf("\nFound example kanji entry: %s\n", form)
			}
		}

		// Count combinations
		if group.HasMultipleDictData() {
			multiDictFiles++

			if hasWord && hasKanji && !hasName {
				wordKanjiFiles++
				if wordKanjiExample == nil {
					wordKanjiExample = &group
					wordKanjiFilename = form
					fmt.Printf("\nFound example word+kanji entry: %s\n", form)
				}
			}

			if hasWord && hasName && !hasKanji {
				wordNameFiles++
			}

			if hasKanji && hasName && !hasWord {
				kanjiNameFiles++
			}

			if hasWord && hasKanji && hasName {
				wordKanjiNameFiles++
				if wordKanjiNameExample == nil {
					wordKanjiNameExample = &group
					wordKanjiNameFilename = form
					fmt.Printf("\nFound example word+kanji+name entry: %s\n", form)
				}
			}
		}

		return nil
	})

	fmt.Printf("\r") // Clear the progress line

	if err != nil && err != filepath.SkipAll {
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

	// Print example entries
	if wordExample != nil {
		fmt.Printf("\n=== Example of word entry: %s ===\n", wordExampleFilename)
		prettyJSON, _ := json.MarshalIndent(wordExample, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("\nNo example word entries found")
	}

	if kanjiExample != nil {
		fmt.Printf("\n=== Example of kanji entry: %s ===\n", kanjiExampleFilename)
		prettyJSON, _ := json.MarshalIndent(kanjiExample, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("\nNo example kanji entries found")
	}

	if nameExample != nil {
		fmt.Printf("\n=== Example of name entry: %s ===\n", nameExampleFilename)
		prettyJSON, _ := json.MarshalIndent(nameExample, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("\nNo example name entries found")
	}

	if wordKanjiExample != nil {
		fmt.Printf("\n=== Example of entry with word+kanji data: %s ===\n", wordKanjiFilename)
		prettyJSON, _ := json.MarshalIndent(wordKanjiExample, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("\nNo entries found with word+kanji data")
	}

	if wordKanjiNameExample != nil {
		fmt.Printf("\n=== Example of entry with word+kanji+name data: %s ===\n", wordKanjiNameFilename)
		prettyJSON, _ := json.MarshalIndent(wordKanjiNameExample, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("\nNo entries found with word+kanji+name data")
	}

	fmt.Printf("\nAnalysis completed in %v\n", elapsedTime)
}
