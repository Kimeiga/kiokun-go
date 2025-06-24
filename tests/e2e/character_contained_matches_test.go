package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/andybalholm/brotli"

	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/kanjidic"
	"kiokun-go/processor"
)

// TestCharacterContainedMatches tests that contained matches work correctly for any character
func TestCharacterContainedMatches(t *testing.T) {
	testCases := []struct {
		char                string
		expectedMinMatches  int
		description         string
	}{
		{"日", 2, "Japanese character for day/sun"},
		{"人", 2, "Japanese character for person"},
		{"水", 2, "Japanese character for water"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Character_%s", tc.char), func(t *testing.T) {
			testCharacterContainedMatches(t, tc.char, tc.expectedMinMatches, tc.description)
		})
	}
}

func testCharacterContainedMatches(t *testing.T, testChar string, expectedMinMatches int, description string) {
	t.Logf("Testing character: %s (%s)", testChar, description)

	// Step 1: Create test data
	entries := createTestData(testChar, t)
	t.Logf("Created %d test entries", len(entries))

	// Step 2: Process entries
	outputDir := fmt.Sprintf("test_output_%s", testChar)
	if err := processEntries(entries, outputDir, t); err != nil {
		t.Fatalf("Error processing entries: %v", err)
	}

	// Step 3: Verify shard distribution
	if err := verifyShardDistribution(testChar, outputDir, t); err != nil {
		t.Fatalf("Error verifying shards: %v", err)
	}

	// Step 4: Test API logic
	containedMatches, err := simulateAPICall(testChar, outputDir, t)
	if err != nil {
		t.Fatalf("Error testing API logic: %v", err)
	}

	totalMatches := 0
	for dictType, count := range containedMatches {
		totalMatches += count
		t.Logf("Dictionary %s: %d matches", dictType, count)
	}

	if totalMatches < expectedMinMatches {
		t.Errorf("Expected at least %d contained matches, but found %d", expectedMinMatches, totalMatches)
	}

	t.Logf("✅ Found %d total contained matches (expected >= %d)", totalMatches, expectedMinMatches)

	// Step 5: Cleanup
	if err := cleanupTestFiles(outputDir); err != nil {
		t.Logf("Warning: Error cleaning up: %v", err)
	}
}

// createTestData creates test dictionary entries for the specified character
func createTestData(testChar string, t *testing.T) []common.Entry {
	var entries []common.Entry

	// Create the character entry itself (Kanjidic)
	charEntry := createKanjidicEntry(testChar)
	entries = append(entries, charEntry)

	// Create words containing the character
	wordsContainingChar := generateWordsContaining(testChar)
	for i, word := range wordsContainingChar {
		wordEntry := createJMdictEntry(word, i+1000000)
		entries = append(entries, wordEntry)
	}

	// Create a word that doesn't contain the character (for comparison)
	nonContainingWord := createJMdictEntry("食べる", 9999999)
	entries = append(entries, nonContainingWord)

	return entries
}

func processEntries(entries []common.Entry, outputDir string, t *testing.T) error {
	proc, err := processor.NewShardedIndexProcessor(outputDir, 1)
	if err != nil {
		return fmt.Errorf("creating processor: %w", err)
	}

	if err := proc.ProcessEntries(entries); err != nil {
		return fmt.Errorf("processing entries: %w", err)
	}

	if err := proc.WriteToFiles(); err != nil {
		return fmt.Errorf("writing files: %w", err)
	}

	return nil
}

func verifyShardDistribution(testChar, outputDir string, t *testing.T) error {
	shardSuffixes := []string{"_non_han", "_han_1char", "_han_2char", "_han_3plus"}
	foundShards := 0

	for _, suffix := range shardSuffixes {
		shardDir := outputDir + suffix
		indexDir := filepath.Join(shardDir, "index")
		
		if _, err := os.Stat(indexDir); err == nil {
			foundShards++
			t.Logf("Found shard: %s", suffix)
		}
	}

	if foundShards == 0 {
		return fmt.Errorf("no shards were created")
	}

	t.Logf("Found %d shards with data", foundShards)
	return nil
}

// simulateAPICall simulates the frontend API call logic
func simulateAPICall(testChar, outputDir string, t *testing.T) (map[string]int, error) {
	// Determine if this is a single character search
	runes := []rune(testChar)
	isSingleCharacter := len(runes) == 1 && isHanCharacter(testChar)
	
	t.Logf("Is single character: %v", isSingleCharacter)
	
	containedMatches := make(map[string]int)

	if isSingleCharacter {
		t.Log("Single character detected - searching across all shards")

		// Search across all shards for contained matches
		shardSuffixes := []string{"_han_1char", "_han_2char", "_han_3plus", "_non_han"}
		
		for _, suffix := range shardSuffixes {
			shardDir := outputDir + suffix
			indexPath := filepath.Join(shardDir, "index", testChar+".json.br")
			
			if _, err := os.Stat(indexPath); os.IsNotExist(err) {
				continue
			}

			// Read and decompress the index file
			indexEntry, err := readIndexFile(indexPath)
			if err != nil {
				continue
			}

			t.Logf("Found index file in shard %s", suffix)

			// Process contained matches
			if indexEntry.C != nil {
				for dictType, ids := range indexEntry.C {
					containedMatches[dictType] += len(ids)
					t.Logf("  Contained matches in %s: %d entries", dictType, len(ids))
				}
			}
		}
	} else {
		t.Log("Multi-character word - checking primary shard only")
		// Handle multi-character case (simplified for this test)
	}

	return containedMatches, nil
}

// Helper functions (simplified versions of the original functions)

func createKanjidicEntry(char string) kanjidic.Kanji {
	meanings := map[string][]string{
		"日": {"day", "sun", "Japan"},
		"水": {"water"},
		"人": {"person", "people"},
	}

	charMeanings := meanings[char]
	if charMeanings == nil {
		charMeanings = []string{"test character"}
	}

	return kanjidic.Kanji{
		Character: char,
		NumericID: "1",
		Meanings:  charMeanings,
		OnYomi:    []string{"TEST"},
		KunYomi:   []string{"てすと"},
		JLPT:      4,
		Grade:     1,
		Stroke:    4,
		Frequency: 1,
	}
}

func createJMdictEntry(word string, id int) jmdict.Word {
	return jmdict.Word{
		ID: strconv.Itoa(id),
		Kanji: []jmdict.KanjiEntry{
			{
				Common: true,
				Text:   word,
			},
		},
		Kana: []jmdict.KanaEntry{
			{
				Common:         true,
				Text:           "てすと",
				AppliesToKanji: []string{"*"},
			},
		},
		Sense: []jmdict.Sense{
			{
				PartOfSpeech: []string{"n"},
				Gloss: []jmdict.Gloss{
					{
						Lang: "eng",
						Text: fmt.Sprintf("test word: %s", word),
					},
				},
			},
		},
	}
}

func generateWordsContaining(testChar string) []string {
	words := []string{
		testChar + "本",     // 2-char word starting with test char
		"今" + testChar,     // 2-char word ending with test char
		testChar + "本語",   // 3-char word starting with test char
		"大" + testChar + "本", // 3-char word with test char in middle
	}

	// Add character-specific words
	switch testChar {
	case "日":
		words = append(words, "日本語", "今日", "毎日")
	case "水":
		words = append(words, "水曜日", "飲み水")
	case "人":
		words = append(words, "人間", "日本人", "大人")
	}

	// Filter to only words that actually contain the character
	var result []string
	for _, word := range words {
		if containsCharacter(word, testChar) {
			result = append(result, word)
		}
	}

	return result
}

func containsCharacter(text, char string) bool {
	for _, r := range text {
		if string(r) == char {
			return true
		}
	}
	return false
}

func isHanCharacter(char string) bool {
	runes := []rune(char)
	if len(runes) != 1 {
		return false
	}
	
	code := runes[0]
	return code >= 0x4e00 && code <= 0x9fff
}

// IndexEntry represents the structure of an index file
type IndexEntry struct {
	E map[string][]int `json:"e,omitempty"` // Exact matches
	C map[string][]int `json:"c,omitempty"` // Contained matches
}

func readIndexFile(path string) (*IndexEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := brotli.NewReader(file)
	
	var indexEntry IndexEntry
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&indexEntry); err != nil {
		return nil, err
	}

	return &indexEntry, nil
}

func cleanupTestFiles(outputDir string) error {
	shardSuffixes := []string{"_non_han", "_han_1char", "_han_2char", "_han_3plus"}
	
	for _, suffix := range shardSuffixes {
		shardDir := outputDir + suffix
		if err := os.RemoveAll(shardDir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing %s: %w", shardDir, err)
		}
	}
	
	return nil
}
