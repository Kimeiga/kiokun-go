#!/bin/bash
# Script to test downloading and converting Chinese dictionary files

set -e  # Exit on error

echo "=== Testing Chinese Dictionary Download and Conversion ==="

# Create directories if they don't exist
mkdir -p dictionaries/chinese_chars/source
mkdir -p dictionaries/chinese_words/source

# Download and convert Chinese character dictionary
echo "Downloading Chinese character dictionary..."
curl -L -o dictionaries/chinese_chars/source/dictionary_char_2024-06-17.jsonl https://data.dong-chinese.com/dump/dictionary_char_2024-06-17.jsonl

echo "Converting Chinese character dictionary from JSONL to JSON..."
go run cmd/jsonl2json/main.go -input=dictionaries/chinese_chars/source/dictionary_char_2024-06-17.jsonl -output=dictionaries/chinese_chars/source/dictionary_char_2024-06-17.json

# Download and convert Chinese word dictionary
echo "Downloading Chinese word dictionary..."
curl -L -o dictionaries/chinese_words/source/dictionary_word_2024-06-17.jsonl https://data.dong-chinese.com/dump/dictionary_word_2024-06-17.jsonl

echo "Converting Chinese word dictionary from JSONL to JSON..."
go run cmd/jsonl2json/main.go -input=dictionaries/chinese_words/source/dictionary_word_2024-06-17.jsonl -output=dictionaries/chinese_words/source/dictionary_word_2024-06-17.json

# Verify files were downloaded and converted
echo "Verifying files..."
ls -la dictionaries/chinese_chars/source/
ls -la dictionaries/chinese_words/source/

echo "=== Test completed successfully ==="
