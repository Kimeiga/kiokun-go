package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/kanjidic"
	"kiokun-go/processor"
)

func main() {
	outputDir := "output_kanji_test"

	fmt.Printf("Debugging Kanjidic processing with output to: %s\n", outputDir)

	// Ensure we're using an absolute path for the output directory
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
		os.Exit(1)
	}
	outputDir = absOutputDir

	// Set dictionaries base path
	common.SetDictionariesBasePath("dictionaries")

	// Create processor
	proc, err := processor.New(outputDir, runtime.NumCPU())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating processor: %v\n", err)
		os.Exit(1)
	}

	// Import only Kanjidic entries
	fmt.Printf("Importing Kanjidic entries directly...\n")

	// Find the kanjidic source file
	sourceDir := filepath.Join("dictionaries", "kanjidic", "source")
	pattern := `^kanjidic2-en-.*\.json$`
	filename, err := common.FindDictionaryFile(sourceDir, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding Kanjidic source file: %v\n", err)
		os.Exit(1)
	}

	// Construct full path
	fullPath := filepath.Join(sourceDir, filename)
	fmt.Printf("Found Kanjidic source file: %s\n", fullPath)

	// Create importer and import entries
	importer := &kanjidic.Importer{}
	kanjiEntries, err := importer.Import(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error importing Kanjidic: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Imported %d Kanjidic entries\n", len(kanjiEntries))

	// Print first 5 entries
	fmt.Printf("\nSample Kanjidic entries:\n")
	for i, entry := range kanjiEntries {
		if i >= 5 {
			break
		}
		kanji := entry.(kanjidic.Kanji)
		fmt.Printf("%d. Character: %s, Meanings: %v\n", i+1, kanji.Character, kanji.Meanings)
	}

	// Process only Kanjidic entries
	fmt.Printf("\nProcessing Kanjidic entries...\n")
	if err := proc.ProcessEntries(kanjiEntries); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing Kanjidic entries: %v\n", err)
		os.Exit(1)
	}

	// Write to files
	fmt.Printf("Writing files to %s...\n", outputDir)
	if err := proc.WriteToFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing files: %v\n", err)
		os.Exit(1)
	}

	// Check a few output files
	fmt.Printf("\nChecking output files for kanji characters...\n")
	// Try to find files for common kanji
	commonKanji := []string{"水", "火", "木", "金", "土"}
	for _, k := range commonKanji {
		filePath := filepath.Join(outputDir, k+".json.br")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("- File for kanji '%s' NOT found\n", k)
		} else {
			fmt.Printf("- File for kanji '%s' found\n", k)
		}
	}

	fmt.Printf("\nDebug completed. Check the %s directory for output files.\n", outputDir)
}
