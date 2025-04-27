#!/bin/bash

echo "=== Building and running dictionary debug tools ==="

# Create directories if they don't exist
mkdir -p cmd/debug_importer cmd/debug_processor_kanji

# Clean up previous test output
rm -rf output_kanji_test
mkdir -p output_kanji_test

# Check if JMNedict entries are loading properly
echo -e "\n=== Debugging JMNedict Importer ==="
go run cmd/debug_importer/main.go

# Check if Kanjidic entries are being processed correctly
echo -e "\n=== Debugging Kanjidic Processing ==="
go run cmd/debug_processor_kanji/main.go

echo -e "\n=== Debug process completed ===" 