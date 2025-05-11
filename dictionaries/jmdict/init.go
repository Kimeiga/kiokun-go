package jmdict

import (
	"path/filepath"

	"kiokun-go/dictionaries/common"
)

func init() {
	sourceDir := filepath.Join("dictionaries", "jmdict", "source")
	pattern := `^jmdict-examples-eng-.*\.json$`

	filename, err := common.FindDictionaryFile(sourceDir, pattern)
	if err != nil {
		// If file not found, register with default name - the importer will handle the error
		filename = "jmdict-examples-eng-3.6.1.json"
	}

	common.RegisterDictionary("jmdict", filename, &Importer{})
}
