package jmdict

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kiokun-go/dictionaries/common"
)

func TestUnmarshalJmdictTypes(t *testing.T) {
	// Test basic unmarshaling of JmdictTypes
	testData := []byte(`{
		"version": "test",
		"languages": ["eng"],
		"commonOnly": false,
		"dictDate": "2023-01-01",
		"dictRevisions": ["1.0"],
		"tags": {"adj-i": "i-adjective"},
		"words": []
	}`)

	result, err := UnmarshalJmdictTypes(testData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JmdictTypes: %v", err)
	}

	// Check basic fields
	if result.Version != "test" {
		t.Errorf("Expected version 'test', got '%s'", result.Version)
	}
	if len(result.Languages) != 1 || result.Languages[0] != "eng" {
		t.Errorf("Languages field incorrect, got %v", result.Languages)
	}
	if result.CommonOnly {
		t.Errorf("Expected CommonOnly to be false")
	}
	if result.DictDate != "2023-01-01" {
		t.Errorf("Expected dictDate '2023-01-01', got '%s'", result.DictDate)
	}
	if len(result.DictRevisions) != 1 || result.DictRevisions[0] != "1.0" {
		t.Errorf("DictRevisions field incorrect, got %v", result.DictRevisions)
	}
	if len(result.Tags) != 1 || result.Tags["adj-i"] != "i-adjective" {
		t.Errorf("Tags field incorrect, got %v", result.Tags)
	}
	if len(result.Words) != 0 {
		t.Errorf("Expected empty words array, got %d items", len(result.Words))
	}
}

func TestJmdictTypesMarshal(t *testing.T) {
	// Test marshaling of JmdictTypes
	dict := JmdictTypes{
		Version:       "test",
		Languages:     []string{"eng"},
		CommonOnly:    false,
		DictDate:      "2023-01-01",
		DictRevisions: []string{"1.0"},
		Tags:          map[string]string{"adj-i": "i-adjective"},
		Words:         []Word{},
	}

	data, err := dict.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal JmdictTypes: %v", err)
	}

	// Unmarshal back to verify
	var result JmdictTypes
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Check basic fields
	if result.Version != "test" {
		t.Errorf("Expected version 'test', got '%s'", result.Version)
	}
}

func TestImporter(t *testing.T) {
	importer := &Importer{}

	// Test Name() function
	if importer.Name() != "jmdict" {
		t.Errorf("Expected importer name 'jmdict', got '%s'", importer.Name())
	}
}

func TestImporterImport(t *testing.T) {
	importer := &Importer{}
	testFile := filepath.Join("testdata", "test_jmdict.json")

	entries, err := importer.Import(testFile)
	if err != nil {
		t.Fatalf("Failed to import test data: %v", err)
	}

	// Verify we got the expected number of entries
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Verify the entries implement the Entry interface
	for i, entry := range entries {
		// Check ID
		id := entry.GetID()
		if id == "" {
			t.Errorf("Entry %d has empty ID", i)
		}

		// Check filename
		filename := entry.GetFilename()
		if filename == "" {
			t.Errorf("Entry %d has empty filename", i)
		}

		// Cast to Word and check specific fields
		word, ok := entry.(Word)
		if !ok {
			t.Errorf("Entry %d is not a Word type", i)
			continue
		}

		// Verify basics of the first entry
		if i == 0 {
			if word.ID != "1000001" {
				t.Errorf("First entry ID expected '1000001', got '%s'", word.ID)
			}

			if len(word.Kanji) == 0 || word.Kanji[0].Text != "食べる" {
				t.Errorf("First entry Kanji text incorrect")
			}

			if len(word.Kana) == 0 || word.Kana[0].Text != "たべる" {
				t.Errorf("First entry Kana text incorrect")
			}

			if len(word.Sense) == 0 {
				t.Errorf("First entry has no senses")
			} else {
				sense := word.Sense[0]
				if len(sense.PartOfSpeech) < 2 || sense.PartOfSpeech[0] != "v1" || sense.PartOfSpeech[1] != "vt" {
					t.Errorf("First entry PartOfSpeech incorrect, got %v", sense.PartOfSpeech)
				}

				if len(sense.Gloss) == 0 || sense.Gloss[0].Text != "to eat" {
					t.Errorf("First entry Gloss text incorrect")
				}

				if len(sense.Examples) == 0 {
					t.Errorf("First entry has no examples")
				} else {
					example := sense.Examples[0]
					if example.Text != "ごはんを食べる" {
						t.Errorf("Example text incorrect, got '%s'", example.Text)
					}

					if len(example.Sentences) == 0 || example.Sentences[0].Text != "I eat rice." {
						t.Errorf("Example sentence text incorrect")
					}
				}
			}
		}
	}
}

func TestWordInterface(t *testing.T) {
	// Test the Word implementation of Entry interface
	word := Word{
		ID: "test-id",
		Kanji: []KanjiEntry{
			{Text: "漢字", Common: true},
		},
		Kana: []KanaEntry{
			{Text: "かんじ", Common: true},
		},
	}

	// Test GetID
	if id := word.GetID(); id != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", id)
	}

	// Test GetFilename - should return first kanji
	if filename := word.GetFilename(); filename != "漢字" {
		t.Errorf("Expected filename '漢字', got '%s'", filename)
	}

	// Test GetFilename with no kanji
	wordNoKanji := Word{
		ID: "test-id",
		Kana: []KanaEntry{
			{Text: "かんじ", Common: true},
		},
	}

	if filename := wordNoKanji.GetFilename(); filename != "かんじ" {
		t.Errorf("Expected filename 'かんじ', got '%s'", filename)
	}

	// Test GetFilename with no kanji or kana
	wordEmpty := Word{
		ID: "test-id",
	}

	if filename := wordEmpty.GetFilename(); filename != "test-id" {
		t.Errorf("Expected filename 'test-id', got '%s'", filename)
	}
}

func TestFieldTypes(t *testing.T) {
	// Read the test file
	testFile := filepath.Join("testdata", "test_jmdict.json")
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	var dict JmdictTypes
	if err := common.ImportJSON(file, &dict); err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	// Test the second word which has more field types
	if len(dict.Words) < 2 {
		t.Fatalf("Not enough test words in the file")
	}

	word := dict.Words[1] // Test the second word (美しい)

	// Check antonym field
	if len(word.Sense) == 0 || len(word.Sense[0].Antonym) == 0 {
		t.Errorf("Expected antonym data not found")
	} else {
		antonym := word.Sense[0].Antonym[0]
		// Check the antonym which should be a Xref with String field
		if antonym.String == nil {
			t.Errorf("Expected antonym to have String field")
		} else if *antonym.String != "醜い" {
			t.Errorf("Expected antonym '醜い' not found, got '%s'", *antonym.String)
		}
	}

	// Check field data
	if len(word.Sense) == 0 || len(word.Sense[0].Field) == 0 {
		t.Errorf("Expected field data not found")
	} else {
		field := word.Sense[0].Field[0]
		if field != "art" {
			t.Errorf("Expected field 'art' not found, got '%s'", field)
		}
	}

	// Check info field
	if len(word.Sense) == 0 || len(word.Sense[0].Info) == 0 {
		t.Errorf("Expected info data not found")
	} else {
		if word.Sense[0].Info[0] != "used to describe natural beauty" {
			t.Errorf("Expected info not found")
		}
	}

	// Check gloss with gender and type
	if len(word.Sense) == 0 || len(word.Sense[0].Gloss) < 2 {
		t.Errorf("Expected gloss data not found")
	} else {
		gloss := word.Sense[0].Gloss[1]
		if gloss.Gender == nil || *gloss.Gender != "neuter" {
			t.Errorf("Expected gender 'neuter' not found")
		}
		if gloss.Type == nil || *gloss.Type != "literal" {
			t.Errorf("Expected type 'literal' not found")
		}
	}

	// Test the third word (日本語) for languageSource
	word = dict.Words[2]
	if len(word.Sense) == 0 || len(word.Sense[0].LanguageSource) == 0 {
		t.Errorf("Expected languageSource data not found")
	} else {
		langSrc := word.Sense[0].LanguageSource[0]
		if langSrc.Lang != "jpn" {
			t.Errorf("Expected lang 'jpn', got '%s'", langSrc.Lang)
		}
		if !langSrc.Full {
			t.Errorf("Expected Full to be true")
		}
		if langSrc.Wasei {
			t.Errorf("Expected Wasei to be false")
		}
		if langSrc.Text == nil || *langSrc.Text != "日本語" {
			t.Errorf("Expected text '日本語', got '%v'", langSrc.Text)
		}
	}
}

func TestInitRegistration(t *testing.T) {
	// This is a test for the init.go file which registers the dictionary
	// Since the init function has already run, we can test that the dictionary was registered correctly

	// Get the registered dictionary from common package
	dictionaries := common.GetRegisteredDictionaries()
	found := false

	for _, dict := range dictionaries {
		if dict.Name == "jmdict" {
			found = true

			// Check that the importer is the correct type
			_, ok := dict.Importer.(*Importer)
			if !ok {
				t.Errorf("Registered jmdict importer is not of type *Importer")
			}

			// Check source directory is set correctly
			expectedSourceDir := filepath.Join("dictionaries", "jmdict", "source")
			if !strings.HasSuffix(dict.SourceDir, expectedSourceDir) {
				t.Errorf("Source directory incorrect. Expected suffix %s, got %s",
					expectedSourceDir, dict.SourceDir)
			}

			break
		}
	}

	if !found {
		t.Errorf("jmdict dictionary not found in registered dictionaries")
	}
}
