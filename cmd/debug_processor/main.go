package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
	"kiokun-go/processor"
)

func main() {
	// Configuration flags
	dictDir := flag.String("dictdir", "dictionaries", "Base directory containing dictionary packages")
	outputDir := flag.String("outdir", "output_test", "Output directory for processed files")
	sampleSize := flag.Int("samples", 5, "Number of samples to process for testing")
	flag.Parse()

	// Ensure we're using an absolute path for the output directory
	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
		os.Exit(1)
	}
	*outputDir = absOutputDir

	fmt.Printf("Starting debugging processor...\n")
	fmt.Printf("Input dictionaries directory: %s\n", *dictDir)
	fmt.Printf("Output directory: %s\n", *outputDir)

	// Debug check for dictionary source files
	fmt.Printf("\n=== Checking dictionary source files ===\n")
	checkDictSourceFile(*dictDir, "jmdict")
	checkDictSourceFile(*dictDir, "jmnedict")
	checkDictSourceFile(*dictDir, "kanjidic")

	// Set dictionaries base path
	common.SetDictionariesBasePath(*dictDir)

	// Create processor
	proc, err := processor.New(*outputDir, runtime.NumCPU())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating processor: %v\n", err)
		os.Exit(1)
	}

	// Import all dictionaries with debug info
	fmt.Printf("\n=== Importing dictionaries with debug info ===\n")
	entries, err := common.ImportAllDictionaries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error importing dictionaries: %v\n", err)
		os.Exit(1)
	}

	// Count entries by dictionary type
	var jmdictCount, jmnedictCount, kanjidicCount, otherCount int
	for _, entry := range entries {
		switch entry.(type) {
		case jmdict.Word:
			jmdictCount++
		case jmnedict.Name:
			jmnedictCount++
		case kanjidic.Kanji:
			kanjidicCount++
		default:
			otherCount++
			fmt.Printf("Unknown entry type: %v\n", reflect.TypeOf(entry))
		}
	}

	fmt.Printf("\n=== Dictionary import summary ===\n")
	fmt.Printf("- Total entries: %d\n", len(entries))
	fmt.Printf("- JMdict entries: %d\n", jmdictCount)
	fmt.Printf("- JMNedict entries: %d\n", jmnedictCount)
	fmt.Printf("- Kanjidic entries: %d\n", kanjidicCount)
	fmt.Printf("- Unknown entry types: %d\n", otherCount)

	// Print sample of entries from each dictionary type
	countJMdict := 0
	countJMNedict := 0
	countKanjidic := 0

	fmt.Printf("\n=== Sample dictionary entries before processing ===\n")
	for i, entry := range entries {
		switch e := entry.(type) {
		case jmdict.Word:
			if countJMdict < *sampleSize {
				fmt.Printf("JMdict entry %d: %s\n", countJMdict, summarizeJMdictWord(e))
				countJMdict++
			}
		case jmnedict.Name:
			if countJMNedict < *sampleSize {
				fmt.Printf("JMNedict entry %d: %s\n", countJMNedict, summarizeJMNedictName(e))
				countJMNedict++
			}
		case kanjidic.Kanji:
			if countKanjidic < *sampleSize {
				fmt.Printf("Kanjidic entry %d: %s\n", countKanjidic, e.Character)
				countKanjidic++
			}
		}

		if countJMdict >= *sampleSize && countJMNedict >= *sampleSize && countKanjidic >= *sampleSize {
			break
		}

		// Only check the first 1000 entries to avoid excessive output
		if i >= 1000 {
			break
		}
	}

	fmt.Printf("\nFound samples - JMdict: %d, JMNedict: %d, Kanjidic: %d\n",
		countJMdict, countJMNedict, countKanjidic)
	fmt.Printf("Total entries to process: %d\n", len(entries))
	fmt.Printf("Processing first 1000 entries only for debugging...\n")

	// Process a limited number of entries for debugging
	debugEntries := entries
	if len(entries) > 1000 {
		debugEntries = entries[:1000]
	}

	if err := proc.ProcessEntries(debugEntries); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing entries: %v\n", err)
		os.Exit(1)
	}

	// Write to files
	fmt.Printf("Writing files to %s...\n", *outputDir)
	if err := proc.WriteToFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Completed debugging processor\n")
	fmt.Printf("Now you can run analyzer on the %s directory\n", *outputDir)
}

// Helper function to check dictionary source files
func checkDictSourceFile(baseDir, dictName string) {
	sourceDir := filepath.Join(baseDir, dictName, "source")
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf("Error reading %s source directory: %v\n", dictName, err)
		return
	}

	if len(files) == 0 {
		fmt.Printf("%s: No source files found\n", dictName)
		return
	}

	fmt.Printf("%s source files:\n", dictName)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, _ := file.Info()
		size := float64(info.Size()) / 1024.0 / 1024.0 // Size in MB
		fmt.Printf("  - %s (%.2f MB)\n", file.Name(), size)
	}
}

// Helper function to summarize JMdict word entries
func summarizeJMdictWord(word jmdict.Word) string {
	forms := []string{}
	for _, k := range word.Kanji {
		forms = append(forms, k.Text)
	}
	for _, k := range word.Kana {
		forms = append(forms, k.Text)
	}
	return fmt.Sprintf("ID:%s Forms:%v", word.ID, forms)
}

// Helper function to summarize JMNedict name entries
func summarizeJMNedictName(name jmnedict.Name) string {
	return fmt.Sprintf("ID:%s Kanji:%v Reading:%v", name.ID, name.Kanji, name.Reading)
}
