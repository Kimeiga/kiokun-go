package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DictionaryFile represents a dictionary file to download and setup
type DictionaryFile struct {
	Name        string
	URL         string
	TargetPath  string
	IsZip       bool
	IsJSONL     bool
	JSONLTarget string // Target JSON file path if converting from JSONL
}

// SetupDictionaryFiles downloads and sets up all required dictionary files if they don't exist
func SetupDictionaryFiles(logf LogFunc) error {
	logf("Checking dictionary files...\n")

	files := []DictionaryFile{
		// JMdict
		{
			Name:       "JMdict",
			URL:        "https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/jmdict-examples-eng-3.6.1+20250505122413.json.zip",
			TargetPath: "dictionaries/jmdict/source/jmdict-examples-eng-3.6.1+20250505122413.json",
			IsZip:      true,
		},
		// JMNedict
		{
			Name:       "JMNedict",
			URL:        "https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/jmnedict-all-3.6.1+20250505122413.json.zip",
			TargetPath: "dictionaries/jmnedict/source/jmnedict-all-3.6.1+20250505122413.json",
			IsZip:      true,
		},
		// Kanjidic
		{
			Name:       "Kanjidic",
			URL:        "https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/kanjidic2-en-3.6.1+20250505122413.json.zip",
			TargetPath: "dictionaries/kanjidic/source/kanjidic2-en-3.6.1+20250505122413.json",
			IsZip:      true,
		},
		// Chinese Characters
		{
			Name:        "Chinese Characters",
			URL:         "https://data.dong-chinese.com/dump/dictionary_char_2024-06-17.jsonl",
			TargetPath:  "dictionaries/chinese_chars/source/dictionary_char_2024-06-17.jsonl",
			IsZip:       false,
			IsJSONL:     true,
			JSONLTarget: "dictionaries/chinese_chars/source/dictionary_char_2024-06-17.json",
		},
		// Chinese Words
		{
			Name:        "Chinese Words",
			URL:         "https://data.dong-chinese.com/dump/dictionary_word_2024-06-17.jsonl",
			TargetPath:  "dictionaries/chinese_words/source/dictionary_word_2024-06-17.jsonl",
			IsZip:       false,
			IsJSONL:     true,
			JSONLTarget: "dictionaries/chinese_words/source/dictionary_word_2024-06-17.json",
		},
		// IDS
		{
			Name:       "IDS",
			URL:        "https://raw.githubusercontent.com/cjkvi/cjkvi-ids/master/ids.txt",
			TargetPath: "dictionaries/ids/source/IDS-UCS-Basic.txt",
			IsZip:      false,
		},
		// IDS Ext A
		{
			Name:       "IDS Ext A",
			URL:        "https://raw.githubusercontent.com/cjkvi/cjkvi-ids/master/ids-ext-a.txt",
			TargetPath: "dictionaries/ids_ext_a/source/IDS-UCS-Ext-A.txt",
			IsZip:      false,
		},
	}

	for _, file := range files {
		// Check if final target file exists
		targetFile := file.TargetPath
		if file.IsJSONL && file.JSONLTarget != "" {
			targetFile = file.JSONLTarget
		}

		if _, err := os.Stat(targetFile); err == nil {
			logf("✓ %s already exists\n", file.Name)
			continue
		}

		logf("⬇ Downloading %s...\n", file.Name)

		// Create directory
		dir := filepath.Dir(file.TargetPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}

		// Download file
		if err := downloadFile(file.URL, file.TargetPath); err != nil {
			return fmt.Errorf("failed to download %s: %v", file.Name, err)
		}

		// Handle zip files
		if file.IsZip {
			logf("📦 Extracting %s...\n", file.Name)
			if err := extractZip(file.TargetPath, filepath.Dir(file.TargetPath)); err != nil {
				return fmt.Errorf("failed to extract %s: %v", file.Name, err)
			}
			// Remove the zip file after extraction
			os.Remove(file.TargetPath)
		}

		// Handle JSONL conversion
		if file.IsJSONL && file.JSONLTarget != "" {
			logf("🔄 Converting %s from JSONL to JSON...\n", file.Name)
			if err := convertJSONLToJSON(file.TargetPath, file.JSONLTarget); err != nil {
				return fmt.Errorf("failed to convert %s: %v", file.Name, err)
			}
		}

		logf("✅ %s setup complete\n", file.Name)
	}

	logf("All dictionary files are ready!\n")
	return nil
}

// downloadFile downloads a file from URL to the specified path
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractZip extracts a zip file to the specified directory
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Skip directories and hidden files
		if f.FileInfo().IsDir() || strings.HasPrefix(f.Name, ".") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)
		outFile, err := os.Create(path)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// convertJSONLToJSON converts a JSONL file to JSON using the existing converter
func convertJSONLToJSON(input, output string) error {
	cmd := exec.Command("go", "run", "cmd/jsonl2json/main.go", "-input="+input, "-output="+output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
