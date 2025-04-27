package kanjidic

import (
	"path/filepath"

	"kiokun-go/dictionaries/common"
)

func init() {
	sourceDir := filepath.Join("dictionaries", "kanjidic", "source")
	pattern := `kanjidic2-en-\d+\.\d+\.\d+\.json$`

	filename, err := common.FindDictionaryFile(sourceDir, pattern)
	if err != nil {
		// If file not found, register with default name - the importer will handle the error
		filename = "kanjidic2-en-3.6.1.json"
	}

	common.RegisterDictionary("kanjidic", filename, &Importer{})
}
