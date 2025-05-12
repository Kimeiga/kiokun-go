package kanjidic

// KanjidicTypes represents the root structure of the Kanjidic file
type KanjidicTypes struct {
	Version       string            `json:"version"`
	Languages     []string          `json:"languages"` // 3-letter codes like "eng"
	CommonOnly    bool              `json:"commonOnly"`
	DictDate      string            `json:"dictDate"`
	DictRevisions []string          `json:"dictRevisions"`
	Tags          map[string]string `json:"tags"`
	Characters    []Kanji           `json:"characters"`
}

// Kanji represents a single kanji character entry in Kanjidic
type Kanji struct {
	Character string   `json:"c"`
	NumericID string   `json:"id,omitempty"` // Added numeric ID
	Meanings  []string `json:"m,omitempty"`
	OnYomi    []string `json:"on,omitempty"`
	KunYomi   []string `json:"kun,omitempty"`
	JLPT      int      `json:"jlpt,omitempty"`
	Grade     int      `json:"grade,omitempty"`
	Stroke    int      `json:"stroke"`
	Frequency int      `json:"freq,omitempty"`
	Radicals  []string `json:"rad,omitempty"`
	IDS       string   `json:"ids,omitempty"` // Ideographic Description Sequence
}

// GetID returns the unique identifier for this kanji
func (k Kanji) GetID() string {
	// If NumericID is set, use it; otherwise fall back to Character
	if k.NumericID != "" {
		return k.NumericID
	}
	return k.Character
}

// GetFilename returns the filename to use for this kanji
func (k Kanji) GetFilename() string {
	// If NumericID is set, use it; otherwise fall back to Character
	if k.NumericID != "" {
		return k.NumericID
	}
	return k.Character
}

// Kanjidic2 represents the root dictionary object
type Kanjidic2 struct {
	Version         string      `json:"version"`
	Languages       []string    `json:"languages"` // 2-letter language codes
	FileVersion     int         `json:"fileVersion"`
	DatabaseVersion string      `json:"databaseVersion"`
	Characters      []Character `json:"characters"`
}

// Character represents a kanji character entry
type Character struct {
	Literal              string                `json:"literal"`
	Codepoints           []Codepoint           `json:"codepoints"`
	Radicals             []Radical             `json:"radicals"`
	Misc                 Misc                  `json:"misc"`
	DictionaryReferences []DictionaryReference `json:"dictionaryReferences"`
	QueryCodes           []QueryCode           `json:"queryCodes"`
	ReadingMeaning       *ReadingMeaning       `json:"readingMeaning"`
}

// Codepoint represents character encodings
type Codepoint struct {
	Type  string `json:"type"` // jis208, jis212, jis213, ucs
	Value string `json:"value"`
}

// Radical represents radical information
type Radical struct {
	Type  string `json:"type"` // classical, nelson_c
	Value int    `json:"value"`
}

// Misc represents miscellaneous kanji information
type Misc struct {
	Grade        *int      `json:"grade"`
	StrokeCounts []int     `json:"strokeCounts"`
	Variants     []Variant `json:"variants"`
	Frequency    *int      `json:"frequency"`
	RadicalNames []string  `json:"radicalNames"`
	JlptLevel    *int      `json:"jlptLevel"`
}

// Variant represents kanji variants
type Variant struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DictionaryReference represents references to other dictionaries
type DictionaryReference struct {
	Type      string     `json:"type"`
	Morohashi *Morohashi `json:"morohashi"`
	Value     string     `json:"value"`
}

// Morohashi represents Morohashi dictionary specific information
type Morohashi struct {
	Volume int `json:"volume"`
	Page   int `json:"page"`
}

// QueryCode represents various dictionary query codes
type QueryCode struct {
	Type                  string  `json:"type"`
	SkipMisclassification *string `json:"skipMisclassification,omitempty"`
	Value                 string  `json:"value"`
}

// ReadingMeaning represents readings and meanings of kanji
type ReadingMeaning struct {
	Groups []ReadingMeaningGroup `json:"groups"`
	Nanori []string              `json:"nanori"`
}

// ReadingMeaningGroup represents a group of readings and meanings
type ReadingMeaningGroup struct {
	Readings []Reading `json:"readings"`
	Meanings []Meaning `json:"meanings"`
}

// Reading represents a kanji reading
type Reading struct {
	Type   string  `json:"type"`
	OnType *string `json:"onType"`
	Status *string `json:"status"`
	Value  string  `json:"value"`
}

// Meaning represents a kanji meaning
type Meaning struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}
