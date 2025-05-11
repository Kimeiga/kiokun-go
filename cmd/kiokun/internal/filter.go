package internal

import (
	"strings"
	"unicode"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
)

// FilterEntries filters dictionary entries based on the configuration
func FilterEntries(entries *DictionaryEntries, config *Config, logf LogFunc) *DictionaryEntries {
	// If no filtering is needed, return the original entries
	if config.OutputMode == OutputAll && !config.TestMode && config.LimitEntries <= 0 {
		return entries
	}

	// Create a copy of the entries to avoid modifying the original
	result := &DictionaryEntries{
		JMdict:       entries.JMdict,
		JMNedict:     entries.JMNedict,
		Kanjidic:     entries.Kanjidic,
		ChineseChars: entries.ChineseChars,
		ChineseWords: entries.ChineseWords,
	}

	// Filter entries based on output mode
	if config.OutputMode != OutputAll {
		result = filterByOutputMode(result, config.OutputMode, logf)
	}

	// Apply test mode filtering if enabled
	if config.TestMode {
		result = filterForTestMode(result, logf)
	}

	// Apply entry limit if specified
	if config.LimitEntries > 0 {
		result = limitEntries(result, config.LimitEntries, logf)
	}

	return result
}

// filterByOutputMode filters entries based on the output mode
func filterByOutputMode(entries *DictionaryEntries, mode OutputMode, logf LogFunc) *DictionaryEntries {
	// Helper function to check if a string contains only Han characters
	isHanOnly := func(s string) bool {
		for _, r := range s {
			if !unicode.Is(unicode.Han, r) {
				return false
			}
		}
		return true
	}

	// Helper function to check if an entry should be included based on the mode
	shouldIncludeEntry := func(entry common.Entry) bool {
		// Get primary text representation for filtering
		var primaryText string
		switch e := entry.(type) {
		case jmdict.Word:
			if len(e.Kanji) > 0 {
				primaryText = e.Kanji[0].Text
			} else if len(e.Kana) > 0 {
				primaryText = e.Kana[0].Text
			} else {
				primaryText = e.ID
			}
		case jmnedict.Name:
			// Use the Name's primary text
			if len(e.Kanji) > 0 {
				primaryText = e.Kanji[0]
			} else if len(e.Reading) > 0 {
				primaryText = e.Reading[0]
			} else {
				primaryText = e.ID
			}
		case kanjidic.Kanji:
			// Use the Kanji character
			primaryText = e.Character
		case chinese_chars.ChineseCharEntry:
			// Use the traditional character
			primaryText = e.Traditional
		case chinese_words.ChineseWordEntry:
			// Use the traditional word
			primaryText = e.Traditional
		default:
			// If we don't know how to filter this type, include it by default
			return true
		}

		// Check if the text contains only Han characters
		isHan := isHanOnly(primaryText)
		charCount := len([]rune(primaryText)) // Get correct Unicode character count

		// Apply filtering based on mode
		switch mode {
		case OutputNonHanOnly:
			return !isHan
		case OutputHanOnly:
			return isHan
		case OutputHan1Char:
			return isHan && charCount == 1
		case OutputHan2Char:
			return isHan && charCount == 2
		case OutputHan3Plus:
			return isHan && charCount >= 3
		default:
			return true
		}
	}

	// Filter each dictionary
	filteredJmdict := make([]common.Entry, 0, len(entries.JMdict))
	filteredJmnedict := make([]common.Entry, 0, len(entries.JMNedict))
	filteredKanjidic := make([]common.Entry, 0, len(entries.Kanjidic))
	filteredChineseChars := make([]common.Entry, 0, len(entries.ChineseChars))
	filteredChineseWords := make([]common.Entry, 0, len(entries.ChineseWords))

	for _, entry := range entries.JMdict {
		if shouldIncludeEntry(entry) {
			filteredJmdict = append(filteredJmdict, entry)
		}
	}
	for _, entry := range entries.JMNedict {
		if shouldIncludeEntry(entry) {
			filteredJmnedict = append(filteredJmnedict, entry)
		}
	}
	for _, entry := range entries.Kanjidic {
		if shouldIncludeEntry(entry) {
			filteredKanjidic = append(filteredKanjidic, entry)
		}
	}
	for _, entry := range entries.ChineseChars {
		if shouldIncludeEntry(entry) {
			filteredChineseChars = append(filteredChineseChars, entry)
		}
	}
	for _, entry := range entries.ChineseWords {
		if shouldIncludeEntry(entry) {
			filteredChineseWords = append(filteredChineseWords, entry)
		}
	}

	logf("Filtered entries - JMdict: %d -> %d, JMNedict: %d -> %d, Kanjidic: %d -> %d, Chinese Chars: %d -> %d, Chinese Words: %d -> %d\n",
		len(entries.JMdict), len(filteredJmdict),
		len(entries.JMNedict), len(filteredJmnedict),
		len(entries.Kanjidic), len(filteredKanjidic),
		len(entries.ChineseChars), len(filteredChineseChars),
		len(entries.ChineseWords), len(filteredChineseWords))

	return &DictionaryEntries{
		JMdict:       filteredJmdict,
		JMNedict:     filteredJmnedict,
		Kanjidic:     filteredKanjidic,
		ChineseChars: filteredChineseChars,
		ChineseWords: filteredChineseWords,
	}
}

// filterForTestMode prioritizes entries that have overlap between Chinese and Japanese dictionaries
func filterForTestMode(entries *DictionaryEntries, logf LogFunc) *DictionaryEntries {
	logf("Test mode enabled - prioritizing entries with overlap between Chinese and Japanese dictionaries\n")

	// Create maps to track characters in each dictionary
	japaneseChars := make(map[string]bool)
	chineseChars := make(map[string]bool)

	// Collect Japanese characters
	for _, entry := range entries.Kanjidic {
		kanji, ok := entry.(kanjidic.Kanji)
		if ok {
			japaneseChars[kanji.Character] = true
		}
	}

	// Collect Chinese characters
	for _, entry := range entries.ChineseChars {
		char, ok := entry.(chinese_chars.ChineseCharEntry)
		if ok {
			chineseChars[char.Traditional] = true
			if char.Simplified != char.Traditional {
				chineseChars[char.Simplified] = true
			}
		}
	}

	// Find common characters
	var commonCharacters []string
	for char := range japaneseChars {
		if chineseChars[char] {
			commonCharacters = append(commonCharacters, char)
		}
	}

	logf("Found %d common characters between Chinese and Japanese dictionaries\n", len(commonCharacters))

	if len(commonCharacters) == 0 {
		// No common characters found, return original entries
		return entries
	}

	// Filter entries to prioritize those with common characters
	var prioritizedJmdictEntries []common.Entry
	var prioritizedJmnedictEntries []common.Entry
	var prioritizedKanjidicEntries []common.Entry
	var prioritizedChineseCharsEntries []common.Entry
	var prioritizedChineseWordsEntries []common.Entry

	// Helper function to check if an entry contains a common character
	containsCommonChar := func(text string) bool {
		for _, commonChar := range commonCharacters {
			if strings.Contains(text, commonChar) {
				return true
			}
		}
		return false
	}

	// Filter Kanjidic entries
	for _, entry := range entries.Kanjidic {
		kanji, ok := entry.(kanjidic.Kanji)
		if ok {
			for _, commonChar := range commonCharacters {
				if kanji.Character == commonChar {
					prioritizedKanjidicEntries = append(prioritizedKanjidicEntries, entry)
					break
				}
			}
		}
	}

	// Filter Chinese character entries
	for _, entry := range entries.ChineseChars {
		char, ok := entry.(chinese_chars.ChineseCharEntry)
		if ok {
			for _, commonChar := range commonCharacters {
				if char.Traditional == commonChar || char.Simplified == commonChar {
					prioritizedChineseCharsEntries = append(prioritizedChineseCharsEntries, entry)
					break
				}
			}
		}
	}

	// Filter JMdict entries
	for _, entry := range entries.JMdict {
		word, ok := entry.(jmdict.Word)
		if ok {
			// Check if any kanji form contains a common character
			found := false
			for _, kanji := range word.Kanji {
				if containsCommonChar(kanji.Text) {
					found = true
					break
				}
			}
			if found {
				prioritizedJmdictEntries = append(prioritizedJmdictEntries, entry)
			}
		}
	}

	// Filter JMNedict entries
	for _, entry := range entries.JMNedict {
		name, ok := entry.(jmnedict.Name)
		if ok {
			// Check if any kanji form contains a common character
			found := false
			for _, kanji := range name.Kanji {
				if containsCommonChar(kanji) {
					found = true
					break
				}
			}
			if found {
				prioritizedJmnedictEntries = append(prioritizedJmnedictEntries, entry)
			}
		}
	}

	// Filter Chinese word entries
	for _, entry := range entries.ChineseWords {
		word, ok := entry.(chinese_words.ChineseWordEntry)
		if ok {
			if containsCommonChar(word.Traditional) || containsCommonChar(word.Simplified) {
				prioritizedChineseWordsEntries = append(prioritizedChineseWordsEntries, entry)
			}
		}
	}

	// If we have prioritized entries, use them
	if len(prioritizedKanjidicEntries) > 0 || len(prioritizedChineseCharsEntries) > 0 {
		logf("Using prioritized entries - JMdict: %d, JMNedict: %d, Kanjidic: %d, Chinese Chars: %d, Chinese Words: %d\n",
			len(prioritizedJmdictEntries), len(prioritizedJmnedictEntries),
			len(prioritizedKanjidicEntries), len(prioritizedChineseCharsEntries),
			len(prioritizedChineseWordsEntries))

		// In test mode, always include all Chinese character entries to ensure overlap
		if len(entries.ChineseChars) > len(prioritizedChineseCharsEntries) {
			logf("In test mode, including all Chinese character entries to ensure overlap\n")
			// Keep the original Chinese character entries
			prioritizedChineseCharsEntries = entries.ChineseChars
		}

		return &DictionaryEntries{
			JMdict:       prioritizedJmdictEntries,
			JMNedict:     prioritizedJmnedictEntries,
			Kanjidic:     prioritizedKanjidicEntries,
			ChineseChars: prioritizedChineseCharsEntries,
			ChineseWords: prioritizedChineseWordsEntries,
		}
	}

	// No prioritized entries found, return original entries
	return entries
}

// limitEntries limits the number of entries from each dictionary
func limitEntries(entries *DictionaryEntries, limit int, logf LogFunc) *DictionaryEntries {
	totalEntries := len(entries.JMdict) + len(entries.JMNedict) + len(entries.Kanjidic) +
		len(entries.ChineseChars) + len(entries.ChineseWords)

	if limit >= totalEntries {
		// No need to limit
		return entries
	}

	logf("Limiting to %d entries (out of %d total)\n", limit, totalEntries)

	// Calculate proportions
	jmdictProportion := float64(len(entries.JMdict)) / float64(totalEntries)
	jmnedictProportion := float64(len(entries.JMNedict)) / float64(totalEntries)
	kanjidicProportion := float64(len(entries.Kanjidic)) / float64(totalEntries)
	chineseCharsProportion := float64(len(entries.ChineseChars)) / float64(totalEntries)
	chineseWordsProportion := float64(len(entries.ChineseWords)) / float64(totalEntries)

	// Calculate limits for each dictionary
	jmdictLimit := int(float64(limit) * jmdictProportion)
	jmnedictLimit := int(float64(limit) * jmnedictProportion)
	kanjidicLimit := int(float64(limit) * kanjidicProportion)
	chineseCharsLimit := int(float64(limit) * chineseCharsProportion)
	chineseWordsLimit := int(float64(limit) * chineseWordsProportion)

	// Adjust for rounding errors
	remaining := limit - jmdictLimit - jmnedictLimit - kanjidicLimit -
		chineseCharsLimit - chineseWordsLimit
	if remaining > 0 && len(entries.Kanjidic) > kanjidicLimit {
		kanjidicLimit += remaining
	}

	// Apply limits
	limitedJmdict := entries.JMdict
	limitedJmnedict := entries.JMNedict
	limitedKanjidic := entries.Kanjidic
	limitedChineseChars := entries.ChineseChars
	limitedChineseWords := entries.ChineseWords

	if jmdictLimit < len(entries.JMdict) {
		limitedJmdict = entries.JMdict[:jmdictLimit]
	}
	if jmnedictLimit < len(entries.JMNedict) {
		limitedJmnedict = entries.JMNedict[:jmnedictLimit]
	}
	if kanjidicLimit < len(entries.Kanjidic) {
		limitedKanjidic = entries.Kanjidic[:kanjidicLimit]
	}
	if chineseCharsLimit < len(entries.ChineseChars) {
		limitedChineseChars = entries.ChineseChars[:chineseCharsLimit]
	}
	if chineseWordsLimit < len(entries.ChineseWords) {
		limitedChineseWords = entries.ChineseWords[:chineseWordsLimit]
	}

	return &DictionaryEntries{
		JMdict:       limitedJmdict,
		JMNedict:     limitedJmnedict,
		Kanjidic:     limitedKanjidic,
		ChineseChars: limitedChineseChars,
		ChineseWords: limitedChineseWords,
	}
}
