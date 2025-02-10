package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Command line flags
	inputFile := flag.String("input", "jmdict-eng-3.5.0.json", "Input JMDict JSON file")
	outputDir := flag.String("output", "dictionary", "Output directory for word files")
	numWorkers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines")
	flag.Parse()

	// Remove existing directory if it exists - using rm -rf for speed
	if _, err := os.Stat(*outputDir); err == nil {
		cmd := exec.Command("rm", "-rf", *outputDir)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing existing directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Create output directory
	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Read and parse input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var dict JmdictTypes
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dict); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	totalWords := len(dict.Words)
	fmt.Printf("Processing %d words with %d workers...\n", totalWords, *numWorkers)

	// Create worker pool
	var wg sync.WaitGroup
	wordChan := make(chan Word, totalWords)
	var processedCount int64 = 0

	// Start progress monitoring goroutine
	done := make(chan bool)
	quit := make(chan bool) // New channel for clean shutdown
	go showProgress(totalWords, &processedCount, done, quit)

	// Start workers
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, wordChan, *outputDir, &processedCount)
	}

	// Send words to workers
	for _, word := range dict.Words {
		wordChan <- word
	}
	close(wordChan)

	// Wait for all workers to finish
	wg.Wait()
	done <- true
	<-quit // Wait for progress goroutine to finish

	fmt.Printf("\nSuccessfully processed %d words\n", totalWords)
}

func showProgress(total int, current *int64, done, quit chan bool) {
	start := time.Now()
	for {
		select {
		case <-done:
			elapsed := time.Since(start)
			rate := float64(total) / elapsed.Seconds()
			fmt.Printf("\rCompleted 100%% (%d/%d) in %.1fs (%.1f words/sec)          \n",
				total, total, elapsed.Seconds(), rate)
			quit <- true
			return
		default:
			processed := atomic.LoadInt64(current)
			percent := float64(processed) * 100 / float64(total)
			elapsed := time.Since(start)
			rate := float64(processed) / elapsed.Seconds()
			fmt.Printf("\rProgress: %.1f%% (%d/%d) - %.1f words/sec",
				percent, processed, total, rate)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func worker(wg *sync.WaitGroup, words <-chan Word, outputDir string, processedCount *int64) {
	defer wg.Done()

	for word := range words {
		// Use first kanji or kana as filename, falling back to ID if neither exists
		var filename string
		if len(word.Kanji) > 0 {
			filename = word.Kanji[0].Text
		} else if len(word.Kana) > 0 {
			filename = word.Kana[0].Text
		} else {
			filename = word.ID
		}

		// Create filename with Unicode support
		filename = filepath.Join(outputDir, filename+".json.gz")

		// Create output file with UTF-8 support
		file, err := os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError creating file %s: %v\n", filename, err)
			continue
		}

		// Create gzip writer
		gzWriter := gzip.NewWriter(file)

		// Encode word to JSON and write to gzip
		encoder := json.NewEncoder(gzWriter)
		encoder.SetEscapeHTML(false) // Preserve Unicode characters
		if err := encoder.Encode(word); err != nil {
			fmt.Fprintf(os.Stderr, "\nError encoding word %s: %v\n", word.ID, err)
			file.Close()
			continue
		}

		// Close writers
		if err := gzWriter.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "\nError closing gzip writer for %s: %v\n", filename, err)
		}
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "\nError closing file %s: %v\n", filename, err)
		}

		atomic.AddInt64(processedCount, 1)
	}
}
