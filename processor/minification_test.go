package processor

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"

	"github.com/andybalholm/brotli"
)

// TestMinification tests the entire minification process:
// 1. Importing dictionary entries
// 2. Processing them
// 3. Writing them to files
// 4. Verifying the minification was applied correctly
func TestMinification(t *testing.T) {
	// Create test directory
	testDir := filepath.Join(os.TempDir(), "kiokun_minification_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Set dictionaries base path
	dictDir := filepath.Join("..", "dictionaries")
	common.SetDictionariesBasePath(dictDir)

	// Create processor
	proc, err := New(testDir, 1)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Test with JMdict entries
	t.Run("JMdict", func(t *testing.T) {
		testJMdictMinification(t, proc, dictDir, testDir)
	})

	// Test with Chinese dictionary entries
	t.Run("Chinese", func(t *testing.T) {
		testChineseMinification(t, proc, dictDir, testDir)
	})

	// Verify the minification
	t.Run("Verification", func(t *testing.T) {
		verifyMinification(t, testDir)
	})
}

// testJMdictMinification tests minification of JMdict entries
func testJMdictMinification(t *testing.T, proc *DictionaryProcessor, dictDir, testDir string) {
	// Find JMdict source file
	jmdictSourceDir := filepath.Join(dictDir, "jmdict", "source")
	pattern := `^jmdict-.*\.json$`
	jmdictFile, err := common.FindDictionaryFile(jmdictSourceDir, pattern)
	if err != nil {
		t.Logf("JMdict file not found, skipping test: %v", err)
		return
	}

	// Import JMdict entries
	jmdictPath := filepath.Join(jmdictSourceDir, jmdictFile)
	t.Logf("Using JMdict file: %s", jmdictPath)
	
	file, err := os.Open(jmdictPath)
	if err != nil {
		t.Fatalf("Failed to open JMdict file: %v", err)
	}
	defer file.Close()

	var dict jmdict.JmdictTypes
	if err := common.ImportJSON(file, &dict); err != nil {
		t.Fatalf("Failed to parse JMdict: %v", err)
	}

	// Limit to first 100 entries for testing
	maxEntries := 100
	if len(dict.Words) > maxEntries {
		dict.Words = dict.Words[:maxEntries]
	}
	t.Logf("Testing with %d JMdict entries", len(dict.Words))

	// Create entries for processing
	entries := make([]common.Entry, len(dict.Words))
	for i, word := range dict.Words {
		entries[i] = word
	}

	// Process entries
	if err := proc.ProcessEntries(entries); err != nil {
		t.Fatalf("Failed to process JMdict entries: %v", err)
	}

	// Write to files
	if err := proc.WriteToFiles(); err != nil {
		t.Fatalf("Failed to write files: %v", err)
	}
}

// testChineseMinification tests minification of Chinese dictionary entries
func testChineseMinification(t *testing.T, proc *DictionaryProcessor, dictDir, testDir string) {
	// Find Chinese character dictionary source file
	chineseCharsSourceDir := filepath.Join(dictDir, "chinese_chars", "source")
	pattern := `^dictionary_char_.*\.json$`
	chineseCharsFile, err := common.FindDictionaryFile(chineseCharsSourceDir, pattern)
	if err != nil {
		t.Logf("Chinese character file not found, skipping test: %v", err)
		return
	}

	// Import Chinese character entries
	chineseCharsPath := filepath.Join(chineseCharsSourceDir, chineseCharsFile)
	t.Logf("Using Chinese character file: %s", chineseCharsPath)
	
	file, err := os.Open(chineseCharsPath)
	if err != nil {
		t.Logf("Failed to open Chinese character file, skipping test: %v", err)
		return
	}
	defer file.Close()

	var charEntries []chinese_chars.ChineseCharEntry
	if err := json.NewDecoder(file).Decode(&charEntries); err != nil {
		t.Logf("Failed to parse Chinese character file, skipping test: %v", err)
		return
	}

	// Limit to first 100 entries for testing
	maxEntries := 100
	if len(charEntries) > maxEntries {
		charEntries = charEntries[:maxEntries]
	}
	t.Logf("Testing with %d Chinese character entries", len(charEntries))

	// Create entries for processing
	entries := make([]common.Entry, len(charEntries))
	for i, char := range charEntries {
		// If ID is not set, use traditional character as ID
		if char.ID == "" {
			char.ID = char.Traditional
		}
		entries[i] = char
	}

	// Process entries
	if err := proc.ProcessEntries(entries); err != nil {
		t.Fatalf("Failed to process Chinese character entries: %v", err)
	}

	// Find Chinese word dictionary source file
	chineseWordsSourceDir := filepath.Join(dictDir, "chinese_words", "source")
	pattern = `^dictionary_word_.*\.json$`
	chineseWordsFile, err := common.FindDictionaryFile(chineseWordsSourceDir, pattern)
	if err != nil {
		t.Logf("Chinese word file not found, skipping that part of test: %v", err)
	} else {
		// Import Chinese word entries
		chineseWordsPath := filepath.Join(chineseWordsSourceDir, chineseWordsFile)
		t.Logf("Using Chinese word file: %s", chineseWordsPath)
		
		file, err := os.Open(chineseWordsPath)
		if err != nil {
			t.Logf("Failed to open Chinese word file, skipping that part of test: %v", err)
		} else {
			defer file.Close()

			var wordEntries []chinese_words.ChineseWordEntry
			if err := json.NewDecoder(file).Decode(&wordEntries); err != nil {
				t.Logf("Failed to parse Chinese word file, skipping that part of test: %v", err)
			} else {
				// Limit to first 100 entries for testing
				if len(wordEntries) > maxEntries {
					wordEntries = wordEntries[:maxEntries]
				}
				t.Logf("Testing with %d Chinese word entries", len(wordEntries))

				// Create entries for processing
				wordEntriesCommon := make([]common.Entry, len(wordEntries))
				for i, word := range wordEntries {
					// If ID is not set, use traditional word as ID
					if word.ID == "" {
						word.ID = word.Traditional
					}
					wordEntriesCommon[i] = word
				}

				// Process entries
				if err := proc.ProcessEntries(wordEntriesCommon); err != nil {
					t.Fatalf("Failed to process Chinese word entries: %v", err)
				}
			}
		}
	}

	// Write to files
	if err := proc.WriteToFiles(); err != nil {
		t.Fatalf("Failed to write files: %v", err)
	}
}

// verifyMinification checks that entries were properly minified
func verifyMinification(t *testing.T, testDir string) {
	// Count of files with issues
	var filesWithWildcards, filesWithEmptyArrays, totalFiles int
	
	// Walk through the output directory and check files
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Only check .json.br files
		if !strings.HasSuffix(path, ".json.br") {
			return nil
		}
		
		totalFiles++
		
		// Read the file
		file, err := os.Open(path)
		if err != nil {
			t.Errorf("Failed to open file %s: %v", path, err)
			return nil
		}
		defer file.Close()
		
		// Decompress the file
		br := brotli.NewReader(file)
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, br); err != nil {
			t.Errorf("Failed to decompress file %s: %v", path, err)
			return nil
		}
		
		// Parse the JSON
		var group WordGroup
		if err := json.Unmarshal(buf.Bytes(), &group); err != nil {
			t.Errorf("Failed to parse JSON from file %s: %v", path, err)
			return nil
		}
		
		// Check for wildcards in JMdict entries
		for _, word := range group.WordJapanese {
			// Check senses for wildcards
			for _, sense := range word.Sense {
				// Check if appliesToKanji contains "*"
				for _, applies := range sense.AppliesToKanji {
					if applies == "*" {
						filesWithWildcards++
						t.Errorf("Found wildcard in appliesToKanji in file %s", path)
					}
				}
				
				// Check if appliesToKana contains "*"
				for _, applies := range sense.AppliesToKana {
					if applies == "*" {
						filesWithWildcards++
						t.Errorf("Found wildcard in appliesToKana in file %s", path)
					}
				}
			}
			
			// Check kana entries for wildcards
			for _, kana := range word.Kana {
				// Check if appliesToKanji contains "*"
				for _, applies := range kana.AppliesToKanji {
					if applies == "*" {
						filesWithWildcards++
						t.Errorf("Found wildcard in kana.appliesToKanji in file %s", path)
					}
				}
			}
		}
		
		// Check for empty arrays in Chinese character entries
		for _, char := range group.CharChinese {
			// Check if Definitions is an empty array
			if len(char.Definitions) == 0 && char.Definitions != nil {
				filesWithEmptyArrays++
				t.Errorf("Found empty Definitions array in file %s", path)
			}
			
			// Check if Pinyin is an empty array
			if len(char.Pinyin) == 0 && char.Pinyin != nil {
				filesWithEmptyArrays++
				t.Errorf("Found empty Pinyin array in file %s", path)
			}
		}
		
		// Check for empty arrays in Chinese word entries
		for _, word := range group.WordChinese {
			// Check if Definitions is an empty array
			if len(word.Definitions) == 0 && word.Definitions != nil {
				filesWithEmptyArrays++
				t.Errorf("Found empty Definitions array in file %s", path)
			}
			
			// Check if Pinyin is an empty array
			if len(word.Pinyin) == 0 && word.Pinyin != nil {
				filesWithEmptyArrays++
				t.Errorf("Found empty Pinyin array in file %s", path)
			}
			
			// Check if Frequency is an empty map
			if word.Frequency != nil && len(word.Frequency) == 0 {
				filesWithEmptyArrays++
				t.Errorf("Found empty Frequency map in file %s", path)
			}
		}
		
		return nil
	})
	
	if err != nil {
		t.Errorf("Error walking output directory: %v", err)
	}
	
	t.Logf("Verified %d files", totalFiles)
	
	// We expect no files with wildcards or empty arrays
	if filesWithWildcards > 0 {
		t.Errorf("Found %d files with wildcards, expected 0", filesWithWildcards)
	}
	
	if filesWithEmptyArrays > 0 {
		t.Errorf("Found %d files with empty arrays, expected 0", filesWithEmptyArrays)
	}
}
