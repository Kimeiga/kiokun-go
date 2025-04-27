package chinese_chars

import (
	"path/filepath"

	"kiokun-go/dictionaries/common"
)

func init() {
	sourceDir := filepath.Join("dictionaries", "chinese_chars", "source")
	pattern := `^dictionary_char_.*\.json$`

	filename, err := common.FindDictionaryFile(sourceDir, pattern)
	if err != nil {
		// If file not found, register with default name - the importer will handle the error
		filename = "dictionary_char_2024-06-17.json"
	}

	common.RegisterDictionary("chinese_chars", filename, &Importer{})
}
