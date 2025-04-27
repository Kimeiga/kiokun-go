package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"kiokun-go/dictionaries/jmnedict"
)

func main() {
	// Direct access to JMNedict source file
	sourceFile := filepath.Join("dictionaries", "jmnedict", "source", "jmnedict-all-3.6.1.json")

	fmt.Printf("Debugging JMNedict importer with file: %s\n", sourceFile)

	// Check if file exists
	fileInfo, err := os.Stat(sourceFile)
	if err != nil {
		fmt.Printf("Error accessing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("File size: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))

	// Open the file directly
	file, err := os.Open(sourceFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Try to decode the JSON
	var dict jmnedict.JMNedictTypes
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dict); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)

		// Rewind the file for further debugging
		file.Seek(0, 0)

		// Read the first 200 bytes to check format
		buf := make([]byte, 200)
		n, _ := file.Read(buf)
		fmt.Printf("First %d bytes of file:\n%s\n", n, buf[:n])

		// Check the JMNedict struct definition
		fmt.Printf("JMNedictTypes struct type: %v\n", reflect.TypeOf(dict))

		os.Exit(1)
	}

	// Success - print some stats
	fmt.Printf("Successfully parsed JMNedict file\n")
	fmt.Printf("Number of names: %d\n", len(dict.Names))

	// Print the actual data structure
	fmt.Printf("JMNedictTypes struct contents:\n")
	fmt.Printf("Version: %s\n", dict.Version)
	fmt.Printf("Languages: %v\n", dict.Languages)
	fmt.Printf("CommonOnly: %v\n", dict.CommonOnly)
	fmt.Printf("DictDate: %s\n", dict.DictDate)
	fmt.Printf("DictRevisions: %v\n", dict.DictRevisions)
	fmt.Printf("Tags count: %d\n", len(dict.Tags))

	// Print first 5 names as examples
	if len(dict.Names) > 0 {
		fmt.Printf("\nSample entries:\n")
		for i, name := range dict.Names {
			if i >= 5 {
				break
			}
			fmt.Printf("%d. ID: %s, Kanji: %v, Reading: %v\n", i+1, name.ID, name.Kanji, name.Reading)
		}
	} else {
		fmt.Printf("\nWARNING: No name entries found in the JMNedict file!\n")

		// Try to manually extract some of the structure to see what's there
		file.Seek(0, 0)
		var rawData map[string]interface{}
		if err := json.NewDecoder(file).Decode(&rawData); err != nil {
			fmt.Printf("Error decoding raw JSON: %v\n", err)
		} else {
			fmt.Printf("\nRaw JSON structure:\n")
			for key, value := range rawData {
				fmt.Printf("- Key: %s, Type: %T\n", key, value)

				// If we find an array that might contain the names, analyze first element
				if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
					fmt.Printf("  First element type: %T\n", arr[0])
					if obj, ok := arr[0].(map[string]interface{}); ok {
						fmt.Printf("  First element keys: %v\n", getMapKeys(obj))
					}
				}
			}
		}
	}
}

// Helper function to get keys from a map
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
