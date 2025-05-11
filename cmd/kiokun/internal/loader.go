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
	var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries []common.Entry

	for _, dict := range dictConfigs {
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
		}

		logf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(startTime).Seconds())
	}

	return &DictionaryEntries{
		JMdict:       jmdictEntries,
		JMNedict:     jmnedictEntries,
		Kanjidic:     kanjidicEntries,
		ChineseChars: chineseCharsEntries,
		ChineseWords: chineseWordsEntries,
	}, nil
}
