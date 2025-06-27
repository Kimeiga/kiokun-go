package main

import (
	"fmt"
	"os"

	// Import for side effects (dictionary registration)
	_ "kiokun-go/dictionaries/chinese_chars"
	_ "kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/ids"
	_ "kiokun-go/dictionaries/ids"
	_ "kiokun-go/dictionaries/jmdict"
	_ "kiokun-go/dictionaries/jmnedict"
	_ "kiokun-go/dictionaries/kanjidic"

	// Import local package functions
	. "kiokun-go/cmd/kiokun/internal"
)

func main() {
	// Parse configuration
	config, logf, err := ParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup dictionary files if they don't exist
	if err := SetupDictionaryFiles(logf); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up dictionary files: %v\n", err)
		os.Exit(1)
	}

	logf("Starting initialization...\n")
	logf("Input dictionaries directory: %s\n", config.DictDir)
	logf("Output directory: %s\n", config.OutputDir)
	logf("Using %d processing workers and %d file writers\n", config.Workers, config.FileWriters)
	if config.OutputMode != OutputAll {
		logf("Filtering mode: %s\n", config.OutputMode)
	}

	// Load dictionaries
	entries, err := LoadDictionaries(config, logf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading dictionaries: %v\n", err)
		os.Exit(1)
	}

	// Create IDS lookup map
	idsMap := make(map[string]string)
	for _, entry := range entries.IDS {
		if idsEntry, ok := entry.(ids.IDSEntry); ok {
			idsMap[idsEntry.Character] = idsEntry.IDS
		}
	}

	logf("Created IDS lookup map with %d entries\n", len(idsMap))

	// Filter entries
	filteredEntries := FilterEntries(entries, config, logf)

	// Process entries with IDS map
	if err := ProcessEntriesWithIDS(filteredEntries, config, logf, idsMap); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing entries: %v\n", err)
		os.Exit(1)
	}

	logf("Successfully processed dictionary files\n")
}
