package chinese_chars

import (
	"encoding/json"
	"os"

	"kiokun-go/dictionaries/common"
)

// Importer handles importing Chinese character dictionary
type Importer struct{}

// Name returns the name of this importer
func (i *Importer) Name() string {
	return "chinese_chars"
}

// Import reads and processes the Chinese character dictionary
func (i *Importer) Import(path string) ([]common.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the JSON file
	var charEntries []ChineseCharEntry
	if err := json.NewDecoder(file).Decode(&charEntries); err != nil {
		return nil, err
	}

	// Convert to common.Entry interface
	entries := make([]common.Entry, len(charEntries))
	for i, entry := range charEntries {
		// If ID is not set, use traditional character as ID
		if entry.ID == "" {
			entry.ID = entry.Traditional
		}
		entries[i] = entry
	}

	return entries, nil
}
