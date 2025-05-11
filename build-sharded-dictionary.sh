#!/bin/bash
# Script to build the dictionary in shards locally

set -e  # Exit on error

echo "=== Building Dictionary in Shards ==="

# Create directories if they don't exist
mkdir -p dictionaries/jmdict/source
mkdir -p dictionaries/jmnedict/source

# Download JMdict file with examples if it doesn't exist
if [ ! -f dictionaries/jmdict/source/jmdict-examples-eng-3.6.1.json ]; then
  echo "Downloading JMdict file with examples..."
  curl -L -o dictionaries/jmdict/source/jmdict-examples-eng-3.6.1.json.zip https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/jmdict-examples-eng-3.6.1+20250505122413.json.zip
  
  # Unzip the JMdict file
  echo "Unzipping JMdict file..."
  unzip -o dictionaries/jmdict/source/jmdict-examples-eng-3.6.1.json.zip -d dictionaries/jmdict/source/
fi

# Download JMNedict file if it doesn't exist
if [ ! -f dictionaries/jmnedict/source/jmnedict-all-3.6.1.json ]; then
  echo "Downloading JMNedict file..."
  curl -L -o dictionaries/jmnedict/source/jmnedict-all-3.6.1.json.zip https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/jmnedict-all-3.6.1+20250505122413.json.zip
  
  # Unzip the JMNedict file
  echo "Unzipping JMNedict file..."
  unzip -o dictionaries/jmnedict/source/jmnedict-all-3.6.1.json.zip -d dictionaries/jmnedict/source/
fi

# Verify files were downloaded
echo "Verifying downloaded files..."
ls -la dictionaries/jmdict/source/
ls -la dictionaries/jmnedict/source/
ls -la dictionaries/kanjidic/source/
ls -la dictionaries/chinese_chars/source/
ls -la dictionaries/chinese_words/source/

# Build each shard
build_shard() {
  local mode=$1
  local output_dir="output_${mode}"
  
  echo "=== Building $mode shard ==="
  go run cmd/kiokun/main.go --writers 16 --mode=$mode --silent
  echo "Generated $(find $output_dir -type f | wc -l) files in $output_dir"
  echo "=== Completed $mode shard ==="
  echo
}

# Build all shards
build_shard "han-1char"
build_shard "han-2char"
build_shard "han-3plus"
build_shard "non-han"

echo "=== All shards built successfully ==="
echo "Total files generated:"
echo "han-1char: $(find output_han_1char -type f | wc -l) files"
echo "han-2char: $(find output_han_2char -type f | wc -l) files"
echo "han-3plus: $(find output_han_3plus -type f | wc -l) files"
echo "non-han: $(find output_non_han -type f | wc -l) files"
