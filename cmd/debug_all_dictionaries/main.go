package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	_ "kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	_ "kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
	_ "kiokun-go/dictionaries/kanjidic"
)

func main() {
	fmt.Println("Debugging all dictionary importers")

	// Get all registered dictionaries
	dictConfigs := common.GetRegisteredDictionaries()

	// Import each dictionary
	var jmdictEntries, jmnedictEntries, kanjidicEntries []common.Entry

	for _, dict := range dictConfigs {
		// Construct full path
		inputPath := filepath.Join(dict.SourceDir, dict.InputFile)

		fmt.Printf("Importing %s from %s...\n", dict.Name, inputPath)

		// Import this dictionary
		entries, err := dict.Importer.Import(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing %s: %v\n", dict.Name, err)
			os.Exit(1)
		}

		// Store entries by dictionary type
		switch dict.Name {
		case "jmdict":
			jmdictEntries = entries
		case "jmnedict":
			jmnedictEntries = entries
		case "kanjidic":
			kanjidicEntries = entries
		}

		fmt.Printf("Imported %s: %d entries\n", dict.Name, len(entries))
	}

	// Display statistics
	fmt.Printf("\nDictionary entry counts:\n")
	fmt.Printf("- JMdict entries: %d\n", len(jmdictEntries))
	fmt.Printf("- JMNedict entries: %d\n", len(jmnedictEntries))
	fmt.Printf("- Kanjidic entries: %d\n", len(kanjidicEntries))

	// Sample entries from each dictionary
	printSampleEntries("JMdict", jmdictEntries, 2)
	printSampleEntries("JMNedict", jmnedictEntries, 2)
	printSampleEntries("Kanjidic", kanjidicEntries, 2)
}

func printSampleEntries(dictName string, entries []common.Entry, count int) {
	if len(entries) == 0 {
		fmt.Printf("\nNo entries in %s\n", dictName)
		return
	}

	fmt.Printf("\nSample entries from %s:\n", dictName)

	maxSamples := count
	if len(entries) < maxSamples {
		maxSamples = len(entries)
	}

	for i := 0; i < maxSamples; i++ {
		fmt.Printf("%d. ID: %s, Filename: %s\n", i+1, entries[i].GetID(), entries[i].GetFilename())

		switch entry := entries[i].(type) {
		case jmdict.Word:
			if len(entry.Kanji) > 0 {
				fmt.Printf("   Kanji: %s\n", entry.Kanji[0].Text)
			}
			if len(entry.Kana) > 0 {
				fmt.Printf("   Reading: %s\n", entry.Kana[0].Text)
			}

		case jmnedict.Name:
			if len(entry.Kanji) > 0 {
				fmt.Printf("   Kanji: %s\n", entry.Kanji[0])
			}
			if len(entry.Reading) > 0 {
				fmt.Printf("   Reading: %s\n", entry.Reading[0])
			}

		case kanjidic.Kanji:
			fmt.Printf("   Character: %s\n", entry.Character)
			fmt.Printf("   JLPT Level: %d\n", entry.JLPT)
			fmt.Printf("   Meanings: %v\n", entry.Meanings)
		}
	}
}
