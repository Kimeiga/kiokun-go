package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/andybalholm/brotli"
	"github.com/ulikunitz/xz"
)

type WordGroup struct {
	WordJapanese []Word `json:"w_j"`
}

// OutputMode determines which words to output
type OutputMode string

const (
	OutputAll        OutputMode = "all"      // Output all words
	OutputHanOnly    OutputMode = "han-only" // Output words with only Han characters
	OutputNonHanOnly OutputMode = "non-han"  // Output words with at least one non-Han character
)

func main() {
	// Command line flags
	inputFile := flag.String("input", "jmdict-eng-3.5.0.json.xz", "Input JMDict JSON file")
	unzipped := flag.Bool("unzipped", false, "Output uncompressed JSON files")
	silent := flag.Bool("silent", false, "Disable progress output")
	outputModeFlag := flag.String("mode", "all", "Output mode: 'all', 'han-only' (Han/Kanji characters only), or 'non-han' (words with at least one kana or ASCII character)")
	flag.Parse()

	outputMode := OutputMode(*outputModeFlag)
	if outputMode != OutputAll && outputMode != OutputHanOnly && outputMode != OutputNonHanOnly {
		fmt.Fprintf(os.Stderr, "Invalid output mode: %s\n", *outputModeFlag)
		os.Exit(1)
	}

	logf := func(format string, a ...interface{}) {
		if !*silent {
			fmt.Printf(format, a...)
		}
	}

	logf("Starting initialization...\n")

	// Determine output directory based on unzipped flag and output mode
	outputDir := "output"
	if outputMode == OutputHanOnly {
		outputDir += "_han"
	} else if outputMode == OutputNonHanOnly {
		outputDir += "_non_han"
	}
	if *unzipped {
		outputDir += "_unzipped"
	}
	logf("Using output directory: %s\n", outputDir)

	// Just create the directory if it doesn't exist
	logf("Ensuring output directory exists...\n")
	start := time.Now()
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}
	logf("Directory creation took %.2fs\n", time.Since(start).Seconds())

	// Read and decompress input file
	logf("Opening input file %s...\n", *inputFile)
	start = time.Now()
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	logf("File open took %.2fs\n", time.Since(start).Seconds())

	var reader io.Reader
	if strings.HasSuffix(*inputFile, ".xz") {
		logf("Creating XZ reader...\n")
		start = time.Now()
		xzReader, err := xz.NewReader(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating XZ reader: %v\n", err)
			os.Exit(1)
		}
		reader = xzReader
		logf("XZ reader creation took %.2fs\n", time.Since(start).Seconds())
	} else {
		reader = file
	}

	logf("Parsing JSON...\n")
	start = time.Now()
	var dict JmdictTypes
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&dict); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	logf("JSON parsing took %.2fs\n", time.Since(start).Seconds())

	// Helper function to check if a string contains only Han characters
	isHanOnly := func(s string) bool {
		for _, r := range s {
			if !unicode.Is(unicode.Han, r) {
				return false
			}
		}
		return true
	}

	// Helper function to determine if a word matches our filter criteria
	wordMatchesFilter := func(word Word) bool {
		if outputMode == OutputAll {
			return true
		}

		// Get the primary text representation for filtering
		var primaryText string
		if len(word.Kanji) > 0 {
			primaryText = word.Kanji[0].Text
		} else if len(word.Kana) > 0 {
			primaryText = word.Kana[0].Text
		} else {
			primaryText = word.ID
		}

		isHan := isHanOnly(primaryText)

		if outputMode == OutputHanOnly && isHan {
			return true
		}
		if outputMode == OutputNonHanOnly && !isHan {
			return true
		}
		return false
	}

	// Filter words based on output mode
	filteredWords := make([]Word, 0)
	for _, word := range dict.Words {
		if wordMatchesFilter(word) {
			filteredWords = append(filteredWords, word)
		}
	}

	totalWords := len(filteredWords)
	logf("Processing %d words after filtering (mode: %s)...\n", totalWords, outputMode)

	// Create a map to group words by filename
	wordGroups := make(map[string]*WordGroup)
	processed := 0

	// First pass: group all words
	logf("Grouping words...\n")
	groupStart := time.Now()
	for _, word := range filteredWords {
		// Use first kanji or kana as filename, falling back to ID if neither exists
		var filename string
		if len(word.Kanji) > 0 {
			filename = word.Kanji[0].Text
		} else if len(word.Kana) > 0 {
			filename = word.Kana[0].Text
		} else {
			filename = word.ID
		}

		// Add word to appropriate group
		if group, exists := wordGroups[filename]; exists {
			group.WordJapanese = append(group.WordJapanese, word)
		} else {
			wordGroups[filename] = &WordGroup{
				WordJapanese: []Word{word},
			}
		}

		processed++
		if !*silent && processed%1000 == 0 {
			elapsed := time.Since(groupStart)
			rate := float64(processed) / elapsed.Seconds()
			logf("\rGrouping: %.1f%% (%d/%d) - %.1f words/sec",
				float64(processed)*100/float64(totalWords), processed, totalWords, rate)
		}
	}
	logf("\nGrouping completed in %.2fs\n", time.Since(groupStart).Seconds())

	// Second pass: write all groups to files
	logf("Writing %d files...\n", len(wordGroups))
	writeStart := time.Now()
	processed = 0
	totalFiles := len(wordGroups)

	for filename, group := range wordGroups {
		var fullPath string
		if *unzipped {
			fullPath = filepath.Join(outputDir, filename+".json")
		} else {
			fullPath = filepath.Join(outputDir, filename+".json.br")
		}

		// Create output file
		file, err := os.Create(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError creating file %s: %v\n", fullPath, err)
			continue
		}

		err = func() error {
			defer file.Close()

			var writer interface {
				Write(p []byte) (n int, err error)
				Close() error
			}

			if *unzipped {
				writer = file
			} else {
				writer = brotli.NewWriter(file)
			}
			defer writer.Close()

			encoder := json.NewEncoder(writer)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(group); err != nil {
				return fmt.Errorf("failed to encode group: %v", err)
			}

			return nil
		}()

		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError processing %s: %v\n", fullPath, err)
		}

		processed++
		if !*silent && processed%100 == 0 {
			elapsed := time.Since(writeStart)
			rate := float64(processed) / elapsed.Seconds()
			logf("\rWriting: %.1f%% (%d/%d) - %.1f files/sec",
				float64(processed)*100/float64(totalFiles), processed, totalFiles, rate)
		}
	}

	logf("\nCompleted in %.1fs (%.1f files/sec)\n",
		time.Since(writeStart).Seconds(),
		float64(totalFiles)/time.Since(writeStart).Seconds())
}
