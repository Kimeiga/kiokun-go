package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"unicode"

	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	_ "kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	_ "kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
	_ "kiokun-go/dictionaries/kanjidic"

	// Import Chinese dictionaries
	"kiokun-go/dictionaries/chinese_chars"
	_ "kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	_ "kiokun-go/dictionaries/chinese_words"
	"kiokun-go/processor"
)

// OutputMode determines which words to output
type OutputMode string

const (
	OutputAll        OutputMode = "all"       // Output all words
	OutputHanOnly    OutputMode = "han-only"  // Output words with only Han characters (legacy mode)
	OutputHan1Char   OutputMode = "han-1char" // Output words with exactly 1 Han character
	OutputHan2Char   OutputMode = "han-2char" // Output words with exactly 2 Han characters
	OutputHan3Plus   OutputMode = "han-3plus" // Output words with 3 or more Han characters
	OutputNonHanOnly OutputMode = "non-han"   // Output words with at least one non-Han character
)

func main() {
	// Configuration flags
	dictDir := flag.String("dictdir", "dictionaries", "Base directory containing dictionary packages")
	outputDir := flag.String("outdir", "output", "Output directory for processed files")
	workers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines for batch processing")
	fileWriters := flag.Int("writers", runtime.NumCPU(), "Number of parallel workers for file writing")
	silent := flag.Bool("silent", false, "Disable progress output")
	devMode := flag.Bool("dev", false, "Development mode - use /tmp for faster I/O")
	limitEntries := flag.Int("limit", 0, "Limit the number of entries to process (0 = no limit)")
	batchSize := flag.Int("batch", 10000, "Process entries in batches of this size")
	outputModeFlag := flag.String("mode", "all", "Output mode: 'all', 'han-only' (legacy), 'han-1char', 'han-2char', 'han-3plus', or 'non-han'")
	flag.Parse()

	// Parse and validate the output mode
	outputMode := OutputMode(*outputModeFlag)
	if outputMode != OutputAll &&
		outputMode != OutputHanOnly &&
		outputMode != OutputHan1Char &&
		outputMode != OutputHan2Char &&
		outputMode != OutputHan3Plus &&
		outputMode != OutputNonHanOnly {
		fmt.Fprintf(os.Stderr, "Invalid output mode: %s\n", *outputModeFlag)
		fmt.Fprintf(os.Stderr, "Valid modes: all, han-only, han-1char, han-2char, han-3plus, non-han\n")
		os.Exit(1)
	}

	// Modify output directory based on mode
	if outputMode == OutputHanOnly {
		*outputDir = *outputDir + "_han"
	} else if outputMode == OutputHan1Char {
		*outputDir = *outputDir + "_han_1char"
	} else if outputMode == OutputHan2Char {
		*outputDir = *outputDir + "_han_2char"
	} else if outputMode == OutputHan3Plus {
		*outputDir = *outputDir + "_han_3plus"
	} else if outputMode == OutputNonHanOnly {
		*outputDir = *outputDir + "_non_han"
	}

	logf := func(format string, a ...interface{}) {
		if !*silent {
			fmt.Printf(format, a...)
		}
	}

	// If dev mode is enabled, use /tmp directory for output
	if *devMode {
		tmpDir := filepath.Join("/tmp", "kiokun-output")
		logf("Development mode enabled: using %s for output\n", tmpDir)
		*outputDir = tmpDir

		// Create the tmp directory if it doesn't exist
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating tmp directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Ensure we're using an absolute path for the output directory
	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
		os.Exit(1)
	}
	*outputDir = absOutputDir

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	logf("Starting initialization...\n")
	logf("Input dictionaries directory: %s\n", *dictDir)
	logf("Output directory: %s\n", *outputDir)
	logf("Using %d processing workers and %d file writers\n", *workers, *fileWriters)
	if outputMode != OutputAll {
		logf("Filtering mode: %s\n", outputMode)
	}

	// Set dictionaries base path - find workspace root
	workspaceRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// If we're in cmd/kiokun, go up two levels
	if filepath.Base(workspaceRoot) == "kiokun" && filepath.Base(filepath.Dir(workspaceRoot)) == "cmd" {
		workspaceRoot = filepath.Dir(filepath.Dir(workspaceRoot))
	}

	// Resolve dictionaries path
	dictPath := filepath.Join(workspaceRoot, *dictDir)
	logf("Using dictionary path: %s\n", dictPath)
	common.SetDictionariesBasePath(dictPath)

	// Create dictionary processor with parallel file writing
	proc, err := processor.New(*outputDir, *fileWriters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating processor: %v\n", err)
		os.Exit(1)
	}

	// Import all dictionaries
	logf("Importing dictionaries...\n")

	// Get all registered dictionaries
	dictConfigs := common.GetRegisteredDictionaries()

	// Import each dictionary
	var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries []common.Entry

	for _, dict := range dictConfigs {
		// Construct full path
		inputPath := filepath.Join(dict.SourceDir, dict.InputFile)

		// Import this dictionary
		logf("Importing %s from %s...\n", dict.Name, inputPath)
		startTime := time.Now()

		entries, err := dict.Importer.Import(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing %s: %v\n", dict.Name, err)
			os.Exit(1)
		}

		// Store entries by dictionary type
		switch dict.Name {
		case "jmdict":
			jmdictEntries = entries
		case "jmnedict":
			jmnedictEntries = entries
		case "kanjidic":
			kanjidicEntries = entries
		case "chinese_chars":
			chineseCharsEntries = entries
		case "chinese_words":
			chineseWordsEntries = entries
		}

		logf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(startTime).Seconds())
	}

	// Helper function to check if a string contains only Han characters
	isHanOnly := func(s string) bool {
		for _, r := range s {
			if !unicode.Is(unicode.Han, r) {
				return false
			}
		}
		return true
	}

	// Helper function to check if an entry should be included based on the mode
	shouldIncludeEntry := func(entry common.Entry) bool {
		if outputMode == OutputAll {
			return true
		}

		// Get primary text representation for filtering
		var primaryText string
		switch e := entry.(type) {
		case jmdict.Word:
			if len(e.Kanji) > 0 {
				primaryText = e.Kanji[0].Text
			} else if len(e.Kana) > 0 {
				primaryText = e.Kana[0].Text
			} else {
				primaryText = e.ID
			}
		case jmnedict.Name:
			// Use the Name's primary text
			if len(e.Kanji) > 0 {
				primaryText = e.Kanji[0]
			} else if len(e.Reading) > 0 {
				primaryText = e.Reading[0]
			} else {
				primaryText = e.ID
			}
		case kanjidic.Kanji:
			// Use the Kanji character
			primaryText = e.Character
		case chinese_chars.ChineseCharEntry:
			// Use the traditional character
			primaryText = e.Traditional
		case chinese_words.ChineseWordEntry:
			// Use the traditional word
			primaryText = e.Traditional
		default:
			// If we don't know how to filter this type, include it by default
			return true
		}

		// Check if the text contains only Han characters
		isHan := isHanOnly(primaryText)
		charCount := len([]rune(primaryText)) // Get correct Unicode character count

		// Apply filtering based on mode
		switch outputMode {
		case OutputNonHanOnly:
			return !isHan
		case OutputHanOnly:
			return isHan
		case OutputHan1Char:
			return isHan && charCount == 1
		case OutputHan2Char:
			return isHan && charCount == 2
		case OutputHan3Plus:
			return isHan && charCount >= 3
		default:
			return true
		}
	}

	// Filter entries based on output mode
	if outputMode != OutputAll {
		filteredJmdict := make([]common.Entry, 0, len(jmdictEntries))
		filteredJmnedict := make([]common.Entry, 0, len(jmnedictEntries))
		filteredKanjidic := make([]common.Entry, 0, len(kanjidicEntries))
		filteredChineseChars := make([]common.Entry, 0, len(chineseCharsEntries))
		filteredChineseWords := make([]common.Entry, 0, len(chineseWordsEntries))

		for _, entry := range jmdictEntries {
			if shouldIncludeEntry(entry) {
				filteredJmdict = append(filteredJmdict, entry)
			}
		}
		for _, entry := range jmnedictEntries {
			if shouldIncludeEntry(entry) {
				filteredJmnedict = append(filteredJmnedict, entry)
			}
		}
		for _, entry := range kanjidicEntries {
			if shouldIncludeEntry(entry) {
				filteredKanjidic = append(filteredKanjidic, entry)
			}
		}
		for _, entry := range chineseCharsEntries {
			if shouldIncludeEntry(entry) {
				filteredChineseChars = append(filteredChineseChars, entry)
			}
		}
		for _, entry := range chineseWordsEntries {
			if shouldIncludeEntry(entry) {
				filteredChineseWords = append(filteredChineseWords, entry)
			}
		}

		logf("Filtered entries - JMdict: %d -> %d, JMNedict: %d -> %d, Kanjidic: %d -> %d, Chinese Chars: %d -> %d, Chinese Words: %d -> %d\n",
			len(jmdictEntries), len(filteredJmdict),
			len(jmnedictEntries), len(filteredJmnedict),
			len(kanjidicEntries), len(filteredKanjidic),
			len(chineseCharsEntries), len(filteredChineseChars),
			len(chineseWordsEntries), len(filteredChineseWords))

		jmdictEntries = filteredJmdict
		jmnedictEntries = filteredJmnedict
		kanjidicEntries = filteredKanjidic
		chineseCharsEntries = filteredChineseChars
		chineseWordsEntries = filteredChineseWords
	}

	// Apply entry limit if specified
	var allEntries []common.Entry

	if *limitEntries > 0 {
		totalEntries := len(jmdictEntries) + len(jmnedictEntries) + len(kanjidicEntries) +
			len(chineseCharsEntries) + len(chineseWordsEntries)
		if *limitEntries < totalEntries {
			logf("Limiting to %d entries (out of %d total)\n", *limitEntries, totalEntries)

			// Calculate proportions
			jmdictProportion := float64(len(jmdictEntries)) / float64(totalEntries)
			jmnedictProportion := float64(len(jmnedictEntries)) / float64(totalEntries)
			kanjidicProportion := float64(len(kanjidicEntries)) / float64(totalEntries)
			chineseCharsProportion := float64(len(chineseCharsEntries)) / float64(totalEntries)
			chineseWordsProportion := float64(len(chineseWordsEntries)) / float64(totalEntries)

			// Calculate limits for each dictionary
			jmdictLimit := int(float64(*limitEntries) * jmdictProportion)
			jmnedictLimit := int(float64(*limitEntries) * jmnedictProportion)
			kanjidicLimit := int(float64(*limitEntries) * kanjidicProportion)
			chineseCharsLimit := int(float64(*limitEntries) * chineseCharsProportion)
			chineseWordsLimit := int(float64(*limitEntries) * chineseWordsProportion)

			// Adjust for rounding errors
			remaining := *limitEntries - jmdictLimit - jmnedictLimit - kanjidicLimit -
				chineseCharsLimit - chineseWordsLimit
			if remaining > 0 && len(kanjidicEntries) > kanjidicLimit {
				kanjidicLimit += remaining
			}

			// Apply limits
			if jmdictLimit < len(jmdictEntries) {
				jmdictEntries = jmdictEntries[:jmdictLimit]
			}
			if jmnedictLimit < len(jmnedictEntries) {
				jmnedictEntries = jmnedictEntries[:jmnedictLimit]
			}
			if kanjidicLimit < len(kanjidicEntries) {
				kanjidicEntries = kanjidicEntries[:kanjidicLimit]
			}
			if chineseCharsLimit < len(chineseCharsEntries) {
				chineseCharsEntries = chineseCharsEntries[:chineseCharsLimit]
			}
			if chineseWordsLimit < len(chineseWordsEntries) {
				chineseWordsEntries = chineseWordsEntries[:chineseWordsLimit]
			}
		}
	}

	// Combine all entries
	allEntries = append(allEntries, jmdictEntries...)
	allEntries = append(allEntries, jmnedictEntries...)
	allEntries = append(allEntries, kanjidicEntries...)
	allEntries = append(allEntries, chineseCharsEntries...)
	allEntries = append(allEntries, chineseWordsEntries...)

	// Count entries by type for debugging
	jmdictCount := 0
	jmnedictCount := 0
	kanjidicCount := 0
	chineseCharsCount := 0
	chineseWordsCount := 0

	for _, entry := range allEntries {
		switch entry.(type) {
		case jmdict.Word:
			jmdictCount++
		case jmnedict.Name:
			jmnedictCount++
		case kanjidic.Kanji:
			kanjidicCount++
		case chinese_chars.ChineseCharEntry:
			chineseCharsCount++
		case chinese_words.ChineseWordEntry:
			chineseWordsCount++
		default:
			logf("Unknown entry type: %T\n", entry)
		}
	}
	logf("Processing %d entries (%d JMdict, %d JMNedict, %d Kanjidic, %d Chinese Chars, %d Chinese Words)\n",
		len(allEntries), jmdictCount, jmnedictCount, kanjidicCount, chineseCharsCount, chineseWordsCount)

	// Process entries in batches with progress reporting
	logf("Processing entries in batches...\n")
	totalEntries := len(allEntries)
	processStart := time.Now()

	for start := 0; start < totalEntries; start += *batchSize {
		end := start + *batchSize
		if end > totalEntries {
			end = totalEntries
		}

		batchEntries := allEntries[start:end]

		// Only update progress every 10 batches to reduce output
		if start%((*batchSize)*10) == 0 || end == totalEntries {
			logf("\rProcessing entries %d-%d of %d (%.1f%%)...",
				start+1, end, totalEntries, float64(end)/float64(totalEntries)*100)
		}

		// Process entries sequentially in a batch
		for _, entry := range batchEntries {
			if err := proc.ProcessEntries([]common.Entry{entry}); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing entry %s: %v\n", entry.GetID(), err)
				os.Exit(1)
			}
		}
	}
	processDuration := time.Since(processStart)
	logf("\rProcessed all %d entries (100%%) in %.2f seconds (%.1f entries/sec)\n",
		totalEntries, processDuration.Seconds(), float64(totalEntries)/processDuration.Seconds())

	// Write all processed entries to files
	logf("Writing files to %s...\n", *outputDir)
	if err := proc.WriteToFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing files: %v\n", err)
		os.Exit(1)
	}

	logf("Successfully processed dictionary files\n")
}
