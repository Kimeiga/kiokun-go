package jmnedict

// JMNedictTypes represents the root structure of the JMnedict file
type JMNedictTypes struct {
	Version       string            `json:"version"`
	Languages     []string          `json:"languages"` // 3-letter codes like "eng"
	CommonOnly    bool              `json:"commonOnly"`
	DictDate      string            `json:"dictDate"`
	DictRevisions []string          `json:"dictRevisions"`
	Tags          map[string]string `json:"tags"`
	Names         []Name            `json:"names"`
}

// Name represents a single name entry in JMnedict
type Name struct {
	ID       string   `json:"id"`
	Kanji    []string `json:"k"`
	Reading  []string `json:"r"`
	Meanings []string `json:"m"`
	Type     []string `json:"type,omitempty"`
}

// GetID returns the unique identifier for this name
func (n Name) GetID() string {
	return n.ID
}

// GetFilename returns the filename to use for this name
func (n Name) GetFilename() string {
	var key string
	if len(n.Kanji) > 0 {
		key = n.Kanji[0]
	} else if len(n.Reading) > 0 {
		key = n.Reading[0]
	}
	return key
}

// JMnedict represents the root dictionary object
type JMnedict struct {
	Version       string            `json:"version"`
	Languages     []string          `json:"languages"` // 3-letter language codes
	DictDate      string            `json:"dictDate"`
	DictRevisions []string          `json:"dictRevisions"`
	Tags          map[string]string `json:"tags"`
	Words         []Word            `json:"words"`
}

// Word represents a dictionary entry
type Word struct {
	ID          string        `json:"id"`
	Kanji       []Kanji       `json:"kanji"`
	Kana        []Kana        `json:"kana"`
	Translation []Translation `json:"translation"`
}

// Kanji represents non-kana writings
type Kanji struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

// Kana represents kana-only writings
type Kana struct {
	Text           string   `json:"text"`
	Tags           []string `json:"tags"`
	AppliesToKanji []string `json:"appliesToKanji"`
}

// Translation represents translations and related data
type Translation struct {
	Type        []string             `json:"type"`
	Related     [][]string           `json:"related"`
	Translation []TranslationDetails `json:"translation"`
}

// TranslationDetails represents individual translation details
type TranslationDetails struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}
