package ids

// IDSEntry represents a single IDS (Ideographic Description Sequence) entry
type IDSEntry struct {
	ID          string `json:"id"`          // Codepoint (e.g., U+4E00)
	Character   string `json:"character"`   // The character itself
	IDS         string `json:"ids"`         // Ideographic Description Sequence
	ApparentIDS string `json:"apparentIds"` // Apparent IDS (optional)
}

// GetID returns the entry ID
func (e IDSEntry) GetID() string {
	return e.ID
}

// GetFilename returns the filename for this entry
func (e IDSEntry) GetFilename() string {
	// Use ID as the filename
	return e.ID
}
