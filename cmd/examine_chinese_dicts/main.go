package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Paths to the Chinese dictionary files
	charDictPath := filepath.Join("dictionaries", "chinese_chars", "source", "dictionary_char_2024-06-17.json")
	wordDictPath := filepath.Join("dictionaries", "chinese_words", "source", "dictionary_word_2024-06-17.json")

	// Examine Chinese character dictionary
	fmt.Println("Examining Chinese character dictionary...")
	examineCharDict(charDictPath)

	// Examine Chinese word dictionary
	fmt.Println("\nExamining Chinese word dictionary...")
	examineWordDict(wordDictPath)
}

func examineCharDict(path string) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Parse the JSON file
	var charEntries []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&charEntries); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Print the total number of entries
	fmt.Printf("Total entries: %d\n", len(charEntries))

	// Examine a few entries
	fmt.Println("Sample entries:")
	for i, entry := range charEntries {
		if i >= 5 {
			break
		}
		fmt.Printf("Entry %d:\n", i+1)
		fmt.Printf("  Raw entry: %v\n", entry)

		// Print specific fields if they exist
		if id, ok := entry["id"]; ok {
			fmt.Printf("  ID: %v\n", id)
		} else if id, ok := entry["_id"]; ok {
			fmt.Printf("  _ID: %v\n", id)
		}

		if trad, ok := entry["traditional"]; ok {
			fmt.Printf("  Traditional: %v\n", trad)
		} else if trad, ok := entry["trad"]; ok {
			fmt.Printf("  Trad: %v\n", trad)
		} else if char, ok := entry["char"]; ok {
			fmt.Printf("  Char: %v\n", char)
		}

		if simp, ok := entry["simplified"]; ok {
			fmt.Printf("  Simplified: %v\n", simp)
		} else if simp, ok := entry["simp"]; ok {
			fmt.Printf("  Simp: %v\n", simp)
		}

		if defs, ok := entry["definitions"]; ok {
			fmt.Printf("  Definitions: %v\n", defs)
		} else if gloss, ok := entry["gloss"]; ok {
			fmt.Printf("  Gloss: %v\n", gloss)
		}

		if pinyin, ok := entry["pinyin"]; ok {
			fmt.Printf("  Pinyin: %v\n", pinyin)
		}

		if stroke, ok := entry["strokeCount"]; ok {
			fmt.Printf("  StrokeCount: %v\n", stroke)
		}

		fmt.Println()
	}

	// Check for entries with non-string IDs or empty IDs
	nonStringIDs := 0
	emptyIDs := 0
	for _, entry := range charEntries {
		if id, ok := entry["id"]; ok {
			if id == "" {
				emptyIDs++
			}
			if _, ok := id.(string); !ok {
				nonStringIDs++
			}
		} else if id, ok := entry["_id"]; ok {
			if id == "" {
				emptyIDs++
			}
			if _, ok := id.(string); !ok {
				nonStringIDs++
			}
		} else {
			emptyIDs++
		}
	}
	fmt.Printf("Entries with empty IDs: %d\n", emptyIDs)
	fmt.Printf("Entries with non-string IDs: %d\n", nonStringIDs)

	// Check for entries where ID != Traditional
	idNotTraditional := 0
	for _, entry := range charEntries {
		var id, trad interface{}
		var hasID, hasTrad bool

		if val, ok := entry["id"]; ok {
			id = val
			hasID = true
		} else if val, ok := entry["_id"]; ok {
			id = val
			hasID = true
		}

		if val, ok := entry["traditional"]; ok {
			trad = val
			hasTrad = true
		} else if val, ok := entry["trad"]; ok {
			trad = val
			hasTrad = true
		} else if val, ok := entry["char"]; ok {
			trad = val
			hasTrad = true
		}

		if hasID && hasTrad && id != trad {
			idNotTraditional++
		}
	}
	fmt.Printf("Entries where ID != Traditional: %d\n", idNotTraditional)
}

func examineWordDict(path string) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Parse the JSON file
	var wordEntries []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&wordEntries); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Print the total number of entries
	fmt.Printf("Total entries: %d\n", len(wordEntries))

	// Examine a few entries
	fmt.Println("Sample entries:")
	for i, entry := range wordEntries {
		if i >= 5 {
			break
		}
		fmt.Printf("Entry %d:\n", i+1)
		fmt.Printf("  Raw entry: %v\n", entry)

		// Print specific fields if they exist
		if id, ok := entry["id"]; ok {
			fmt.Printf("  ID: %v\n", id)
		} else if id, ok := entry["_id"]; ok {
			fmt.Printf("  _ID: %v\n", id)
		}

		if trad, ok := entry["traditional"]; ok {
			fmt.Printf("  Traditional: %v\n", trad)
		} else if trad, ok := entry["trad"]; ok {
			fmt.Printf("  Trad: %v\n", trad)
		}

		if simp, ok := entry["simplified"]; ok {
			fmt.Printf("  Simplified: %v\n", simp)
		} else if simp, ok := entry["simp"]; ok {
			fmt.Printf("  Simp: %v\n", simp)
		}

		if defs, ok := entry["definitions"]; ok {
			fmt.Printf("  Definitions: %v\n", defs)
		} else if items, ok := entry["items"]; ok {
			fmt.Printf("  Items: %v\n", items)
		} else if gloss, ok := entry["gloss"]; ok {
			fmt.Printf("  Gloss: %v\n", gloss)
		}

		if pinyin, ok := entry["pinyin"]; ok {
			fmt.Printf("  Pinyin: %v\n", pinyin)
		} else if pinyinSearch, ok := entry["pinyinSearchString"]; ok {
			fmt.Printf("  PinyinSearchString: %v\n", pinyinSearch)
		}

		if hsk, ok := entry["hskLevel"]; ok {
			fmt.Printf("  HskLevel: %v\n", hsk)
		} else if stats, ok := entry["statistics"]; ok {
			if statsMap, ok := stats.(map[string]interface{}); ok {
				if hsk, ok := statsMap["hskLevel"]; ok {
					fmt.Printf("  Statistics.HskLevel: %v\n", hsk)
				}
			}
		}

		fmt.Println()
	}

	// Check for entries with non-string IDs or empty IDs
	nonStringIDs := 0
	emptyIDs := 0
	for _, entry := range wordEntries {
		if id, ok := entry["id"]; ok {
			if id == "" {
				emptyIDs++
			}
			if _, ok := id.(string); !ok {
				nonStringIDs++
			}
		} else if id, ok := entry["_id"]; ok {
			if id == "" {
				emptyIDs++
			}
			if _, ok := id.(string); !ok {
				nonStringIDs++
			}
		} else {
			emptyIDs++
		}
	}
	fmt.Printf("Entries with empty IDs: %d\n", emptyIDs)
	fmt.Printf("Entries with non-string IDs: %d\n", nonStringIDs)

	// Check for entries where ID != Traditional
	idNotTraditional := 0
	for _, entry := range wordEntries {
		var id, trad interface{}
		var hasID, hasTrad bool

		if val, ok := entry["id"]; ok {
			id = val
			hasID = true
		} else if val, ok := entry["_id"]; ok {
			id = val
			hasID = true
		}

		if val, ok := entry["traditional"]; ok {
			trad = val
			hasTrad = true
		} else if val, ok := entry["trad"]; ok {
			trad = val
			hasTrad = true
		}

		if hasID && hasTrad && id != trad {
			idNotTraditional++
		}
	}
	fmt.Printf("Entries where ID != Traditional: %d\n", idNotTraditional)
}
