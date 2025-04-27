package chinese_words

import (
	"path/filepath"

	"kiokun-go/dictionaries/common"
)

func init() {
	sourceDir := filepath.Join("dictionaries", "chinese_words", "source")
	pattern := `^dictionary_word_.*\.json$`

	filename, err := common.FindDictionaryFile(sourceDir, pattern)
	if err != nil {
		// If file not found, register with default name - the importer will handle the error
		filename = "dictionary_word_2024-06-17.json"
	}

	common.RegisterDictionary("chinese_words", filename, &Importer{})
}
