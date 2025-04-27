package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Parse command line arguments
	inputFile := flag.String("input", "", "Input JSONL file path (required)")
	outputFile := flag.String("output", "", "Output JSON file path (optional)")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Input file is required")
		fmt.Println("Usage: go run main.go -input=<jsonl-file> [-output=<json-file>]")
		os.Exit(1)
	}

	// If output file not specified, use input filename with .json extension
	if *outputFile == "" {
		base := filepath.Base(*inputFile)
		ext := filepath.Ext(base)
		nameWithoutExt := strings.TrimSuffix(base, ext)
		*outputFile = nameWithoutExt + ".json"
	}

	// Open the input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a scanner for reading the file line by line
	scanner := bufio.NewScanner(file)

	// Prepare an array to hold all JSON objects
	var jsonObjects []json.RawMessage

	// Process each line
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse the JSON line
		var rawJson json.RawMessage
		err := json.Unmarshal([]byte(line), &rawJson)
		if err != nil {
			fmt.Printf("Error parsing JSON at line %d: %v\n", lineCount, err)
			continue
		}

		jsonObjects = append(jsonObjects, rawJson)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Marshal the array of objects to JSON
	outputJSON, err := json.Marshal(jsonObjects)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Write the JSON to the output file
	err = os.WriteFile(*outputFile, outputJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %d JSON lines to %s\n", len(jsonObjects), *outputFile)
}
