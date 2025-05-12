package chinese_words

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

	// Check if the file is JSONL or JSON
	isJSONL := strings.HasSuffix(path, ".jsonl")

	fmt.Printf("DEBUG: Importing Chinese word dictionary from %s (isJSONL: %v)\n", path, isJSONL)

	// Parse the file based on its format
	var rawEntries []map[string]interface{}

	if isJSONL {
		// Parse JSONL (one JSON object per line)
		scanner := bufio.NewScanner(file)
		lineCount := 0
		ribenFound := false

		// Set a larger buffer size for the scanner
		const maxScanTokenSize = 1024 * 1024 // 1MB
		buf := make([]byte, maxScanTokenSize)
		scanner.Buffer(buf, maxScanTokenSize)

		for scanner.Scan() {
			lineCount++
			line := scanner.Text()

			// Debug: Check for "日本" in the raw line
			if strings.Contains(line, "日本") && !strings.Contains(line, "日本國誌") && !strings.Contains(line, "日本国志") {
				fmt.Printf("DEBUG: Found '日本' in line %d: %s\n", lineCount, line[:100]+"...")
				ribenFound = true
			}

			var entry map[string]interface{}
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				return nil, fmt.Errorf("error parsing JSONL line %d: %v", lineCount, err)
			}

			// Debug: Check for "日本" in the parsed entry
			if ribenFound {
				if word, ok := entry["simp"]; ok {
					if wordStr, ok := word.(string); ok && wordStr == "日本" {
						fmt.Printf("DEBUG: Found entry with simp='日本': %+v\n", entry)
					}
				}
				if word, ok := entry["trad"]; ok {
					if wordStr, ok := word.(string); ok && wordStr == "日本" {
						fmt.Printf("DEBUG: Found entry with trad='日本': %+v\n", entry)
					}
				}
			}

			rawEntries = append(rawEntries, entry)
		}

		fmt.Printf("DEBUG: Processed %d lines from JSONL file\n", lineCount)

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading JSONL file: %v", err)
		}
	} else {
		// Parse JSON array
		if err := json.NewDecoder(file).Decode(&rawEntries); err != nil {
			return nil, err
		}
		fmt.Printf("DEBUG: Processed JSON array with %d entries\n", len(rawEntries))
	}

	// Debug: Count entries containing "日本"
	var ribenEntries []map[string]interface{}
	for _, entry := range rawEntries {
		// Check all possible field names for "日本"
		fieldsToCheck := []string{"word", "simplified", "simp", "trad", "traditional"}

		for _, field := range fieldsToCheck {
			if word, ok := entry[field]; ok {
				if wordStr, ok := word.(string); ok && wordStr == "日本" {
					// Only add if not already added
					alreadyAdded := false
					for _, e := range ribenEntries {
						// Compare by _id since maps can't be directly compared
						if eID, ok := e["_id"]; ok {
							if entryID, ok := entry["_id"]; ok && eID == entryID {
								alreadyAdded = true
								break
							}
						}
					}

					if !alreadyAdded {
						ribenEntries = append(ribenEntries, entry)
					}
				}
			}
		}
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

		// Map traditional form
		// First try "trad" (JSONL format)
		if trad, ok := rawEntry["trad"]; ok {
			if tradStr, ok := trad.(string); ok {
				entry.Traditional = tradStr
				// If no simplified form is specified, use the traditional form
				entry.Simplified = tradStr
			}
		} else if word, ok := rawEntry["word"]; ok {
			// Then try "word" (JSON format)
			if wordStr, ok := word.(string); ok {
				entry.Traditional = wordStr
				// If no simplified form is specified, use the traditional form
				entry.Simplified = wordStr
			}
		}

		// Map simplified form
		// First try "simp" (JSONL format)
		if simp, ok := rawEntry["simp"]; ok {
			if simpStr, ok := simp.(string); ok {
				entry.Simplified = simpStr
			}
		} else if simp, ok := rawEntry["simplified"]; ok {
			// Then try "simplified" (JSON format)
			if simpStr, ok := simp.(string); ok {
				entry.Simplified = simpStr
			}
		}

		// Map pinyin
		// First try JSONL format (items array with pinyin field)
		if items, ok := rawEntry["items"]; ok {
			if itemsArr, ok := items.([]interface{}); ok && len(itemsArr) > 0 {
				for _, item := range itemsArr {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if pinyin, ok := itemMap["pinyin"]; ok {
							if pinyinStr, ok := pinyin.(string); ok {
								entry.Pinyin = append(entry.Pinyin, pinyinStr)
							}
						}

						// Extract definitions from items
						if defs, ok := itemMap["definitions"]; ok {
							if defsArr, ok := defs.([]interface{}); ok {
								for _, def := range defsArr {
									if defStr, ok := def.(string); ok {
										entry.Definitions = append(entry.Definitions, defStr)
									}
								}
							}
						}
					}
				}
			}
		} else {
			// Try JSON format
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

			// Map definitions (JSON format)
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
		}

		// If no definitions found yet, try the gloss field (JSONL format)
		if len(entry.Definitions) == 0 {
			if gloss, ok := rawEntry["gloss"]; ok {
				if glossStr, ok := gloss.(string); ok {
					entry.Definitions = []string{glossStr}
				}
			}
		}

		// Map HSK level
		// First try statistics.hskLevel (JSONL format)
		if stats, ok := rawEntry["statistics"]; ok {
			if statsMap, ok := stats.(map[string]interface{}); ok {
				if hskLevel, ok := statsMap["hskLevel"]; ok {
					if hskFloat, ok := hskLevel.(float64); ok {
						entry.HskLevel = int(hskFloat)
					}
				}
			}
		} else if hsk, ok := rawEntry["hsk"]; ok {
			// Then try hsk (JSON format)
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
