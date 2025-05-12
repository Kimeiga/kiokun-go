package internal

import (
	"fmt"
	"path/filepath"
	"time"

	"kiokun-go/dictionaries/common"
)

// DictionaryEntries holds entries from all dictionaries
type DictionaryEntries struct {
	JMdict       []common.Entry
	JMNedict     []common.Entry
	Kanjidic     []common.Entry
	ChineseChars []common.Entry
	ChineseWords []common.Entry
	IDS          []common.Entry
}

// LoadDictionaries loads all dictionaries and returns their entries
func LoadDictionaries(config *Config, logf LogFunc) (*DictionaryEntries, error) {
	// Resolve dictionaries path
	dictPath := filepath.Join(config.WorkspaceRoot, config.DictDir)
	logf("Using dictionary path: %s\n", dictPath)
	common.SetDictionariesBasePath(dictPath)

	// Import all dictionaries
	logf("Importing dictionaries...\n")

	// Get all registered dictionaries
	dictConfigs := common.GetRegisteredDictionaries()

	// Import each dictionary
	var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries, idsEntries []common.Entry

	// Check if any specific dictionary is selected
	onlySpecificDict := config.OnlyJMdict || config.OnlyJMNedict || config.OnlyKanjidic ||
		config.OnlyChineseChars || config.OnlyChineseWords || config.OnlyIDS

	for _, dict := range dictConfigs {
		// Skip dictionaries that are not selected when using specific dictionary flags
		if onlySpecificDict {
			switch dict.Name {
			case "jmdict":
				if !config.OnlyJMdict {
					logf("Skipping %s (not selected)\n", dict.Name)
					continue
				}
			case "jmnedict":
				if !config.OnlyJMNedict {
					logf("Skipping %s (not selected)\n", dict.Name)
					continue
				}
			case "kanjidic":
				if !config.OnlyKanjidic {
					logf("Skipping %s (not selected)\n", dict.Name)
					continue
				}
			case "chinese_chars":
				if !config.OnlyChineseChars {
					logf("Skipping %s (not selected)\n", dict.Name)
					continue
				}
			case "chinese_words":
				if !config.OnlyChineseWords {
					logf("Skipping %s (not selected)\n", dict.Name)
					continue
				}
			case "ids", "ids_ext_a":
				// Always load IDS dictionaries for character composition data
				// but only if we're processing Kanjidic or Chinese Chars
				if !config.OnlyIDS && !config.OnlyKanjidic && !config.OnlyChineseChars {
					logf("Skipping %s (not selected)\n", dict.Name)
					continue
				}
			}
		}

		// Construct full path
		inputPath := filepath.Join(dict.SourceDir, dict.InputFile)

		// Import this dictionary
		logf("Importing %s from %s...\n", dict.Name, inputPath)
		startTime := time.Now()

		entries, err := dict.Importer.Import(inputPath)
		if err != nil {
			return nil, fmt.Errorf("error importing %s: %v", dict.Name, err)
		}

		// Store entries by dictionary type
		switch dict.Name {
		case "jmdict":
			jmdictEntries = entries
		case "jmnedict":
			jmnedictEntries = entries
		case "kanjidic":
			kanjidicEntries = entries
		case "chinese_chars":
			chineseCharsEntries = entries
		case "chinese_words":
			chineseWordsEntries = entries
		case "ids", "ids_ext_a":
			// Append IDS entries from different files
			idsEntries = append(idsEntries, entries...)
		}

		logf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(startTime).Seconds())
	}

	return &DictionaryEntries{
		JMdict:       jmdictEntries,
		JMNedict:     jmnedictEntries,
		Kanjidic:     kanjidicEntries,
		ChineseChars: chineseCharsEntries,
		ChineseWords: chineseWordsEntries,
		IDS:          idsEntries,
	}, nil
}
