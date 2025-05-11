package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/kanjidic"

	"github.com/andybalholm/brotli"
)

// TestIndexProcessor tests the index-based processor with exact and contained-in matches
func TestIndexProcessor(t *testing.T) {
	// Create a temporary output directory
	testDir := filepath.Join(os.TempDir(), "kiokun_index_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create an index processor
	proc, err := NewIndexProcessor(testDir, 1)
	if err != nil {
		t.Fatalf("Failed to create index processor: %v", err)
	}

	// Create test entries for characters that should exist in both Chinese and Japanese
	testChars := []string{"人", "日", "月"} // Person, Day/Sun, Month/Moon

	// Create compound words that contain these characters
	testWords := []struct {
		kanji    string
		kana     string
		meaning  string
		contains string // which test character it contains
	}{
		{"人間", "にんげん", "human being", "人"},
		{"日本", "にほん", "Japan", "日"},
		{"月曜日", "げつようび", "Monday", "月"},
	}

	// Process each test character
	for _, char := range testChars {
		// Create a Japanese kanji entry
		kanji := kanjidic.Kanji{
			Character: char,
			Meanings:  []string{"Test meaning for " + char},
			OnYomi:    []string{"on"},
			KunYomi:   []string{"kun"},
			Stroke:    5,
		}

		// Create a Chinese character entry
		chineseChar := chinese_chars.ChineseCharEntry{
			ID:          "test_char_" + char,
			Traditional: char,
			Simplified:  char,
			Definitions: []string{"Test definition for " + char},
			Pinyin:      []string{"test"},
		}

		// Process both entries
		entries := []common.Entry{kanji, chineseChar}
		if err := proc.ProcessEntries(entries); err != nil {
			t.Fatalf("Failed to process entries for %s: %v", char, err)
		}
	}

	// Process compound words
	for i, word := range testWords {
		// Create a JMdict entry
		jmdictEntry := jmdict.Word{
			ID: "test_word_" + word.kanji,
			Kanji: []jmdict.KanjiEntry{
				{Text: word.kanji, Common: true},
			},
			Kana: []jmdict.KanaEntry{
				{Text: word.kana, Common: true, AppliesToKanji: []string{"*"}},
			},
			Sense: []jmdict.Sense{
				{
					PartOfSpeech:   []string{"n"},
					AppliesToKanji: []string{"*"},
					AppliesToKana:  []string{"*"},
					Gloss: []jmdict.Gloss{
						{Lang: "eng", Text: word.meaning},
					},
				},
			},
		}

		// Create a Chinese word entry
		chineseWord := chinese_words.ChineseWordEntry{
			ID:          "test_cword_" + word.kanji,
			Traditional: word.kanji,
			Simplified:  word.kanji,
			Definitions: []string{"Test definition for " + word.kanji},
			Pinyin:      []string{"test" + fmt.Sprintf("%d", i)},
		}

		// Process both entries
		entries := []common.Entry{jmdictEntry, chineseWord}
		if err := proc.ProcessEntries(entries); err != nil {
			t.Fatalf("Failed to process entries for %s: %v", word.kanji, err)
		}
	}

	// Write to files
	if err := proc.WriteToFiles(); err != nil {
		t.Fatalf("Failed to write files: %v", err)
	}

	// Verify that index files were created for single characters
	for _, char := range testChars {
		indexPath := filepath.Join(testDir, "index", char+".json.br")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			t.Errorf("Expected index file %s does not exist", indexPath)
			continue
		}

		// Read and check the index file
		indexEntry, err := readIndexEntry(indexPath)
		if err != nil {
			t.Errorf("Failed to read index file %s: %v", indexPath, err)
			continue
		}

		// Check that the index entry has exact matches for single characters
		if len(indexEntry.De) == 0 {
			t.Errorf("Index entry for %s has no Kanjidic exact matches", char)
		}
		if len(indexEntry.Ce) == 0 {
			t.Errorf("Index entry for %s has no Chinese character exact matches", char)
		}

		// Check that the index entry has contained-in matches for compound words
		// Find words that contain this character
		var containingWords []string
		for _, word := range testWords {
			if word.contains == char {
				containingWords = append(containingWords, word.kanji)
			}
		}

		if len(containingWords) > 0 {
			if len(indexEntry.Jc) == 0 {
				t.Errorf("Index entry for %s should have JMdict contained-in matches", char)
			}
			if len(indexEntry.Wc) == 0 {
				t.Errorf("Index entry for %s should have Chinese word contained-in matches", char)
			}
		}

		t.Logf("Index entry for %s: %+v", char, indexEntry)

		// Verify that dictionary files were created for exact matches
		for _, kanjiID := range indexEntry.De {
			kanjiPath := filepath.Join(testDir, "d", fmt.Sprintf("%d.json.br", kanjiID))
			if _, err := os.Stat(kanjiPath); os.IsNotExist(err) {
				t.Errorf("Expected kanjidic file %s does not exist", kanjiPath)
			}
		}

		for _, charID := range indexEntry.Ce {
			charPath := filepath.Join(testDir, "c", fmt.Sprintf("%d.json.br", charID))
			if _, err := os.Stat(charPath); os.IsNotExist(err) {
				t.Errorf("Expected Chinese character file %s does not exist", charPath)
			}
		}

		// Verify that dictionary files were created for contained-in matches
		for _, wordID := range indexEntry.Jc {
			wordPath := filepath.Join(testDir, "j", fmt.Sprintf("%d.json.br", wordID))
			if _, err := os.Stat(wordPath); os.IsNotExist(err) {
				t.Errorf("Expected JMdict file %s does not exist", wordPath)
			}
		}

		for _, wordID := range indexEntry.Wc {
			wordPath := filepath.Join(testDir, "w", fmt.Sprintf("%d.json.br", wordID))
			if _, err := os.Stat(wordPath); os.IsNotExist(err) {
				t.Errorf("Expected Chinese word file %s does not exist", wordPath)
			}
		}
	}

	// Verify that index files were created for compound words
	for _, word := range testWords {
		indexPath := filepath.Join(testDir, "index", word.kanji+".json.br")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			t.Errorf("Expected index file %s does not exist", indexPath)
			continue
		}

		// Read and check the index file
		indexEntry, err := readIndexEntry(indexPath)
		if err != nil {
			t.Errorf("Failed to read index file %s: %v", indexPath, err)
			continue
		}

		// Check that the index entry has exact matches for compound words
		if len(indexEntry.Je) == 0 {
			t.Errorf("Index entry for %s has no JMdict exact matches", word.kanji)
		}
		if len(indexEntry.We) == 0 {
			t.Errorf("Index entry for %s has no Chinese word exact matches", word.kanji)
		}

		// Compound words should not have contained-in matches for themselves
		if len(indexEntry.Jc) > 0 {
			t.Errorf("Index entry for %s should not have JMdict contained-in matches", word.kanji)
		}
		if len(indexEntry.Wc) > 0 {
			t.Errorf("Index entry for %s should not have Chinese word contained-in matches", word.kanji)
		}

		t.Logf("Index entry for %s: %+v", word.kanji, indexEntry)

		// Verify that dictionary files were created for exact matches
		for _, wordID := range indexEntry.Je {
			wordPath := filepath.Join(testDir, "j", fmt.Sprintf("%d.json.br", wordID))
			if _, err := os.Stat(wordPath); os.IsNotExist(err) {
				t.Errorf("Expected JMdict file %s does not exist", wordPath)
			}
		}

		for _, wordID := range indexEntry.We {
			wordPath := filepath.Join(testDir, "w", fmt.Sprintf("%d.json.br", wordID))
			if _, err := os.Stat(wordPath); os.IsNotExist(err) {
				t.Errorf("Expected Chinese word file %s does not exist", wordPath)
			}
		}
	}
}

// readIndexEntry reads an IndexEntry from a compressed JSON file
func readIndexEntry(filePath string) (*IndexEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decompress the file
	br := brotli.NewReader(file)
	decoder := json.NewDecoder(br)

	var entry IndexEntry
	if err := decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}
