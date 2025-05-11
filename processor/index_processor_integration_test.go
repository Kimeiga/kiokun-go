package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"

	"github.com/andybalholm/brotli"
)

// TestIndexProcessorIntegration tests the index-based processor with a larger dataset
// and verifies the content of both index files and dictionary files
func TestIndexProcessorIntegration(t *testing.T) {
	// Create a temporary output directory
	testDir := filepath.Join(os.TempDir(), "kiokun_index_integration_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create an index processor
	proc, err := NewIndexProcessor(testDir, 4)
	if err != nil {
		t.Fatalf("Failed to create index processor: %v", err)
	}

	// Create a larger test dataset
	var allEntries []common.Entry

	// Add Japanese kanji entries
	kanjiChars := []string{"人", "日", "月", "水", "火", "木", "金", "土", "山", "川"}
	for i, char := range kanjiChars {
		kanji := kanjidic.Kanji{
			Character: char,
			Meanings:  []string{"Test meaning for " + char},
			OnYomi:    []string{"on" + fmt.Sprintf("%d", i)},
			KunYomi:   []string{"kun" + fmt.Sprintf("%d", i)},
			Stroke:    5 + i,
		}
		allEntries = append(allEntries, kanji)
	}

	// Add Chinese character entries
	for i, char := range kanjiChars {
		chineseChar := chinese_chars.ChineseCharEntry{
			ID:          "test_char_" + char,
			Traditional: char,
			Simplified:  char,
			Definitions: []string{"Test definition for " + char},
			Pinyin:      []string{"test" + fmt.Sprintf("%d", i)},
		}
		allEntries = append(allEntries, chineseChar)
	}

	// Add JMdict entries (words)
	jmdictWords := []struct {
		id    string
		kanji string
		kana  string
		sense string
	}{
		{"jmdict_1", "人間", "にんげん", "human being"},
		{"jmdict_2", "日本", "にほん", "Japan"},
		{"jmdict_3", "月曜日", "げつようび", "Monday"},
		{"jmdict_4", "火山", "かざん", "volcano"},
		{"jmdict_5", "木曜日", "もくようび", "Thursday"},
	}

	for _, word := range jmdictWords {
		jmdictEntry := jmdict.Word{
			ID: word.id,
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
						{Lang: "eng", Text: word.sense},
					},
				},
			},
		}
		allEntries = append(allEntries, jmdictEntry)
	}

	// Add JMNedict entries (names)
	jmnedictNames := []struct {
		id      string
		kanji   string
		reading string
		type_   string
	}{
		{"jmnedict_1", "山田", "やまだ", "surname"},
		{"jmnedict_2", "川上", "かわかみ", "surname"},
		{"jmnedict_3", "金子", "かねこ", "surname"},
		{"jmnedict_4", "東京", "とうきょう", "place"},
		{"jmnedict_5", "大阪", "おおさか", "place"},
	}

	for _, name := range jmnedictNames {
		jmnedictEntry := jmnedict.Name{
			ID:      name.id,
			Kanji:   []string{name.kanji},
			Reading: []string{name.reading},
			Type:    []string{name.type_},
		}
		allEntries = append(allEntries, jmnedictEntry)
	}

	// Process all entries
	t.Logf("Processing %d test entries", len(allEntries))
	if err := proc.ProcessEntries(allEntries); err != nil {
		t.Fatalf("Failed to process entries: %v", err)
	}

	// Write to files
	if err := proc.WriteToFiles(); err != nil {
		t.Fatalf("Failed to write files: %v", err)
	}

	// Verify index files for kanji characters
	t.Logf("Verifying index files for kanji characters")
	for _, char := range kanjiChars {
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

		// Check that the index entry has both Japanese and Chinese IDs
		if len(indexEntry.De) == 0 {
			t.Errorf("Index entry for %s has no Kanjidic exact matches", char)
		}
		if len(indexEntry.Ce) == 0 {
			t.Errorf("Index entry for %s has no Chinese character exact matches", char)
		}

		t.Logf("Index entry for %s: %+v", char, indexEntry)

		// Verify that dictionary files were created and contain correct data
		for _, kanjiID := range indexEntry.De {
			kanjiPath := filepath.Join(testDir, "d", fmt.Sprintf("%d.json.br", kanjiID))
			if _, err := os.Stat(kanjiPath); os.IsNotExist(err) {
				t.Errorf("Expected kanjidic file %s does not exist", kanjiPath)
				continue
			}

			// Read and verify the kanji file
			kanjiEntry, err := readKanjiEntry(kanjiPath)
			if err != nil {
				t.Errorf("Failed to read kanji file %s: %v", kanjiPath, err)
				continue
			}

			if kanjiEntry.Character != char {
				t.Errorf("Kanji file %s has incorrect character: expected %s, got %s",
					kanjiPath, char, kanjiEntry.Character)
			}

			t.Logf("Kanji entry for %d: Character=%s, Meanings=%v",
				kanjiID, kanjiEntry.Character, kanjiEntry.Meanings)
		}

		for _, charID := range indexEntry.Ce {
			charPath := filepath.Join(testDir, "c", fmt.Sprintf("%d.json.br", charID))
			if _, err := os.Stat(charPath); os.IsNotExist(err) {
				t.Errorf("Expected Chinese character file %s does not exist", charPath)
				continue
			}

			// Read and verify the Chinese character file
			chineseEntry, err := readChineseCharEntry(charPath)
			if err != nil {
				t.Errorf("Failed to read Chinese character file %s: %v", charPath, err)
				continue
			}

			if chineseEntry.Traditional != char {
				t.Errorf("Chinese character file %s has incorrect character: expected %s, got %s",
					charPath, char, chineseEntry.Traditional)
			}

			t.Logf("Chinese character entry for %d: Traditional=%s, Definitions=%v",
				charID, chineseEntry.Traditional, chineseEntry.Definitions)
		}
	}

	// Verify index files for compound words
	t.Logf("Verifying index files for compound words")
	for _, word := range jmdictWords {
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

		// Check that the index entry has JMdict IDs
		if len(indexEntry.J) == 0 {
			t.Errorf("Index entry for %s has no JMdict IDs", word.kanji)
		}

		t.Logf("Index entry for %s: %+v", word.kanji, indexEntry)

		// Verify that JMdict files were created and contain correct data
		for _, jmdictID := range indexEntry.Je {
			jmdictPath := filepath.Join(testDir, "j", fmt.Sprintf("%d.json.br", jmdictID))
			if _, err := os.Stat(jmdictPath); os.IsNotExist(err) {
				t.Errorf("Expected JMdict file %s does not exist", jmdictPath)
				continue
			}

			// Read and verify the JMdict file
			jmdictEntry, err := readJMdictEntry(jmdictPath)
			if err != nil {
				t.Errorf("Failed to read JMdict file %s: %v", jmdictPath, err)
				continue
			}

			if len(jmdictEntry.Kanji) == 0 || jmdictEntry.Kanji[0].Text != word.kanji {
				t.Errorf("JMdict file %s has incorrect kanji: expected %s",
					jmdictPath, word.kanji)
			}

			t.Logf("JMdict entry for %d: ID=%s, Kanji=%s",
				jmdictID, jmdictEntry.ID, jmdictEntry.Kanji[0].Text)
		}
	}

	// Verify index files for names
	t.Logf("Verifying index files for names")
	for _, name := range jmnedictNames {
		indexPath := filepath.Join(testDir, "index", name.kanji+".json.br")
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

		// Check that the index entry has JMNedict IDs
		if len(indexEntry.N) == 0 {
			t.Errorf("Index entry for %s has no JMNedict IDs", name.kanji)
		}

		t.Logf("Index entry for %s: %+v", name.kanji, indexEntry)

		// Verify that JMNedict files were created and contain correct data
		for _, jmnedictID := range indexEntry.Ne {
			jmnedictPath := filepath.Join(testDir, "n", fmt.Sprintf("%d.json.br", jmnedictID))
			if _, err := os.Stat(jmnedictPath); os.IsNotExist(err) {
				t.Errorf("Expected JMNedict file %s does not exist", jmnedictPath)
				continue
			}

			// Read and verify the JMNedict file
			jmnedictEntry, err := readJMNedictEntry(jmnedictPath)
			if err != nil {
				t.Errorf("Failed to read JMNedict file %s: %v", jmnedictPath, err)
				continue
			}

			if len(jmnedictEntry.Kanji) == 0 || jmnedictEntry.Kanji[0] != name.kanji {
				t.Errorf("JMNedict file %s has incorrect kanji: expected %s",
					jmnedictPath, name.kanji)
			}

			t.Logf("JMNedict entry for %d: ID=%s, Kanji=%s",
				jmnedictID, jmnedictEntry.ID, jmnedictEntry.Kanji[0])
		}
	}

	// Verify that reading-only entries are also indexed
	t.Logf("Verifying reading-only entries")
	for _, word := range jmdictWords {
		indexPath := filepath.Join(testDir, "index", word.kana+".json.br")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			t.Errorf("Expected index file for reading %s does not exist", word.kana)
			continue
		}

		// Read and check the index file
		indexEntry, err := readIndexEntry(indexPath)
		if err != nil {
			t.Errorf("Failed to read index file for reading %s: %v", word.kana, err)
			continue
		}

		// Check that the index entry has JMdict exact matches
		if len(indexEntry.Je) == 0 {
			t.Errorf("Index entry for reading %s has no JMdict exact matches", word.kana)
		}

		t.Logf("Index entry for reading %s: %+v", word.kana, indexEntry)
	}
}

// readKanjiEntry reads a Kanji entry from a compressed JSON file
func readKanjiEntry(filePath string) (*kanjidic.Kanji, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decompress the file
	br := brotli.NewReader(file)
	decoder := json.NewDecoder(br)

	var entry kanjidic.Kanji
	if err := decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// readChineseCharEntry reads a Chinese character entry from a compressed JSON file
func readChineseCharEntry(filePath string) (*chinese_chars.ChineseCharEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decompress the file
	br := brotli.NewReader(file)
	decoder := json.NewDecoder(br)

	var entry chinese_chars.ChineseCharEntry
	if err := decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// readJMdictEntry reads a JMdict entry from a compressed JSON file
func readJMdictEntry(filePath string) (*jmdict.Word, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decompress the file
	br := brotli.NewReader(file)
	decoder := json.NewDecoder(br)

	var entry jmdict.Word
	if err := decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// readJMNedictEntry reads a JMNedict entry from a compressed JSON file
func readJMNedictEntry(filePath string) (*jmnedict.Name, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decompress the file
	br := brotli.NewReader(file)
	decoder := json.NewDecoder(br)

	var entry jmnedict.Name
	if err := decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}
