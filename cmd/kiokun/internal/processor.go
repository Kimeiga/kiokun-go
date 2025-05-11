package internal

import (
	"fmt"
	"time"

	"kiokun-go/dictionaries/common"
	"kiokun-go/processor"
)

// ProcessEntries processes dictionary entries and writes them to files
func ProcessEntries(entries *DictionaryEntries, config *Config, logf LogFunc) error {
	// Always use the index-based processor
	logf("Using index-based processor with separate files for each dictionary\n")
	proc, err := processor.NewIndexProcessor(config.OutputDir, config.FileWriters)

	if err != nil {
		return fmt.Errorf("error creating processor: %v", err)
	}

	// Calculate total entries and pre-allocate the slice
	totalEntries := len(entries.JMdict) + len(entries.JMNedict) + len(entries.Kanjidic) +
		len(entries.ChineseChars) + len(entries.ChineseWords)

	// Pre-allocate the slice to avoid reallocations
	allEntries := make([]common.Entry, 0, totalEntries)

	// Append entries directly to avoid function call overhead
	allEntries = append(allEntries, entries.JMdict...)
	allEntries = append(allEntries, entries.JMNedict...)
	allEntries = append(allEntries, entries.Kanjidic...)
	allEntries = append(allEntries, entries.ChineseChars...)
	allEntries = append(allEntries, entries.ChineseWords...)

	// Log entry counts directly from the source slices to avoid type assertions
	logf("Processing %d entries (%d JMdict, %d JMNedict, %d Kanjidic, %d Chinese Chars, %d Chinese Words)\n",
		totalEntries, len(entries.JMdict), len(entries.JMNedict), len(entries.Kanjidic),
		len(entries.ChineseChars), len(entries.ChineseWords))

	// Process entries in batches with progress reporting
	logf("Processing entries in batches...\n")
	processStart := time.Now()

	// If the total entries is small enough, process them all at once
	if totalEntries <= config.BatchSize*10 {
		logf("\rProcessing all %d entries...", totalEntries)
		if err := proc.ProcessEntries(allEntries); err != nil {
			return fmt.Errorf("error processing entries: %v", err)
		}
	} else {
		// Process in batches for larger datasets
		for start := 0; start < totalEntries; start += config.BatchSize {
			end := start + config.BatchSize
			if end > totalEntries {
				end = totalEntries
			}

			batchEntries := allEntries[start:end]

			// Only update progress every 10 batches to reduce output
			if start%(config.BatchSize*10) == 0 || end == totalEntries {
				logf("\rProcessing entries %d-%d of %d (%.1f%%)...",
					start+1, end, totalEntries, float64(end)/float64(totalEntries)*100)
			}

			// Process entries in a batch
			if err := proc.ProcessEntries(batchEntries); err != nil {
				return fmt.Errorf("error processing batch: %v", err)
			}
		}
	}
	processDuration := time.Since(processStart)
	logf("\rProcessed all %d entries (100%%) in %.2f seconds (%.1f entries/sec)\n",
		totalEntries, processDuration.Seconds(), float64(totalEntries)/processDuration.Seconds())

	// Write all processed entries to files
	logf("Writing files to %s...\n", config.OutputDir)
	if err := proc.WriteToFiles(); err != nil {
		return fmt.Errorf("error writing files: %v", err)
	}

	return nil
}
