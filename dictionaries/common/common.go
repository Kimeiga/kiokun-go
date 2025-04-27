package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// Entry represents a generic dictionary entry
type Entry interface {
	GetID() string
	GetFilename() string
}

// DictionaryImporter defines the interface for dictionary importers
type DictionaryImporter interface {
	Name() string
	Import(inputPath string) ([]Entry, error)
}

// DictionaryConfig holds configuration for a dictionary
type DictionaryConfig struct {
	Name      string
	SourceDir string
	InputFile string
	Importer  DictionaryImporter
}

var (
	registeredDictionaries []DictionaryConfig
	dictionariesBasePath   = "dictionaries" // default path
)

// SetDictionariesBasePath sets the base path for all dictionary operations
func SetDictionariesBasePath(path string) {
	dictionariesBasePath = path
	// Update source directories for all registered dictionaries
	for i := range registeredDictionaries {
		registeredDictionaries[i].SourceDir = filepath.Join(dictionariesBasePath, registeredDictionaries[i].Name, "source")
	}
}

// RegisterDictionary adds a dictionary to the registry
func RegisterDictionary(name, inputFile string, importer DictionaryImporter) {
	sourceDir := filepath.Join(dictionariesBasePath, name, "source")
	config := DictionaryConfig{
		Name:      name,
		SourceDir: sourceDir,
		InputFile: inputFile,
		Importer:  importer,
	}
	registeredDictionaries = append(registeredDictionaries, config)
}

// GetRegisteredDictionaries returns a copy of the registered dictionaries slice
// This function is primarily used for testing purposes
func GetRegisteredDictionaries() []DictionaryConfig {
	result := make([]DictionaryConfig, len(registeredDictionaries))
	copy(result, registeredDictionaries)
	return result
}

// ImportAllDictionaries imports all registered dictionaries
func ImportAllDictionaries() ([]Entry, error) {
	var allEntries []Entry

	for _, dict := range registeredDictionaries {
		inputPath := filepath.Join(dict.SourceDir, dict.InputFile)
		file, err := os.Open(inputPath)
		if err != nil {
			return nil, fmt.Errorf("error importing %s: %v", dict.Name, err)
		}
		defer file.Close()

		entries, err := dict.Importer.Import(inputPath)
		if err != nil {
			return nil, fmt.Errorf("error importing %s: %v", dict.Name, err)
		}
		allEntries = append(allEntries, entries...)
	}

	return allEntries, nil
}

// FindDictionaryFile finds a dictionary file in the given directory using regex pattern
func FindDictionaryFile(dir, pattern string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && re.MatchString(entry.Name()) {
			return entry.Name(), nil
		}
	}

	return "", os.ErrNotExist
}

// ImportJSON is a helper function to decode JSON data
func ImportJSON(reader io.Reader, v interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(v)
}
