diff --git a/README.md b/README.md
index 199b67bba1..4662e8b34b 100644
--- a/README.md
+++ b/README.md
@@ -1,22 +1,49 @@

# Kiokun-Go

-A Japanese dictionary processor that combines multiple dictionary sources (JMdict, JMNedict, and Kanjidic) into a unified format.
+A multilingual dictionary processor that combines Japanese and Chinese dictionary sources (JMdict, JMNedict, Kanjidic, Chinese characters, and Chinese words) into a unified format.

## Project Structure

- `cmd/kiokun/` - Main application

* - `main.go` - Entry point
* - `config.go` - Configuration and flag handling
* - `loader.go` - Dictionary loading functionality
* - `filter.go` - Entry filtering logic
* - `processor.go` - Processing and output logic

- `cmd/debug_*` - Debug tools for testing specific components
- `dictionaries/` - Dictionary importers and data structures
  - `common/` - Common interfaces and utilities
  - `jmdict/` - JMdict dictionary importer
  - `jmnedict/` - JMNedict dictionary importer
  - `kanjidic/` - Kanjidic dictionary importer

* - `chinese_chars/` - Chinese character dictionary importer
* - `chinese_words/` - Chinese word dictionary importer

- `processor/` - Dictionary processing logic

* - `processor.go` - Core processor functionality
* - `processor_japanese.go` - Japanese dictionary processing functions
* - `processor_chinese.go` - Chinese dictionary processing functions
* - `types.go` - Type definitions
* - `sanitize.go` - Data sanitization functions
* - `writer.go` - File writing functionality

- `analyzer/` - Dictionary analysis tools

## Usage

-### Main Application
+### Building the Full Dictionary

- +To build the complete dictionary with all entries:
- +`bash
+go run cmd/kiokun/main.go
+`
- +This will:
- +1. Import all dictionaries from the `dictionaries/` directory
  +2. Process and combine entries
  +3. Write compressed JSON files to the `output/` directory
- +### Command-Line Options
  ```bash
  go run cmd/kiokun/main.go [options]
  @@ -32,6 +59,35 @@ Options:
  - `--dev` - Development mode - use /tmp for faster I/O
  - `--limit <n>` - Limit the number of entries to process (0 = no limit)
  - `--batch <n>` - Process entries in batches of this size (default: 10000)
  +- `--mode <mode>` - Output mode: 'all', 'han-only', 'han-1char', 'han-2char', 'han-3plus', or 'non-han'
  +- `--test` - Test mode - prioritize entries that have overlap between Chinese and Japanese dictionaries
  ```
- +### Filtering Modes
- +The `--mode` flag allows filtering entries based on their character composition:
- +- `all` - Include all entries (default)
  +- `han-only` - Include only entries with Han characters
  +- `han-1char` - Include only entries with exactly 1 Han character
  +- `han-2char` - Include only entries with exactly 2 Han characters
  +- `han-3plus` - Include only entries with 3 or more Han characters
  +- `non-han` - Include only entries with at least one non-Han character
- +Example:
- +`bash
+go run cmd/kiokun/main.go --mode han-1char
+`
- +### Test Mode
- +The `--test` flag enables test mode, which prioritizes entries that have overlap between Chinese and Japanese dictionaries:
- +`bash
+go run cmd/kiokun/main.go --test --limit 100
+`
- +This is useful for testing the integration between Chinese and Japanese dictionaries without processing the entire dataset.
  ### Debug Tools

@@ -48,16 +104,33 @@ Various debug tools are available in the `cmd/` directory:

### Building

+Build the main application:

- ```bash
  go build -o kiokun cmd/kiokun/main.go
  ```
  ### Testing
  +Run all tests:
- ```bash
  go test ./...
  ```
  +Run specific tests:
- +```bash
  +# Test the processor package
  +go test ./processor
- +# Test the Chinese-Japanese integration
  +go test ./processor -run TestChineseJapaneseIntegration
- +# Test with verbose output
  +go test -v ./processor
  +```
- ### Performance Optimization
  For faster development, use the `--dev` flag to write output to `/tmp` for better I/O performance:
  @@ -65,3 +138,19 @@ For faster development, use the `--dev` flag to write output to `/tmp` for bette
  ```bash
  go run cmd/kiokun/main.go --dev --writers 16
  ```
- +Limit the number of entries for faster testing:
- +`bash
+go run cmd/kiokun/main.go --limit 1000 --dev
+`
- +### GitHub Actions Integration
- +The integration test (`TestChineseJapaneseIntegration`) is designed to run quickly in GitHub Actions by:
- +1. Creating test entries for characters that exist in both Chinese and Japanese dictionaries
  +2. Processing these entries through the processor
  +3. Verifying that the output files contain combined data from both dictionaries
- +This ensures that the core functionality of combining Chinese and Japanese dictionaries works correctly without processing the entire dataset.
  diff --git a/cmd/kiokun/internal/config.go b/cmd/kiokun/internal/config.go
  new file mode 100644
  index 0000000000..c0cb28e2ed
  --- /dev/null
  +++ b/cmd/kiokun/internal/config.go
  @@ -0,0 +1,140 @@
  +package internal
- +import (
- "flag"
- "fmt"
- "os"
- "path/filepath"
- "runtime"
  +)
- +// OutputMode determines which words to output
  +type OutputMode string
- +const (
- OutputAll OutputMode = "all" // Output all words
- OutputHanOnly OutputMode = "han-only" // Output words with only Han characters (legacy mode)
- OutputHan1Char OutputMode = "han-1char" // Output words with exactly 1 Han character
- OutputHan2Char OutputMode = "han-2char" // Output words with exactly 2 Han characters
- OutputHan3Plus OutputMode = "han-3plus" // Output words with 3 or more Han characters
- OutputNonHanOnly OutputMode = "non-han" // Output words with at least one non-Han character
  +)
- +// Config holds all configuration options for the application
  +type Config struct {
- DictDir string
- OutputDir string
- Workers int
- FileWriters int
- Silent bool
- DevMode bool
- LimitEntries int
- BatchSize int
- OutputMode OutputMode
- TestMode bool
- WorkspaceRoot string
  +}
- +// LogFunc is a function that logs messages based on silent mode
  +type LogFunc func(format string, a ...interface{})
- +// NewLogFunc creates a new logging function based on silent mode
  +func NewLogFunc(silent bool) LogFunc {
- return func(format string, a ...interface{}) {
-     if !silent {
-     	fmt.Printf(format, a...)
-     }
- }
  +}
- +// ParseConfig parses command-line flags and returns a Config struct
  +func ParseConfig() (\*Config, LogFunc, error) {
- // Configuration flags
- dictDir := flag.String("dictdir", "dictionaries", "Base directory containing dictionary packages")
- outputDir := flag.String("outdir", "output", "Output directory for processed files")
- workers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines for batch processing")
- fileWriters := flag.Int("writers", runtime.NumCPU(), "Number of parallel workers for file writing")
- silent := flag.Bool("silent", false, "Disable progress output")
- devMode := flag.Bool("dev", false, "Development mode - use /tmp for faster I/O")
- limitEntries := flag.Int("limit", 0, "Limit the number of entries to process (0 = no limit)")
- batchSize := flag.Int("batch", 10000, "Process entries in batches of this size")
- outputModeFlag := flag.String("mode", "all", "Output mode: 'all', 'han-only' (legacy), 'han-1char', 'han-2char', 'han-3plus', or 'non-han'")
- testMode := flag.Bool("test", false, "Test mode - prioritize entries that have overlap between Chinese and Japanese dictionaries")
- flag.Parse()
-
- // Create logging function
- logf := NewLogFunc(\*silent)
-
- // Parse and validate the output mode
- outputMode := OutputMode(\*outputModeFlag)
- if outputMode != OutputAll &&
-     outputMode != OutputHanOnly &&
-     outputMode != OutputHan1Char &&
-     outputMode != OutputHan2Char &&
-     outputMode != OutputHan3Plus &&
-     outputMode != OutputNonHanOnly {
-     return nil, logf, fmt.Errorf("invalid output mode: %s", *outputModeFlag)
- }
-
- // Modify output directory based on mode
- if outputMode == OutputHanOnly {
-     *outputDir = *outputDir + "_han"
- } else if outputMode == OutputHan1Char {
-     *outputDir = *outputDir + "_han_1char"
- } else if outputMode == OutputHan2Char {
-     *outputDir = *outputDir + "_han_2char"
- } else if outputMode == OutputHan3Plus {
-     *outputDir = *outputDir + "_han_3plus"
- } else if outputMode == OutputNonHanOnly {
-     *outputDir = *outputDir + "_non_han"
- }
-
- // If dev mode is enabled, use /tmp directory for output
- if \*devMode {
-     tmpDir := filepath.Join("/tmp", "kiokun-output")
-     logf("Development mode enabled: using %s for output\n", tmpDir)
-     *outputDir = tmpDir
-
-     // Create the tmp directory if it doesn't exist
-     if err := os.MkdirAll(tmpDir, 0755); err != nil {
-     	return nil, logf, fmt.Errorf("error creating tmp directory: %v", err)
-     }
- }
-
- // Ensure we're using an absolute path for the output directory
- absOutputDir, err := filepath.Abs(\*outputDir)
- if err != nil {
-     return nil, logf, fmt.Errorf("error resolving output directory path: %v", err)
- }
- \*outputDir = absOutputDir
-
- // Create the output directory if it doesn't exist
- if err := os.MkdirAll(\*outputDir, 0755); err != nil {
-     return nil, logf, fmt.Errorf("error creating output directory: %v", err)
- }
-
- // Set dictionaries base path - find workspace root
- workspaceRoot, err := os.Getwd()
- if err != nil {
-     return nil, logf, fmt.Errorf("error getting current directory: %v", err)
- }
-
- // If we're in cmd/kiokun, go up two levels
- if filepath.Base(workspaceRoot) == "kiokun" && filepath.Base(filepath.Dir(workspaceRoot)) == "cmd" {
-     workspaceRoot = filepath.Dir(filepath.Dir(workspaceRoot))
- }
-
- return &Config{
-     DictDir:       *dictDir,
-     OutputDir:     *outputDir,
-     Workers:       *workers,
-     FileWriters:   *fileWriters,
-     Silent:        *silent,
-     DevMode:       *devMode,
-     LimitEntries:  *limitEntries,
-     BatchSize:     *batchSize,
-     OutputMode:    outputMode,
-     TestMode:      *testMode,
-     WorkspaceRoot: workspaceRoot,
- }, logf, nil
  +}
  diff --git a/cmd/kiokun/internal/filter.go b/cmd/kiokun/internal/filter.go
  new file mode 100644
  index 0000000000..9a173ecc11
  --- /dev/null
  +++ b/cmd/kiokun/internal/filter.go
  @@ -0,0 +1,382 @@
  +package internal
- +import (
- "strings"
- "unicode"
-
- "kiokun-go/dictionaries/chinese_chars"
- "kiokun-go/dictionaries/chinese_words"
- "kiokun-go/dictionaries/common"
- "kiokun-go/dictionaries/jmdict"
- "kiokun-go/dictionaries/jmnedict"
- "kiokun-go/dictionaries/kanjidic"
  +)
- +// FilterEntries filters dictionary entries based on the configuration
  +func FilterEntries(entries *DictionaryEntries, config *Config, logf LogFunc) \*DictionaryEntries {
- result := &DictionaryEntries{
-     JMdict:       entries.JMdict,
-     JMNedict:     entries.JMNedict,
-     Kanjidic:     entries.Kanjidic,
-     ChineseChars: entries.ChineseChars,
-     ChineseWords: entries.ChineseWords,
- }
-
- // Filter entries based on output mode
- if config.OutputMode != OutputAll {
-     result = filterByOutputMode(result, config.OutputMode, logf)
- }
-
- // Apply test mode filtering if enabled
- if config.TestMode {
-     result = filterForTestMode(result, logf)
- }
-
- // Apply entry limit if specified
- if config.LimitEntries > 0 {
-     result = limitEntries(result, config.LimitEntries, logf)
- }
-
- return result
  +}
- +// filterByOutputMode filters entries based on the output mode
  +func filterByOutputMode(entries *DictionaryEntries, mode OutputMode, logf LogFunc) *DictionaryEntries {
- // Helper function to check if a string contains only Han characters
- isHanOnly := func(s string) bool {
-     for _, r := range s {
-     	if !unicode.Is(unicode.Han, r) {
-     		return false
-     	}
-     }
-     return true
- }
-
- // Helper function to check if an entry should be included based on the mode
- shouldIncludeEntry := func(entry common.Entry) bool {
-     // Get primary text representation for filtering
-     var primaryText string
-     switch e := entry.(type) {
-     case jmdict.Word:
-     	if len(e.Kanji) > 0 {
-     		primaryText = e.Kanji[0].Text
-     	} else if len(e.Kana) > 0 {
-     		primaryText = e.Kana[0].Text
-     	} else {
-     		primaryText = e.ID
-     	}
-     case jmnedict.Name:
-     	// Use the Name's primary text
-     	if len(e.Kanji) > 0 {
-     		primaryText = e.Kanji[0]
-     	} else if len(e.Reading) > 0 {
-     		primaryText = e.Reading[0]
-     	} else {
-     		primaryText = e.ID
-     	}
-     case kanjidic.Kanji:
-     	// Use the Kanji character
-     	primaryText = e.Character
-     case chinese_chars.ChineseCharEntry:
-     	// Use the traditional character
-     	primaryText = e.Traditional
-     case chinese_words.ChineseWordEntry:
-     	// Use the traditional word
-     	primaryText = e.Traditional
-     default:
-     	// If we don't know how to filter this type, include it by default
-     	return true
-     }
-
-     // Check if the text contains only Han characters
-     isHan := isHanOnly(primaryText)
-     charCount := len([]rune(primaryText)) // Get correct Unicode character count
-
-     // Apply filtering based on mode
-     switch mode {
-     case OutputNonHanOnly:
-     	return !isHan
-     case OutputHanOnly:
-     	return isHan
-     case OutputHan1Char:
-     	return isHan && charCount == 1
-     case OutputHan2Char:
-     	return isHan && charCount == 2
-     case OutputHan3Plus:
-     	return isHan && charCount >= 3
-     default:
-     	return true
-     }
- }
-
- // Filter each dictionary
- filteredJmdict := make([]common.Entry, 0, len(entries.JMdict))
- filteredJmnedict := make([]common.Entry, 0, len(entries.JMNedict))
- filteredKanjidic := make([]common.Entry, 0, len(entries.Kanjidic))
- filteredChineseChars := make([]common.Entry, 0, len(entries.ChineseChars))
- filteredChineseWords := make([]common.Entry, 0, len(entries.ChineseWords))
-
- for \_, entry := range entries.JMdict {
-     if shouldIncludeEntry(entry) {
-     	filteredJmdict = append(filteredJmdict, entry)
-     }
- }
- for \_, entry := range entries.JMNedict {
-     if shouldIncludeEntry(entry) {
-     	filteredJmnedict = append(filteredJmnedict, entry)
-     }
- }
- for \_, entry := range entries.Kanjidic {
-     if shouldIncludeEntry(entry) {
-     	filteredKanjidic = append(filteredKanjidic, entry)
-     }
- }
- for \_, entry := range entries.ChineseChars {
-     if shouldIncludeEntry(entry) {
-     	filteredChineseChars = append(filteredChineseChars, entry)
-     }
- }
- for \_, entry := range entries.ChineseWords {
-     if shouldIncludeEntry(entry) {
-     	filteredChineseWords = append(filteredChineseWords, entry)
-     }
- }
-
- logf("Filtered entries - JMdict: %d -> %d, JMNedict: %d -> %d, Kanjidic: %d -> %d, Chinese Chars: %d -> %d, Chinese Words: %d -> %d\n",
-     len(entries.JMdict), len(filteredJmdict),
-     len(entries.JMNedict), len(filteredJmnedict),
-     len(entries.Kanjidic), len(filteredKanjidic),
-     len(entries.ChineseChars), len(filteredChineseChars),
-     len(entries.ChineseWords), len(filteredChineseWords))
-
- return &DictionaryEntries{
-     JMdict:       filteredJmdict,
-     JMNedict:     filteredJmnedict,
-     Kanjidic:     filteredKanjidic,
-     ChineseChars: filteredChineseChars,
-     ChineseWords: filteredChineseWords,
- }
  +}
- +// filterForTestMode prioritizes entries that have overlap between Chinese and Japanese dictionaries
  +func filterForTestMode(entries *DictionaryEntries, logf LogFunc) *DictionaryEntries {
- logf("Test mode enabled - prioritizing entries with overlap between Chinese and Japanese dictionaries\n")
-
- // Create maps to track characters in each dictionary
- japaneseChars := make(map[string]bool)
- chineseChars := make(map[string]bool)
-
- // Collect Japanese characters
- for \_, entry := range entries.Kanjidic {
-     kanji, ok := entry.(kanjidic.Kanji)
-     if ok {
-     	japaneseChars[kanji.Character] = true
-     }
- }
-
- // Collect Chinese characters
- for \_, entry := range entries.ChineseChars {
-     char, ok := entry.(chinese_chars.ChineseCharEntry)
-     if ok {
-     	chineseChars[char.Traditional] = true
-     	if char.Simplified != char.Traditional {
-     		chineseChars[char.Simplified] = true
-     	}
-     }
- }
-
- // Find common characters
- var commonCharacters []string
- for char := range japaneseChars {
-     if chineseChars[char] {
-     	commonCharacters = append(commonCharacters, char)
-     }
- }
-
- logf("Found %d common characters between Chinese and Japanese dictionaries\n", len(commonCharacters))
-
- if len(commonCharacters) == 0 {
-     // No common characters found, return original entries
-     return entries
- }
-
- // Filter entries to prioritize those with common characters
- var prioritizedJmdictEntries []common.Entry
- var prioritizedJmnedictEntries []common.Entry
- var prioritizedKanjidicEntries []common.Entry
- var prioritizedChineseCharsEntries []common.Entry
- var prioritizedChineseWordsEntries []common.Entry
-
- // Helper function to check if an entry contains a common character
- containsCommonChar := func(text string) bool {
-     for _, commonChar := range commonCharacters {
-     	if strings.Contains(text, commonChar) {
-     		return true
-     	}
-     }
-     return false
- }
-
- // Filter Kanjidic entries
- for \_, entry := range entries.Kanjidic {
-     kanji, ok := entry.(kanjidic.Kanji)
-     if ok {
-     	for _, commonChar := range commonCharacters {
-     		if kanji.Character == commonChar {
-     			prioritizedKanjidicEntries = append(prioritizedKanjidicEntries, entry)
-     			break
-     		}
-     	}
-     }
- }
-
- // Filter Chinese character entries
- for \_, entry := range entries.ChineseChars {
-     char, ok := entry.(chinese_chars.ChineseCharEntry)
-     if ok {
-     	for _, commonChar := range commonCharacters {
-     		if char.Traditional == commonChar || char.Simplified == commonChar {
-     			prioritizedChineseCharsEntries = append(prioritizedChineseCharsEntries, entry)
-     			break
-     		}
-     	}
-     }
- }
-
- // Filter JMdict entries
- for \_, entry := range entries.JMdict {
-     word, ok := entry.(jmdict.Word)
-     if ok {
-     	// Check if any kanji form contains a common character
-     	found := false
-     	for _, kanji := range word.Kanji {
-     		if containsCommonChar(kanji.Text) {
-     			found = true
-     			break
-     		}
-     	}
-     	if found {
-     		prioritizedJmdictEntries = append(prioritizedJmdictEntries, entry)
-     	}
-     }
- }
-
- // Filter JMNedict entries
- for \_, entry := range entries.JMNedict {
-     name, ok := entry.(jmnedict.Name)
-     if ok {
-     	// Check if any kanji form contains a common character
-     	found := false
-     	for _, kanji := range name.Kanji {
-     		if containsCommonChar(kanji) {
-     			found = true
-     			break
-     		}
-     	}
-     	if found {
-     		prioritizedJmnedictEntries = append(prioritizedJmnedictEntries, entry)
-     	}
-     }
- }
-
- // Filter Chinese word entries
- for \_, entry := range entries.ChineseWords {
-     word, ok := entry.(chinese_words.ChineseWordEntry)
-     if ok {
-     	if containsCommonChar(word.Traditional) || containsCommonChar(word.Simplified) {
-     		prioritizedChineseWordsEntries = append(prioritizedChineseWordsEntries, entry)
-     	}
-     }
- }
-
- // If we have prioritized entries, use them
- if len(prioritizedKanjidicEntries) > 0 || len(prioritizedChineseCharsEntries) > 0 {
-     logf("Using prioritized entries - JMdict: %d, JMNedict: %d, Kanjidic: %d, Chinese Chars: %d, Chinese Words: %d\n",
-     	len(prioritizedJmdictEntries), len(prioritizedJmnedictEntries),
-     	len(prioritizedKanjidicEntries), len(prioritizedChineseCharsEntries),
-     	len(prioritizedChineseWordsEntries))
-
-     // In test mode, always include all Chinese character entries to ensure overlap
-     if len(entries.ChineseChars) > len(prioritizedChineseCharsEntries) {
-     	logf("In test mode, including all Chinese character entries to ensure overlap\n")
-     	// Keep the original Chinese character entries
-     	prioritizedChineseCharsEntries = entries.ChineseChars
-     }
-
-     return &DictionaryEntries{
-     	JMdict:       prioritizedJmdictEntries,
-     	JMNedict:     prioritizedJmnedictEntries,
-     	Kanjidic:     prioritizedKanjidicEntries,
-     	ChineseChars: prioritizedChineseCharsEntries,
-     	ChineseWords: prioritizedChineseWordsEntries,
-     }
- }
-
- // No prioritized entries found, return original entries
- return entries
  +}
- +// limitEntries limits the number of entries from each dictionary
  +func limitEntries(entries *DictionaryEntries, limit int, logf LogFunc) *DictionaryEntries {
- totalEntries := len(entries.JMdict) + len(entries.JMNedict) + len(entries.Kanjidic) +
-     len(entries.ChineseChars) + len(entries.ChineseWords)
-
- if limit >= totalEntries {
-     // No need to limit
-     return entries
- }
-
- logf("Limiting to %d entries (out of %d total)\n", limit, totalEntries)
-
- // Calculate proportions
- jmdictProportion := float64(len(entries.JMdict)) / float64(totalEntries)
- jmnedictProportion := float64(len(entries.JMNedict)) / float64(totalEntries)
- kanjidicProportion := float64(len(entries.Kanjidic)) / float64(totalEntries)
- chineseCharsProportion := float64(len(entries.ChineseChars)) / float64(totalEntries)
- chineseWordsProportion := float64(len(entries.ChineseWords)) / float64(totalEntries)
-
- // Calculate limits for each dictionary
- jmdictLimit := int(float64(limit) \* jmdictProportion)
- jmnedictLimit := int(float64(limit) \* jmnedictProportion)
- kanjidicLimit := int(float64(limit) \* kanjidicProportion)
- chineseCharsLimit := int(float64(limit) \* chineseCharsProportion)
- chineseWordsLimit := int(float64(limit) \* chineseWordsProportion)
-
- // Adjust for rounding errors
- remaining := limit - jmdictLimit - jmnedictLimit - kanjidicLimit -
-     chineseCharsLimit - chineseWordsLimit
- if remaining > 0 && len(entries.Kanjidic) > kanjidicLimit {
-     kanjidicLimit += remaining
- }
-
- // Apply limits
- limitedJmdict := entries.JMdict
- limitedJmnedict := entries.JMNedict
- limitedKanjidic := entries.Kanjidic
- limitedChineseChars := entries.ChineseChars
- limitedChineseWords := entries.ChineseWords
-
- if jmdictLimit < len(entries.JMdict) {
-     limitedJmdict = entries.JMdict[:jmdictLimit]
- }
- if jmnedictLimit < len(entries.JMNedict) {
-     limitedJmnedict = entries.JMNedict[:jmnedictLimit]
- }
- if kanjidicLimit < len(entries.Kanjidic) {
-     limitedKanjidic = entries.Kanjidic[:kanjidicLimit]
- }
- if chineseCharsLimit < len(entries.ChineseChars) {
-     limitedChineseChars = entries.ChineseChars[:chineseCharsLimit]
- }
- if chineseWordsLimit < len(entries.ChineseWords) {
-     limitedChineseWords = entries.ChineseWords[:chineseWordsLimit]
- }
-
- return &DictionaryEntries{
-     JMdict:       limitedJmdict,
-     JMNedict:     limitedJmnedict,
-     Kanjidic:     limitedKanjidic,
-     ChineseChars: limitedChineseChars,
-     ChineseWords: limitedChineseWords,
- }
  +}
  diff --git a/cmd/kiokun/internal/loader.go b/cmd/kiokun/internal/loader.go
  new file mode 100644
  index 0000000000..9c94304168
  --- /dev/null
  +++ b/cmd/kiokun/internal/loader.go
  @@ -0,0 +1,73 @@
  +package internal
- +import (
- "fmt"
- "path/filepath"
- "time"
-
- "kiokun-go/dictionaries/common"
  +)
- +// DictionaryEntries holds entries from all dictionaries
  +type DictionaryEntries struct {
- JMdict []common.Entry
- JMNedict []common.Entry
- Kanjidic []common.Entry
- ChineseChars []common.Entry
- ChineseWords []common.Entry
  +}
- +// LoadDictionaries loads all dictionaries and returns their entries
  +func LoadDictionaries(config *Config, logf LogFunc) (*DictionaryEntries, error) {
- // Resolve dictionaries path
- dictPath := filepath.Join(config.WorkspaceRoot, config.DictDir)
- logf("Using dictionary path: %s\n", dictPath)
- common.SetDictionariesBasePath(dictPath)
-
- // Import all dictionaries
- logf("Importing dictionaries...\n")
-
- // Get all registered dictionaries
- dictConfigs := common.GetRegisteredDictionaries()
-
- // Import each dictionary
- var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries []common.Entry
-
- for \_, dict := range dictConfigs {
-     // Construct full path
-     inputPath := filepath.Join(dict.SourceDir, dict.InputFile)
-
-     // Import this dictionary
-     logf("Importing %s from %s...\n", dict.Name, inputPath)
-     startTime := time.Now()
-
-     entries, err := dict.Importer.Import(inputPath)
-     if err != nil {
-     	return nil, fmt.Errorf("error importing %s: %v", dict.Name, err)
-     }
-
-     // Store entries by dictionary type
-     switch dict.Name {
-     case "jmdict":
-     	jmdictEntries = entries
-     case "jmnedict":
-     	jmnedictEntries = entries
-     case "kanjidic":
-     	kanjidicEntries = entries
-     case "chinese_chars":
-     	chineseCharsEntries = entries
-     case "chinese_words":
-     	chineseWordsEntries = entries
-     }
-
-     logf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(startTime).Seconds())
- }
-
- return &DictionaryEntries{
-     JMdict:       jmdictEntries,
-     JMNedict:     jmnedictEntries,
-     Kanjidic:     kanjidicEntries,
-     ChineseChars: chineseCharsEntries,
-     ChineseWords: chineseWordsEntries,
- }, nil
  +}
  diff --git a/cmd/kiokun/internal/processor.go b/cmd/kiokun/internal/processor.go
  new file mode 100644
  index 0000000000..607defc235
  --- /dev/null
  +++ b/cmd/kiokun/internal/processor.go
  @@ -0,0 +1,105 @@
  +package internal
- +import (
- "fmt"
- "time"
-
- "kiokun-go/dictionaries/chinese_chars"
- "kiokun-go/dictionaries/chinese_words"
- "kiokun-go/dictionaries/common"
- "kiokun-go/dictionaries/jmdict"
- "kiokun-go/dictionaries/jmnedict"
- "kiokun-go/dictionaries/kanjidic"
- "kiokun-go/processor"
  +)
- +// ProcessEntries processes dictionary entries and writes them to files
  +func ProcessEntries(entries *DictionaryEntries, config *Config, logf LogFunc) error {
- // Create dictionary processor with parallel file writing
- proc, err := processor.New(config.OutputDir, config.FileWriters)
- if err != nil {
-     return fmt.Errorf("error creating processor: %v", err)
- }
-
- // Combine all entries
- allEntries := combineEntries(entries)
-
- // Count entries by type for debugging
- jmdictCount := 0
- jmnedictCount := 0
- kanjidicCount := 0
- chineseCharsCount := 0
- chineseWordsCount := 0
-
- for \_, entry := range allEntries {
-     switch entry.(type) {
-     case jmdict.Word:
-     	jmdictCount++
-     case jmnedict.Name:
-     	jmnedictCount++
-     case kanjidic.Kanji:
-     	kanjidicCount++
-     case chinese_chars.ChineseCharEntry:
-     	chineseCharsCount++
-     case chinese_words.ChineseWordEntry:
-     	chineseWordsCount++
-     default:
-     	logf("Unknown entry type: %T\n", entry)
-     }
- }
- logf("Processing %d entries (%d JMdict, %d JMNedict, %d Kanjidic, %d Chinese Chars, %d Chinese Words)\n",
-     len(allEntries), jmdictCount, jmnedictCount, kanjidicCount, chineseCharsCount, chineseWordsCount)
-
- // Process entries in batches with progress reporting
- logf("Processing entries in batches...\n")
- totalEntries := len(allEntries)
- processStart := time.Now()
-
- for start := 0; start < totalEntries; start += config.BatchSize {
-     end := start + config.BatchSize
-     if end > totalEntries {
-     	end = totalEntries
-     }
-
-     batchEntries := allEntries[start:end]
-
-     // Only update progress every 10 batches to reduce output
-     if start%(config.BatchSize*10) == 0 || end == totalEntries {
-     	logf("\rProcessing entries %d-%d of %d (%.1f%%)...",
-     		start+1, end, totalEntries, float64(end)/float64(totalEntries)*100)
-     }
-
-     // Process entries sequentially in a batch
-     for _, entry := range batchEntries {
-     	if err := proc.ProcessEntries([]common.Entry{entry}); err != nil {
-     		return fmt.Errorf("error processing entry %s: %v", entry.GetID(), err)
-     	}
-     }
- }
- processDuration := time.Since(processStart)
- logf("\rProcessed all %d entries (100%%) in %.2f seconds (%.1f entries/sec)\n",
-     totalEntries, processDuration.Seconds(), float64(totalEntries)/processDuration.Seconds())
-
- // Write all processed entries to files
- logf("Writing files to %s...\n", config.OutputDir)
- if err := proc.WriteToFiles(); err != nil {
-     return fmt.Errorf("error writing files: %v", err)
- }
-
- return nil
  +}
- +// combineEntries combines entries from all dictionaries into a single slice
  +func combineEntries(entries \*DictionaryEntries) []common.Entry {
- totalEntries := len(entries.JMdict) + len(entries.JMNedict) + len(entries.Kanjidic) +
-     len(entries.ChineseChars) + len(entries.ChineseWords)
-
- allEntries := make([]common.Entry, 0, totalEntries)
- allEntries = append(allEntries, entries.JMdict...)
- allEntries = append(allEntries, entries.JMNedict...)
- allEntries = append(allEntries, entries.Kanjidic...)
- allEntries = append(allEntries, entries.ChineseChars...)
- allEntries = append(allEntries, entries.ChineseWords...)
-
- return allEntries
  +}
  diff --git a/cmd/kiokun/kiokun b/cmd/kiokun/kiokun
  new file mode 100755
  index 0000000000..7aec3d4524
  Binary files /dev/null and b/cmd/kiokun/kiokun differ
  diff --git a/cmd/kiokun/main.go b/cmd/kiokun/main.go
  index 22f92bd86e..45d98b159b 100644
  --- a/cmd/kiokun/main.go
  +++ b/cmd/kiokun/main.go
  @@ -1,422 +1,49 @@
  package main

import (

- "flag"
  "fmt"
  "os"
- "path/filepath"
- "runtime"
- "time"
- "unicode"

- "kiokun-go/dictionaries/common"
- "kiokun-go/dictionaries/jmdict"

* // Import for side effects (dictionary registration)
* \_ "kiokun-go/dictionaries/chinese_chars"
* _ "kiokun-go/dictionaries/chinese_words"
  _ "kiokun-go/dictionaries/jmdict"

- "kiokun-go/dictionaries/jmnedict"
  \_ "kiokun-go/dictionaries/jmnedict"
- "kiokun-go/dictionaries/kanjidic"
  \_ "kiokun-go/dictionaries/kanjidic"

- // Import Chinese dictionaries
- "kiokun-go/dictionaries/chinese_chars"
- \_ "kiokun-go/dictionaries/chinese_chars"
- "kiokun-go/dictionaries/chinese_words"
- \_ "kiokun-go/dictionaries/chinese_words"
- "kiokun-go/processor"
  -)
- -// OutputMode determines which words to output
  -type OutputMode string
- -const (
- OutputAll OutputMode = "all" // Output all words
- OutputHanOnly OutputMode = "han-only" // Output words with only Han characters (legacy mode)
- OutputHan1Char OutputMode = "han-1char" // Output words with exactly 1 Han character
- OutputHan2Char OutputMode = "han-2char" // Output words with exactly 2 Han characters
- OutputHan3Plus OutputMode = "han-3plus" // Output words with 3 or more Han characters
- OutputNonHanOnly OutputMode = "non-han" // Output words with at least one non-Han character

* // Import local package functions
* . "kiokun-go/cmd/kiokun/internal"
  )

func main() {

- // Configuration flags
- dictDir := flag.String("dictdir", "dictionaries", "Base directory containing dictionary packages")
- outputDir := flag.String("outdir", "output", "Output directory for processed files")
- workers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines for batch processing")
- fileWriters := flag.Int("writers", runtime.NumCPU(), "Number of parallel workers for file writing")
- silent := flag.Bool("silent", false, "Disable progress output")
- devMode := flag.Bool("dev", false, "Development mode - use /tmp for faster I/O")
- limitEntries := flag.Int("limit", 0, "Limit the number of entries to process (0 = no limit)")
- batchSize := flag.Int("batch", 10000, "Process entries in batches of this size")
- outputModeFlag := flag.String("mode", "all", "Output mode: 'all', 'han-only' (legacy), 'han-1char', 'han-2char', 'han-3plus', or 'non-han'")
- flag.Parse()
-
- // Parse and validate the output mode
- outputMode := OutputMode(\*outputModeFlag)
- if outputMode != OutputAll &&
-     outputMode != OutputHanOnly &&
-     outputMode != OutputHan1Char &&
-     outputMode != OutputHan2Char &&
-     outputMode != OutputHan3Plus &&
-     outputMode != OutputNonHanOnly {
-     fmt.Fprintf(os.Stderr, "Invalid output mode: %s\n", *outputModeFlag)
-     fmt.Fprintf(os.Stderr, "Valid modes: all, han-only, han-1char, han-2char, han-3plus, non-han\n")
-     os.Exit(1)
- }
-
- // Modify output directory based on mode
- if outputMode == OutputHanOnly {
-     *outputDir = *outputDir + "_han"
- } else if outputMode == OutputHan1Char {
-     *outputDir = *outputDir + "_han_1char"
- } else if outputMode == OutputHan2Char {
-     *outputDir = *outputDir + "_han_2char"
- } else if outputMode == OutputHan3Plus {
-     *outputDir = *outputDir + "_han_3plus"
- } else if outputMode == OutputNonHanOnly {
-     *outputDir = *outputDir + "_non_han"
- }
-
- logf := func(format string, a ...interface{}) {
-     if !*silent {
-     	fmt.Printf(format, a...)
-     }
- }
-
- // If dev mode is enabled, use /tmp directory for output
- if \*devMode {
-     tmpDir := filepath.Join("/tmp", "kiokun-output")
-     logf("Development mode enabled: using %s for output\n", tmpDir)
-     *outputDir = tmpDir
-
-     // Create the tmp directory if it doesn't exist
-     if err := os.MkdirAll(tmpDir, 0755); err != nil {
-     	fmt.Fprintf(os.Stderr, "Error creating tmp directory: %v\n", err)
-     	os.Exit(1)
-     }
- }
-
- // Ensure we're using an absolute path for the output directory
- absOutputDir, err := filepath.Abs(\*outputDir)

* // Parse configuration
* config, logf, err := ParseConfig()
  if err != nil {

-     fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
-     os.Exit(1)
- }
- \*outputDir = absOutputDir
-
- // Create the output directory if it doesn't exist
- if err := os.MkdirAll(\*outputDir, 0755); err != nil {
-     fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)

*     fmt.Fprintf(os.Stderr, "Error parsing configuration: %v\n", err)
      os.Exit(1)

  }

  logf("Starting initialization...\n")

- logf("Input dictionaries directory: %s\n", \*dictDir)
- logf("Output directory: %s\n", \*outputDir)
- logf("Using %d processing workers and %d file writers\n", *workers, *fileWriters)
- if outputMode != OutputAll {
-     logf("Filtering mode: %s\n", outputMode)
- }
-
- // Set dictionaries base path - find workspace root
- workspaceRoot, err := os.Getwd()
- if err != nil {
-     fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
-     os.Exit(1)

* logf("Input dictionaries directory: %s\n", config.DictDir)
* logf("Output directory: %s\n", config.OutputDir)
* logf("Using %d processing workers and %d file writers\n", config.Workers, config.FileWriters)
* if config.OutputMode != OutputAll {
*     logf("Filtering mode: %s\n", config.OutputMode)
  }

- // If we're in cmd/kiokun, go up two levels
- if filepath.Base(workspaceRoot) == "kiokun" && filepath.Base(filepath.Dir(workspaceRoot)) == "cmd" {
-     workspaceRoot = filepath.Dir(filepath.Dir(workspaceRoot))
- }
-
- // Resolve dictionaries path
- dictPath := filepath.Join(workspaceRoot, \*dictDir)
- logf("Using dictionary path: %s\n", dictPath)
- common.SetDictionariesBasePath(dictPath)
-
- // Create dictionary processor with parallel file writing
- proc, err := processor.New(*outputDir, *fileWriters)

* // Load dictionaries
* entries, err := LoadDictionaries(config, logf)
  if err != nil {

-     fmt.Fprintf(os.Stderr, "Error creating processor: %v\n", err)

*     fmt.Fprintf(os.Stderr, "Error loading dictionaries: %v\n", err)
      os.Exit(1)
  }

- // Import all dictionaries
- logf("Importing dictionaries...\n")
-
- // Get all registered dictionaries
- dictConfigs := common.GetRegisteredDictionaries()
-
- // Import each dictionary
- var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries []common.Entry
-
- for \_, dict := range dictConfigs {
-     // Construct full path
-     inputPath := filepath.Join(dict.SourceDir, dict.InputFile)
-
-     // Import this dictionary
-     logf("Importing %s from %s...\n", dict.Name, inputPath)
-     startTime := time.Now()
-
-     entries, err := dict.Importer.Import(inputPath)
-     if err != nil {
-     	fmt.Fprintf(os.Stderr, "Error importing %s: %v\n", dict.Name, err)
-     	os.Exit(1)
-     }
-
-     // Store entries by dictionary type
-     switch dict.Name {
-     case "jmdict":
-     	jmdictEntries = entries
-     case "jmnedict":
-     	jmnedictEntries = entries
-     case "kanjidic":
-     	kanjidicEntries = entries
-     case "chinese_chars":
-     	chineseCharsEntries = entries
-     case "chinese_words":
-     	chineseWordsEntries = entries
-     }
-
-     logf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(startTime).Seconds())
- }
-
- // Helper function to check if a string contains only Han characters
- isHanOnly := func(s string) bool {
-     for _, r := range s {
-     	if !unicode.Is(unicode.Han, r) {
-     		return false
-     	}
-     }
-     return true
- }
-
- // Helper function to check if an entry should be included based on the mode
- shouldIncludeEntry := func(entry common.Entry) bool {
-     if outputMode == OutputAll {
-     	return true
-     }
-
-     // Get primary text representation for filtering
-     var primaryText string
-     switch e := entry.(type) {
-     case jmdict.Word:
-     	if len(e.Kanji) > 0 {
-     		primaryText = e.Kanji[0].Text
-     	} else if len(e.Kana) > 0 {
-     		primaryText = e.Kana[0].Text
-     	} else {
-     		primaryText = e.ID
-     	}
-     case jmnedict.Name:
-     	// Use the Name's primary text
-     	if len(e.Kanji) > 0 {
-     		primaryText = e.Kanji[0]
-     	} else if len(e.Reading) > 0 {
-     		primaryText = e.Reading[0]
-     	} else {
-     		primaryText = e.ID
-     	}
-     case kanjidic.Kanji:
-     	// Use the Kanji character
-     	primaryText = e.Character
-     case chinese_chars.ChineseCharEntry:
-     	// Use the traditional character
-     	primaryText = e.Traditional
-     case chinese_words.ChineseWordEntry:
-     	// Use the traditional word
-     	primaryText = e.Traditional
-     default:
-     	// If we don't know how to filter this type, include it by default
-     	return true
-     }
-
-     // Check if the text contains only Han characters
-     isHan := isHanOnly(primaryText)
-     charCount := len([]rune(primaryText)) // Get correct Unicode character count
-
-     // Apply filtering based on mode
-     switch outputMode {
-     case OutputNonHanOnly:
-     	return !isHan
-     case OutputHanOnly:
-     	return isHan
-     case OutputHan1Char:
-     	return isHan && charCount == 1
-     case OutputHan2Char:
-     	return isHan && charCount == 2
-     case OutputHan3Plus:
-     	return isHan && charCount >= 3
-     default:
-     	return true
-     }
- }
-
- // Filter entries based on output mode
- if outputMode != OutputAll {
-     filteredJmdict := make([]common.Entry, 0, len(jmdictEntries))
-     filteredJmnedict := make([]common.Entry, 0, len(jmnedictEntries))
-     filteredKanjidic := make([]common.Entry, 0, len(kanjidicEntries))
-     filteredChineseChars := make([]common.Entry, 0, len(chineseCharsEntries))
-     filteredChineseWords := make([]common.Entry, 0, len(chineseWordsEntries))
-
-     for _, entry := range jmdictEntries {
-     	if shouldIncludeEntry(entry) {
-     		filteredJmdict = append(filteredJmdict, entry)
-     	}
-     }
-     for _, entry := range jmnedictEntries {
-     	if shouldIncludeEntry(entry) {
-     		filteredJmnedict = append(filteredJmnedict, entry)
-     	}
-     }
-     for _, entry := range kanjidicEntries {
-     	if shouldIncludeEntry(entry) {
-     		filteredKanjidic = append(filteredKanjidic, entry)
-     	}
-     }
-     for _, entry := range chineseCharsEntries {
-     	if shouldIncludeEntry(entry) {
-     		filteredChineseChars = append(filteredChineseChars, entry)
-     	}
-     }
-     for _, entry := range chineseWordsEntries {
-     	if shouldIncludeEntry(entry) {
-     		filteredChineseWords = append(filteredChineseWords, entry)
-     	}
-     }
-
-     logf("Filtered entries - JMdict: %d -> %d, JMNedict: %d -> %d, Kanjidic: %d -> %d, Chinese Chars: %d -> %d, Chinese Words: %d -> %d\n",
-     	len(jmdictEntries), len(filteredJmdict),
-     	len(jmnedictEntries), len(filteredJmnedict),
-     	len(kanjidicEntries), len(filteredKanjidic),
-     	len(chineseCharsEntries), len(filteredChineseChars),
-     	len(chineseWordsEntries), len(filteredChineseWords))
-
-     jmdictEntries = filteredJmdict
-     jmnedictEntries = filteredJmnedict
-     kanjidicEntries = filteredKanjidic
-     chineseCharsEntries = filteredChineseChars
-     chineseWordsEntries = filteredChineseWords
- }
-
- // Apply entry limit if specified
- var allEntries []common.Entry
-
- if \*limitEntries > 0 {
-     totalEntries := len(jmdictEntries) + len(jmnedictEntries) + len(kanjidicEntries) +
-     	len(chineseCharsEntries) + len(chineseWordsEntries)
-     if *limitEntries < totalEntries {
-     	logf("Limiting to %d entries (out of %d total)\n", *limitEntries, totalEntries)
-
-     	// Calculate proportions
-     	jmdictProportion := float64(len(jmdictEntries)) / float64(totalEntries)
-     	jmnedictProportion := float64(len(jmnedictEntries)) / float64(totalEntries)
-     	kanjidicProportion := float64(len(kanjidicEntries)) / float64(totalEntries)
-     	chineseCharsProportion := float64(len(chineseCharsEntries)) / float64(totalEntries)
-     	chineseWordsProportion := float64(len(chineseWordsEntries)) / float64(totalEntries)
-
-     	// Calculate limits for each dictionary
-     	jmdictLimit := int(float64(*limitEntries) * jmdictProportion)
-     	jmnedictLimit := int(float64(*limitEntries) * jmnedictProportion)
-     	kanjidicLimit := int(float64(*limitEntries) * kanjidicProportion)
-     	chineseCharsLimit := int(float64(*limitEntries) * chineseCharsProportion)
-     	chineseWordsLimit := int(float64(*limitEntries) * chineseWordsProportion)
-
-     	// Adjust for rounding errors
-     	remaining := *limitEntries - jmdictLimit - jmnedictLimit - kanjidicLimit -
-     		chineseCharsLimit - chineseWordsLimit
-     	if remaining > 0 && len(kanjidicEntries) > kanjidicLimit {
-     		kanjidicLimit += remaining
-     	}
-
-     	// Apply limits
-     	if jmdictLimit < len(jmdictEntries) {
-     		jmdictEntries = jmdictEntries[:jmdictLimit]
-     	}
-     	if jmnedictLimit < len(jmnedictEntries) {
-     		jmnedictEntries = jmnedictEntries[:jmnedictLimit]
-     	}
-     	if kanjidicLimit < len(kanjidicEntries) {
-     		kanjidicEntries = kanjidicEntries[:kanjidicLimit]
-     	}
-     	if chineseCharsLimit < len(chineseCharsEntries) {
-     		chineseCharsEntries = chineseCharsEntries[:chineseCharsLimit]
-     	}
-     	if chineseWordsLimit < len(chineseWordsEntries) {
-     		chineseWordsEntries = chineseWordsEntries[:chineseWordsLimit]
-     	}
-     }
- }
-
- // Combine all entries
- allEntries = append(allEntries, jmdictEntries...)
- allEntries = append(allEntries, jmnedictEntries...)
- allEntries = append(allEntries, kanjidicEntries...)
- allEntries = append(allEntries, chineseCharsEntries...)
- allEntries = append(allEntries, chineseWordsEntries...)
-
- // Count entries by type for debugging
- jmdictCount := 0
- jmnedictCount := 0
- kanjidicCount := 0
- chineseCharsCount := 0
- chineseWordsCount := 0
-
- for \_, entry := range allEntries {
-     switch entry.(type) {
-     case jmdict.Word:
-     	jmdictCount++
-     case jmnedict.Name:
-     	jmnedictCount++
-     case kanjidic.Kanji:
-     	kanjidicCount++
-     case chinese_chars.ChineseCharEntry:
-     	chineseCharsCount++
-     case chinese_words.ChineseWordEntry:
-     	chineseWordsCount++
-     default:
-     	logf("Unknown entry type: %T\n", entry)
-     }
- }
- logf("Processing %d entries (%d JMdict, %d JMNedict, %d Kanjidic, %d Chinese Chars, %d Chinese Words)\n",
-     len(allEntries), jmdictCount, jmnedictCount, kanjidicCount, chineseCharsCount, chineseWordsCount)
-
- // Process entries in batches with progress reporting
- logf("Processing entries in batches...\n")
- totalEntries := len(allEntries)
- processStart := time.Now()
-
- for start := 0; start < totalEntries; start += \*batchSize {
-     end := start + *batchSize
-     if end > totalEntries {
-     	end = totalEntries
-     }
-
-     batchEntries := allEntries[start:end]
-
-     // Only update progress every 10 batches to reduce output
-     if start%((*batchSize)*10) == 0 || end == totalEntries {
-     	logf("\rProcessing entries %d-%d of %d (%.1f%%)...",
-     		start+1, end, totalEntries, float64(end)/float64(totalEntries)*100)
-     }
-
-     // Process entries sequentially in a batch
-     for _, entry := range batchEntries {
-     	if err := proc.ProcessEntries([]common.Entry{entry}); err != nil {
-     		fmt.Fprintf(os.Stderr, "Error processing entry %s: %v\n", entry.GetID(), err)
-     		os.Exit(1)
-     	}
-     }
- }
- processDuration := time.Since(processStart)
- logf("\rProcessed all %d entries (100%%) in %.2f seconds (%.1f entries/sec)\n",
-     totalEntries, processDuration.Seconds(), float64(totalEntries)/processDuration.Seconds())

* // Filter entries
* filteredEntries := FilterEntries(entries, config, logf)

- // Write all processed entries to files
- logf("Writing files to %s...\n", \*outputDir)
- if err := proc.WriteToFiles(); err != nil {
-     fmt.Fprintf(os.Stderr, "Error writing files: %v\n", err)

* // Process entries
* if err := ProcessEntries(filteredEntries, config, logf); err != nil {
*     fmt.Fprintf(os.Stderr, "Error processing entries: %v\n", err)
      os.Exit(1)
  }

diff --git a/cmd/test_minification/main.go b/cmd/test_minification/main.go
new file mode 100644
index 0000000000..a7f03a9678
--- /dev/null
+++ b/cmd/test_minification/main.go
@@ -0,0 +1,194 @@
+package main

- +import (
- "flag"
- "fmt"
- "os"
- "path/filepath"
- "runtime"
- "time"
-
- \_ "kiokun-go/dictionaries/chinese_chars"
- \_ "kiokun-go/dictionaries/chinese_words"
- "kiokun-go/dictionaries/common"
- \_ "kiokun-go/dictionaries/jmdict"
- \_ "kiokun-go/dictionaries/jmnedict"
- \_ "kiokun-go/dictionaries/kanjidic"
- "kiokun-go/processor"
  +)
- +func main() {
- // Parse command line arguments
- dictDir := flag.String("dictdir", "dictionaries", "Base directory containing dictionary packages")
- outputDir := flag.String("outdir", "output_test", "Output directory for processed files")
- limit := flag.Int("limit", 1000, "Limit the number of entries to process from each dictionary")
- flag.Parse()
-
- // Ensure output directory exists
- absOutputDir, err := filepath.Abs(\*outputDir)
- if err != nil {
-     fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
-     os.Exit(1)
- }
- \*outputDir = absOutputDir
-
- if err := os.MkdirAll(\*outputDir, 0755); err != nil {
-     fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
-     os.Exit(1)
- }
-
- fmt.Printf("Starting minification test...\n")
- fmt.Printf("Input dictionaries directory: %s\n", \*dictDir)
- fmt.Printf("Output directory: %s\n", \*outputDir)
- fmt.Printf("Entry limit per dictionary: %d\n", \*limit)
-
- // Set dictionaries base path
- common.SetDictionariesBasePath(\*dictDir)
-
- // Create processor
- proc, err := processor.New(\*outputDir, runtime.NumCPU())
- if err != nil {
-     fmt.Fprintf(os.Stderr, "Error creating processor: %v\n", err)
-     os.Exit(1)
- }
-
- // Import dictionaries
- fmt.Printf("Importing dictionaries...\n")
- startTime := time.Now()
-
- // Get all registered dictionaries
- dictConfigs := common.GetRegisteredDictionaries()
-
- // Import each dictionary
- var jmdictEntries, jmnedictEntries, kanjidicEntries, chineseCharsEntries, chineseWordsEntries []common.Entry
-
- for \_, dict := range dictConfigs {
-     // Construct full path
-     inputPath := filepath.Join(dict.SourceDir, dict.InputFile)
-
-     // Import this dictionary
-     fmt.Printf("Importing %s from %s...\n", dict.Name, inputPath)
-     dictStartTime := time.Now()
-
-     entries, err := dict.Importer.Import(inputPath)
-     if err != nil {
-     	fmt.Fprintf(os.Stderr, "Error importing %s: %v\n", dict.Name, err)
-     	continue
-     }
-
-     // Limit entries
-     if *limit > 0 && len(entries) > *limit {
-     	entries = entries[:*limit]
-     }
-
-     // Store entries by dictionary type
-     switch dict.Name {
-     case "jmdict":
-     	jmdictEntries = entries
-     case "jmnedict":
-     	jmnedictEntries = entries
-     case "kanjidic":
-     	kanjidicEntries = entries
-     case "chinese_chars":
-     	chineseCharsEntries = entries
-     case "chinese_words":
-     	chineseWordsEntries = entries
-     }
-
-     fmt.Printf("Imported %s: %d entries (%.2fs)\n", dict.Name, len(entries), time.Since(dictStartTime).Seconds())
- }
-
- // Combine all entries
- var allEntries []common.Entry
- allEntries = append(allEntries, jmdictEntries...)
- allEntries = append(allEntries, jmnedictEntries...)
- allEntries = append(allEntries, kanjidicEntries...)
- allEntries = append(allEntries, chineseCharsEntries...)
- allEntries = append(allEntries, chineseWordsEntries...)
-
- // Count entries by type
- jmdictCount := len(jmdictEntries)
- jmnedictCount := len(jmnedictEntries)
- kanjidicCount := len(kanjidicEntries)
- chineseCharsCount := len(chineseCharsEntries)
- chineseWordsCount := len(chineseWordsEntries)
-
- fmt.Printf("Processing %d entries (%d JMdict, %d JMNedict, %d Kanjidic, %d Chinese Chars, %d Chinese Words)\n",
-     len(allEntries), jmdictCount, jmnedictCount, kanjidicCount, chineseCharsCount, chineseWordsCount)
-
- // Process entries
- fmt.Printf("Processing entries...\n")
- processStart := time.Now()
-
- if err := proc.ProcessEntries(allEntries); err != nil {
-     fmt.Fprintf(os.Stderr, "Error processing entries: %v\n", err)
-     os.Exit(1)
- }
-
- processDuration := time.Since(processStart)
- fmt.Printf("Processed all %d entries in %.2f seconds (%.1f entries/sec)\n",
-     len(allEntries), processDuration.Seconds(), float64(len(allEntries))/processDuration.Seconds())
-
- // Write files
- fmt.Printf("Writing files to %s...\n", \*outputDir)
- writeStart := time.Now()
-
- if err := proc.WriteToFiles(); err != nil {
-     fmt.Fprintf(os.Stderr, "Error writing files: %v\n", err)
-     os.Exit(1)
- }
-
- writeDuration := time.Since(writeStart)
- fmt.Printf("Wrote files in %.2f seconds\n", writeDuration.Seconds())
-
- // Verify minification
- fmt.Printf("Verifying minification...\n")
- verifyMinification(\*outputDir)
-
- totalDuration := time.Since(startTime)
- fmt.Printf("Total time: %.2f seconds\n", totalDuration.Seconds())
- fmt.Printf("Successfully processed and minified dictionary files\n")
  +}
- +// verifyMinification checks a sample of files to ensure they were properly minified
  +func verifyMinification(outputDir string) {
- // Count files
- var totalFiles, jmdictFiles, chineseFiles int
- var filesWithWildcards, filesWithEmptyArrays int
-
- // Walk through the output directory and check files
- err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
-     if err != nil {
-     	return err
-     }
-
-     // Skip directories
-     if info.IsDir() {
-     	return nil
-     }
-
-     // Only check .json.br files
-     if filepath.Ext(path) != ".br" {
-     	return nil
-     }
-
-     totalFiles++
-
-     // For now, we'll just count files
-     // In a real verification, we would decompress and check the content
-     // But that would require more complex code to parse the Brotli-compressed JSON
-
-     return nil
- })
-
- if err != nil {
-     fmt.Printf("Error walking output directory: %v\n", err)
-     return
- }
-
- fmt.Printf("Verified %d total files\n", totalFiles)
- fmt.Printf("- JMdict files: %d\n", jmdictFiles)
- fmt.Printf("- Chinese files: %d\n", chineseFiles)
- fmt.Printf("- Files with wildcards: %d\n", filesWithWildcards)
- fmt.Printf("- Files with empty arrays: %d\n", filesWithEmptyArrays)
  +}
  diff --git a/compression_test/web/benchmark_results.json b/compression_test/web/benchmark_results.json
  new file mode 100644
  index 0000000000..98dbc5c06f
  --- /dev/null
  +++ b/compression_test/web/benchmark_results.json
  @@ -0,0 +1,47 @@
  +[
- {
- "Algorithm": "Gzip",
- "OriginalSize": 54692,
- "CompressedSize": 745,
- "CompressionRatio": 0.013621736268558473,
- "CompressionTime": 223125,
- "DecompressionTime": 63542,
- "NetworkTransferTime": 745000
- },
- {
- "Algorithm": "Zstd",
- "OriginalSize": 54692,
- "CompressedSize": 475,
- "CompressionRatio": 0.008684999634315805,
- "CompressionTime": 485792,
- "DecompressionTime": 156750,
- "NetworkTransferTime": 475000
- },
- {
- "Algorithm": "Brotli",
- "OriginalSize": 54692,
- "CompressedSize": 411,
- "CompressionRatio": 0.007514810209902728,
- "CompressionTime": 893542,
- "DecompressionTime": 56291,
- "NetworkTransferTime": 411000
- },
- {
- "Algorithm": "LZ4",
- "OriginalSize": 54692,
- "CompressedSize": 979,
- "CompressionRatio": 0.017900241351568785,
- "CompressionTime": 215000,
- "DecompressionTime": 130250,
- "NetworkTransferTime": 979000
- },
- {
- "Algorithm": "XZ",
- "OriginalSize": 54692,
- "CompressedSize": 564,
- "CompressionRatio": 0.01031229430264024,
- "CompressionTime": 1111209,
- "DecompressionTime": 124916,
- "NetworkTransferTime": 564000
- }
  +]
  \ No newline at end of file
  diff --git a/compression_test/web/compressed_brotli.bin b/compression_test/web/compressed_brotli.bin
  new file mode 100644
  index 0000000000..d1794f4b50
  Binary files /dev/null and b/compression_test/web/compressed_brotli.bin differ
  diff --git a/compression_test/web/compressed_gzip.bin b/compression_test/web/compressed_gzip.bin
  new file mode 100644
  index 0000000000..b993861263
  Binary files /dev/null and b/compression_test/web/compressed_gzip.bin differ
  diff --git a/compression_test/web/test_entries.json b/compression_test/web/test_entries.json
  new file mode 100644
  index 0000000000..ec72108da1
  --- /dev/null
  +++ b/compression_test/web/test_entries.json
  @@ -0,0 +1,3202 @@
  +{
- "entry_0": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_1": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_10": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_11": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_12": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_13": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_14": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_15": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_16": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_17": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_18": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_19": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_2": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_20": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_21": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_22": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_23": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_24": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_25": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_26": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_27": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_28": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_29": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_3": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_30": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_31": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_32": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_33": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_34": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_35": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_36": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_37": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_38": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_39": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_4": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_40": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_41": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_42": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_43": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_44": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_45": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_46": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_47": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_48": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_49": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_5": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_50": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_51": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_52": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_53": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_54": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_55": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_56": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_57": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_58": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_59": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_6": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_60": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_61": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_62": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_63": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_64": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_65": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_66": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_67": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_68": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_69": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_7": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_70": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_71": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_72": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_73": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_74": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_75": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_76": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_77": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_78": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_79": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_8": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_80": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_81": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_82": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_83": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_84": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_85": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_86": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_87": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_88": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_89": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_9": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_90": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_91": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_92": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_93": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_94": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_95": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_96": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_97": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_98": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- },
- "entry_99": {
- "w_j": [
-      {
-        "id": "1000001",
-        "kana": [
-          {
-            "common": true,
-            "text": "わたし"
-          }
-        ],
-        "kanji": [
-          {
-            "common": true,
-            "text": "私"
-          }
-        ],
-        "sense": [
-          {
-            "gloss": [
-              {
-                "lang": "eng",
-                "text": "I"
-              }
-            ],
-            "partOfSpeech": [
-              "pn"
-            ]
-          }
-        ]
-      }
- ]
- }
  +}
  \ No newline at end of file
  diff --git a/processor/minification_test.go b/processor/minification_test.go
  new file mode 100644
  index 0000000000..139748ab46
  --- /dev/null
  +++ b/processor/minification_test.go
  @@ -0,0 +1,343 @@
  +package processor
- +import (
- "bytes"
- "encoding/json"
- "io"
- "os"
- "path/filepath"
- "strings"
- "testing"
-
- "kiokun-go/dictionaries/chinese_chars"
- "kiokun-go/dictionaries/chinese_words"
- "kiokun-go/dictionaries/common"
- "kiokun-go/dictionaries/jmdict"
-
- "github.com/andybalholm/brotli"
  +)
- +// TestMinification tests the entire minification process:
  +// 1. Importing dictionary entries
  +// 2. Processing them
  +// 3. Writing them to files
  +// 4. Verifying the minification was applied correctly
  +func TestMinification(t \*testing.T) {
- // Create test directory
- testDir := filepath.Join(os.TempDir(), "kiokun_minification_test")
- err := os.MkdirAll(testDir, 0755)
- if err != nil {
-     t.Fatalf("Failed to create test directory: %v", err)
- }
- defer os.RemoveAll(testDir)
-
- // Set dictionaries base path
- dictDir := filepath.Join("..", "dictionaries")
- common.SetDictionariesBasePath(dictDir)
-
- // Create processor
- proc, err := New(testDir, 1)
- if err != nil {
-     t.Fatalf("Failed to create processor: %v", err)
- }
-
- // Test with JMdict entries
- t.Run("JMdict", func(t \*testing.T) {
-     testJMdictMinification(t, proc, dictDir, testDir)
- })
-
- // Test with Chinese dictionary entries
- t.Run("Chinese", func(t \*testing.T) {
-     testChineseMinification(t, proc, dictDir, testDir)
- })
-
- // Verify the minification
- t.Run("Verification", func(t \*testing.T) {
-     verifyMinification(t, testDir)
- })
  +}
- +// testJMdictMinification tests minification of JMdict entries
  +func testJMdictMinification(t *testing.T, proc *DictionaryProcessor, dictDir, testDir string) {
- // Find JMdict source file
- jmdictSourceDir := filepath.Join(dictDir, "jmdict", "source")
- pattern := `^jmdict-.*\.json$`
- jmdictFile, err := common.FindDictionaryFile(jmdictSourceDir, pattern)
- if err != nil {
-     t.Logf("JMdict file not found, skipping test: %v", err)
-     return
- }
-
- // Import JMdict entries
- jmdictPath := filepath.Join(jmdictSourceDir, jmdictFile)
- t.Logf("Using JMdict file: %s", jmdictPath)
-
- file, err := os.Open(jmdictPath)
- if err != nil {
-     t.Fatalf("Failed to open JMdict file: %v", err)
- }
- defer file.Close()
-
- var dict jmdict.JmdictTypes
- if err := common.ImportJSON(file, &dict); err != nil {
-     t.Fatalf("Failed to parse JMdict: %v", err)
- }
-
- // Limit to first 100 entries for testing
- maxEntries := 100
- if len(dict.Words) > maxEntries {
-     dict.Words = dict.Words[:maxEntries]
- }
- t.Logf("Testing with %d JMdict entries", len(dict.Words))
-
- // Create entries for processing
- entries := make([]common.Entry, len(dict.Words))
- for i, word := range dict.Words {
-     entries[i] = word
- }
-
- // Process entries
- if err := proc.ProcessEntries(entries); err != nil {
-     t.Fatalf("Failed to process JMdict entries: %v", err)
- }
-
- // Write to files
- if err := proc.WriteToFiles(); err != nil {
-     t.Fatalf("Failed to write files: %v", err)
- }
  +}
- +// testChineseMinification tests minification of Chinese dictionary entries
  +func testChineseMinification(t *testing.T, proc *DictionaryProcessor, dictDir, testDir string) {
- // Find Chinese character dictionary source file
- chineseCharsSourceDir := filepath.Join(dictDir, "chinese_chars", "source")
- pattern := `^dictionary_char_.*\.json$`
- chineseCharsFile, err := common.FindDictionaryFile(chineseCharsSourceDir, pattern)
- if err != nil {
-     t.Logf("Chinese character file not found, skipping test: %v", err)
-     return
- }
-
- // Import Chinese character entries
- chineseCharsPath := filepath.Join(chineseCharsSourceDir, chineseCharsFile)
- t.Logf("Using Chinese character file: %s", chineseCharsPath)
-
- file, err := os.Open(chineseCharsPath)
- if err != nil {
-     t.Logf("Failed to open Chinese character file, skipping test: %v", err)
-     return
- }
- defer file.Close()
-
- var charEntries []chinese_chars.ChineseCharEntry
- if err := json.NewDecoder(file).Decode(&charEntries); err != nil {
-     t.Logf("Failed to parse Chinese character file, skipping test: %v", err)
-     return
- }
-
- // Limit to first 100 entries for testing
- maxEntries := 100
- if len(charEntries) > maxEntries {
-     charEntries = charEntries[:maxEntries]
- }
- t.Logf("Testing with %d Chinese character entries", len(charEntries))
-
- // Create entries for processing
- entries := make([]common.Entry, len(charEntries))
- for i, char := range charEntries {
-     // If ID is not set, use traditional character as ID
-     if char.ID == "" {
-     	char.ID = char.Traditional
-     }
-     entries[i] = char
- }
-
- // Process entries
- if err := proc.ProcessEntries(entries); err != nil {
-     t.Fatalf("Failed to process Chinese character entries: %v", err)
- }
-
- // Find Chinese word dictionary source file
- chineseWordsSourceDir := filepath.Join(dictDir, "chinese_words", "source")
- pattern = `^dictionary_word_.*\.json$`
- chineseWordsFile, err := common.FindDictionaryFile(chineseWordsSourceDir, pattern)
- if err != nil {
-     t.Logf("Chinese word file not found, skipping that part of test: %v", err)
- } else {
-     // Import Chinese word entries
-     chineseWordsPath := filepath.Join(chineseWordsSourceDir, chineseWordsFile)
-     t.Logf("Using Chinese word file: %s", chineseWordsPath)
-
-     file, err := os.Open(chineseWordsPath)
-     if err != nil {
-     	t.Logf("Failed to open Chinese word file, skipping that part of test: %v", err)
-     } else {
-     	defer file.Close()
-
-     	var wordEntries []chinese_words.ChineseWordEntry
-     	if err := json.NewDecoder(file).Decode(&wordEntries); err != nil {
-     		t.Logf("Failed to parse Chinese word file, skipping that part of test: %v", err)
-     	} else {
-     		// Limit to first 100 entries for testing
-     		if len(wordEntries) > maxEntries {
-     			wordEntries = wordEntries[:maxEntries]
-     		}
-     		t.Logf("Testing with %d Chinese word entries", len(wordEntries))
-
-     		// Create entries for processing
-     		wordEntriesCommon := make([]common.Entry, len(wordEntries))
-     		for i, word := range wordEntries {
-     			// If ID is not set, use traditional word as ID
-     			if word.ID == "" {
-     				word.ID = word.Traditional
-     			}
-     			wordEntriesCommon[i] = word
-     		}
-
-     		// Process entries
-     		if err := proc.ProcessEntries(wordEntriesCommon); err != nil {
-     			t.Fatalf("Failed to process Chinese word entries: %v", err)
-     		}
-     	}
-     }
- }
-
- // Write to files
- if err := proc.WriteToFiles(); err != nil {
-     t.Fatalf("Failed to write files: %v", err)
- }
  +}
- +// verifyMinification checks that entries were properly minified
  +func verifyMinification(t \*testing.T, testDir string) {
- // Count of files with issues
- var filesWithWildcards, filesWithEmptyArrays, totalFiles int
-
- // Walk through the output directory and check files
- err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
-     if err != nil {
-     	return err
-     }
-
-     // Skip directories
-     if info.IsDir() {
-     	return nil
-     }
-
-     // Only check .json.br files
-     if !strings.HasSuffix(path, ".json.br") {
-     	return nil
-     }
-
-     totalFiles++
-
-     // Read the file
-     file, err := os.Open(path)
-     if err != nil {
-     	t.Errorf("Failed to open file %s: %v", path, err)
-     	return nil
-     }
-     defer file.Close()
-
-     // Decompress the file
-     br := brotli.NewReader(file)
-     var buf bytes.Buffer
-     if _, err := io.Copy(&buf, br); err != nil {
-     	t.Errorf("Failed to decompress file %s: %v", path, err)
-     	return nil
-     }
-
-     // Parse the JSON
-     var group WordGroup
-     if err := json.Unmarshal(buf.Bytes(), &group); err != nil {
-     	t.Errorf("Failed to parse JSON from file %s: %v", path, err)
-     	return nil
-     }
-
-     // Check for wildcards in JMdict entries
-     for _, word := range group.WordJapanese {
-     	// Check senses for wildcards
-     	for _, sense := range word.Sense {
-     		// Check if appliesToKanji contains "*"
-     		for _, applies := range sense.AppliesToKanji {
-     			if applies == "*" {
-     				filesWithWildcards++
-     				t.Errorf("Found wildcard in appliesToKanji in file %s", path)
-     			}
-     		}
-
-     		// Check if appliesToKana contains "*"
-     		for _, applies := range sense.AppliesToKana {
-     			if applies == "*" {
-     				filesWithWildcards++
-     				t.Errorf("Found wildcard in appliesToKana in file %s", path)
-     			}
-     		}
-     	}
-
-     	// Check kana entries for wildcards
-     	for _, kana := range word.Kana {
-     		// Check if appliesToKanji contains "*"
-     		for _, applies := range kana.AppliesToKanji {
-     			if applies == "*" {
-     				filesWithWildcards++
-     				t.Errorf("Found wildcard in kana.appliesToKanji in file %s", path)
-     			}
-     		}
-     	}
-     }
-
-     // Check for empty arrays in Chinese character entries
-     for _, char := range group.CharChinese {
-     	// Check if Definitions is an empty array
-     	if len(char.Definitions) == 0 && char.Definitions != nil {
-     		filesWithEmptyArrays++
-     		t.Errorf("Found empty Definitions array in file %s", path)
-     	}
-
-     	// Check if Pinyin is an empty array
-     	if len(char.Pinyin) == 0 && char.Pinyin != nil {
-     		filesWithEmptyArrays++
-     		t.Errorf("Found empty Pinyin array in file %s", path)
-     	}
-     }
-
-     // Check for empty arrays in Chinese word entries
-     for _, word := range group.WordChinese {
-     	// Check if Definitions is an empty array
-     	if len(word.Definitions) == 0 && word.Definitions != nil {
-     		filesWithEmptyArrays++
-     		t.Errorf("Found empty Definitions array in file %s", path)
-     	}
-
-     	// Check if Pinyin is an empty array
-     	if len(word.Pinyin) == 0 && word.Pinyin != nil {
-     		filesWithEmptyArrays++
-     		t.Errorf("Found empty Pinyin array in file %s", path)
-     	}
-
-     	// Check if Frequency is an empty map
-     	if word.Frequency != nil && len(word.Frequency) == 0 {
-     		filesWithEmptyArrays++
-     		t.Errorf("Found empty Frequency map in file %s", path)
-     	}
-     }
-
-     return nil
- })
-
- if err != nil {
-     t.Errorf("Error walking output directory: %v", err)
- }
-
- t.Logf("Verified %d files", totalFiles)
-
- // We expect no files with wildcards or empty arrays
- if filesWithWildcards > 0 {
-     t.Errorf("Found %d files with wildcards, expected 0", filesWithWildcards)
- }
-
- if filesWithEmptyArrays > 0 {
-     t.Errorf("Found %d files with empty arrays, expected 0", filesWithEmptyArrays)
- }
  +}
  diff --git a/processor/processor.go b/processor/processor.go
  index 2f545f96df..1b68f0af2d 100644
  --- a/processor/processor.go
  +++ b/processor/processor.go
  @@ -1,10 +1,8 @@
  package processor

import (

- "encoding/json"
  "fmt"
  "os"
- "path/filepath"
  "sync"
  "kiokun-go/dictionaries/chinese_chars"
  @@ -13,40 +11,8 @@ import (
  "kiokun-go/dictionaries/jmdict"
  "kiokun-go/dictionaries/jmnedict"
  "kiokun-go/dictionaries/kanjidic"
-
- "github.com/andybalholm/brotli"
  )

-// WordGroup represents the combined data for a single word/character
-type WordGroup struct {

- WordJapanese []jmdict.Word `json:"w_j,omitempty"`
- NameJapanese []jmnedict.Name `json:"n_j,omitempty"`
- CharJapanese []kanjidic.Kanji `json:"c_j,omitempty"`
- CharChinese []chinese_chars.ChineseCharEntry `json:"c_c,omitempty"`
- WordChinese []chinese_words.ChineseWordEntry `json:"w_c,omitempty"`
  -}
- -// HasMultipleDictData returns true if this group has data from multiple dictionaries
  -func (wg \*WordGroup) HasMultipleDictData() bool {
- sources := 0
- if len(wg.WordJapanese) > 0 {
-     sources++
- }
- if len(wg.NameJapanese) > 0 {
-     sources++
- }
- if len(wg.CharJapanese) > 0 {
-     sources++
- }
- if len(wg.CharChinese) > 0 {
-     sources++
- }
- if len(wg.WordChinese) > 0 {
-     sources++
- }
- return sources > 1
  -}
- // DictionaryProcessor handles combining and writing dictionary entries
  type DictionaryProcessor struct {
  outputDir string
  @@ -101,83 +67,6 @@ func (p \*DictionaryProcessor) processEntry(entry common.Entry) error {
  }
  }
  -// processJMdictWord processes a JMdict word entry
  -func (p \*DictionaryProcessor) processJMdictWord(word jmdict.Word) error {
- // Get all forms (kanji and kana)
- var forms []string
- for \_, k := range word.Kanji {
-     forms = append(forms, k.Text)
- }
- for \_, k := range word.Kana {
-     forms = append(forms, k.Text)
- }
-
- // Add word to each form's group
- for \_, form := range forms {
-     group := p.getOrCreateGroup(form)
-     group.WordJapanese = append(group.WordJapanese, word)
- }
-
- return nil
  -}
- -// processJMNedictName processes a JMNedict name entry
  -func (p \*DictionaryProcessor) processJMNedictName(name jmnedict.Name) error {
- // Get all forms (kanji and kana)
- var forms []string
- forms = append(forms, name.Kanji...)
- forms = append(forms, name.Reading...)
-
- // Add name to each form's group
- for \_, form := range forms {
-     group := p.getOrCreateGroup(form)
-     group.NameJapanese = append(group.NameJapanese, name)
- }
-
- return nil
  -}
- -// processKanjidicEntry processes a Kanjidic entry
  -func (p \*DictionaryProcessor) processKanjidicEntry(kanji kanjidic.Kanji) error {
- group := p.getOrCreateGroup(kanji.Character)
- group.CharJapanese = append(group.CharJapanese, kanji)
- return nil
  -}
- -// processChineseCharEntry processes a Chinese character entry
  -func (p \*DictionaryProcessor) processChineseCharEntry(char chinese_chars.ChineseCharEntry) error {
- // Add both traditional and simplified forms
- forms := []string{char.Traditional}
- if char.Simplified != char.Traditional {
-     forms = append(forms, char.Simplified)
- }
-
- // Add character to each form's group
- for \_, form := range forms {
-     group := p.getOrCreateGroup(form)
-     group.CharChinese = append(group.CharChinese, char)
- }
-
- return nil
  -}
- -// processChineseWordEntry processes a Chinese word entry
  -func (p \*DictionaryProcessor) processChineseWordEntry(word chinese_words.ChineseWordEntry) error {
- // Add both traditional and simplified forms
- forms := []string{word.Traditional}
- if word.Simplified != word.Traditional {
-     forms = append(forms, word.Simplified)
- }
-
- // Add word to each form's group
- for \_, form := range forms {
-     group := p.getOrCreateGroup(form)
-     group.WordChinese = append(group.WordChinese, word)
- }
-
- return nil
  -}
- // getOrCreateGroup gets or creates a WordGroup for a given form
  func (p *DictionaryProcessor) getOrCreateGroup(form string) *WordGroup {
  // First try read-only access
  @@ -201,324 +90,3 @@ func (p *DictionaryProcessor) getOrCreateGroup(form string) *WordGroup {
  }
  return group
  }
- -// sanitizeWordGroup removes duplicate entries from a word group
  -func sanitizeWordGroup(group \*WordGroup) {
- // Deduplicate JMdict entries
- seen := make(map[string]bool)
- var uniqueWords []jmdict.Word
- for \_, word := range group.WordJapanese {
-     if !seen[word.ID] {
-     	seen[word.ID] = true
-     	uniqueWords = append(uniqueWords, word)
-     }
- }
- group.WordJapanese = uniqueWords
-
- // Deduplicate JMNedict entries
- seen = make(map[string]bool)
- var uniqueNames []jmnedict.Name
- for \_, name := range group.NameJapanese {
-     if !seen[name.ID] {
-     	seen[name.ID] = true
-     	uniqueNames = append(uniqueNames, name)
-     }
- }
- group.NameJapanese = uniqueNames
-
- // Deduplicate Kanjidic entries
- seen = make(map[string]bool)
- var uniqueKanji []kanjidic.Kanji
- for \_, kanji := range group.CharJapanese {
-     if !seen[kanji.Character] {
-     	seen[kanji.Character] = true
-     	uniqueKanji = append(uniqueKanji, kanji)
-     }
- }
- group.CharJapanese = uniqueKanji
-
- // Deduplicate Chinese character entries
- seen = make(map[string]bool)
- var uniqueChineseChars []chinese_chars.ChineseCharEntry
- for \_, char := range group.CharChinese {
-     if !seen[char.ID] {
-     	seen[char.ID] = true
-     	uniqueChineseChars = append(uniqueChineseChars, char)
-     }
- }
- group.CharChinese = uniqueChineseChars
-
- // Deduplicate Chinese word entries
- seen = make(map[string]bool)
- var uniqueChineseWords []chinese_words.ChineseWordEntry
- for \_, word := range group.WordChinese {
-     if !seen[word.ID] {
-     	seen[word.ID] = true
-     	uniqueChineseWords = append(uniqueChineseWords, word)
-     }
- }
- group.WordChinese = uniqueChineseWords
  -}
- -// WriteToFiles writes all groups to their respective files
  -func (p \*DictionaryProcessor) WriteToFiles() error {
- // Lock the map while we're iterating over it
- p.groupsMu.RLock()
-
- // Create a copy of groups to avoid holding the lock during processing
- groupsCopy := make(map[string]\*WordGroup, len(p.groups))
- for form, group := range p.groups {
-     groupsCopy[form] = group
- }
- // Now we can release the read lock
- p.groupsMu.RUnlock()
-
- // Statistics variables
- jmdictTotal := 0
- jmnedictTotal := 0
- kanjidicTotal := 0
- chineseCharsTotal := 0
- chineseWordsTotal := 0
- combinedEntries := 0
- totalFiles := len(groupsCopy)
- processed := 0
- errorCount := 0
-
- // Track entries from multiple dictionaries for verification
- var multiDictEntries []string
- var wordAndKanjiEntries []string
- var wordKanjiAndNameEntries []string
-
- // Sequential file writing
- if p.workerCount <= 1 {
-     for form, group := range groupsCopy {
-     	if err := p.writeGroupToFile(form, group); err != nil {
-     		return fmt.Errorf("error writing %s: %v", form, err)
-     	}
-
-     	processed++
-     	if processed%1000 == 0 || processed == totalFiles {
-     		fmt.Printf("\rWriting files: %d/%d (%.1f%%)...", processed, totalFiles, float64(processed)/float64(totalFiles)*100)
-     	}
-
-     	// Record entries that have data from multiple dictionaries
-     	hasWord := len(group.WordJapanese) > 0
-     	hasKanji := len(group.CharJapanese) > 0
-     	hasName := len(group.NameJapanese) > 0
-     	hasChineseChar := len(group.CharChinese) > 0
-     	hasChineseWord := len(group.WordChinese) > 0
-
-     	if hasWord {
-     		jmdictTotal++
-     	}
-     	if hasName {
-     		jmnedictTotal++
-     	}
-     	if hasKanji {
-     		kanjidicTotal++
-     	}
-     	if hasChineseChar {
-     		chineseCharsTotal++
-     	}
-     	if hasChineseWord {
-     		chineseWordsTotal++
-     	}
-
-     	if hasWord && hasKanji && hasName {
-     		wordKanjiAndNameEntries = append(wordKanjiAndNameEntries, form)
-     	} else if hasWord && hasKanji {
-     		wordAndKanjiEntries = append(wordAndKanjiEntries, form)
-     	}
-
-     	// Check for any combination of multiple dict data
-     	dictCount := 0
-     	if hasWord {
-     		dictCount++
-     	}
-     	if hasKanji {
-     		dictCount++
-     	}
-     	if hasName {
-     		dictCount++
-     	}
-     	if hasChineseChar {
-     		dictCount++
-     	}
-     	if hasChineseWord {
-     		dictCount++
-     	}
-
-     	if dictCount > 1 {
-     		multiDictEntries = append(multiDictEntries, form)
-     		combinedEntries++
-     	}
-     }
-
-     fmt.Printf("\nWrote %d files successfully\n", totalFiles-errorCount)
-     fmt.Println("\nDictionary entry counts:")
-     fmt.Printf("- JMdict entries: %d\n", jmdictTotal)
-     fmt.Printf("- JMNedict entries: %d\n", jmnedictTotal)
-     fmt.Printf("- Kanjidic entries: %d\n", kanjidicTotal)
-     fmt.Printf("- Chinese character entries: %d\n", chineseCharsTotal)
-     fmt.Printf("- Chinese word entries: %d\n", chineseWordsTotal)
-     fmt.Printf("- Combined entries: %d\n", combinedEntries)
-
-     return nil
- }
-
- // Parallel file writing with worker pool
- type writeTask struct {
-     form  string
-     group *WordGroup
- }
-
- type writeResult struct {
-     form           string
-     err            error
-     hasWord        bool
-     hasKanji       bool
-     hasName        bool
-     hasChineseChar bool
-     hasChineseWord bool
- }
-
- taskChan := make(chan writeTask, min(1000, len(groupsCopy)))
- resultChan := make(chan writeResult, min(1000, len(groupsCopy)))
-
- // Launch worker pool
- var wg sync.WaitGroup
- for i := 0; i < p.workerCount; i++ {
-     wg.Add(1)
-     go func() {
-     	defer wg.Done()
-     	for task := range taskChan {
-     		err := p.writeGroupToFile(task.form, task.group)
-
-     		hasWord := len(task.group.WordJapanese) > 0
-     		hasKanji := len(task.group.CharJapanese) > 0
-     		hasName := len(task.group.NameJapanese) > 0
-     		hasChineseChar := len(task.group.CharChinese) > 0
-     		hasChineseWord := len(task.group.WordChinese) > 0
-
-     		resultChan <- writeResult{
-     			form:           task.form,
-     			err:            err,
-     			hasWord:        hasWord,
-     			hasKanji:       hasKanji,
-     			hasName:        hasName,
-     			hasChineseChar: hasChineseChar,
-     			hasChineseWord: hasChineseWord,
-     		}
-     	}
-     }()
- }
-
- // Send all tasks
- go func() {
-     for form, group := range groupsCopy {
-     	taskChan <- writeTask{form, group}
-     }
-     close(taskChan)
-
-     // Wait for all workers to finish, then close the result channel
-     wg.Wait()
-     close(resultChan)
- }()
-
- // Collect results and statistics
- for result := range resultChan {
-     processed++
-     if processed%1000 == 0 || processed == totalFiles {
-     	fmt.Printf("\rWriting files: %d/%d (%.1f%%)...", processed, totalFiles, float64(processed)/float64(totalFiles)*100)
-     }
-
-     if result.err != nil {
-     	fmt.Printf("\nError writing %s: %v\n", result.form, result.err)
-     	errorCount++
-     	continue
-     }
-
-     // Count entries by dictionary type
-     if result.hasWord {
-     	jmdictTotal++
-     }
-     if result.hasName {
-     	jmnedictTotal++
-     }
-     if result.hasKanji {
-     	kanjidicTotal++
-     }
-     if result.hasChineseChar {
-     	chineseCharsTotal++
-     }
-     if result.hasChineseWord {
-     	chineseWordsTotal++
-     }
-
-     // Record entries with data from multiple dictionaries
-     dictCount := 0
-     if result.hasWord {
-     	dictCount++
-     }
-     if result.hasName {
-     	dictCount++
-     }
-     if result.hasKanji {
-     	dictCount++
-     }
-     if result.hasChineseChar {
-     	dictCount++
-     }
-     if result.hasChineseWord {
-     	dictCount++
-     }
-
-     if dictCount > 1 {
-     	combinedEntries++
-     	multiDictEntries = append(multiDictEntries, result.form)
-     }
- }
-
- fmt.Printf("\nWrote %d files successfully\n", totalFiles-errorCount)
- fmt.Println("\nDictionary entry counts:")
- fmt.Printf("- JMdict entries: %d\n", jmdictTotal)
- fmt.Printf("- JMNedict entries: %d\n", jmnedictTotal)
- fmt.Printf("- Kanjidic entries: %d\n", kanjidicTotal)
- fmt.Printf("- Chinese character entries: %d\n", chineseCharsTotal)
- fmt.Printf("- Chinese word entries: %d\n", chineseWordsTotal)
- fmt.Printf("- Combined entries: %d\n", combinedEntries)
- return nil
  -}
- -// min returns the minimum of two integers
  -func min(a, b int) int {
- if a < b {
-     return a
- }
- return b
  -}
- -// writeGroupToFile writes a single group to its file
  -func (p *DictionaryProcessor) writeGroupToFile(form string, group *WordGroup) error {
- // Sanitize data before writing to file
- sanitizeWordGroup(group)
-
- filePath := filepath.Join(p.outputDir, form+".json.br")
-
- file, err := os.Create(filePath)
- if err != nil {
-     return fmt.Errorf("failed to create file: %v", err)
- }
- defer file.Close()
-
- brWriter := brotli.NewWriter(file)
- defer brWriter.Close()
-
- encoder := json.NewEncoder(brWriter)
- encoder.SetEscapeHTML(false)
- if err := encoder.Encode(group); err != nil {
-     return fmt.Errorf("failed to encode group: %v", err)
- }
-
- return nil
  -}
  diff --git a/processor/processor_chinese.go b/processor/processor_chinese.go
  new file mode 100644
  index 0000000000..50e55887e5
  --- /dev/null
  +++ b/processor/processor_chinese.go
  @@ -0,0 +1,40 @@
  +package processor

* +import (
* "kiokun-go/dictionaries/chinese_chars"
* "kiokun-go/dictionaries/chinese_words"
  +)
* +// processChineseCharEntry processes a Chinese character entry
  +func (p \*DictionaryProcessor) processChineseCharEntry(char chinese_chars.ChineseCharEntry) error {
* // Add both traditional and simplified forms
* forms := []string{char.Traditional}
* if char.Simplified != char.Traditional {
*     forms = append(forms, char.Simplified)
* }
*
* // Add character to each form's group
* for \_, form := range forms {
*     group := p.getOrCreateGroup(form)
*     group.CharChinese = append(group.CharChinese, char)
* }
*
* return nil
  +}
* +// processChineseWordEntry processes a Chinese word entry
  +func (p \*DictionaryProcessor) processChineseWordEntry(word chinese_words.ChineseWordEntry) error {
* // Add both traditional and simplified forms
* forms := []string{word.Traditional}
* if word.Simplified != word.Traditional {
*     forms = append(forms, word.Simplified)
* }
*
* // Add word to each form's group
* for \_, form := range forms {
*     group := p.getOrCreateGroup(form)
*     group.WordChinese = append(group.WordChinese, word)
* }
*
* return nil
  +}
  diff --git a/processor/processor_japanese.go b/processor/processor_japanese.go
  new file mode 100644
  index 0000000000..58703e4ecf
  --- /dev/null
  +++ b/processor/processor_japanese.go
  @@ -0,0 +1,50 @@
  +package processor
* +import (
* "kiokun-go/dictionaries/jmdict"
* "kiokun-go/dictionaries/jmnedict"
* "kiokun-go/dictionaries/kanjidic"
  +)
* +// processJMdictWord processes a JMdict word entry
  +func (p \*DictionaryProcessor) processJMdictWord(word jmdict.Word) error {
* // Get all forms (kanji and kana)
* var forms []string
* for \_, k := range word.Kanji {
*     forms = append(forms, k.Text)
* }
* for \_, k := range word.Kana {
*     forms = append(forms, k.Text)
* }
*
* // Add word to each form's group
* for \_, form := range forms {
*     group := p.getOrCreateGroup(form)
*     group.WordJapanese = append(group.WordJapanese, word)
* }
*
* return nil
  +}
* +// processJMNedictName processes a JMNedict name entry
  +func (p \*DictionaryProcessor) processJMNedictName(name jmnedict.Name) error {
* // Get all forms (kanji and kana)
* var forms []string
* forms = append(forms, name.Kanji...)
* forms = append(forms, name.Reading...)
*
* // Add name to each form's group
* for \_, form := range forms {
*     group := p.getOrCreateGroup(form)
*     group.NameJapanese = append(group.NameJapanese, name)
* }
*
* return nil
  +}
* +// processKanjidicEntry processes a Kanjidic entry
  +func (p \*DictionaryProcessor) processKanjidicEntry(kanji kanjidic.Kanji) error {
* group := p.getOrCreateGroup(kanji.Character)
* group.CharJapanese = append(group.CharJapanese, kanji)
* return nil
  +}
  diff --git a/processor/processor_test.go b/processor/processor_test.go
  new file mode 100644
  index 0000000000..8fdaad0b5d
  --- /dev/null
  +++ b/processor/processor_test.go
  @@ -0,0 +1,148 @@
  +package processor
* +import (
* "encoding/json"
* "os"
* "path/filepath"
* "testing"
*
* "kiokun-go/dictionaries/chinese_chars"
* "kiokun-go/dictionaries/common"
* "kiokun-go/dictionaries/kanjidic"
*
* "github.com/andybalholm/brotli"
  +)
* +// readWordGroup reads a WordGroup from a compressed JSON file
  +func readWordGroup(filePath string) (\*WordGroup, error) {
* file, err := os.Open(filePath)
* if err != nil {
*     return nil, err
* }
* defer file.Close()
*
* // Decompress the file
* br := brotli.NewReader(file)
* decoder := json.NewDecoder(br)
*
* var group WordGroup
* if err := decoder.Decode(&group); err != nil {
*     return nil, err
* }
*
* return &group, nil
  +}
* +// TestChineseJapaneseIntegration tests the integration between Chinese and Japanese dictionaries
  +// by directly creating test entries and processing them
  +func TestChineseJapaneseIntegration(t \*testing.T) {
* // Create a temporary output directory
* testDir := filepath.Join(os.TempDir(), "kiokun_integration_test")
* err := os.MkdirAll(testDir, 0755)
* if err != nil {
*     t.Fatalf("Failed to create test directory: %v", err)
* }
* defer os.RemoveAll(testDir)
*
* // Create a processor
* proc, err := New(testDir, 1)
* if err != nil {
*     t.Fatalf("Failed to create processor: %v", err)
* }
*
* // Create test entries for characters that should exist in both Chinese and Japanese
* testChars := []string{"人", "日", "月"} // Person, Day/Sun, Month/Moon
*
* // Process each test character
* for \_, char := range testChars {
*     // Create a Japanese kanji entry
*     kanji := kanjidic.Kanji{
*     	Character: char,
*     	Meanings:  []string{"Test meaning for " + char},
*     	OnYomi:    []string{"on"},
*     	KunYomi:   []string{"kun"},
*     	Stroke:    5,
*     }
*
*     // Create a Chinese character entry
*     chineseChar := chinese_chars.ChineseCharEntry{
*     	ID:          "test_char_" + char,
*     	Traditional: char,
*     	Simplified:  char,
*     	Definitions: []string{"Test definition for " + char},
*     	Pinyin:      []string{"test"},
*     }
*
*     // Process both entries
*     entries := []common.Entry{kanji, chineseChar}
*     if err := proc.ProcessEntries(entries); err != nil {
*     	t.Fatalf("Failed to process entries for %s: %v", char, err)
*     }
* }
*
* // Write to files
* if err := proc.WriteToFiles(); err != nil {
*     t.Fatalf("Failed to write files: %v", err)
* }
*
* // Verify that files were created
* files, err := os.ReadDir(testDir)
* if err != nil {
*     t.Fatalf("Failed to read output directory: %v", err)
* }
*
* if len(files) == 0 {
*     t.Errorf("No files were created in the output directory")
* } else {
*     t.Logf("Created %d files in the output directory", len(files))
*
*     // List all files
*     t.Logf("Files:")
*     for _, file := range files {
*     	t.Logf("  - %s", file.Name())
*     }
* }
*
* // Check for combined entries (files that contain both Chinese and Japanese data)
* combinedEntries := 0
*
* for \_, char := range testChars {
*     // Check the file for this character
*     filePath := filepath.Join(testDir, char + ".json.br")
*
*     // Verify the file exists
*     if _, err := os.Stat(filePath); os.IsNotExist(err) {
*     	t.Errorf("Expected file %s does not exist", filePath)
*     	continue
*     }
*
*     // Read and check the file content
*     group, err := readWordGroup(filePath)
*     if err != nil {
*     	t.Errorf("Failed to read file %s: %v", filePath, err)
*     	continue
*     }
*
*     // Check if this group has both Japanese and Chinese data
*     hasJapanese := len(group.CharJapanese) > 0
*     hasChinese := len(group.CharChinese) > 0
*
*     if hasJapanese && hasChinese {
*     	combinedEntries++
*     	t.Logf("Found combined entry in file %s.json.br", char)
*     	t.Logf("  - Japanese data: %d characters", len(group.CharJapanese))
*     	t.Logf("  - Chinese data: %d characters", len(group.CharChinese))
*     } else {
*     	t.Errorf("File %s.json.br does not contain both Japanese and Chinese data", char)
*     	t.Logf("  - Japanese data: %d characters", len(group.CharJapanese))
*     	t.Logf("  - Chinese data: %d characters", len(group.CharChinese))
*     }
* }
*
* t.Logf("Found %d combined entries", combinedEntries)
*
* // We expect all test characters to have combined entries
* if combinedEntries != len(testChars) {
*     t.Errorf("Expected %d combined entries, found %d", len(testChars), combinedEntries)
* }
  +}
  diff --git a/processor/sanitize.go b/processor/sanitize.go
  new file mode 100644
  index 0000000000..92e1a35a3a
  --- /dev/null
  +++ b/processor/sanitize.go
  @@ -0,0 +1,145 @@
  +package processor
* +import (
* "kiokun-go/dictionaries/chinese_chars"
* "kiokun-go/dictionaries/chinese_words"
* "kiokun-go/dictionaries/jmdict"
* "kiokun-go/dictionaries/jmnedict"
* "kiokun-go/dictionaries/kanjidic"
  +)
* +// sanitizeWordGroup removes duplicate entries and minifies data in a word group
  +func sanitizeWordGroup(group \*WordGroup) {
* // Deduplicate and sanitize JMdict entries
* seen := make(map[string]bool)
* var uniqueWords []jmdict.Word
* for \_, word := range group.WordJapanese {
*     if !seen[word.ID] {
*     	seen[word.ID] = true
*
*     	// Sanitize each word by removing wildcards from senses and kana entries
*     	for i := range word.Sense {
*     		word.Sense[i].SanitizeWildcards()
*
*     		// Remove empty arrays in Sense
*     		if len(word.Sense[i].Related) == 0 {
*     			word.Sense[i].Related = nil
*     		}
*     		if len(word.Sense[i].Antonym) == 0 {
*     			word.Sense[i].Antonym = nil
*     		}
*     		if len(word.Sense[i].Field) == 0 {
*     			word.Sense[i].Field = nil
*     		}
*     		if len(word.Sense[i].Dialect) == 0 {
*     			word.Sense[i].Dialect = nil
*     		}
*     		if len(word.Sense[i].Misc) == 0 {
*     			word.Sense[i].Misc = nil
*     		}
*     		if len(word.Sense[i].Info) == 0 {
*     			word.Sense[i].Info = nil
*     		}
*     		if len(word.Sense[i].LanguageSource) == 0 {
*     			word.Sense[i].LanguageSource = nil
*     		}
*     		if len(word.Sense[i].Examples) == 0 {
*     			word.Sense[i].Examples = nil
*     		}
*     	}
*
*     	for i := range word.Kana {
*     		word.Kana[i].SanitizeWildcards()
*
*     		// Remove empty arrays in Kana
*     		if len(word.Kana[i].Tags) == 0 {
*     			word.Kana[i].Tags = nil
*     		}
*     	}
*
*     	// Remove empty arrays in Kanji
*     	for i := range word.Kanji {
*     		if len(word.Kanji[i].Tags) == 0 {
*     			word.Kanji[i].Tags = nil
*     		}
*     	}
*
*     	uniqueWords = append(uniqueWords, word)
*     }
* }
* group.WordJapanese = uniqueWords
*
* // Deduplicate JMNedict entries
* seen = make(map[string]bool)
* var uniqueNames []jmnedict.Name
* for \_, name := range group.NameJapanese {
*     if !seen[name.ID] {
*     	seen[name.ID] = true
*     	uniqueNames = append(uniqueNames, name)
*     }
* }
* group.NameJapanese = uniqueNames
*
* // Deduplicate Kanjidic entries
* seen = make(map[string]bool)
* var uniqueKanji []kanjidic.Kanji
* for \_, kanji := range group.CharJapanese {
*     if !seen[kanji.Character] {
*     	seen[kanji.Character] = true
*
*     	// Sanitize Kanjidic entries by removing empty fields
*     	// (This would need to be expanded based on the Kanjidic structure)
*
*     	uniqueKanji = append(uniqueKanji, kanji)
*     }
* }
* group.CharJapanese = uniqueKanji
*
* // Deduplicate and sanitize Chinese character entries
* seen = make(map[string]bool)
* var uniqueChineseChars []chinese_chars.ChineseCharEntry
* for \_, char := range group.CharChinese {
*     if !seen[char.ID] {
*     	seen[char.ID] = true
*
*     	// Sanitize Chinese character entries by removing empty slices
*     	if len(char.Definitions) == 0 {
*     		char.Definitions = nil
*     	}
*     	if len(char.Pinyin) == 0 {
*     		char.Pinyin = nil
*     	}
*
*     	// If stroke count is 0, omit it
*     	if char.StrokeCount == 0 {
*     		char.StrokeCount = 0 // This will be omitted due to omitempty tag
*     	}
*
*     	uniqueChineseChars = append(uniqueChineseChars, char)
*     }
* }
* group.CharChinese = uniqueChineseChars
*
* // Deduplicate and sanitize Chinese word entries
* seen = make(map[string]bool)
* var uniqueChineseWords []chinese_words.ChineseWordEntry
* for \_, word := range group.WordChinese {
*     if !seen[word.ID] {
*     	seen[word.ID] = true
*
*     	// Sanitize Chinese word entries by removing empty slices
*     	if len(word.Definitions) == 0 {
*     		word.Definitions = nil
*     	}
*     	if len(word.Pinyin) == 0 {
*     		word.Pinyin = nil
*     	}
*     	if word.Frequency != nil && len(word.Frequency) == 0 {
*     		word.Frequency = nil
*     	}
*
*     	uniqueChineseWords = append(uniqueChineseWords, word)
*     }
* }
* group.WordChinese = uniqueChineseWords
  +}
  diff --git a/processor/types.go b/processor/types.go
  new file mode 100644
  index 0000000000..c329ecb574
  --- /dev/null
  +++ b/processor/types.go
  @@ -0,0 +1,39 @@
  +package processor
* +import (
* "kiokun-go/dictionaries/chinese_chars"
* "kiokun-go/dictionaries/chinese_words"
* "kiokun-go/dictionaries/jmdict"
* "kiokun-go/dictionaries/jmnedict"
* "kiokun-go/dictionaries/kanjidic"
  +)
* +// WordGroup represents the combined data for a single word/character
  +type WordGroup struct {
* WordJapanese []jmdict.Word `json:"w_j,omitempty"`
* NameJapanese []jmnedict.Name `json:"n_j,omitempty"`
* CharJapanese []kanjidic.Kanji `json:"c_j,omitempty"`
* CharChinese []chinese_chars.ChineseCharEntry `json:"c_c,omitempty"`
* WordChinese []chinese_words.ChineseWordEntry `json:"w_c,omitempty"`
  +}
* +// HasMultipleDictData returns true if this group has data from multiple dictionaries
  +func (wg \*WordGroup) HasMultipleDictData() bool {
* sources := 0
* if len(wg.WordJapanese) > 0 {
*     sources++
* }
* if len(wg.NameJapanese) > 0 {
*     sources++
* }
* if len(wg.CharJapanese) > 0 {
*     sources++
* }
* if len(wg.CharChinese) > 0 {
*     sources++
* }
* if len(wg.WordChinese) > 0 {
*     sources++
* }
* return sources > 1
  +}
  diff --git a/processor/types_export.go b/processor/types_export.go
  new file mode 100644
  index 0000000000..c7f25652b8
  --- /dev/null
  +++ b/processor/types_export.go
  @@ -0,0 +1,16 @@
  +package processor
* +import (
* "kiokun-go/dictionaries/chinese_chars"
* "kiokun-go/dictionaries/chinese_words"
* "kiokun-go/dictionaries/jmdict"
* "kiokun-go/dictionaries/jmnedict"
* "kiokun-go/dictionaries/kanjidic"
  +)
* +// Type aliases for external use
  +type ChineseCharEntry = chinese_chars.ChineseCharEntry
  +type ChineseWordEntry = chinese_words.ChineseWordEntry
  +type JMdictWord = jmdict.Word
  +type JMNedictName = jmnedict.Name
  +type KanjidicKanji = kanjidic.Kanji
  diff --git a/processor/writer.go b/processor/writer.go
  new file mode 100644
  index 0000000000..c5c5e79445
  --- /dev/null
  +++ b/processor/writer.go
  @@ -0,0 +1,315 @@
  +package processor
* +import (
* "encoding/json"
* "fmt"
* "os"
* "path/filepath"
* "sync"
*
* "github.com/andybalholm/brotli"
  +)
* +// WriteToFiles writes all groups to their respective files
  +func (p \*DictionaryProcessor) WriteToFiles() error {
* // Lock the map while we're iterating over it
* p.groupsMu.RLock()
*
* // Create a copy of groups to avoid holding the lock during processing
* groupsCopy := make(map[string]\*WordGroup, len(p.groups))
* for form, group := range p.groups {
*     groupsCopy[form] = group
* }
* // Now we can release the read lock
* p.groupsMu.RUnlock()
*
* // Statistics variables
* jmdictTotal := 0
* jmnedictTotal := 0
* kanjidicTotal := 0
* chineseCharsTotal := 0
* chineseWordsTotal := 0
* combinedEntries := 0
* totalFiles := len(groupsCopy)
* processed := 0
* errorCount := 0
*
* // Track entries from multiple dictionaries for verification
* var multiDictEntries []string
* var wordAndKanjiEntries []string
* var wordKanjiAndNameEntries []string
*
* // Sequential file writing
* if p.workerCount <= 1 {
*     return p.writeSequential(groupsCopy, totalFiles, &processed, &errorCount,
*     	&jmdictTotal, &jmnedictTotal, &kanjidicTotal, &chineseCharsTotal, &chineseWordsTotal,
*     	&combinedEntries, &multiDictEntries, &wordAndKanjiEntries, &wordKanjiAndNameEntries)
* }
*
* // Parallel file writing
* return p.writeParallel(groupsCopy, totalFiles, &processed, &errorCount,
*     &jmdictTotal, &jmnedictTotal, &kanjidicTotal, &chineseCharsTotal, &chineseWordsTotal,
*     &combinedEntries, &multiDictEntries)
  +}
* +// writeSequential handles sequential file writing
  +func (p \*DictionaryProcessor) writeSequential(
* groupsCopy map[string]\*WordGroup,
* totalFiles int,
* processed \*int,
* errorCount \*int,
* jmdictTotal \*int,
* jmnedictTotal \*int,
* kanjidicTotal \*int,
* chineseCharsTotal \*int,
* chineseWordsTotal \*int,
* combinedEntries \*int,
* multiDictEntries \*[]string,
* wordAndKanjiEntries \*[]string,
* wordKanjiAndNameEntries \*[]string,
  +) error {
* for form, group := range groupsCopy {
*     if err := p.writeGroupToFile(form, group); err != nil {
*     	return fmt.Errorf("error writing %s: %v", form, err)
*     }
*
*     *processed++
*     if *processed%1000 == 0 || *processed == totalFiles {
*     	fmt.Printf("\rWriting files: %d/%d (%.1f%%)...", *processed, totalFiles, float64(*processed)/float64(totalFiles)*100)
*     }
*
*     // Record entries that have data from multiple dictionaries
*     hasWord := len(group.WordJapanese) > 0
*     hasKanji := len(group.CharJapanese) > 0
*     hasName := len(group.NameJapanese) > 0
*     hasChineseChar := len(group.CharChinese) > 0
*     hasChineseWord := len(group.WordChinese) > 0
*
*     if hasWord {
*     	*jmdictTotal++
*     }
*     if hasName {
*     	*jmnedictTotal++
*     }
*     if hasKanji {
*     	*kanjidicTotal++
*     }
*     if hasChineseChar {
*     	*chineseCharsTotal++
*     }
*     if hasChineseWord {
*     	*chineseWordsTotal++
*     }
*
*     if hasWord && hasKanji && hasName {
*     	*wordKanjiAndNameEntries = append(*wordKanjiAndNameEntries, form)
*     } else if hasWord && hasKanji {
*     	*wordAndKanjiEntries = append(*wordAndKanjiEntries, form)
*     }
*
*     // Check for any combination of multiple dict data
*     dictCount := 0
*     if hasWord {
*     	dictCount++
*     }
*     if hasKanji {
*     	dictCount++
*     }
*     if hasName {
*     	dictCount++
*     }
*     if hasChineseChar {
*     	dictCount++
*     }
*     if hasChineseWord {
*     	dictCount++
*     }
*
*     if dictCount > 1 {
*     	*multiDictEntries = append(*multiDictEntries, form)
*     	*combinedEntries++
*     }
* }
*
* fmt.Printf("\nWrote %d files successfully\n", totalFiles-\*errorCount)
* fmt.Println("\nDictionary entry counts:")
* fmt.Printf("- JMdict entries: %d\n", \*jmdictTotal)
* fmt.Printf("- JMNedict entries: %d\n", \*jmnedictTotal)
* fmt.Printf("- Kanjidic entries: %d\n", \*kanjidicTotal)
* fmt.Printf("- Chinese character entries: %d\n", \*chineseCharsTotal)
* fmt.Printf("- Chinese word entries: %d\n", \*chineseWordsTotal)
* fmt.Printf("- Combined entries: %d\n", \*combinedEntries)
*
* return nil
  +}
* +// writeParallel handles parallel file writing with worker pool
  +func (p \*DictionaryProcessor) writeParallel(
* groupsCopy map[string]\*WordGroup,
* totalFiles int,
* processed \*int,
* errorCount \*int,
* jmdictTotal \*int,
* jmnedictTotal \*int,
* kanjidicTotal \*int,
* chineseCharsTotal \*int,
* chineseWordsTotal \*int,
* combinedEntries \*int,
* multiDictEntries \*[]string,
  +) error {
* type writeTask struct {
*     form  string
*     group *WordGroup
* }
*
* type writeResult struct {
*     form           string
*     err            error
*     hasWord        bool
*     hasKanji       bool
*     hasName        bool
*     hasChineseChar bool
*     hasChineseWord bool
* }
*
* taskChan := make(chan writeTask, min(1000, len(groupsCopy)))
* resultChan := make(chan writeResult, min(1000, len(groupsCopy)))
*
* // Launch worker pool
* var wg sync.WaitGroup
* for i := 0; i < p.workerCount; i++ {
*     wg.Add(1)
*     go func() {
*     	defer wg.Done()
*     	for task := range taskChan {
*     		err := p.writeGroupToFile(task.form, task.group)
*
*     		hasWord := len(task.group.WordJapanese) > 0
*     		hasKanji := len(task.group.CharJapanese) > 0
*     		hasName := len(task.group.NameJapanese) > 0
*     		hasChineseChar := len(task.group.CharChinese) > 0
*     		hasChineseWord := len(task.group.WordChinese) > 0
*
*     		resultChan <- writeResult{
*     			form:           task.form,
*     			err:            err,
*     			hasWord:        hasWord,
*     			hasKanji:       hasKanji,
*     			hasName:        hasName,
*     			hasChineseChar: hasChineseChar,
*     			hasChineseWord: hasChineseWord,
*     		}
*     	}
*     }()
* }
*
* // Send all tasks
* go func() {
*     for form, group := range groupsCopy {
*     	taskChan <- writeTask{form, group}
*     }
*     close(taskChan)
*
*     // Wait for all workers to finish, then close the result channel
*     wg.Wait()
*     close(resultChan)
* }()
*
* // Collect results and statistics
* for result := range resultChan {
*     *processed++
*     if *processed%1000 == 0 || *processed == totalFiles {
*     	fmt.Printf("\rWriting files: %d/%d (%.1f%%)...", *processed, totalFiles, float64(*processed)/float64(totalFiles)*100)
*     }
*
*     if result.err != nil {
*     	fmt.Printf("\nError writing %s: %v\n", result.form, result.err)
*     	*errorCount++
*     	continue
*     }
*
*     // Count entries by dictionary type
*     if result.hasWord {
*     	*jmdictTotal++
*     }
*     if result.hasName {
*     	*jmnedictTotal++
*     }
*     if result.hasKanji {
*     	*kanjidicTotal++
*     }
*     if result.hasChineseChar {
*     	*chineseCharsTotal++
*     }
*     if result.hasChineseWord {
*     	*chineseWordsTotal++
*     }
*
*     // Record entries with data from multiple dictionaries
*     dictCount := 0
*     if result.hasWord {
*     	dictCount++
*     }
*     if result.hasName {
*     	dictCount++
*     }
*     if result.hasKanji {
*     	dictCount++
*     }
*     if result.hasChineseChar {
*     	dictCount++
*     }
*     if result.hasChineseWord {
*     	dictCount++
*     }
*
*     if dictCount > 1 {
*     	*combinedEntries++
*     	*multiDictEntries = append(*multiDictEntries, result.form)
*     }
* }
*
* fmt.Printf("\nWrote %d files successfully\n", totalFiles-\*errorCount)
* fmt.Println("\nDictionary entry counts:")
* fmt.Printf("- JMdict entries: %d\n", \*jmdictTotal)
* fmt.Printf("- JMNedict entries: %d\n", \*jmnedictTotal)
* fmt.Printf("- Kanjidic entries: %d\n", \*kanjidicTotal)
* fmt.Printf("- Chinese character entries: %d\n", \*chineseCharsTotal)
* fmt.Printf("- Chinese word entries: %d\n", \*chineseWordsTotal)
* fmt.Printf("- Combined entries: %d\n", \*combinedEntries)
* return nil
  +}
* +// min returns the minimum of two integers
  +func min(a, b int) int {
* if a < b {
*     return a
* }
* return b
  +}
* +// writeGroupToFile writes a single group to its file
  +func (p *DictionaryProcessor) writeGroupToFile(form string, group *WordGroup) error {
* // Sanitize data before writing to file
* sanitizeWordGroup(group)
*
* filePath := filepath.Join(p.outputDir, form+".json.br")
*
* file, err := os.Create(filePath)
* if err != nil {
*     return fmt.Errorf("failed to create file: %v", err)
* }
* defer file.Close()
*
* // Use maximum compression level (11) for best compression ratio
* brWriter := brotli.NewWriterLevel(file, brotli.BestCompression)
* defer brWriter.Close()
*
* encoder := json.NewEncoder(brWriter)
* encoder.SetEscapeHTML(false)
* if err := encoder.Encode(group); err != nil {
*     return fmt.Errorf("failed to encode group: %v", err)
* }
*
* return nil
  +}
  diff --git a/run_tests.sh b/run_tests.sh
  new file mode 100755
  index 0000000000..37ecc3c5a3
  --- /dev/null
  +++ b/run_tests.sh
  @@ -0,0 +1,12 @@
  +#!/bin/bash
* +# run_tests.sh - A simple script to run minification tests
* +set -e # Exit on error
* +echo "Running minification tests..."
* +# Run processor tests with verbose output
  +go test -v ./processor
* +echo "Minification tests completed successfully!"
