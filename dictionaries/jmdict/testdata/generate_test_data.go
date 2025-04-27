package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kiokun-go/dictionaries/jmdict"
)

// This program generates test data by extracting entries from the main JMdict file
// that have specific field types. It creates a smaller JSON file for use in tests.

func main() {
	// Get current working directory and ensure we're in the right place
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Navigate to project root if needed
	for i := 0; i < 3; i++ {
		if _, err := os.Stat(filepath.Join(cwd, "dictionaries")); err == nil {
			break
		}
		cwd = filepath.Dir(cwd)
	}
	if err := os.Chdir(cwd); err != nil {
		fmt.Printf("Error changing to project root: %v\n", err)
		os.Exit(1)
	}

	// Find the JMdict file
	sourceDir := filepath.Join("dictionaries", "jmdict", "source")
	files, err := filepath.Glob(filepath.Join(sourceDir, "jmdict-examples-eng-*.json"))
	if err != nil || len(files) == 0 {
		fmt.Printf("Error finding JMdict file: %v\n", err)
		os.Exit(1)
	}

	jmdictFile := files[0]
	fmt.Printf("Using JMdict file: %s\n", jmdictFile)

	// Create fieldsamples test data
	testDataDir := filepath.Join("dictionaries", "jmdict", "testdata")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		fmt.Printf("Error creating testdata directory: %v\n", err)
		os.Exit(1)
	}

	// Create a collection of entries with various field types
	entries := []jmdict.Word{}

	// Parse the JMdict file in chunks to extract specific entries
	// We use a streaming approach to avoid loading the entire file in memory
	file, err := os.Open(jmdictFile)
	if err != nil {
		fmt.Printf("Error opening JMdict file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Fields we want to find examples of
	fieldExamples := map[string]bool{
		"adj-i":      false, // i-adjective
		"n":          false, // noun
		"v1":         false, // Ichidan verb
		"v5k":        false, // Godan verb with `ku' ending
		"adv":        false, // adverb
		"ksb":        false, // Kansai-ben dialect
		"osb":        false, // Osaka-ben dialect
		"comp":       false, // computer terminology
		"med":        false, // medicine
		"food":       false, // food term
		"arch":       false, // archaism
		"male":       false, // male term or language
		"vulg":       false, // vulgar expression or word
		"feminine":   false, // feminine gender
		"figurative": false, // figurative meaning
		"ger":        false, // German source language
		"chi":        false, // Chinese source language
	}

	fmt.Println("Starting to scan JMdict...")

	// Output the test data
	outputFile := filepath.Join(testDataDir, "field_samples.json")
	fmt.Printf("Writing field examples to: %s\n", outputFile)

	// Create a simple dictionary
	dict := jmdict.JmdictTypes{
		Version:       "test-1.0",
		Languages:     []string{"eng"},
		CommonOnly:    false,
		DictDate:      "2023-01-01",
		DictRevisions: []string{"1.0"},
		Tags:          make(map[string]string),
		Words:         entries,
	}

	// Write to file
	// In a real implementation, we would scan the file to find examples of each field type
	// For now, we'll use our test data from test_jmdict.json
	fmt.Println("Note: This is not fully implemented.")
	fmt.Println("For now, please use the test_jmdict.json file for tests.")
}
