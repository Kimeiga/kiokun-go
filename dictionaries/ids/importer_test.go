package ids

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImporter_Import(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test_ids.txt")

	// Create test data
	testData := `U+4E00	一	⿻一丨
U+4E01	丁	⿻一亅
U+4E02	丂	⿱丄㇉
U+4E03	七	⿻一𠃌
U+4E04	丄	⿱一一	@apparent=⿱一丨
;; This is a comment
`
	err := os.WriteFile(testFilePath, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create an importer
	importer := &Importer{}

	// Import the test file
	entries, err := importer.Import(testFilePath)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// Check the number of entries
	expectedEntries := 5
	if len(entries) != expectedEntries {
		t.Errorf("Expected %d entries, got %d", expectedEntries, len(entries))
	}

	// Check the first entry
	if len(entries) > 0 {
		entry, ok := entries[0].(IDSEntry)
		if !ok {
			t.Fatalf("Entry is not an IDSEntry")
		}

		if entry.ID != "U+4E00" {
			t.Errorf("Expected ID 'U+4E00', got '%s'", entry.ID)
		}

		if entry.Character != "一" {
			t.Errorf("Expected Character '一', got '%s'", entry.Character)
		}

		if entry.IDS != "⿻一丨" {
			t.Errorf("Expected IDS '⿻一丨', got '%s'", entry.IDS)
		}
	}

	// Check the entry with apparent IDS
	if len(entries) > 4 {
		entry, ok := entries[4].(IDSEntry)
		if !ok {
			t.Fatalf("Entry is not an IDSEntry")
		}

		if entry.ID != "U+4E04" {
			t.Errorf("Expected ID 'U+4E04', got '%s'", entry.ID)
		}

		if entry.ApparentIDS != "⿱一丨" {
			t.Errorf("Expected ApparentIDS '⿱一丨', got '%s'", entry.ApparentIDS)
		}
	}
}
