# JSONL to JSON Converter

A simple utility to convert JSON Lines (JSONL) files to standard JSON arrays.

## Purpose

This tool converts JSONL format (one JSON object per line) to a standard JSON array format that's suitable for use with schema generation tools like quicktype.io.

## Usage

```bash
# Build the tool
go build -o jsonl2json

# Run with required input file
./jsonl2json -input=path/to/file.jsonl

# Specify custom output file (optional)
./jsonl2json -input=path/to/file.jsonl -output=path/to/output.json
```

If no output file is specified, the tool will use the input filename with a `.json` extension.

## Example

Converting a Chinese dictionary:

```bash
# From project root
go run cmd/jsonl2json/main.go -input=dictionaries/chinese_characters/source/dictionary_char_2024-06-17.jsonl
```

This will create `dictionary_char_2024-06-17.json` in the current directory.
