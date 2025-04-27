package jmnedict

import (
	"path/filepath"

	"kiokun-go/dictionaries/common"
)

func init() {
	sourceDir := filepath.Join("dictionaries", "jmnedict", "source")
	pattern := `^jmnedict-all-.*\.json$`

	filename, err := common.FindDictionaryFile(sourceDir, pattern)
	if err != nil {
		// If file not found, register with default name - the importer will handle the error
		filename = "jmnedict-all-.json"
	}

	common.RegisterDictionary("jmnedict", filename, &Importer{})
}
