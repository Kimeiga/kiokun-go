package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/processor"
)

func main() {
	// Create a temporary output directory
	tempDir := filepath.Join(os.TempDir(), "kiokun-debug")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		os.Exit(1)
	}

	// Create an index processor
	proc, err := processor.NewIndexProcessor(tempDir, 4)
	if err != nil {
		fmt.Printf("Error creating index processor: %v\n", err)
		os.Exit(1)
	}

	// Create a few sample Chinese character entries
	charEntries := []common.Entry{
		chinese_chars.ChineseCharEntry{
			ID:          "char1",
			Traditional: "人",
			Simplified:  "人",
			Definitions: []string{"person", "people"},
			Pinyin:      []string{"rén"},
			StrokeCount: 2,
		},
		chinese_chars.ChineseCharEntry{
			ID:          "char2",
			Traditional: "山",
			Simplified:  "山",
			Definitions: []string{"mountain"},
			Pinyin:      []string{"shān"},
			StrokeCount: 3,
		},
	}

	// Create a few sample Chinese word entries
	wordEntries := []common.Entry{
		chinese_words.ChineseWordEntry{
			ID:          "word1",
			Traditional: "人山人海",
			Simplified:  "人山人海",
			Definitions: []string{"huge crowd of people"},
			Pinyin:      []string{"rén shān rén hǎi"},
		},
		chinese_words.ChineseWordEntry{
			ID:          "word2",
			Traditional: "山水",
			Simplified:  "山水",
			Definitions: []string{"landscape", "scenery"},
			Pinyin:      []string{"shān shuǐ"},
		},
	}

	// Process the entries
	fmt.Println("Processing Chinese character entries...")
	if err := proc.ProcessEntries(charEntries); err != nil {
		fmt.Printf("Error processing character entries: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing Chinese word entries...")
	if err := proc.ProcessEntries(wordEntries); err != nil {
		fmt.Printf("Error processing word entries: %v\n", err)
		os.Exit(1)
	}

	// Write to files
	fmt.Println("Writing to files...")
	if err := proc.WriteToFiles(); err != nil {
		fmt.Printf("Error writing to files: %v\n", err)
		os.Exit(1)
	}

	// Print the index entries for debugging
	fmt.Println("\nIndex entries:")
	for key, entry := range proc.GetIndex() {
		fmt.Printf("Key: %s\n", key)
		
		// Print exact matches
		if entry.E != nil {
			fmt.Println("  Exact matches:")
			for dictType, ids := range entry.E {
				fmt.Printf("    %s: %v\n", dictType, ids)
			}
		}
		
		// Print contained-in matches
		if entry.C != nil {
			fmt.Println("  Contained-in matches:")
			for dictType, ids := range entry.C {
				fmt.Printf("    %s: %v\n", dictType, ids)
			}
		}
		
		fmt.Println()
	}

	fmt.Println("Debug complete. Check the output above for any issues.")
}
