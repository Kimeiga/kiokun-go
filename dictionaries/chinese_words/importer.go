package chinese_words

import (
	"encoding/json"
	"os"
	"strings"

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

	// Convert to ChineseWordEntry structs
	entries := make([]common.Entry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry := ChineseWordEntry{}

		// Map _id to ID
		if id, ok := rawEntry["_id"]; ok {
			if idStr, ok := id.(string); ok {
				entry.ID = idStr
			}
		}

		// Map trad to Traditional
		if trad, ok := rawEntry["trad"]; ok {
			if tradStr, ok := trad.(string); ok {
				entry.Traditional = tradStr
			}
		}

		// Map simp to Simplified
		if simp, ok := rawEntry["simp"]; ok {
			if simpStr, ok := simp.(string); ok {
				entry.Simplified = simpStr
			}
		}

		// Map definitions or gloss
		if defs, ok := rawEntry["definitions"]; ok {
			if defsArr, ok := defs.([]interface{}); ok {
				entry.Definitions = make([]string, len(defsArr))
				for j, def := range defsArr {
					if defStr, ok := def.(string); ok {
						entry.Definitions[j] = defStr
					}
				}
			}
		} else if gloss, ok := rawEntry["gloss"]; ok {
			if glossStr, ok := gloss.(string); ok {
				entry.Definitions = []string{glossStr}
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
		} else if pinyinSearch, ok := rawEntry["pinyinSearchString"]; ok {
			if pinyinStr, ok := pinyinSearch.(string); ok && pinyinStr != "" {
				// Split by spaces to get individual pinyin values
				entry.Pinyin = strings.Fields(pinyinStr)
			}
		}

		// Map HSK level
		if stats, ok := rawEntry["statistics"]; ok {
			if statsMap, ok := stats.(map[string]interface{}); ok {
				if hsk, ok := statsMap["hskLevel"]; ok {
					if hskFloat, ok := hsk.(float64); ok {
						entry.HskLevel = int(hskFloat)
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

		// If Simplified is not set, use Traditional
		if entry.Simplified == "" {
			entry.Simplified = entry.Traditional
		}

		entries[i] = entry
	}

	return entries, nil
}
