package processor

import (
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
)

// ShardedIndexProcessor processes dictionary entries and builds sharded indexes
type ShardedIndexProcessor struct {
	baseDir          string
	shardDirs        map[ShardType]string
	indexDirs        map[ShardType]string
	jmdictDirs       map[ShardType]string
	jmnedictDirs     map[ShardType]string
	kanjidicDirs     map[ShardType]string
	chineseCharsDirs map[ShardType]string
	chineseWordsDirs map[ShardType]string
	indexes          map[ShardType]map[string]*IndexEntry
	writtenEntries   map[ShardType]map[string]bool
	fileWriters      int
	idsMap           map[string]string // Map of character to IDS
	mu               sync.Mutex
}

// NewShardedIndexProcessor creates a new sharded index processor
func NewShardedIndexProcessor(baseDir string, fileWriters int) (*ShardedIndexProcessor, error) {
	// Create the processor
	p := &ShardedIndexProcessor{
		baseDir:          baseDir,
		shardDirs:        make(map[ShardType]string),
		indexDirs:        make(map[ShardType]string),
		jmdictDirs:       make(map[ShardType]string),
		jmnedictDirs:     make(map[ShardType]string),
		kanjidicDirs:     make(map[ShardType]string),
		chineseCharsDirs: make(map[ShardType]string),
		chineseWordsDirs: make(map[ShardType]string),
		indexes:          make(map[ShardType]map[string]*IndexEntry),
		writtenEntries:   make(map[ShardType]map[string]bool),
		fileWriters:      fileWriters,
		idsMap:           make(map[string]string),
	}

	// Initialize indexes and writtenEntries for each shard
	for _, shardType := range []ShardType{ShardNonHan, ShardHan1Char, ShardHan2Char, ShardHan3Plus} {
		p.indexes[shardType] = make(map[string]*IndexEntry)
		p.writtenEntries[shardType] = make(map[string]bool)
	}

	// Create the output directories
	if err := p.createDirectories(); err != nil {
		return nil, err
	}

	return p, nil
}

// SetIDSMap sets the IDS map for the processor
func (p *ShardedIndexProcessor) SetIDSMap(idsMap map[string]string) {
	p.idsMap = idsMap
}

// createDirectories creates the necessary output directories for each shard
func (p *ShardedIndexProcessor) createDirectories() error {
	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(p.baseDir, 0755); err != nil {
		return err
	}

	// Create directories for each shard
	for _, shardType := range []ShardType{ShardNonHan, ShardHan1Char, ShardHan2Char, ShardHan3Plus} {
		// Create the shard directory
		shardDir := GetOutputDirForShard(p.baseDir, shardType)
		p.shardDirs[shardType] = shardDir
		if err := os.MkdirAll(shardDir, 0755); err != nil {
			return err
		}

		// Create subdirectories with one-letter names
		indexDir := filepath.Join(shardDir, "index")
		jmdictDir := filepath.Join(shardDir, "j")       // j for JMdict
		jmnedictDir := filepath.Join(shardDir, "n")     // n for JMNedict
		kanjidicDir := filepath.Join(shardDir, "d")     // d for Kanjidic
		chineseCharsDir := filepath.Join(shardDir, "c") // c for Chinese Characters
		chineseWordsDir := filepath.Join(shardDir, "w") // w for Chinese Words

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
		p.indexDirs[shardType] = indexDir
		p.jmdictDirs[shardType] = jmdictDir
		p.jmnedictDirs[shardType] = jmnedictDir
		p.kanjidicDirs[shardType] = kanjidicDir
		p.chineseCharsDirs[shardType] = chineseCharsDir
		p.chineseWordsDirs[shardType] = chineseWordsDir
	}

	return nil
}

// ProcessEntries processes a slice of entries
func (p *ShardedIndexProcessor) ProcessEntries(entries []common.Entry) error {
	for _, entry := range entries {
		if err := p.processEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

// processEntry processes a single entry
func (p *ShardedIndexProcessor) processEntry(entry common.Entry) error {
	// Get the entry ID
	originalID := entry.GetID()

	// Get the shard type
	shardType := GetShardType(entry)

	// Special logging for "日" character (check by ID since type assertion might not work)
	if originalID == "4057102" {
		fmt.Printf("🌞 PROCESSOR: Processing '日' entry - originalID: %s, shardType: %d\n", originalID, shardType)
	}

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
		indexEntry, ok := p.indexes[shardType][key]
		if !ok {
			indexEntry = &IndexEntry{
				E: make(map[string][]int64),
				C: make(map[string][]int64),
			}
			p.indexes[shardType][key] = indexEntry

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

		} else if key == "日本" {
			// Debug: Print when ID already exists in exact match list for "日本"
			fmt.Printf("DEBUG: ID %d already exists in exact match list for key '日本' with dictionary type '%s'\n", idInt, dictType)
		}
	}

	// Process contained-in matches
	for _, key := range containedMatches {
		// Get or create the index entry
		indexEntry, ok := p.indexes[shardType][key]
		if !ok {
			indexEntry = &IndexEntry{
				E: make(map[string][]int64),
				C: make(map[string][]int64),
			}
			p.indexes[shardType][key] = indexEntry
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

	// Add IDS data to single Han character entries
	var updatedEntry common.Entry
	switch e := entry.(type) {
	case kanjidic.Kanji:
		// For Kanjidic entries, add IDS data if available
		if ids, ok := p.idsMap[e.Character]; ok {
			// Create a copy of the entry with IDS data
			entryCopy := e
			entryCopy.IDS = ids
			updatedEntry = entryCopy
		}
	case chinese_chars.ChineseCharEntry:
		// For Chinese character entries, add IDS data if available
		if ids, ok := p.idsMap[e.Traditional]; ok {
			// Create a copy of the entry with IDS data
			entryCopy := e
			entryCopy.IDS = ids
			updatedEntry = entryCopy
		}
	}

	// Use the updated entry if available
	if updatedEntry != nil {
		entry = updatedEntry
	}

	// Write the entry to its dictionary file if not already written
	p.mu.Lock()
	alreadyWritten := p.writtenEntries[shardType][id]
	if !alreadyWritten {
		p.writtenEntries[shardType][id] = true
	}
	p.mu.Unlock()

	if !alreadyWritten {
		return p.writeEntryToFile(entry, shardType)
	}

	return nil
}

// writeEntryToFile writes an entry to its dictionary file in the appropriate shard
func (p *ShardedIndexProcessor) writeEntryToFile(entry common.Entry, shardType ShardType) error {
	var dir string

	// Get the original ID
	originalID := entry.GetID()

	// Create the sharded ID by prepending the shard type
	shardedID := fmt.Sprintf("%d%s", shardType, originalID)

	// Special logging for "日" character (check by ID since type assertion might not work)
	if originalID == "4057102" {
		fmt.Printf("🌞 WRITE_FILE: Writing '日' entry - originalID: %s, shardedID: %s, shardType: %d\n", originalID, shardedID, shardType)
	}

	// Determine the directory based on entry type
	switch entry.(type) {
	case jmdict.Word:
		dir = p.jmdictDirs[shardType]
	case jmnedict.Name:
		dir = p.jmnedictDirs[shardType]
	case kanjidic.Kanji:
		dir = p.kanjidicDirs[shardType]
	case chinese_chars.ChineseCharEntry:
		dir = p.chineseCharsDirs[shardType]
	case chinese_words.ChineseWordEntry:
		dir = p.chineseWordsDirs[shardType]
	default:
		return fmt.Errorf("unknown entry type: %T", entry)
	}

	// Write the entry to a file
	filePath := filepath.Join(dir, shardedID+".json.br")

	// Special logging for "日" character (check by ID since type assertion might not work)
	if originalID == "4057102" {
		fmt.Printf("🌞 FINAL_FILE: Writing '日' entry to file: %s\n", filePath)
	}

	return writeCompressedJSON(filePath, entry)
}

// WriteToFiles writes all index entries to files for each shard
func (p *ShardedIndexProcessor) WriteToFiles() error {
	// Count total files to write across all shards
	totalFiles := 0
	totalDictFiles := 0
	for shardType := range p.indexes {
		totalFiles += len(p.indexes[shardType])
		totalDictFiles += len(p.writtenEntries[shardType])
	}

	fmt.Printf("Writing %d index files and %d dictionary files across all shards...\n", totalFiles, totalDictFiles)

	// Process each shard
	for shardType, index := range p.indexes {
		shardFiles := len(index)
		shardDictFiles := len(p.writtenEntries[shardType])

		fmt.Printf("Shard %d: Writing %d index files and %d dictionary files...\n",
			shardType, shardFiles, shardDictFiles)

		// Create channels for worker pool
		type job struct {
			key   string
			entry *IndexEntry
		}

		jobs := make(chan job, shardFiles)
		results := make(chan error, shardFiles)

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
						fmt.Printf("\rShard %d: Writing index files: %d/%d (%.1f%%)...",
							shardType, current, shardFiles, float64(current)/float64(shardFiles)*100)
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
					filename := filepath.Join(p.indexDirs[shardType], j.key+".json.br")
					err := writeCompressedJSON(filename, j.entry)

					mu.Lock()
					completed++
					mu.Unlock()

					results <- err
				}
			}()
		}

		// Send jobs to workers
		for key, entry := range index {

			// Optimize the index entry before writing
			optimizeIndexEntry(entry)
			jobs <- job{key, entry}
		}
		close(jobs)

		// Collect results
		var errors []error
		for i := 0; i < shardFiles; i++ {
			err := <-results
			if err != nil {
				errors = append(errors, err)
			}
		}

		// Stop progress reporting
		close(done)
		fmt.Printf("\rShard %d: Wrote %d index files successfully\n", shardType, shardFiles)
	}

	// Print statistics for all shards
	p.printStatistics()

	return nil
}

// printStatistics prints statistics about the index
func (p *ShardedIndexProcessor) printStatistics() {
	// Aggregate statistics across all shards
	var totalJmdictExactCount, totalJmdictContainedCount int
	var totalJmnedictExactCount, totalJmnedictContainedCount int
	var totalKanjidicExactCount, totalKanjidicContainedCount int
	var totalChineseCharsExactCount, totalChineseCharsContainedCount int
	var totalChineseWordsExactCount, totalChineseWordsContainedCount int
	var totalUniqueEntries int

	// Process each shard
	for shardType, index := range p.indexes {
		var jmdictExactCount, jmdictContainedCount int
		var jmnedictExactCount, jmnedictContainedCount int
		var kanjidicExactCount, kanjidicContainedCount int
		var chineseCharsExactCount, chineseCharsContainedCount int
		var chineseWordsExactCount, chineseWordsContainedCount int

		for _, entry := range index {
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

		// Print statistics for this shard
		fmt.Printf("\nShard %d statistics:\n", shardType)
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
		fmt.Printf("- Total exact matches: %d\n", jmdictExactCount+jmnedictExactCount+kanjidicExactCount+chineseCharsExactCount+chineseWordsExactCount)
		fmt.Printf("- Total contained-in matches: %d\n", jmdictContainedCount+jmnedictContainedCount+kanjidicContainedCount+chineseCharsContainedCount+chineseWordsContainedCount)
		fmt.Printf("- Total unique dictionary entries: %d\n", len(p.writtenEntries[shardType]))

		// Accumulate totals
		totalJmdictExactCount += jmdictExactCount
		totalJmdictContainedCount += jmdictContainedCount
		totalJmnedictExactCount += jmnedictExactCount
		totalJmnedictContainedCount += jmnedictContainedCount
		totalKanjidicExactCount += kanjidicExactCount
		totalKanjidicContainedCount += kanjidicContainedCount
		totalChineseCharsExactCount += chineseCharsExactCount
		totalChineseCharsContainedCount += chineseCharsContainedCount
		totalChineseWordsExactCount += chineseWordsExactCount
		totalChineseWordsContainedCount += chineseWordsContainedCount
		totalUniqueEntries += len(p.writtenEntries[shardType])
	}

	// Print overall statistics
	fmt.Printf("\nOverall statistics across all shards:\n")
	fmt.Printf("- JMdict exact matches: %d\n", totalJmdictExactCount)
	fmt.Printf("- JMdict contained-in matches: %d\n", totalJmdictContainedCount)
	fmt.Printf("- JMNedict exact matches: %d\n", totalJmnedictExactCount)
	fmt.Printf("- JMNedict contained-in matches: %d\n", totalJmnedictContainedCount)
	fmt.Printf("- Kanjidic exact matches: %d\n", totalKanjidicExactCount)
	fmt.Printf("- Kanjidic contained-in matches: %d\n", totalKanjidicContainedCount)
	fmt.Printf("- Chinese character exact matches: %d\n", totalChineseCharsExactCount)
	fmt.Printf("- Chinese character contained-in matches: %d\n", totalChineseCharsContainedCount)
	fmt.Printf("- Chinese word exact matches: %d\n", totalChineseWordsExactCount)
	fmt.Printf("- Chinese word contained-in matches: %d\n", totalChineseWordsContainedCount)
	fmt.Printf("- Total exact matches: %d\n", totalJmdictExactCount+totalJmnedictExactCount+totalKanjidicExactCount+totalChineseCharsExactCount+totalChineseWordsExactCount)
	fmt.Printf("- Total contained-in matches: %d\n", totalJmdictContainedCount+totalJmnedictContainedCount+totalKanjidicContainedCount+totalChineseCharsContainedCount+totalChineseWordsContainedCount)
	fmt.Printf("- Total unique dictionary entries: %d\n", totalUniqueEntries)
}
