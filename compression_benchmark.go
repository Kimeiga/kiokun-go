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
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// ... existing code until main() ...

func main() {
	// Create test directory
	testDir := "test_data"
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

	// Create JavaScript test file
	jsTest := `<!DOCTYPE html>
<html>
<head>
    <title>Compression Test</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        pre { background: #f5f5f5; padding: 10px; border-radius: 4px; }
        .success { color: green; }
        .error { color: red; }
    </style>
</head>
<body>
    <h2>Decompression Test Results</h2>
    <pre id="results"></pre>

    <script>
        async function runTests() {
            const results = document.getElementById('results');
            const algorithms = ['gzip', 'brotli'];
            const testData = await fetch('test_entries.json').then(r => r.text());
            
            results.textContent = 'Running tests...\n\n';
            
            for (const algo of algorithms) {
                try {
                    results.textContent += `Testing ${algo}...\n`;
                    
                    const start = performance.now();
                    const response = await fetch(`compressed_${algo}.bin`, {
                        headers: {
                            'Accept-Encoding': algo === 'brotli' ? 'br' : 'gzip'
                        }
                    });
                    
                    const decompressed = await response.text();
                    const end = performance.now();
                    const time = end - start;
                    
                    results.textContent += `${algo}:\n`;
                    results.textContent += `  Time: ${time.toFixed(2)}ms\n`;
                    
                    // Verify decompression
                    if (decompressed === testData) {
                        results.textContent += '  Verification: ✓ Success\n\n';
                    } else {
                        results.textContent += '  Verification: ✗ Failed (content mismatch)\n\n';
                    }
                } catch (error) {
                    results.textContent += `  Error: ${error.message}\n\n`;
                }
            }
        }

        runTests().catch(error => {
            document.getElementById('results').textContent += `\nTest error: ${error.message}`;
        });
    </script>
</body>
</html>`

	// Save JavaScript test file
	if err := os.WriteFile(filepath.Join(testDir, "test.html"), []byte(jsTest), 0644); err != nil {
		fmt.Printf("Error writing JavaScript test: %v\n", err)
		return
	}

	// Save compressed versions for JavaScript testing
	for _, result := range results {
		if result.Algorithm == "Gzip" || result.Algorithm == "Brotli" {
			compressed, _, err := eval("compress" + result.Algorithm)(data)
			if err != nil {
				fmt.Printf("Error compressing for JavaScript test: %v\n", err)
				continue
			}
			filename := filepath.Join(testDir, "compressed_"+strings.ToLower(result.Algorithm)+".bin")
			if err := os.WriteFile(filename, compressed, 0644); err != nil {
				fmt.Printf("Error writing compressed file: %v\n", err)
			}
		}
	}

	fmt.Println("\nTest files created in 'test_data' directory")
	fmt.Println("To run browser tests, serve the directory with:")
	fmt.Println("python3 -m http.server 8000")
	fmt.Println("Then visit: http://localhost:8000/test_data/test.html")
} 