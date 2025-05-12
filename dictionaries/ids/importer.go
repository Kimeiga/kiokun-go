package ids

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"

	"kiokun-go/dictionaries/common"
)

// Importer handles importing IDS data
type Importer struct{}

// Name returns the name of this importer
func (i *Importer) Name() string {
	return "ids"
}

// Import reads and processes the IDS file
func (i *Importer) Import(path string) ([]common.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []common.Entry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, ";;") || strings.TrimSpace(line) == "" {
			continue
		}

		// Parse the line
		fields := strings.Split(line, "\t")
		if len(fields) < 3 {
			continue
		}

		codepoint := fields[0]
		character := fields[1]
		ids := fields[2]

		// Check if the character field contains exactly one character
		if utf8.RuneCountInString(character) != 1 {
			continue
		}

		// Check for apparent IDS
		apparentIDS := ""
		for i := 3; i < len(fields); i++ {
			if strings.HasPrefix(fields[i], "@apparent=") {
				apparentIDS = strings.TrimPrefix(fields[i], "@apparent=")
				break
			}
		}

		entry := IDSEntry{
			ID:          codepoint,
			Character:   character,
			IDS:         ids,
			ApparentIDS: apparentIDS,
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
