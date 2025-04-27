package kanjidic

import (
	"os"

	"kiokun-go/dictionaries/common"
)

type Importer struct{}

func (i *Importer) Name() string {
	return "Kanjidic"
}

func (i *Importer) Import(path string) ([]common.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse as Kanjidic2 which matches the actual file format
	var dict Kanjidic2
	if err := common.ImportJSON(file, &dict); err != nil {
		return nil, err
	}

	// Convert each Character in the file to our simplified Kanji type
	entries := make([]common.Entry, len(dict.Characters))
	for i, char := range dict.Characters {
		// Extract readings from the ReadingMeaning structure
		kanji := Kanji{
			Character: char.Literal,
			Meanings:  []string{},
			OnYomi:    []string{},
			KunYomi:   []string{},
			Stroke:    0,
		}

		// Extract stroke count from the first element if available
		if char.Misc.StrokeCounts != nil && len(char.Misc.StrokeCounts) > 0 {
			kanji.Stroke = char.Misc.StrokeCounts[0]
		}

		// Extract grade if available
		if char.Misc.Grade != nil {
			kanji.Grade = *char.Misc.Grade
		}

		// Extract JLPT level if available
		if char.Misc.JlptLevel != nil {
			kanji.JLPT = *char.Misc.JlptLevel
		}

		// Extract frequency if available
		if char.Misc.Frequency != nil {
			kanji.Frequency = *char.Misc.Frequency
		}

		// Extract readings and meanings
		if char.ReadingMeaning != nil {
			for _, group := range char.ReadingMeaning.Groups {
				// Extract on and kun readings
				for _, reading := range group.Readings {
					if reading.Type == "ja_on" {
						kanji.OnYomi = append(kanji.OnYomi, reading.Value)
					} else if reading.Type == "ja_kun" {
						kanji.KunYomi = append(kanji.KunYomi, reading.Value)
					}
				}

				// Extract English meanings
				for _, meaning := range group.Meanings {
					if meaning.Lang == "en" {
						kanji.Meanings = append(kanji.Meanings, meaning.Value)
					}
				}
			}

			// Add nanori readings (name readings) to a separate field if needed
			if len(char.ReadingMeaning.Nanori) > 0 {
				kanji.Radicals = char.ReadingMeaning.Nanori // Reusing radicals field for nanori readings
			}
		}

		entries[i] = kanji
	}

	return entries, nil
}
