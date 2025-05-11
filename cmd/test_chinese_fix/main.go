package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/processor"
)

func main() {
	// Create a temporary output directory
	tempDir := filepath.Join(os.TempDir(), "kiokun-test-chinese-fix")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		os.Exit(1)
	}

	// Import a small number of Chinese entries
	fmt.Println("Importing Chinese entries...")
	charImporter := &chinese_chars.Importer{}
	wordImporter := &chinese_words.Importer{}

	charPath := filepath.Join("dictionaries", "chinese_chars", "source", "dictionary_char_2024-06-17.json")
	wordPath := filepath.Join("dictionaries", "chinese_words", "source", "dictionary_word_2024-06-17.json")

	// Import with a limit of 100 entries
	charEntries, err := importWithLimit(charImporter, charPath, 100)
	if err != nil {
		fmt.Printf("Error importing Chinese characters: %v\n", err)
		os.Exit(1)
	}

	wordEntries, err := importWithLimit(wordImporter, wordPath, 100)
	if err != nil {
		fmt.Printf("Error importing Chinese words: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Imported %d Chinese characters and %d Chinese words\n", len(charEntries), len(wordEntries))

	// Create an index processor
	proc, err := processor.NewIndexProcessor(tempDir, 4)
	if err != nil {
		fmt.Printf("Error creating index processor: %v\n", err)
		os.Exit(1)
	}

	// Process the entries
	fmt.Println("Processing Chinese character entries...")
	if err := proc.ProcessEntries(charEntries); err != nil {
		fmt.Printf("Error processing character entries: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing Chinese word entries...")
	if err := proc.ProcessEntries(wordEntries); err != nil {
		fmt.Printf("Error processing word entries: %v\n", err)
		os.Exit(1)
	}

	// Write to files
	fmt.Println("Writing to files...")
	if err := proc.WriteToFiles(); err != nil {
		fmt.Printf("Error writing to files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Test completed successfully!")
}

// importWithLimit imports entries with a limit
func importWithLimit(importer common.DictionaryImporter, path string, limit int) ([]common.Entry, error) {
	allEntries, err := importer.Import(path)
	if err != nil {
		return nil, err
	}

	if limit > 0 && limit < len(allEntries) {
		return allEntries[:limit], nil
	}
	return allEntries, nil
}
