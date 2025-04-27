package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	githubAPI = "https://api.github.com/repos/scriptin/jmdict-simplified/releases/latest"
)

type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var requiredFiles = map[string]struct {
	prefix   string
	destDir  string
	jsonName string
}{
	"jmdict": {
		prefix:   "jmdict-examples-eng-",
		destDir:  "dictionaries/jmdict/source/",
		jsonName: "jmdict-examples-eng-",
	},
	"jmnedict": {
		prefix:   "jmnedict-all-",
		destDir:  "dictionaries/jmnedict/source/",
		jsonName: "jmnedict-all-",
	},
	"kanjidic": {
		prefix:   "kanjidic2-en-",
		destDir:  "dictionaries/kanjidic/source/",
		jsonName: "kanjidic2-en-",
	},
}

func main() {
	useXZ := flag.Bool("xz", false, "Compress JSON files with XZ after extraction")
	flag.Parse()

	// Create HTTP client with reasonable timeout
	client := &http.Client{
		Timeout: 10 * time.Minute, // Long timeout for large files
	}

	// Fetch latest release info
	fmt.Println("Fetching latest release information...")
	release, err := getLatestRelease(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching release info: %v\n", err)
		os.Exit(1)
	}

	// Clean up old dictionary files
	for _, dict := range requiredFiles {
		if err := cleanupDirectory(dict.destDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error cleaning up %s: %v\n", dict.destDir, err)
		}
	}

	// Process each required file
	for dictName, dict := range requiredFiles {
		// Ensure destination directory exists
		if err := os.MkdirAll(dict.destDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dict.destDir, err)
			continue
		}

		// Find matching asset
		var matchingAsset *GitHubAsset
		for i := range release.Assets {
			if strings.HasPrefix(release.Assets[i].Name, dict.prefix) &&
				strings.HasSuffix(release.Assets[i].Name, ".zip") {
				matchingAsset = &release.Assets[i]
				break
			}
		}
		if matchingAsset == nil {
			fmt.Fprintf(os.Stderr, "Could not find asset starting with %s\n", dict.prefix)
			continue
		}

		fmt.Printf("Processing %s dictionary: %s\n", dictName, matchingAsset.Name)

		// Download and process the file
		if err := downloadAndExtract(client, matchingAsset.BrowserDownloadURL, dict.destDir, dict.jsonName, *useXZ); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", matchingAsset.Name, err)
			continue
		}
	}

	fmt.Println("Dictionary update complete!")
}

func cleanupDirectory(dir string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Remove all files in the directory
	for _, entry := range entries {
		if !entry.IsDir() { // Only remove files, not subdirectories
			path := filepath.Join(dir, entry.Name())
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %v", path, err)
			}
			fmt.Printf("Removed old file: %s\n", path)
		}
	}

	return nil
}

func getLatestRelease(client *http.Client) (*GitHubRelease, error) {
	req, err := http.NewRequest("GET", githubAPI, nil)
	if err != nil {
		return nil, err
	}

	// Add User-Agent header to avoid GitHub API limitations
	req.Header.Set("User-Agent", "Kiokun-Dictionary-Updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadAndExtract(client *http.Client, url, destDir, jsonPrefix string, useXZ bool) error {
	// Download file
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status: %s", resp.Status)
	}

	// Create temporary file for ZIP
	tmpFile, err := os.CreateTemp("", "dict-*.zip")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Copy download to temporary file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	// Open ZIP file
	zipReader, err := zip.OpenReader(tmpPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Extract JSON file
	for _, file := range zipReader.File {
		if !strings.HasSuffix(file.Name, ".json") || !strings.HasPrefix(file.Name, jsonPrefix) {
			continue
		}

		// Open file in ZIP
		rc, err := file.Open()
		if err != nil {
			return err
		}

		// Create destination file
		destPath := filepath.Join(destDir, file.Name)
		destFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return err
		}

		// Copy contents
		if _, err := io.Copy(destFile, rc); err != nil {
			rc.Close()
			destFile.Close()
			return err
		}

		rc.Close()
		destFile.Close()

		fmt.Printf("Extracted %s\n", destPath)

		// Compress with XZ if requested
		if useXZ {
			if err := compressWithXZ(destPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: XZ compression failed for %s: %v\n", destPath, err)
			} else {
				// Remove original JSON file
				os.Remove(destPath)
			}
		}
	}

	return nil
}

func compressWithXZ(filepath string) error {
	// Run xz command to compress the file
	cmd := exec.Command("xz", "-9", "-f", filepath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("xz compression failed: %v\nOutput: %s", err, output)
	}
	return nil
}
