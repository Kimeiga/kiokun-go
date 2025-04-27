package jmnedict

import (
	"os"

	"kiokun-go/dictionaries/common"
)

type Importer struct{}

func (i *Importer) Name() string {
	return "JMNedict"
}

func (i *Importer) Import(path string) ([]common.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dict JMnedict
	if err := common.ImportJSON(file, &dict); err != nil {
		return nil, err
	}

	// Convert []Word to []common.Entry
	entries := make([]common.Entry, len(dict.Words))
	for i, word := range dict.Words {
		// Convert Word to Name for compatibility with existing code
		name := Name{
			ID: word.ID,
		}

		// Extract kanji texts
		kanjiTexts := make([]string, len(word.Kanji))
		for j, k := range word.Kanji {
			kanjiTexts[j] = k.Text
		}
		name.Kanji = kanjiTexts

		// Extract kana texts as readings
		readings := make([]string, len(word.Kana))
		for j, k := range word.Kana {
			readings[j] = k.Text
		}
		name.Reading = readings

		// Extract meanings from translations
		var meanings []string
		for _, trans := range word.Translation {
			for _, detail := range trans.Translation {
				if detail.Lang == "eng" { // Assuming we want English meanings
					meanings = append(meanings, detail.Text)
				}
			}
			// Add type information
			name.Type = trans.Type
		}
		name.Meanings = meanings

		entries[i] = name
	}

	return entries, nil
}
