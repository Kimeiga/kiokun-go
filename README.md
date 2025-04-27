# Kiokun-Go

A Japanese dictionary processor that combines multiple dictionary sources (JMdict, JMNedict, and Kanjidic) into a unified format.

## Project Structure

- `cmd/kiokun/` - Main application
- `cmd/debug_*` - Debug tools for testing specific components
- `dictionaries/` - Dictionary importers and data structures
  - `common/` - Common interfaces and utilities
  - `jmdict/` - JMdict dictionary importer
  - `jmnedict/` - JMNedict dictionary importer
  - `kanjidic/` - Kanjidic dictionary importer
- `processor/` - Dictionary processing logic
- `analyzer/` - Dictionary analysis tools

## Usage

### Main Application

```bash
go run cmd/kiokun/main.go [options]
```

Options:

- `--dictdir <dir>` - Base directory containing dictionary packages (default: "dictionaries")
- `--outdir <dir>` - Output directory for processed files (default: "output")
- `--workers <n>` - Number of worker goroutines for batch processing (default: CPU count)
- `--writers <n>` - Number of parallel workers for file writing (default: CPU count)
- `--silent` - Disable progress output
- `--dev` - Development mode - use /tmp for faster I/O
- `--limit <n>` - Limit the number of entries to process (0 = no limit)
- `--batch <n>` - Process entries in batches of this size (default: 10000)

### Debug Tools

Various debug tools are available in the `cmd/` directory:

- `cmd/debug_all_dictionaries/` - Test importing all dictionaries
- `cmd/debug_kanjidic_importer/` - Test the Kanjidic importer
- `cmd/debug_processor_kanji/` - Test processing Kanjidic entries
- `cmd/debug_importer/` - Test dictionary importers
- `cmd/debug_analyzer/` - Test dictionary analysis
- `cmd/debug_processor/` - Test dictionary processing

## Development

### Building

```bash
go build -o kiokun cmd/kiokun/main.go
```

### Testing

```bash
go test ./...
```

### Performance Optimization

For faster development, use the `--dev` flag to write output to `/tmp` for better I/O performance:

```bash
go run cmd/kiokun/main.go --dev --writers 16
```
