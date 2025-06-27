package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/andybalholm/brotli"
)

// Test data structure similar to Chinese word entry
type TestEntry struct {
	ID          string   `json:"id"`
	Traditional string   `json:"traditional"`
	Simplified  string   `json:"simplified"`
	Definitions []string `json:"definitions,omitempty"`
	Pinyin      []string `json:"pinyin,omitempty"`
	HskLevel    int      `json:"hskLevel,omitempty"`
}

// Old writeCompressedJSON function (with the bug)
func writeCompressedJSONOld(filename string, obj interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bw := brotli.NewWriter(file)
	defer bw.Close()

	encoder := json.NewEncoder(bw)
	err = encoder.Encode(obj)
	return err
}

// New writeCompressedJSON function (with the fix)
func writeCompressedJSONNew(filename string, obj interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bw := brotli.NewWriter(file)
	
	encoder := json.NewEncoder(bw)
	if err := encoder.Encode(obj); err != nil {
		bw.Close() // Close on error to clean up
		return err
	}
	
	// Explicitly close the brotli writer to flush buffers
	return bw.Close()
}

// Function to read and verify compressed JSON
func readCompressedJSON(filename string, obj interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	br := brotli.NewReader(file)
	decoder := json.NewDecoder(br)
	return decoder.Decode(obj)
}

func main() {
	// Create test data similar to the 日 entry
	testEntry := TestEntry{
		ID:          "57102",
		Traditional: "日",
		Simplified:  "日",
		Definitions: []string{"sun", "day", "daytime"},
		Pinyin:      []string{"rì"},
		HskLevel:    1,
	}

	// Create test directory
	testDir := "test_brotli_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Test old function
	oldFile := filepath.Join(testDir, "old_157102.json.br")
	fmt.Println("Testing old writeCompressedJSON function...")
	if err := writeCompressedJSONOld(oldFile, testEntry); err != nil {
		fmt.Printf("Error writing with old function: %v\n", err)
	} else {
		fmt.Println("Old function wrote file successfully")
		
		// Try to read it back
		var readEntry TestEntry
		if err := readCompressedJSON(oldFile, &readEntry); err != nil {
			fmt.Printf("Error reading old file: %v\n", err)
		} else {
			fmt.Printf("Old file read successfully: %+v\n", readEntry)
		}
	}

	// Test new function
	newFile := filepath.Join(testDir, "new_157102.json.br")
	fmt.Println("\nTesting new writeCompressedJSON function...")
	if err := writeCompressedJSONNew(newFile, testEntry); err != nil {
		fmt.Printf("Error writing with new function: %v\n", err)
	} else {
		fmt.Println("New function wrote file successfully")
		
		// Try to read it back
		var readEntry TestEntry
		if err := readCompressedJSON(newFile, &readEntry); err != nil {
			fmt.Printf("Error reading new file: %v\n", err)
		} else {
			fmt.Printf("New file read successfully: %+v\n", readEntry)
		}
	}

	// Compare file sizes
	oldStat, _ := os.Stat(oldFile)
	newStat, _ := os.Stat(newFile)
	fmt.Printf("\nFile sizes:\n")
	fmt.Printf("Old function: %d bytes\n", oldStat.Size())
	fmt.Printf("New function: %d bytes\n", newStat.Size())
}
