package processor

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"

	"github.com/andybalholm/brotli"
)

// IndexEntry represents an entry in the index with minimal field names
// Using single-letter field names for maximum compression
type IndexEntry struct {
	// Exact matches (when the key exactly matches the entry)
	E map[string][]int64 `json:"e,omitempty"` // Exact matches by dictionary type (j, n, d, c, w)

	// Contained-in matches (when the key is contained within the entry)
	C map[string][]int64 `json:"c,omitempty"` // Contained-in matches by dictionary type (j, n, d, c, w)
}

// IndexProcessor processes dictionary entries and builds an index
type IndexProcessor struct {
	baseDir         string
	indexDir        string
	jmdictDir       string
	jmnedictDir     string
	kanjidicDir     string
	chineseCharsDir string
	chineseWordsDir string
	index           map[string]*IndexEntry
	writtenEntries  map[string]bool
	fileWriters     int
	mu              sync.Mutex
}

// NewIndexProcessor creates a new index-based processor
func NewIndexProcessor(baseDir string, fileWriters int) (*IndexProcessor, error) {
	// Create the processor
	p := &IndexProcessor{
		baseDir:        baseDir,
		index:          make(map[string]*IndexEntry),
		writtenEntries: make(map[string]bool),
		fileWriters:    fileWriters,
	}

	// Create the output directories
	if err := p.createDirectories(); err != nil {
		return nil, err
	}

	return p, nil
}

// createDirectories creates the necessary output directories
func (p *IndexProcessor) createDirectories() error {
	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(p.baseDir, 0755); err != nil {
		return err
	}

	// Create subdirectories with one-letter names
	indexDir := filepath.Join(p.baseDir, "index")
	jmdictDir := filepath.Join(p.baseDir, "j")       // j for JMdict
	jmnedictDir := filepath.Join(p.baseDir, "n")     // n for JMNedict
	kanjidicDir := filepath.Join(p.baseDir, "d")     // d for Kanjidic
	chineseCharsDir := filepath.Join(p.baseDir, "c") // c for Chinese Characters
	chineseWordsDir := filepath.Join(p.baseDir, "w") // w for Chinese Words

	// Create each directory
	for _, dir := range []string{
		indexDir,
		jmdictDir,
		jmnedictDir,
		kanjidicDir,
		chineseCharsDir,
		chineseWordsDir,
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Store the directory paths
	p.indexDir = indexDir
	p.jmdictDir = jmdictDir
	p.jmnedictDir = jmnedictDir
	p.kanjidicDir = kanjidicDir
	p.chineseCharsDir = chineseCharsDir
	p.chineseWordsDir = chineseWordsDir

	return nil
}

// getEntryKeys returns all possible keys for an entry
func getEntryKeys(entry common.Entry) []string {
	var keys []string

	switch e := entry.(type) {
	case jmdict.Word:
		// Add all kanji forms
		for _, k := range e.Kanji {
			keys = append(keys, k.Text)
		}
		// Add all kana forms
		for _, k := range e.Kana {
			keys = append(keys, k.Text)
		}
		// If no forms, use ID
		if len(keys) == 0 {
			keys = append(keys, e.ID)
		}
	case jmnedict.Name:
		// Add all kanji forms
		keys = append(keys, e.Kanji...)
		// Add all reading forms
		keys = append(keys, e.Reading...)
		// If no forms, use ID
		if len(keys) == 0 {
			keys = append(keys, e.ID)
		}
	case kanjidic.Kanji:
		// Add the character
		keys = append(keys, e.Character)
	case chinese_chars.ChineseCharEntry:
		// Add the traditional form
		keys = append(keys, e.Traditional)
		// Add the simplified form if different
		if e.Simplified != e.Traditional {
			keys = append(keys, e.Simplified)
		}
	case chinese_words.ChineseWordEntry:
		// Add the traditional form
		keys = append(keys, e.Traditional)
		// Add the simplified form if different
		if e.Simplified != e.Traditional {
			keys = append(keys, e.Simplified)
		}
	default:
		// For unknown entry types, use ID
		keys = append(keys, entry.GetID())
	}

	return keys
}

// ProcessEntries processes a slice of entries
func (p *IndexProcessor) ProcessEntries(entries []common.Entry) error {
	for _, entry := range entries {
		if err := p.processEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

// processEntry processes a single entry
func (p *IndexProcessor) processEntry(entry common.Entry) error {
	// Get the entry ID
	originalID := entry.GetID()

	// Get the shard type
	shardType := GetShardType(entry)

	// Create the sharded ID by prepending the shard type
	id := fmt.Sprintf("%d%s", shardType, originalID)

	// Convert string ID to int64 where possible
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		// If ID is not a number, use a hash of the string as the ID
		h := fnv.New64a()
		h.Write([]byte(id))
		idInt = int64(h.Sum64())
	}

	// Determine exact matches and contained-in matches based on entry type
	var exactMatches, containedMatches []string

	switch e := entry.(type) {
	case jmdict.Word:
		// For JMdict words, exact matches are the kanji and kana forms
		for _, k := range e.Kanji {
			exactMatches = append(exactMatches, k.Text)
		}
		for _, k := range e.Kana {
			exactMatches = append(exactMatches, k.Text)
		}

		// If no forms, use ID as exact match
		if len(exactMatches) == 0 {
			exactMatches = append(exactMatches, e.ID)
		}

		// For multi-character entries, each character is a contained-in match
		for _, form := range exactMatches {
			for _, char := range form {
				// Skip non-CJK characters
				if !isHanCharacter(string(char)) {
					continue
				}
				containedMatches = append(containedMatches, string(char))
			}
		}

	case jmnedict.Name:
		// For JMNedict names, exact matches are the kanji and reading forms
		exactMatches = append(exactMatches, e.Kanji...)
		exactMatches = append(exactMatches, e.Reading...)

		// If no forms, use ID as exact match
		if len(exactMatches) == 0 {
			exactMatches = append(exactMatches, e.ID)
		}

		// For multi-character entries, each character is a contained-in match
		for _, form := range exactMatches {
			for _, char := range form {
				// Skip non-CJK characters
				if !isHanCharacter(string(char)) {
					continue
				}
				containedMatches = append(containedMatches, string(char))
			}
		}

	case kanjidic.Kanji:
		// For Kanjidic entries, the character is an exact match
		// There are no contained-in matches for single characters
		exactMatches = append(exactMatches, e.Character)

	case chinese_chars.ChineseCharEntry:
		// For Chinese character entries, the traditional and simplified forms are exact matches
		// There are no contained-in matches for single characters
		exactMatches = append(exactMatches, e.Traditional)
		if e.Simplified != e.Traditional {
			exactMatches = append(exactMatches, e.Simplified)
		}

	case chinese_words.ChineseWordEntry:
		// For Chinese word entries, the traditional and simplified forms are exact matches
		exactMatches = append(exactMatches, e.Traditional)
		if e.Simplified != e.Traditional {
			exactMatches = append(exactMatches, e.Simplified)
		}

		// For multi-character entries, each character is a contained-in match
		for _, form := range exactMatches {
			for _, char := range form {
				// Skip non-CJK characters
				if !isHanCharacter(string(char)) {
					continue
				}
				containedMatches = append(containedMatches, string(char))
			}
		}

	default:
		// For unknown entry types, use ID as exact match
		exactMatches = append(exactMatches, entry.GetID())
	}

	// Remove duplicates from containedMatches
	containedMatches = removeDuplicates(containedMatches)

	// Remove exact matches from contained-in matches
	containedMatches = removeExactMatches(containedMatches, exactMatches)

	// Get the dictionary type code
	var dictType string
	switch entry.(type) {
	case jmdict.Word:
		dictType = "j"
	case jmnedict.Name:
		dictType = "n"
	case kanjidic.Kanji:
		dictType = "d"
	case chinese_chars.ChineseCharEntry:
		dictType = "c"
	case chinese_words.ChineseWordEntry:
		dictType = "w"
	}

	// Add the entry to the index for each key
	p.mu.Lock()

	// Process exact matches
	for _, key := range exactMatches {
		// Get or create the index entry
		indexEntry, ok := p.index[key]
		if !ok {
			indexEntry = &IndexEntry{
				E: make(map[string][]int64),
				C: make(map[string][]int64),
			}
			p.index[key] = indexEntry
		}

		// Initialize maps if nil
		if indexEntry.E == nil {
			indexEntry.E = make(map[string][]int64)
		}

		// Add to the exact match list
		exists := false
		for _, existingID := range indexEntry.E[dictType] {
			if existingID == idInt {
				exists = true
				break
			}
		}
		if !exists {
			indexEntry.E[dictType] = append(indexEntry.E[dictType], idInt)
		}
	}

	// Process contained-in matches
	for _, key := range containedMatches {
		// Get or create the index entry
		indexEntry, ok := p.index[key]
		if !ok {
			indexEntry = &IndexEntry{
				E: make(map[string][]int64),
				C: make(map[string][]int64),
			}
			p.index[key] = indexEntry
		}

		// Initialize maps if nil
		if indexEntry.C == nil {
			indexEntry.C = make(map[string][]int64)
		}

		// Add to the contained-in match list
		exists := false
		for _, existingID := range indexEntry.C[dictType] {
			if existingID == idInt {
				exists = true
				break
			}
		}
		if !exists {
			indexEntry.C[dictType] = append(indexEntry.C[dictType], idInt)
		}
	}

	p.mu.Unlock()

	// Write the entry to its dictionary file if not already written
	p.mu.Lock()
	alreadyWritten := p.writtenEntries[id]
	if !alreadyWritten {
		p.writtenEntries[id] = true
	}
	p.mu.Unlock()

	if !alreadyWritten {
		return p.writeEntryToFile(entry)
	}

	return nil
}

// writeEntryToFile writes an entry to its dictionary file
func (p *IndexProcessor) writeEntryToFile(entry common.Entry) error {
	var dir string

	// Get the original ID
	originalID := entry.GetID()

	// Get the shard type
	shardType := GetShardType(entry)

	// Create the sharded ID by prepending the shard type
	shardedID := fmt.Sprintf("%d%s", shardType, originalID)

	// Determine the directory based on entry type
	switch entry.(type) {
	case jmdict.Word:
		dir = p.jmdictDir
	case jmnedict.Name:
		dir = p.jmnedictDir
	case kanjidic.Kanji:
		dir = p.kanjidicDir
	case chinese_chars.ChineseCharEntry:
		dir = p.chineseCharsDir
	case chinese_words.ChineseWordEntry:
		dir = p.chineseWordsDir
	default:
		return fmt.Errorf("unknown entry type: %T", entry)
	}

	// Write the entry to a file
	filePath := filepath.Join(dir, shardedID+".json.br")
	return writeCompressedJSON(filePath, entry)
}

// WriteToFiles writes all index entries to files
func (p *IndexProcessor) WriteToFiles() error {
	// Count total files to write
	totalFiles := len(p.index)
	totalDictFiles := len(p.writtenEntries)

	fmt.Printf("Writing %d index files and %d dictionary files...\n", totalFiles, totalDictFiles)

	// Create channels for worker pool
	type job struct {
		key   string
		entry *IndexEntry
	}

	jobs := make(chan job, totalFiles)
	results := make(chan error, totalFiles)

	// Progress tracking
	var (
		mu        sync.Mutex
		completed int
	)

	// Progress reporting goroutine
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				mu.Lock()
				current := completed
				mu.Unlock()

				if current > 0 {
					fmt.Printf("\rWriting index files: %d/%d (%.1f%%)...",
						current, totalFiles, float64(current)/float64(totalFiles)*100)
				}
			case <-done:
				return
			}
		}
	}()

	// Create worker pool
	for w := 0; w < p.fileWriters; w++ {
		go func() {
			for j := range jobs {
				filename := filepath.Join(p.indexDir, j.key+".json.br")
				err := writeCompressedJSON(filename, j.entry)

				mu.Lock()
				completed++
				mu.Unlock()

				results <- err
			}
		}()
	}

	// Send jobs to workers
	for key, entry := range p.index {
		// Optimize the index entry before writing
		optimizeIndexEntry(entry)
		jobs <- job{key, entry}
	}
	close(jobs)

	// Collect results
	var errors []error
	for i := 0; i < totalFiles; i++ {
		err := <-results
		if err != nil {
			errors = append(errors, err)
		}
	}

	// Stop progress reporting
	close(done)
	fmt.Printf("\rWrote %d index files successfully\n", totalFiles)

	// Print statistics
	var jmdictExactCount, jmdictContainedCount int
	var jmnedictExactCount, jmnedictContainedCount int
	var kanjidicExactCount, kanjidicContainedCount int
	var chineseCharsExactCount, chineseCharsContainedCount int
	var chineseWordsExactCount, chineseWordsContainedCount int

	for _, entry := range p.index {
		// Count exact matches
		jmdictExactCount += len(entry.E["j"])
		jmnedictExactCount += len(entry.E["n"])
		kanjidicExactCount += len(entry.E["d"])
		chineseCharsExactCount += len(entry.E["c"])
		chineseWordsExactCount += len(entry.E["w"])

		// Count contained-in matches
		jmdictContainedCount += len(entry.C["j"])
		jmnedictContainedCount += len(entry.C["n"])
		kanjidicContainedCount += len(entry.C["d"])
		chineseCharsContainedCount += len(entry.C["c"])
		chineseWordsContainedCount += len(entry.C["w"])
	}

	// Calculate totals
	totalExactCount := jmdictExactCount + jmnedictExactCount + kanjidicExactCount +
		chineseCharsExactCount + chineseWordsExactCount
	totalContainedCount := jmdictContainedCount + jmnedictContainedCount + kanjidicContainedCount +
		chineseCharsContainedCount + chineseWordsContainedCount

	fmt.Printf("\nDictionary entry counts in index:\n")
	fmt.Printf("- JMdict exact matches: %d\n", jmdictExactCount)
	fmt.Printf("- JMdict contained-in matches: %d\n", jmdictContainedCount)
	fmt.Printf("- JMNedict exact matches: %d\n", jmnedictExactCount)
	fmt.Printf("- JMNedict contained-in matches: %d\n", jmnedictContainedCount)
	fmt.Printf("- Kanjidic exact matches: %d\n", kanjidicExactCount)
	fmt.Printf("- Kanjidic contained-in matches: %d\n", kanjidicContainedCount)
	fmt.Printf("- Chinese character exact matches: %d\n", chineseCharsExactCount)
	fmt.Printf("- Chinese character contained-in matches: %d\n", chineseCharsContainedCount)
	fmt.Printf("- Chinese word exact matches: %d\n", chineseWordsExactCount)
	fmt.Printf("- Chinese word contained-in matches: %d\n", chineseWordsContainedCount)
	fmt.Printf("- Total exact matches: %d\n", totalExactCount)
	fmt.Printf("- Total contained-in matches: %d\n", totalContainedCount)
	fmt.Printf("- Total unique dictionary entries: %d\n", totalDictFiles)

	// Check for errors
	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors while writing files", len(errors))
	}

	return nil
}

// optimizeIndexEntry optimizes an index entry to reduce size
// It removes empty arrays and ensures the entry is as small as possible
func optimizeIndexEntry(entry *IndexEntry) {
	// Remove empty dictionary types from exact matches
	if entry.E != nil {
		for dictType, ids := range entry.E {
			if len(ids) == 0 {
				delete(entry.E, dictType)
			}
		}
		if len(entry.E) == 0 {
			entry.E = nil
		}
	}

	// Remove empty dictionary types from contained-in matches
	if entry.C != nil {
		for dictType, ids := range entry.C {
			if len(ids) == 0 {
				delete(entry.C, dictType)
			}
		}
		if len(entry.C) == 0 {
			entry.C = nil
		}
	}
}

// writeCompressedJSON writes an object to a Brotli-compressed JSON file
func writeCompressedJSON(filename string, obj interface{}) error {
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a Brotli writer
	bw := brotli.NewWriter(file)
	defer bw.Close()

	// Encode the object as JSON
	encoder := json.NewEncoder(bw)
	err = encoder.Encode(obj)
	return err
}

// isHanCharacter returns true if the character is a CJK character
func isHanCharacter(char string) bool {
	if len(char) == 0 {
		return false
	}

	r := []rune(char)[0]

	// CJK Unified Ideographs
	if (r >= 0x4E00 && r <= 0x9FFF) ||
		// CJK Unified Ideographs Extension A
		(r >= 0x3400 && r <= 0x4DBF) ||
		// CJK Unified Ideographs Extension B
		(r >= 0x20000 && r <= 0x2A6DF) ||
		// CJK Unified Ideographs Extension C
		(r >= 0x2A700 && r <= 0x2B73F) ||
		// CJK Unified Ideographs Extension D
		(r >= 0x2B740 && r <= 0x2B81F) ||
		// CJK Unified Ideographs Extension E
		(r >= 0x2B820 && r <= 0x2CEAF) ||
		// CJK Unified Ideographs Extension F
		(r >= 0x2CEB0 && r <= 0x2EBEF) ||
		// CJK Compatibility Ideographs
		(r >= 0xF900 && r <= 0xFAFF) {
		return true
	}

	return false
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

// removeExactMatches removes exact matches from contained-in matches
func removeExactMatches(containedMatches, exactMatches []string) []string {
	// Create a map of exact matches for O(1) lookup
	exactMap := make(map[string]bool)
	for _, match := range exactMatches {
		exactMap[match] = true
	}

	// Filter out exact matches from contained-in matches
	var filtered []string
	for _, match := range containedMatches {
		if !exactMap[match] {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

// GetIndex returns a copy of the index map for debugging purposes
func (p *IndexProcessor) GetIndex() map[string]*IndexEntry {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Create a copy of the index map
	indexCopy := make(map[string]*IndexEntry, len(p.index))
	for key, entry := range p.index {
		indexCopy[key] = entry
	}

	return indexCopy
}
