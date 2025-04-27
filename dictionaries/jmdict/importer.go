package jmdict

import (
	"os"

	"kiokun-go/dictionaries/common"
)

type Importer struct{}

func (i *Importer) Name() string {
	return "jmdict"
}

func (i *Importer) Import(path string) ([]common.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dict JmdictTypes
	if err := common.ImportJSON(file, &dict); err != nil {
		return nil, err
	}

	// Convert words to entries
	entries := make([]common.Entry, len(dict.Words))
	for i, word := range dict.Words {
		entries[i] = word
	}

	return entries, nil
}
