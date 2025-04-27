package main

import (
	"encoding/json"
	"fmt"
	"os"

	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
)

func main() {
	// Test JMNedict import
	jmnedictFile, err := os.Open("dictionaries/jmnedict/source/jmnedict-all-3.6.1.json")
	if err != nil {
		fmt.Printf("Error opening JMNedict file: %v\n", err)
		return
	}
	defer jmnedictFile.Close()

	// First, decode into a map to check the structure
	var rawDict map[string]interface{}
	if err := json.NewDecoder(jmnedictFile).Decode(&rawDict); err != nil {
		fmt.Printf("Error decoding JMNedict JSON: %v\n", err)
		return
	}

	fmt.Println("JMNedict file keys:")
	for k := range rawDict {
		fmt.Printf("- %s\n", k)
	}

	if words, ok := rawDict["words"].([]interface{}); ok {
		fmt.Printf("JMNedict has %d words entries\n", len(words))
	}

	// Reset file pointer to beginning
	jmnedictFile.Seek(0, 0)

	// Now decode using the correct struct (JMnedict instead of JMNedictTypes)
	var jmneDict jmnedict.JMnedict
	if err := json.NewDecoder(jmnedictFile).Decode(&jmneDict); err != nil {
		fmt.Printf("Error decoding JMNedict into struct: %v\n", err)
		return
	}
	fmt.Printf("JMNedict successfully parsed with %d words\n", len(jmneDict.Words))

	// Test Kanjidic import
	kanjidicFile, err := os.Open("dictionaries/kanjidic/source/kanjidic2-en-3.6.1.json")
	if err != nil {
		fmt.Printf("Error opening Kanjidic file: %v\n", err)
		return
	}
	defer kanjidicFile.Close()

	var kanjiDict kanjidic.KanjidicTypes
	if err := json.NewDecoder(kanjidicFile).Decode(&kanjiDict); err != nil {
		fmt.Printf("Error decoding Kanjidic JSON: %v\n", err)
		return
	}
	fmt.Printf("Kanjidic successfully parsed with %d characters\n", len(kanjiDict.Characters))
}
