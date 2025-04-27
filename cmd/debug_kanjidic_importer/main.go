package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/kanjidic"
)

func main() {
	// Path to the kanjidic source file
	sourceFile := filepath.Join("dictionaries", "kanjidic", "source", "kanjidic2-en-3.6.1.json")

	fmt.Printf("Testing kanjidic importer with file: %s\n", sourceFile)

	// Check if file exists
	fileInfo, err := os.Stat(sourceFile)
	if err != nil {
		fmt.Printf("Error accessing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("File size: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))

	// Create the importer
	importer := &kanjidic.Importer{}

	// Import the file
	entries, err := importer.Import(sourceFile)
	if err != nil {
		fmt.Printf("Error importing file: %v\n", err)
		os.Exit(1)
	}

	// Print stats
	fmt.Printf("Successfully imported kanjidic file\n")
	fmt.Printf("Number of kanji entries: %d\n\n", len(entries))

	// Print first 5 entries as examples
	if len(entries) > 0 {
		fmt.Printf("Sample entries:\n")
		maxSamples := 5
		if len(entries) < maxSamples {
			maxSamples = len(entries)
		}

		for i := 0; i < maxSamples; i++ {
			if kanji, ok := entries[i].(kanjidic.Kanji); ok {
				fmt.Printf("%d. Character: %s\n", i+1, kanji.Character)
				fmt.Printf("   JLPT Level: %d\n", kanji.JLPT)
				fmt.Printf("   Grade: %d\n", kanji.Grade)
				fmt.Printf("   Stroke Count: %d\n", kanji.Stroke)
				fmt.Printf("   Meanings: %v\n", kanji.Meanings)
				fmt.Printf("   On Readings: %v\n", kanji.OnYomi)
				fmt.Printf("   Kun Readings: %v\n\n", kanji.KunYomi)
			} else {
				fmt.Printf("%d. Entry is not a Kanji type: %T\n", i+1, entries[i])
			}
		}
	} else {
		fmt.Printf("WARNING: No kanji entries found in the kanjidic file!\n")
	}
}
