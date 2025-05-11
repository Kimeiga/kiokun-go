#!/bin/bash
# Script to test the dictionary build process locally

set -e  # Exit on error

echo "=== Testing Dictionary Build Process ==="

# Create directories if they don't exist
mkdir -p dictionaries/jmdict/source
mkdir -p dictionaries/jmnedict/source

# Download JMdict file with examples
echo "Downloading JMdict file with examples..."
curl -L -o dictionaries/jmdict/source/jmdict-examples-eng-3.6.1.json.zip https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/jmdict-examples-eng-3.6.1+20250505122413.json.zip

# Unzip the JMdict file
echo "Unzipping JMdict file..."
unzip -o dictionaries/jmdict/source/jmdict-examples-eng-3.6.1.json.zip -d dictionaries/jmdict/source/
# File is already named correctly, no need to rename

# Download JMNedict file
echo "Downloading JMNedict file..."
curl -L -o dictionaries/jmnedict/source/jmnedict-all-3.6.1.json.zip https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1+20250505122413/jmnedict-all-3.6.1+20250505122413.json.zip

# Unzip the JMNedict file
echo "Unzipping JMNedict file..."
unzip -o dictionaries/jmnedict/source/jmnedict-all-3.6.1.json.zip -d dictionaries/jmnedict/source/
# File is already named correctly, no need to rename

# Verify files were downloaded
echo "Verifying downloaded files..."
ls -la dictionaries/jmdict/source/
ls -la dictionaries/jmnedict/source/
ls -la dictionaries/kanjidic/source/
ls -la dictionaries/chinese_chars/source/
ls -la dictionaries/chinese_words/source/

# Build Dictionary
echo "Building dictionary..."
go run cmd/kiokun/main.go --writers 16 --silent

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Dictionary build successful!"
    echo "Generated $(find output -type f | wc -l) files"
else
    echo "Dictionary build failed!"
    exit 1
fi

echo "=== Test completed successfully ==="
