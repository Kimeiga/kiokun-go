package chinese_chars

import "fmt"

// ChineseCharEntry represents a single Chinese character entry
type ChineseCharEntry struct {
	ID          string   `json:"id"`
	Traditional string   `json:"traditional"`
	Simplified  string   `json:"simplified"`
	Definitions []string `json:"definitions,omitempty"`
	Pinyin      []string `json:"pinyin,omitempty"`
	StrokeCount int      `json:"strokeCount,omitempty"`
}

// GetID returns the entry ID
func (c ChineseCharEntry) GetID() string {
	return c.ID
}

// GetFilename returns the filename for this entry
func (c ChineseCharEntry) GetFilename() string {
	// Use traditional character as the filename
	return fmt.Sprintf("%s.json", c.Traditional)
}
