package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	lz4 "github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// Simulated network speeds (bytes per second)
const (
	NETWORK_SPEED_4G = 1_000_000 // 1MB/s (conservative 4G speed)
)

type CompressionResult struct {
	Algorithm           string
	OriginalSize        int64
	CompressedSize      int64
	CompressionRatio    float64
	CompressionTime     time.Duration
	DecompressionTime   time.Duration
	NetworkTransferTime time.Duration
}

func simulateNetworkTransfer(size int64) time.Duration {
	return time.Duration(float64(size) / float64(NETWORK_SPEED_4G) * float64(time.Second))
}

func compressGzip(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	if _, err := w.Write(data); err != nil {
		return nil, 0, err
	}
	if err := w.Close(); err != nil {
		return nil, 0, err
	}
	return compressed.Bytes(), time.Since(start), nil
}

func decompressGzip(compressed []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, 0, err
	}
	defer r.Close()
	decompressed, err := io.ReadAll(r)
	return decompressed, time.Since(start), err
}

func compressZstd(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	var compressed bytes.Buffer
	w, err := zstd.NewWriter(&compressed)
	if err != nil {
		return nil, 0, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, 0, err
	}
	if err := w.Close(); err != nil {
		return nil, 0, err
	}
	return compressed.Bytes(), time.Since(start), nil
}

func decompressZstd(compressed []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	r, err := zstd.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, 0, err
	}
	defer r.Close()
	decompressed, err := io.ReadAll(r)
	return decompressed, time.Since(start), err
}

func compressBrotli(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	var compressed bytes.Buffer
	w := brotli.NewWriter(&compressed)
	if _, err := w.Write(data); err != nil {
		return nil, 0, err
	}
	if err := w.Close(); err != nil {
		return nil, 0, err
	}
	return compressed.Bytes(), time.Since(start), nil
}

func decompressBrotli(compressed []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	r := brotli.NewReader(bytes.NewReader(compressed))
	decompressed, err := io.ReadAll(r)
	return decompressed, time.Since(start), err
}

func compressLZ4(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	var compressed bytes.Buffer
	w := lz4.NewWriter(&compressed)
	if _, err := w.Write(data); err != nil {
		return nil, 0, err
	}
	if err := w.Close(); err != nil {
		return nil, 0, err
	}
	return compressed.Bytes(), time.Since(start), nil
}

func decompressLZ4(compressed []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	r := lz4.NewReader(bytes.NewReader(compressed))
	decompressed, err := io.ReadAll(r)
	return decompressed, time.Since(start), err
}

func compressXZ(data []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	var compressed bytes.Buffer
	w, err := xz.NewWriter(&compressed)
	if err != nil {
		return nil, 0, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, 0, err
	}
	if err := w.Close(); err != nil {
		return nil, 0, err
	}
	return compressed.Bytes(), time.Since(start), nil
}

func decompressXZ(compressed []byte) ([]byte, time.Duration, error) {
	start := time.Now()
	r, err := xz.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, 0, err
	}
	decompressed, err := io.ReadAll(r)
	return decompressed, time.Since(start), err
}

func benchmarkCompression(data []byte) []CompressionResult {
	results := make([]CompressionResult, 0)
	originalSize := int64(len(data))

	// Test Gzip
	if compressed, compTime, err := compressGzip(data); err == nil {
		if _, decompTime, err := decompressGzip(compressed); err == nil {
			results = append(results, CompressionResult{
				Algorithm:           "Gzip",
				OriginalSize:        originalSize,
				CompressedSize:      int64(len(compressed)),
				CompressionRatio:    float64(len(compressed)) / float64(originalSize),
				CompressionTime:     compTime,
				DecompressionTime:   decompTime,
				NetworkTransferTime: simulateNetworkTransfer(int64(len(compressed))),
			})
		}
	}

	// Test Zstd
	if compressed, compTime, err := compressZstd(data); err == nil {
		if _, decompTime, err := decompressZstd(compressed); err == nil {
			results = append(results, CompressionResult{
				Algorithm:           "Zstd",
				OriginalSize:        originalSize,
				CompressedSize:      int64(len(compressed)),
				CompressionRatio:    float64(len(compressed)) / float64(originalSize),
				CompressionTime:     compTime,
				DecompressionTime:   decompTime,
				NetworkTransferTime: simulateNetworkTransfer(int64(len(compressed))),
			})
		}
	}

	// Test Brotli
	if compressed, compTime, err := compressBrotli(data); err == nil {
		if _, decompTime, err := decompressBrotli(compressed); err == nil {
			results = append(results, CompressionResult{
				Algorithm:           "Brotli",
				OriginalSize:        originalSize,
				CompressedSize:      int64(len(compressed)),
				CompressionRatio:    float64(len(compressed)) / float64(originalSize),
				CompressionTime:     compTime,
				DecompressionTime:   decompTime,
				NetworkTransferTime: simulateNetworkTransfer(int64(len(compressed))),
			})
		}
	}

	// Test LZ4
	if compressed, compTime, err := compressLZ4(data); err == nil {
		if _, decompTime, err := decompressLZ4(compressed); err == nil {
			results = append(results, CompressionResult{
				Algorithm:           "LZ4",
				OriginalSize:        originalSize,
				CompressedSize:      int64(len(compressed)),
				CompressionRatio:    float64(len(compressed)) / float64(originalSize),
				CompressionTime:     compTime,
				DecompressionTime:   decompTime,
				NetworkTransferTime: simulateNetworkTransfer(int64(len(compressed))),
			})
		}
	}

	// Test XZ
	if compressed, compTime, err := compressXZ(data); err == nil {
		if _, decompTime, err := decompressXZ(compressed); err == nil {
			results = append(results, CompressionResult{
				Algorithm:           "XZ",
				OriginalSize:        originalSize,
				CompressedSize:      int64(len(compressed)),
				CompressionRatio:    float64(len(compressed)) / float64(originalSize),
				CompressionTime:     compTime,
				DecompressionTime:   decompTime,
				NetworkTransferTime: simulateNetworkTransfer(int64(len(compressed))),
			})
		}
	}

	return results
}

func main() {
	// Create test directory
	testDir := "web"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("Error creating test directory: %v\n", err)
		return
	}

	// Create test data with multiple entries
	testEntries := make(map[string]interface{})
	baseEntry := map[string]interface{}{
		"w_j": []interface{}{
			map[string]interface{}{
				"id": "1000001",
				"kanji": []interface{}{
					map[string]interface{}{
						"common": true,
						"text":   "私",
					},
				},
				"kana": []interface{}{
					map[string]interface{}{
						"common": true,
						"text":   "わたし",
					},
				},
				"sense": []interface{}{
					map[string]interface{}{
						"partOfSpeech": []string{"pn"},
						"gloss": []interface{}{
							map[string]interface{}{
								"lang": "eng",
								"text": "I",
							},
						},
					},
				},
			},
		},
	}

	// Create 100 entries with slight variations
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("entry_%d", i)
		entry := make(map[string]interface{})
		for k, v := range baseEntry {
			entry[k] = v
		}
		testEntries[key] = entry
	}

	// Convert to JSON
	data, err := json.MarshalIndent(testEntries, "", "  ")
	if err != nil {
		fmt.Printf("Error creating test data: %v\n", err)
		return
	}

	// Save test data
	if err := os.WriteFile(filepath.Join(testDir, "test_entries.json"), data, 0644); err != nil {
		fmt.Printf("Error writing test data: %v\n", err)
		return
	}

	// Run benchmarks
	results := benchmarkCompression(data)

	// Print results
	fmt.Printf("Original size: %d bytes\n\n", len(data))
	fmt.Printf("%-10s | %-15s | %-15s | %-15s | %-15s | %-15s\n",
		"Algorithm",
		"Comp. Size",
		"Comp. Ratio",
		"Comp. Time",
		"Decomp. Time",
		"Network Time")
	fmt.Println(strings.Repeat("-", 90))

	for _, result := range results {
		fmt.Printf("%-10s | %-15d | %-15.2f | %-15s | %-15s | %-15s\n",
			result.Algorithm,
			result.CompressedSize,
			result.CompressionRatio,
			result.CompressionTime,
			result.DecompressionTime,
			result.NetworkTransferTime)
	}

	// Save compressed versions
	for _, result := range results {
		if result.Algorithm == "Gzip" || result.Algorithm == "Brotli" {
			compressed, _, err := eval("compress" + result.Algorithm)(data)
			if err != nil {
				fmt.Printf("Error compressing with %s: %v\n", result.Algorithm, err)
				continue
			}
			filename := filepath.Join(testDir, "compressed_"+strings.ToLower(result.Algorithm)+".bin")
			if err := os.WriteFile(filename, compressed, 0644); err != nil {
				fmt.Printf("Error writing compressed file: %v\n", err)
			}
		}
	}

	// Save benchmark results as JSON
	resultsJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling results: %v\n", err)
		return
	}
	if err := os.WriteFile(filepath.Join(testDir, "benchmark_results.json"), resultsJSON, 0644); err != nil {
		fmt.Printf("Error writing results: %v\n", err)
		return
	}

	fmt.Println("\nTest files created in 'web' directory")
	fmt.Println("Now you can run a development server (e.g., Vite) to test browser decompression")
}

// Helper function to dynamically call compression functions
func eval(name string) func([]byte) ([]byte, time.Duration, error) {
	switch name {
	case "compressGzip":
		return compressGzip
	case "compressBrotli":
		return compressBrotli
	default:
		return nil
	}
}
