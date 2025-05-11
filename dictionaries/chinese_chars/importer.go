package chinese_chars

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

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

	// Parse the JSON file as a generic map to handle MongoDB-style fields
	var rawEntries []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&rawEntries); err != nil {
		return nil, err
	}

	// Create a slice of entries for sorting
	tempEntries := make([]ChineseCharEntry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry := ChineseCharEntry{}

		// Map _id to ID (we'll replace this with sequential IDs later)
		if id, ok := rawEntry["_id"]; ok {
			if idStr, ok := id.(string); ok {
				entry.ID = idStr
			}
		}

		// Map char to Traditional
		if char, ok := rawEntry["char"]; ok {
			if charStr, ok := char.(string); ok {
				entry.Traditional = charStr
				// If no simplified form is specified, use the traditional form
				entry.Simplified = charStr
			}
		}

		// Map simplified if available
		if simp, ok := rawEntry["simpVariants"]; ok {
			if simpArr, ok := simp.([]interface{}); ok && len(simpArr) > 0 {
				if simpStr, ok := simpArr[0].(string); ok {
					entry.Simplified = simpStr
				}
			}
		}

		// Map gloss to Definitions
		if gloss, ok := rawEntry["gloss"]; ok {
			if glossStr, ok := gloss.(string); ok {
				entry.Definitions = []string{glossStr}
			}
		}

		// Map strokeCount
		if stroke, ok := rawEntry["strokeCount"]; ok {
			if strokeFloat, ok := stroke.(float64); ok {
				entry.StrokeCount = int(strokeFloat)
			}
		}

		// Ensure ID is set
		if entry.ID == "" {
			entry.ID = entry.Traditional
		}

		// Ensure Traditional is set
		if entry.Traditional == "" && entry.ID != "" {
			entry.Traditional = entry.ID
		}

		tempEntries[i] = entry
	}

	// Sort entries by Traditional character for consistent IDs
	sort.Slice(tempEntries, func(i, j int) bool {
		return tempEntries[i].Traditional < tempEntries[j].Traditional
	})

	// Assign sequential numeric IDs
	entries := make([]common.Entry, len(tempEntries))
	for i, entry := range tempEntries {
		// Assign sequential ID (starting from 1)
		entry.ID = fmt.Sprintf("%d", i+1)
		entries[i] = entry
	}

	return entries, nil
}
