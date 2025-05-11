package processor

import (
	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
)

// Type aliases for external use
type ChineseCharEntry = chinese_chars.ChineseCharEntry
type ChineseWordEntry = chinese_words.ChineseWordEntry
type JMdictWord = jmdict.Word
type JMNedictName = jmnedict.Name
type KanjidicKanji = kanjidic.Kanji
