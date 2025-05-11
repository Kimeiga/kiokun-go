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

# Build all shards at once with the new sharded processor
echo "=== Building all shards at once ==="
go run cmd/kiokun/main.go --writers 16 --silent

# Check the output directories
for SHARD in "non_han" "han_1char" "han_2char" "han_3plus"; do
    SHARD_DIR="output_${SHARD}"
    if [ -d "$SHARD_DIR" ]; then
        INDEX_COUNT=$(find "$SHARD_DIR/index" -type f -name "*.json.br" | wc -l)
        JMDICT_COUNT=$(find "$SHARD_DIR/j" -type f -name "*.json.br" | wc -l)
        JMNEDICT_COUNT=$(find "$SHARD_DIR/n" -type f -name "*.json.br" | wc -l)
        KANJIDIC_COUNT=$(find "$SHARD_DIR/d" -type f -name "*.json.br" | wc -l)
        CHINESE_CHARS_COUNT=$(find "$SHARD_DIR/c" -type f -name "*.json.br" | wc -l)
        CHINESE_WORDS_COUNT=$(find "$SHARD_DIR/w" -type f -name "*.json.br" | wc -l)
        TOTAL_COUNT=$((JMDICT_COUNT + JMNEDICT_COUNT + KANJIDIC_COUNT + CHINESE_CHARS_COUNT + CHINESE_WORDS_COUNT))

        echo "Shard $SHARD:"
        echo "  - Index files: $INDEX_COUNT"
        echo "  - JMdict entries: $JMDICT_COUNT"
        echo "  - JMNedict entries: $JMNEDICT_COUNT"
        echo "  - Kanjidic entries: $KANJIDIC_COUNT"
        echo "  - Chinese character entries: $CHINESE_CHARS_COUNT"
        echo "  - Chinese word entries: $CHINESE_WORDS_COUNT"
        echo "  - Total dictionary entries: $TOTAL_COUNT"
    else
        echo "Shard $SHARD: Directory not found"
    fi
done

echo "=== All shards built successfully ==="
echo "Total files generated:"
echo "han-1char: $(find output_han_1char -type f 2>/dev/null | wc -l) files"
echo "han-2char: $(find output_han_2char -type f 2>/dev/null | wc -l) files"
echo "han-3plus: $(find output_han_3plus -type f 2>/dev/null | wc -l) files"
echo "non-han: $(find output_non_han -type f 2>/dev/null | wc -l) files"
