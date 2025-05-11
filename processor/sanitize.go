package processor

import (
	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
)

// sanitizeWordGroup removes duplicate entries and minifies data in a word group
func sanitizeWordGroup(group *WordGroup) {
	// Deduplicate and sanitize JMdict entries
	seen := make(map[string]bool)
	var uniqueWords []jmdict.Word
	for _, word := range group.WordJapanese {
		if !seen[word.ID] {
			seen[word.ID] = true

			// Sanitize each word by removing wildcards from senses and kana entries
			for i := range word.Sense {
				word.Sense[i].SanitizeWildcards()

				// Remove empty arrays in Sense
				if len(word.Sense[i].Related) == 0 {
					word.Sense[i].Related = nil
				}
				if len(word.Sense[i].Antonym) == 0 {
					word.Sense[i].Antonym = nil
				}
				if len(word.Sense[i].Field) == 0 {
					word.Sense[i].Field = nil
				}
				if len(word.Sense[i].Dialect) == 0 {
					word.Sense[i].Dialect = nil
				}
				if len(word.Sense[i].Misc) == 0 {
					word.Sense[i].Misc = nil
				}
				if len(word.Sense[i].Info) == 0 {
					word.Sense[i].Info = nil
				}
				if len(word.Sense[i].LanguageSource) == 0 {
					word.Sense[i].LanguageSource = nil
				}
				if len(word.Sense[i].Examples) == 0 {
					word.Sense[i].Examples = nil
				}
			}

			for i := range word.Kana {
				word.Kana[i].SanitizeWildcards()

				// Remove empty arrays in Kana
				if len(word.Kana[i].Tags) == 0 {
					word.Kana[i].Tags = nil
				}
			}

			// Remove empty arrays in Kanji
			for i := range word.Kanji {
				if len(word.Kanji[i].Tags) == 0 {
					word.Kanji[i].Tags = nil
				}
			}

			uniqueWords = append(uniqueWords, word)
		}
	}
	group.WordJapanese = uniqueWords

	// Deduplicate JMNedict entries
	seen = make(map[string]bool)
	var uniqueNames []jmnedict.Name
	for _, name := range group.NameJapanese {
		if !seen[name.ID] {
			seen[name.ID] = true
			uniqueNames = append(uniqueNames, name)
		}
	}
	group.NameJapanese = uniqueNames

	// Deduplicate Kanjidic entries
	seen = make(map[string]bool)
	var uniqueKanji []kanjidic.Kanji
	for _, kanji := range group.CharJapanese {
		if !seen[kanji.Character] {
			seen[kanji.Character] = true

			// Sanitize Kanjidic entries by removing empty fields
			// (This would need to be expanded based on the Kanjidic structure)

			uniqueKanji = append(uniqueKanji, kanji)
		}
	}
	group.CharJapanese = uniqueKanji

	// Deduplicate and sanitize Chinese character entries
	seen = make(map[string]bool)
	var uniqueChineseChars []chinese_chars.ChineseCharEntry
	for _, char := range group.CharChinese {
		if !seen[char.ID] {
			seen[char.ID] = true

			// Sanitize Chinese character entries by removing empty slices
			if len(char.Definitions) == 0 {
				char.Definitions = nil
			}
			if len(char.Pinyin) == 0 {
				char.Pinyin = nil
			}

			// If stroke count is 0, omit it
			if char.StrokeCount == 0 {
				char.StrokeCount = 0 // This will be omitted due to omitempty tag
			}

			uniqueChineseChars = append(uniqueChineseChars, char)
		}
	}
	group.CharChinese = uniqueChineseChars

	// Deduplicate and sanitize Chinese word entries
	seen = make(map[string]bool)
	var uniqueChineseWords []chinese_words.ChineseWordEntry
	for _, word := range group.WordChinese {
		if !seen[word.ID] {
			seen[word.ID] = true

			// Sanitize Chinese word entries by removing empty slices
			if len(word.Definitions) == 0 {
				word.Definitions = nil
			}
			if len(word.Pinyin) == 0 {
				word.Pinyin = nil
			}
			if word.Frequency != nil && len(word.Frequency) == 0 {
				word.Frequency = nil
			}

			uniqueChineseWords = append(uniqueChineseWords, word)
		}
	}
	group.WordChinese = uniqueChineseWords
}
