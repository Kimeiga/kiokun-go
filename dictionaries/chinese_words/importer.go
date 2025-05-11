package chinese_words

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

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

	// Parse the JSON file as a generic map to handle MongoDB-style fields
	var rawEntries []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&rawEntries); err != nil {
		return nil, err
	}

	// Create a slice of entries for sorting
	tempEntries := make([]ChineseWordEntry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry := ChineseWordEntry{}

		// Map _id to ID (we'll replace this with sequential IDs later)
		if id, ok := rawEntry["_id"]; ok {
			if idStr, ok := id.(string); ok {
				entry.ID = idStr
			}
		}

		// Map word to Traditional
		if word, ok := rawEntry["word"]; ok {
			if wordStr, ok := word.(string); ok {
				entry.Traditional = wordStr
				// If no simplified form is specified, use the traditional form
				entry.Simplified = wordStr
			}
		}

		// Map simplified if available
		if simp, ok := rawEntry["simplified"]; ok {
			if simpStr, ok := simp.(string); ok {
				entry.Simplified = simpStr
			}
		}

		// Map pinyin
		if pinyin, ok := rawEntry["pinyin"]; ok {
			if pinyinArr, ok := pinyin.([]interface{}); ok {
				entry.Pinyin = make([]string, len(pinyinArr))
				for j, p := range pinyinArr {
					if pStr, ok := p.(string); ok {
						entry.Pinyin[j] = pStr
					}
				}
			} else if pinyinStr, ok := pinyin.(string); ok {
				entry.Pinyin = []string{pinyinStr}
			}
		}

		// Map definitions
		if defs, ok := rawEntry["definitions"]; ok {
			if defsArr, ok := defs.([]interface{}); ok {
				entry.Definitions = make([]string, len(defsArr))
				for j, d := range defsArr {
					if dStr, ok := d.(string); ok {
						entry.Definitions[j] = dStr
					}
				}
			} else if defStr, ok := defs.(string); ok {
				entry.Definitions = []string{defStr}
			}
		}

		// Map HSK level
		if hsk, ok := rawEntry["hsk"]; ok {
			if hskFloat, ok := hsk.(float64); ok {
				entry.HskLevel = int(hskFloat)
			}
		}

		// Map frequency
		if freq, ok := rawEntry["frequency"]; ok {
			if freqMap, ok := freq.(map[string]interface{}); ok {
				entry.Frequency = make(map[string]int)
				for k, v := range freqMap {
					if vFloat, ok := v.(float64); ok {
						entry.Frequency[k] = int(vFloat)
					}
				}
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

	// Sort entries by Traditional word for consistent IDs
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
