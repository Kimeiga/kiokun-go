package testdata

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"kiokun-go/dictionaries/jmdict"
)

// FieldFinder helps find entries in the JMdict that use specific field types
type FieldFinder struct {
	dictPath string
}

// NewFieldFinder creates a new FieldFinder for the specified dictionary file
func NewFieldFinder(dictPath string) *FieldFinder {
	return &FieldFinder{
		dictPath: dictPath,
	}
}

// FindEntriesWithField processes the dictionary file in chunks to find entries with a specific field
// without loading the entire dictionary into memory
func (ff *FieldFinder) FindEntriesWithField(field, value string, maxResults int) ([]jmdict.Word, error) {
	file, err := os.Open(ff.dictPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dictionary file: %w", err)
	}
	defer file.Close()

	// Read the file header to get past the root object opening
	header := make([]byte, 1024)
	_, err = file.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}

	// Find the start of the "words" array
	wordsStart := strings.Index(string(header), `"words":[`)
	if wordsStart == -1 {
		return nil, fmt.Errorf("couldn't find words array in header")
	}

	// Seek to the position after "words":[
	_, err = file.Seek(int64(wordsStart)+9, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to words array: %w", err)
	}

	// Process in chunks to avoid loading the entire file
	results := []jmdict.Word{}
	decoder := json.NewDecoder(file)

	// Read one word at a time
	for {
		if len(results) >= maxResults {
			break
		}

		var word jmdict.Word
		err := decoder.Decode(&word)
		if err == io.EOF {
			break
		}
		if err != nil {
			// Skip invalid entries
			continue
		}

		// Check if this word has the desired field/value
		if ff.wordHasField(&word, field, value) {
			results = append(results, word)
		}

		// Read past the comma (if any)
		// This is a simplification and may need adjustment based on the exact JSON format
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if token == json.Delim(']') {
			break
		}
	}

	return results, nil
}

// wordHasField checks if a word contains the specified field and value
func (ff *FieldFinder) wordHasField(word *jmdict.Word, field, value string) bool {
	switch field {
	case "partOfSpeech":
		return ff.hasPartOfSpeech(word, value)
	case "dialect":
		return ff.hasDialect(word, value)
	case "field":
		return ff.hasField(word, value)
	case "misc":
		return ff.hasMisc(word, value)
	case "languageSource":
		return ff.hasLanguageSource(word, value)
	case "gloss.gender":
		return ff.hasGlossGender(word, value)
	case "gloss.type":
		return ff.hasGlossType(word, value)
	case "example":
		return ff.hasExample(word)
	default:
		return false
	}
}

// Helper functions to check for specific fields
func (ff *FieldFinder) hasPartOfSpeech(word *jmdict.Word, pos string) bool {
	for _, sense := range word.Sense {
		for _, p := range sense.PartOfSpeech {
			if p == pos {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasDialect(word *jmdict.Word, dialect string) bool {
	for _, sense := range word.Sense {
		for _, d := range sense.Dialect {
			if d == dialect {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasField(word *jmdict.Word, field string) bool {
	for _, sense := range word.Sense {
		for _, f := range sense.Field {
			if f == field {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasMisc(word *jmdict.Word, misc string) bool {
	for _, sense := range word.Sense {
		for _, m := range sense.Misc {
			if m == misc {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasLanguageSource(word *jmdict.Word, lang string) bool {
	for _, sense := range word.Sense {
		for _, ls := range sense.LanguageSource {
			if string(ls.Lang) == lang {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasGlossGender(word *jmdict.Word, gender string) bool {
	for _, sense := range word.Sense {
		for _, g := range sense.Gloss {
			if g.Gender != nil && *g.Gender == gender {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasGlossType(word *jmdict.Word, glossType string) bool {
	for _, sense := range word.Sense {
		for _, g := range sense.Gloss {
			if g.Type != nil && *g.Type == glossType {
				return true
			}
		}
	}
	return false
}

func (ff *FieldFinder) hasExample(word *jmdict.Word) bool {
	for _, sense := range word.Sense {
		if len(sense.Examples) > 0 {
			return true
		}
	}
	return false
}

// WriteEntriesToFile writes the found entries to a test file
func (ff *FieldFinder) WriteEntriesToFile(entries []jmdict.Word, outputPath string) error {
	// Create a JmdictTypes container
	dict := jmdict.JmdictTypes{
		Version:       "test-field-finder",
		Languages:     []string{"eng"},
		CommonOnly:    false,
		DictDate:      "2023-01-01",
		DictRevisions: []string{"1.0"},
		Tags:          make(map[string]string),
		Words:         entries,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(dict, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Write to file
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// SimpleToolToFindExamples is a simple tool to extract examples for testing
func SimpleToolToFindExamples() error {
	sourceDir := filepath.Join("dictionaries", "jmdict", "source")
	pattern := `^jmdict-examples-eng-.*\.json$`

	filename, err := filepath.Glob(filepath.Join(sourceDir, pattern))
	if err != nil || len(filename) == 0 {
		return fmt.Errorf("failed to find dictionary file: %w", err)
	}

	finder := NewFieldFinder(filename[0])

	// Find different types of entries
	partOfSpeechEntries, err := finder.FindEntriesWithField("partOfSpeech", "adj-i", 2)
	if err != nil {
		return err
	}

	fieldEntries, err := finder.FindEntriesWithField("field", "comp", 2)
	if err != nil {
		return err
	}

	dialectEntries, err := finder.FindEntriesWithField("dialect", "ksb", 2)
	if err != nil {
		return err
	}

	miscEntries, err := finder.FindEntriesWithField("misc", "arch", 2)
	if err != nil {
		return err
	}

	exampleEntries, err := finder.FindEntriesWithField("example", "", 2)
	if err != nil {
		return err
	}

	// Combine all entries
	allEntries := []jmdict.Word{}
	allEntries = append(allEntries, partOfSpeechEntries...)
	allEntries = append(allEntries, fieldEntries...)
	allEntries = append(allEntries, dialectEntries...)
	allEntries = append(allEntries, miscEntries...)
	allEntries = append(allEntries, exampleEntries...)

	// Write to test file
	outputPath := filepath.Join("dictionaries", "jmdict", "testdata", "field_examples.json")
	return finder.WriteEntriesToFile(allEntries, outputPath)
}
