package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

// Config holds all configuration options for the application
type Config struct {
	DictDir       string
	OutputDir     string
	Workers       int
	FileWriters   int
	Silent        bool
	DevMode       bool
	LimitEntries  int
	BatchSize     int
	OutputMode    OutputMode
	TestMode      bool
	WorkspaceRoot string
	// UseIndexMode removed - always using index-based approach

	// Dictionary selection flags
	OnlyJMdict       bool
	OnlyJMNedict     bool
	OnlyKanjidic     bool
	OnlyChineseChars bool
	OnlyChineseWords bool
	OnlyIDS          bool
}

// LogFunc is a function that logs messages based on silent mode
type LogFunc func(format string, a ...interface{})

// NewLogFunc creates a new logging function based on silent mode
func NewLogFunc(silent bool) LogFunc {
	return func(format string, a ...interface{}) {
		if !silent {
			fmt.Printf(format, a...)
		}
	}
}

// ParseConfig parses command-line flags and returns a Config struct
func ParseConfig() (*Config, LogFunc, error) {
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
	testMode := flag.Bool("test", false, "Test mode - prioritize entries that have overlap between Chinese and Japanese dictionaries")
	// Index mode flag removed - always using index-based approach

	// Dictionary selection flags
	onlyJMdict := flag.Bool("only-jmdict", false, "Process only JMdict (Japanese words)")
	onlyJMNedict := flag.Bool("only-jmnedict", false, "Process only JMNedict (Japanese names)")
	onlyKanjidic := flag.Bool("only-kanjidic", false, "Process only Kanjidic (Japanese kanji)")
	onlyChineseChars := flag.Bool("only-chinese-chars", false, "Process only Chinese characters")
	onlyChineseWords := flag.Bool("only-chinese-words", false, "Process only Chinese words")
	onlyIDS := flag.Bool("only-ids", false, "Process only IDS (Ideographic Description Sequences)")
	flag.Parse()

	// Create logging function
	logf := NewLogFunc(*silent)

	// Parse and validate the output mode
	outputMode := OutputMode(*outputModeFlag)
	if outputMode != OutputAll &&
		outputMode != OutputHanOnly &&
		outputMode != OutputHan1Char &&
		outputMode != OutputHan2Char &&
		outputMode != OutputHan3Plus &&
		outputMode != OutputNonHanOnly {
		return nil, logf, fmt.Errorf("invalid output mode: %s", *outputModeFlag)
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

	// Log the output directory
	logf("Output directory: %s\n", *outputDir)

	// If dev mode is enabled, use /tmp directory for output
	if *devMode {
		tmpDir := filepath.Join("/tmp", "kiokun-output")
		logf("Development mode enabled: using %s for output\n", tmpDir)
		*outputDir = tmpDir

		// Create the tmp directory if it doesn't exist
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return nil, logf, fmt.Errorf("error creating tmp directory: %v", err)
		}
	}

	// Ensure we're using an absolute path for the output directory
	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		return nil, logf, fmt.Errorf("error resolving output directory path: %v", err)
	}
	*outputDir = absOutputDir

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		return nil, logf, fmt.Errorf("error creating output directory: %v", err)
	}

	// Set dictionaries base path - find workspace root
	workspaceRoot, err := os.Getwd()
	if err != nil {
		return nil, logf, fmt.Errorf("error getting current directory: %v", err)
	}

	// If we're in cmd/kiokun, go up two levels
	if filepath.Base(workspaceRoot) == "kiokun" && filepath.Base(filepath.Dir(workspaceRoot)) == "cmd" {
		workspaceRoot = filepath.Dir(filepath.Dir(workspaceRoot))
	}

	return &Config{
		DictDir:       *dictDir,
		OutputDir:     *outputDir,
		Workers:       *workers,
		FileWriters:   *fileWriters,
		Silent:        *silent,
		DevMode:       *devMode,
		LimitEntries:  *limitEntries,
		BatchSize:     *batchSize,
		OutputMode:    outputMode,
		TestMode:      *testMode,
		WorkspaceRoot: workspaceRoot,

		// Dictionary selection flags
		OnlyJMdict:       *onlyJMdict,
		OnlyJMNedict:     *onlyJMNedict,
		OnlyKanjidic:     *onlyKanjidic,
		OnlyChineseChars: *onlyChineseChars,
		OnlyChineseWords: *onlyChineseWords,
		OnlyIDS:          *onlyIDS,
	}, logf, nil
}
