package processor

import (
	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
)

// WordGroup represents the combined data for a single word/character
type WordGroup struct {
	WordJapanese []jmdict.Word                    `json:"w_j,omitempty"`
	NameJapanese []jmnedict.Name                  `json:"n_j,omitempty"`
	CharJapanese []kanjidic.Kanji                 `json:"c_j,omitempty"`
	CharChinese  []chinese_chars.ChineseCharEntry `json:"c_c,omitempty"`
	WordChinese  []chinese_words.ChineseWordEntry `json:"w_c,omitempty"`
}

// HasMultipleDictData returns true if this group has data from multiple dictionaries
func (wg *WordGroup) HasMultipleDictData() bool {
	sources := 0
	if len(wg.WordJapanese) > 0 {
		sources++
	}
	if len(wg.NameJapanese) > 0 {
		sources++
	}
	if len(wg.CharJapanese) > 0 {
		sources++
	}
	if len(wg.CharChinese) > 0 {
		sources++
	}
	if len(wg.WordChinese) > 0 {
		sources++
	}
	return sources > 1
}
