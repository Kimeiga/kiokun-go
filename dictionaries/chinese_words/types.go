package chinese_words

import "fmt"

// ChineseWordEntry represents a single Chinese word entry
type ChineseWordEntry struct {
	ID          string         `json:"id"`
	Traditional string         `json:"traditional"`
	Simplified  string         `json:"simplified"`
	Pinyin      []string       `json:"pinyin,omitempty"`
	Definitions []string       `json:"definitions,omitempty"`
	HskLevel    int            `json:"hskLevel,omitempty"`
	Frequency   map[string]int `json:"frequency,omitempty"`
}

// GetID returns the entry ID
func (w ChineseWordEntry) GetID() string {
	return w.ID
}

// GetFilename returns the filename for this entry
func (w ChineseWordEntry) GetFilename() string {
	// Use traditional word as the filename
	return fmt.Sprintf("%s.json", w.Traditional)
}
