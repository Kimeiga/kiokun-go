package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/andybalholm/brotli"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

func main() {
	// Source file
	sourcePath := filepath.Join("..", "test_data", "jmdict-examples-eng-3.6.1.json")

	// Create test_data directory in web-test/src
	testDataDir := filepath.Join("..", "web-test", "src", "test_data")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		fmt.Printf("Error creating test data directory: %v\n", err)
		return
	}

	// Read the source JSON file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Printf("Error reading source file: %v\n", err)
		return
	}

	// Save raw JSON
	rawPath := filepath.Join(testDataDir, "jmdict.json")
	if err := os.WriteFile(rawPath, data, 0644); err != nil {
		fmt.Printf("Error writing raw JSON: %v\n", err)
		return
	}

	// Compress with gzip
	var gzipBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipBuf)
	if _, err := gzipWriter.Write(data); err != nil {
		fmt.Printf("Error compressing with gzip: %v\n", err)
		return
	}
	gzipWriter.Close()

	gzipPath := filepath.Join(testDataDir, "jmdict.json.gz")
	if err := os.WriteFile(gzipPath, gzipBuf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing gzip file: %v\n", err)
		return
	}

	// Compress with brotli
	var brotliBuf bytes.Buffer
	brotliWriter := brotli.NewWriter(&brotliBuf)
	if _, err := brotliWriter.Write(data); err != nil {
		fmt.Printf("Error compressing with brotli: %v\n", err)
		return
	}
	brotliWriter.Close()

	brotliPath := filepath.Join(testDataDir, "jmdict.json.br")
	if err := os.WriteFile(brotliPath, brotliBuf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing brotli file: %v\n", err)
		return
	}

	// Compress with zstd
	var zstdBuf bytes.Buffer
	zstdWriter, err := zstd.NewWriter(&zstdBuf)
	if err != nil {
		fmt.Printf("Error creating zstd writer: %v\n", err)
		return
	}
	if _, err := zstdWriter.Write(data); err != nil {
		fmt.Printf("Error compressing with zstd: %v\n", err)
		return
	}
	zstdWriter.Close()

	zstdPath := filepath.Join(testDataDir, "jmdict.json.zst")
	if err := os.WriteFile(zstdPath, zstdBuf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing zstd file: %v\n", err)
		return
	}

	// Compress with lz4
	var lz4Buf bytes.Buffer
	lz4Writer := lz4.NewWriter(&lz4Buf)
	if _, err := lz4Writer.Write(data); err != nil {
		fmt.Printf("Error compressing with lz4: %v\n", err)
		return
	}
	lz4Writer.Close()

	lz4Path := filepath.Join(testDataDir, "jmdict.json.lz4")
	if err := os.WriteFile(lz4Path, lz4Buf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing lz4 file: %v\n", err)
		return
	}

	// Compress with xz
	var xzBuf bytes.Buffer
	xzWriter, err := xz.NewWriter(&xzBuf)
	if err != nil {
		fmt.Printf("Error creating xz writer: %v\n", err)
		return
	}
	if _, err := xzWriter.Write(data); err != nil {
		fmt.Printf("Error compressing with xz: %v\n", err)
		return
	}
	xzWriter.Close()

	xzPath := filepath.Join(testDataDir, "jmdict.json.xz")
	if err := os.WriteFile(xzPath, xzBuf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing xz file: %v\n", err)
		return
	}

	// Compress with snappy
	snappyData := snappy.Encode(nil, data)
	snappyPath := filepath.Join(testDataDir, "jmdict.json.snappy")
	if err := os.WriteFile(snappyPath, snappyData, 0644); err != nil {
		fmt.Printf("Error writing snappy file: %v\n", err)
		return
	}

	fmt.Printf("Successfully processed test file\n")
	fmt.Printf("  Raw size: %d bytes\n", len(data))
	fmt.Printf("  Gzip size: %d bytes\n", gzipBuf.Len())
	fmt.Printf("  Brotli size: %d bytes\n", brotliBuf.Len())
	fmt.Printf("  Zstd size: %d bytes\n", zstdBuf.Len())
	fmt.Printf("  LZ4 size: %d bytes\n", lz4Buf.Len())
	fmt.Printf("  XZ size: %d bytes\n", xzBuf.Len())
	fmt.Printf("  Snappy size: %d bytes\n", len(snappyData))

	fmt.Printf("\nFiles written to %s:\n", testDataDir)
	fmt.Printf("  - jmdict.json\n")
	fmt.Printf("  - jmdict.json.gz\n")
	fmt.Printf("  - jmdict.json.br\n")
	fmt.Printf("  - jmdict.json.zst\n")
	fmt.Printf("  - jmdict.json.lz4\n")
	fmt.Printf("  - jmdict.json.xz\n")
	fmt.Printf("  - jmdict.json.snappy\n")
}
