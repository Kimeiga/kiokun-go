package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"

	"github.com/andybalholm/brotli"
)

// WordGroup represents the combined data for a single word/character
type WordGroup struct {
	WordJapanese []jmdict.Word                    `json:"w_j,omitempty"`
	NameJapanese []jmnedict.Name                  `json:"n_j,omitempty"`
	CharJapanese []kanjidic.Kanji                 `json:"c_j,omitempty"`
	CharChinese  []chinese_chars.ChineseCharEntry `json:"c_c,omitempty"`
	WordChinese  []chinese_words.ChineseWordEntry `json:"w_c,omitempty"`
}

// HasMultipleDictData returns true if this group has data from multiple dictionaries
func (wg *WordGroup) HasMultipleDictData() bool {
	sources := 0
	if len(wg.WordJapanese) > 0 {
		sources++
	}
	if len(wg.NameJapanese) > 0 {
		sources++
	}
	if len(wg.CharJapanese) > 0 {
		sources++
	}
	if len(wg.CharChinese) > 0 {
		sources++
	}
	if len(wg.WordChinese) > 0 {
		sources++
	}
	return sources > 1
}

// DictionaryProcessor handles combining and writing dictionary entries
type DictionaryProcessor struct {
	outputDir   string
	groups      map[string]*WordGroup // Cache of word groups
	groupsMu    sync.RWMutex          // Mutex to protect the groups map
	workerCount int                   // Number of worker goroutines for file writing
}

// New creates a new processor instance
func New(outputDir string, workerCount int) (*DictionaryProcessor, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}

	// Default to 1 worker (sequential) if an invalid count is provided
	if workerCount < 1 {
		workerCount = 1
	}

	return &DictionaryProcessor{
		outputDir:   outputDir,
		groups:      make(map[string]*WordGroup),
		workerCount: workerCount,
	}, nil
}

// ProcessEntries processes a batch of entries
func (p *DictionaryProcessor) ProcessEntries(entries []common.Entry) error {
	for _, entry := range entries {
		if err := p.processEntry(entry); err != nil {
			return fmt.Errorf("error processing entry: %v", err)
		}
	}
	return nil
}

// processEntry handles a single entry
func (p *DictionaryProcessor) processEntry(entry common.Entry) error {
	switch e := entry.(type) {
	case jmdict.Word:
		return p.processJMdictWord(e)
	case jmnedict.Name:
		return p.processJMNedictName(e)
	case kanjidic.Kanji:
		return p.processKanjidicEntry(e)
	case chinese_chars.ChineseCharEntry:
		return p.processChineseCharEntry(e)
	case chinese_words.ChineseWordEntry:
		return p.processChineseWordEntry(e)
	default:
		return fmt.Errorf("unknown entry type: %T", entry)
	}
}

// processJMdictWord processes a JMdict word entry
func (p *DictionaryProcessor) processJMdictWord(word jmdict.Word) error {
	// Get all forms (kanji and kana)
	var forms []string
	for _, k := range word.Kanji {
		forms = append(forms, k.Text)
	}
	for _, k := range word.Kana {
		forms = append(forms, k.Text)
	}

	// Add word to each form's group
	for _, form := range forms {
		group := p.getOrCreateGroup(form)
		group.WordJapanese = append(group.WordJapanese, word)
	}

	return nil
}

// processJMNedictName processes a JMNedict name entry
func (p *DictionaryProcessor) processJMNedictName(name jmnedict.Name) error {
	// Get all forms (kanji and kana)
	var forms []string
	forms = append(forms, name.Kanji...)
	forms = append(forms, name.Reading...)

	// Add name to each form's group
	for _, form := range forms {
		group := p.getOrCreateGroup(form)
		group.NameJapanese = append(group.NameJapanese, name)
	}

	return nil
}

// processKanjidicEntry processes a Kanjidic entry
func (p *DictionaryProcessor) processKanjidicEntry(kanji kanjidic.Kanji) error {
	group := p.getOrCreateGroup(kanji.Character)
	group.CharJapanese = append(group.CharJapanese, kanji)
	return nil
}

// processChineseCharEntry processes a Chinese character entry
func (p *DictionaryProcessor) processChineseCharEntry(char chinese_chars.ChineseCharEntry) error {
	// Add both traditional and simplified forms
	forms := []string{char.Traditional}
	if char.Simplified != char.Traditional {
		forms = append(forms, char.Simplified)
	}

	// Add character to each form's group
	for _, form := range forms {
		group := p.getOrCreateGroup(form)
		group.CharChinese = append(group.CharChinese, char)
	}

	return nil
}

// processChineseWordEntry processes a Chinese word entry
func (p *DictionaryProcessor) processChineseWordEntry(word chinese_words.ChineseWordEntry) error {
	// Add both traditional and simplified forms
	forms := []string{word.Traditional}
	if word.Simplified != word.Traditional {
		forms = append(forms, word.Simplified)
	}

	// Add word to each form's group
	for _, form := range forms {
		group := p.getOrCreateGroup(form)
		group.WordChinese = append(group.WordChinese, word)
	}

	return nil
}

// getOrCreateGroup gets or creates a WordGroup for a given form
func (p *DictionaryProcessor) getOrCreateGroup(form string) *WordGroup {
	// First try read-only access
	p.groupsMu.RLock()
	group, exists := p.groups[form]
	p.groupsMu.RUnlock()

	if exists {
		return group
	}

	// If not found, acquire write lock and check again
	p.groupsMu.Lock()
	defer p.groupsMu.Unlock()

	// Check again after acquiring the write lock
	group, exists = p.groups[form]
	if !exists {
		group = &WordGroup{}
		p.groups[form] = group
	}
	return group
}

// sanitizeWordGroup removes duplicate entries from a word group
func sanitizeWordGroup(group *WordGroup) {
	// Deduplicate JMdict entries
	seen := make(map[string]bool)
	var uniqueWords []jmdict.Word
	for _, word := range group.WordJapanese {
		if !seen[word.ID] {
			seen[word.ID] = true
			uniqueWords = append(uniqueWords, word)
		}
	}
	group.WordJapanese = uniqueWords

	// Deduplicate JMNedict entries
	seen = make(map[string]bool)
	var uniqueNames []jmnedict.Name
	for _, name := range group.NameJapanese {
		if !seen[name.ID] {
			seen[name.ID] = true
			uniqueNames = append(uniqueNames, name)
		}
	}
	group.NameJapanese = uniqueNames

	// Deduplicate Kanjidic entries
	seen = make(map[string]bool)
	var uniqueKanji []kanjidic.Kanji
	for _, kanji := range group.CharJapanese {
		if !seen[kanji.Character] {
			seen[kanji.Character] = true
			uniqueKanji = append(uniqueKanji, kanji)
		}
	}
	group.CharJapanese = uniqueKanji

	// Deduplicate Chinese character entries
	seen = make(map[string]bool)
	var uniqueChineseChars []chinese_chars.ChineseCharEntry
	for _, char := range group.CharChinese {
		if !seen[char.ID] {
			seen[char.ID] = true
			uniqueChineseChars = append(uniqueChineseChars, char)
		}
	}
	group.CharChinese = uniqueChineseChars

	// Deduplicate Chinese word entries
	seen = make(map[string]bool)
	var uniqueChineseWords []chinese_words.ChineseWordEntry
	for _, word := range group.WordChinese {
		if !seen[word.ID] {
			seen[word.ID] = true
			uniqueChineseWords = append(uniqueChineseWords, word)
		}
	}
	group.WordChinese = uniqueChineseWords
}

// WriteToFiles writes all groups to their respective files
func (p *DictionaryProcessor) WriteToFiles() error {
	// Lock the map while we're iterating over it
	p.groupsMu.RLock()

	// Create a copy of groups to avoid holding the lock during processing
	groupsCopy := make(map[string]*WordGroup, len(p.groups))
	for form, group := range p.groups {
		groupsCopy[form] = group
	}
	// Now we can release the read lock
	p.groupsMu.RUnlock()

	// Statistics variables
	jmdictTotal := 0
	jmnedictTotal := 0
	kanjidicTotal := 0
	chineseCharsTotal := 0
	chineseWordsTotal := 0
	combinedEntries := 0
	totalFiles := len(groupsCopy)
	processed := 0
	errorCount := 0

	// Track entries from multiple dictionaries for verification
	var multiDictEntries []string
	var wordAndKanjiEntries []string
	var wordKanjiAndNameEntries []string

	// Sequential file writing
	if p.workerCount <= 1 {
		for form, group := range groupsCopy {
			if err := p.writeGroupToFile(form, group); err != nil {
				return fmt.Errorf("error writing %s: %v", form, err)
			}

			processed++
			if processed%1000 == 0 || processed == totalFiles {
				fmt.Printf("\rWriting files: %d/%d (%.1f%%)...", processed, totalFiles, float64(processed)/float64(totalFiles)*100)
			}

			// Record entries that have data from multiple dictionaries
			hasWord := len(group.WordJapanese) > 0
			hasKanji := len(group.CharJapanese) > 0
			hasName := len(group.NameJapanese) > 0
			hasChineseChar := len(group.CharChinese) > 0
			hasChineseWord := len(group.WordChinese) > 0

			if hasWord {
				jmdictTotal++
			}
			if hasName {
				jmnedictTotal++
			}
			if hasKanji {
				kanjidicTotal++
			}
			if hasChineseChar {
				chineseCharsTotal++
			}
			if hasChineseWord {
				chineseWordsTotal++
			}

			if hasWord && hasKanji && hasName {
				wordKanjiAndNameEntries = append(wordKanjiAndNameEntries, form)
			} else if hasWord && hasKanji {
				wordAndKanjiEntries = append(wordAndKanjiEntries, form)
			}

			// Check for any combination of multiple dict data
			dictCount := 0
			if hasWord {
				dictCount++
			}
			if hasKanji {
				dictCount++
			}
			if hasName {
				dictCount++
			}
			if hasChineseChar {
				dictCount++
			}
			if hasChineseWord {
				dictCount++
			}

			if dictCount > 1 {
				multiDictEntries = append(multiDictEntries, form)
				combinedEntries++
			}
		}

		fmt.Printf("\nWrote %d files successfully\n", totalFiles-errorCount)
		fmt.Println("\nDictionary entry counts:")
		fmt.Printf("- JMdict entries: %d\n", jmdictTotal)
		fmt.Printf("- JMNedict entries: %d\n", jmnedictTotal)
		fmt.Printf("- Kanjidic entries: %d\n", kanjidicTotal)
		fmt.Printf("- Chinese character entries: %d\n", chineseCharsTotal)
		fmt.Printf("- Chinese word entries: %d\n", chineseWordsTotal)
		fmt.Printf("- Combined entries: %d\n", combinedEntries)

		return nil
	}

	// Parallel file writing with worker pool
	type writeTask struct {
		form  string
		group *WordGroup
	}

	type writeResult struct {
		form           string
		err            error
		hasWord        bool
		hasKanji       bool
		hasName        bool
		hasChineseChar bool
		hasChineseWord bool
	}

	taskChan := make(chan writeTask, min(1000, len(groupsCopy)))
	resultChan := make(chan writeResult, min(1000, len(groupsCopy)))

	// Launch worker pool
	var wg sync.WaitGroup
	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				err := p.writeGroupToFile(task.form, task.group)

				hasWord := len(task.group.WordJapanese) > 0
				hasKanji := len(task.group.CharJapanese) > 0
				hasName := len(task.group.NameJapanese) > 0
				hasChineseChar := len(task.group.CharChinese) > 0
				hasChineseWord := len(task.group.WordChinese) > 0

				resultChan <- writeResult{
					form:           task.form,
					err:            err,
					hasWord:        hasWord,
					hasKanji:       hasKanji,
					hasName:        hasName,
					hasChineseChar: hasChineseChar,
					hasChineseWord: hasChineseWord,
				}
			}
		}()
	}

	// Send all tasks
	go func() {
		for form, group := range groupsCopy {
			taskChan <- writeTask{form, group}
		}
		close(taskChan)

		// Wait for all workers to finish, then close the result channel
		wg.Wait()
		close(resultChan)
	}()

	// Collect results and statistics
	for result := range resultChan {
		processed++
		if processed%1000 == 0 || processed == totalFiles {
			fmt.Printf("\rWriting files: %d/%d (%.1f%%)...", processed, totalFiles, float64(processed)/float64(totalFiles)*100)
		}

		if result.err != nil {
			fmt.Printf("\nError writing %s: %v\n", result.form, result.err)
			errorCount++
			continue
		}

		// Count entries by dictionary type
		if result.hasWord {
			jmdictTotal++
		}
		if result.hasName {
			jmnedictTotal++
		}
		if result.hasKanji {
			kanjidicTotal++
		}
		if result.hasChineseChar {
			chineseCharsTotal++
		}
		if result.hasChineseWord {
			chineseWordsTotal++
		}

		// Record entries with data from multiple dictionaries
		dictCount := 0
		if result.hasWord {
			dictCount++
		}
		if result.hasName {
			dictCount++
		}
		if result.hasKanji {
			dictCount++
		}
		if result.hasChineseChar {
			dictCount++
		}
		if result.hasChineseWord {
			dictCount++
		}

		if dictCount > 1 {
			combinedEntries++
			multiDictEntries = append(multiDictEntries, result.form)
		}
	}

	fmt.Printf("\nWrote %d files successfully\n", totalFiles-errorCount)
	fmt.Println("\nDictionary entry counts:")
	fmt.Printf("- JMdict entries: %d\n", jmdictTotal)
	fmt.Printf("- JMNedict entries: %d\n", jmnedictTotal)
	fmt.Printf("- Kanjidic entries: %d\n", kanjidicTotal)
	fmt.Printf("- Chinese character entries: %d\n", chineseCharsTotal)
	fmt.Printf("- Chinese word entries: %d\n", chineseWordsTotal)
	fmt.Printf("- Combined entries: %d\n", combinedEntries)
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// writeGroupToFile writes a single group to its file
func (p *DictionaryProcessor) writeGroupToFile(form string, group *WordGroup) error {
	// Sanitize data before writing to file
	sanitizeWordGroup(group)

	filePath := filepath.Join(p.outputDir, form+".json.br")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	brWriter := brotli.NewWriter(file)
	defer brWriter.Close()

	encoder := json.NewEncoder(brWriter)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(group); err != nil {
		return fmt.Errorf("failed to encode group: %v", err)
	}

	return nil
}
