package main

import (
	"fmt"
	"os"

	// Import for side effects (dictionary registration)
	_ "kiokun-go/dictionaries/chinese_chars"
	_ "kiokun-go/dictionaries/chinese_words"
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

	// Filter entries
	filteredEntries := FilterEntries(entries, config, logf)

	// Process entries
	if err := ProcessEntries(filteredEntries, config, logf); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing entries: %v\n", err)
		os.Exit(1)
	}

	logf("Successfully processed dictionary files\n")
}
