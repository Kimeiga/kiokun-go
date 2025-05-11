package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run preprocess_json.go <json-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Read the JSON file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON array
	var entries []map[string]interface{}
	if err := json.Unmarshal(data, &entries); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Process each entry
	for i, entry := range entries {
		// Convert strokeCount from string to int if it exists
		if strokeCount, ok := entry["strokeCount"]; ok {
			switch sc := strokeCount.(type) {
			case string:
				// Convert string to int
				if intVal, err := strconv.Atoi(sc); err == nil {
					entry["strokeCount"] = intVal
				} else {
					fmt.Printf("Warning: Could not convert strokeCount '%s' to int at entry %d\n", sc, i)
				}
			case float64:
				// Already a number, convert to int
				entry["strokeCount"] = int(sc)
			}
		}
	}

	// Marshal back to JSON
	outputData, err := json.Marshal(entries)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Write back to the file
	if err := ioutil.WriteFile(filePath, outputData, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully preprocessed %d entries in %s\n", len(entries), filePath)
}
