package chinese_chars

// ChineseCharEntry represents a single Chinese character entry
type ChineseCharEntry struct {
	ID          string   `json:"id"`
	Traditional string   `json:"traditional"`
	Simplified  string   `json:"simplified"`
	Definitions []string `json:"definitions,omitempty"`
	Pinyin      []string `json:"pinyin,omitempty"`
	StrokeCount int      `json:"strokeCount,omitempty"`
	IDS         string   `json:"ids,omitempty"` // Ideographic Description Sequence
}

// GetID returns the entry ID
func (c ChineseCharEntry) GetID() string {
	return c.ID
}

// GetFilename returns the filename for this entry
func (c ChineseCharEntry) GetFilename() string {
	// Use ID as the filename
	return c.ID
}
