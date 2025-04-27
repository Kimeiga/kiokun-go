package chinese_words

import (
	"encoding/json"
	"os"

	"kiokun-go/dictionaries/common"
)

// Importer handles importing Chinese word dictionary
type Importer struct{}

// Name returns the name of this importer
func (i *Importer) Name() string {
	return "chinese_words"
}

// Import reads and processes the Chinese word dictionary
func (i *Importer) Import(path string) ([]common.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the JSON file
	var wordEntries []ChineseWordEntry
	if err := json.NewDecoder(file).Decode(&wordEntries); err != nil {
		return nil, err
	}

	// Convert to common.Entry interface
	entries := make([]common.Entry, len(wordEntries))
	for i, entry := range wordEntries {
		// If ID is not set, use traditional word as ID
		if entry.ID == "" {
			entry.ID = entry.Traditional
		}
		entries[i] = entry
	}

	return entries, nil
}
