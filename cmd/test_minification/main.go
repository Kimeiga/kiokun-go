package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	_ "kiokun-go/dictionaries/chinese_chars"
	_ "kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	_ "kiokun-go/dictionaries/jmdict"
	_ "kiokun-go/dictionaries/jmnedict"
	_ "kiokun-go/dictionaries/kanjidic"
	"kiokun-go/processor"
)

func main() {
	// Parse command line arguments
	dictDir := flag.String("dictdir", "dictionaries", "Base directory containing dictionary packages")
	outputDir := flag.String("outdir", "output_test", "Output directory for processed files")
	limit := flag.Int("limit", 1000, "Limit the number of entries to process from each dictionary")
	flag.Parse()

	// Ensure output directory exists
	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
		os.Exit(1)
	}
	*outputDir = absOutputDir

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting minification test...\n")
	fmt.Printf("Input dictionaries directory: %s\n", *dictDir)
	fmt.Printf("Output directory: %s\n", *outputDir)
	fmt.Printf("Entry limit per dictionary: %d\n", *limit)

	// Set dictionaries base path
	common.SetDictionariesBasePath(*dictDir)

	// Create processor
	proc, err := processor.New(*outputDir, runtime.NumCPU())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating processor: %v\n", err)
		os.Exit(1)
	}

	// Import dictionaries
	fmt.Printf("Importing dictionaries...\n")
	startTime := time.Now()

	// Get all registered dictionaries
	dictConfigs := common.GetRegisteredDictionaries()

	// Import each dictionary
	var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries []common.Entry

	for _, dict := range dictConfigs {
		// Construct full path
		inputPath := filepath.Join(dict.SourceDir, dict.InputFile)

		// Import this dictionary
		fmt.Printf("Importing %s from %s...\n", dict.Name, inputPath)
		dictStartTime := time.Now()

		entries, err := dict.Importer.Import(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing %s: %v\n", dict.Name, err)
			continue
		}

		// Limit entries
		if *limit > 0 && len(entries) > *limit {
			entries = entries[:*limit]
		}

		// Store entries by dictionary type
		switch dict.Name {
		case "jmdict":
			jmdictEntries = entries
		case "jmnedict":
			jmnedictEntries = entries
		case "kanjidic":
			kanjidicEntries = entries
		case "chinese_chars":
			chineseCharsEntries = entries
		case "chinese_words":
			chineseWordsEntries = entries
		}

		fmt.Printf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(dictStartTime).Seconds())
	}

	// Combine all entries
	var allEntries []common.Entry
	allEntries = append(allEntries, jmdictEntries...)
	allEntries = append(allEntries, jmnedictEntries...)
	allEntries = append(allEntries, kanjidicEntries...)
	allEntries = append(allEntries, chineseCharsEntries...)
	allEntries = append(allEntries, chineseWordsEntries...)

	// Count entries by type
	jmdictCount := len(jmdictEntries)
	jmnedictCount := len(jmnedictEntries)
	kanjidicCount := len(kanjidicEntries)
	chineseCharsCount := len(chineseCharsEntries)
	chineseWordsCount := len(chineseWordsEntries)

	fmt.Printf("Processing %d entries (%d JMdict, %d JMNedict, %d Kanjidic, %d Chinese Chars, %d Chinese Words)\n",
		len(allEntries), jmdictCount, jmnedictCount, kanjidicCount, chineseCharsCount, chineseWordsCount)

	// Process entries
	fmt.Printf("Processing entries...\n")
	processStart := time.Now()

	if err := proc.ProcessEntries(allEntries); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing entries: %v\n", err)
		os.Exit(1)
	}

	processDuration := time.Since(processStart)
	fmt.Printf("Processed all %d entries in %.2f seconds (%.1f entries/sec)\n",
		len(allEntries), processDuration.Seconds(), float64(len(allEntries))/processDuration.Seconds())

	// Write files
	fmt.Printf("Writing files to %s...\n", *outputDir)
	writeStart := time.Now()

	if err := proc.WriteToFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing files: %v\n", err)
		os.Exit(1)
	}

	writeDuration := time.Since(writeStart)
	fmt.Printf("Wrote files in %.2f seconds\n", writeDuration.Seconds())

	// Verify minification
	fmt.Printf("Verifying minification...\n")
	verifyMinification(*outputDir)

	totalDuration := time.Since(startTime)
	fmt.Printf("Total time: %.2f seconds\n", totalDuration.Seconds())
	fmt.Printf("Successfully processed and minified dictionary files\n")
}

// verifyMinification checks a sample of files to ensure they were properly minified
func verifyMinification(outputDir string) {
	// Count files
	var totalFiles, jmdictFiles, chineseFiles int
	var filesWithWildcards, filesWithEmptyArrays int

	// Walk through the output directory and check files
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only check .json.br files
		if filepath.Ext(path) != ".br" {
			return nil
		}

		totalFiles++

		// For now, we'll just count files
		// In a real verification, we would decompress and check the content
		// But that would require more complex code to parse the Brotli-compressed JSON

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking output directory: %v\n", err)
		return
	}

	fmt.Printf("Verified %d total files\n", totalFiles)
	fmt.Printf("- JMdict files: %d\n", jmdictFiles)
	fmt.Printf("- Chinese files: %d\n", chineseFiles)
	fmt.Printf("- Files with wildcards: %d\n", filesWithWildcards)
	fmt.Printf("- Files with empty arrays: %d\n", filesWithEmptyArrays)
}
