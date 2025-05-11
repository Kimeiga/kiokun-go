# Kiokun-Go

A multilingual dictionary processor for Japanese and Chinese dictionary sources (JMdict, JMNedict, Kanjidic, Chinese characters, and Chinese words). It creates an optimized index-based structure for efficient lookup and retrieval.

## Recent Updates

### Index Structure Optimization

We've optimized the index structure to use single-letter field names for maximum compression:

```json
{
  "e": {
    "j": [1234, 5678], // JMdict exact matches
    "n": [7890], // JMNedict exact matches
    "d": [3691], // Kanjidic exact matches
    "c": [1592], // Chinese Characters exact matches
    "w": [7531] // Chinese Words exact matches
  },
  "c": {
    "j": [9012, 3456], // JMdict contained-in matches
    "n": [1357, 2468], // JMNedict contained-in matches
    "w": [9517, 8642] // Chinese Words contained-in matches
  }
}
```

This structure:

- Distinguishes between exact matches and contained-in matches
- Uses minimal field names for optimal compression
- Supports pagination for contained-in matches

### Directory Structure Optimization

We've optimized the directory structure to use one-letter names:

```
output/
├── index/
│   ├── 人.json.br
│   ├── 日.json.br
│   └── ...
├── j/        # JMdict entries
├── n/        # JMNedict entries
├── d/        # Kanjidic entries
├── c/        # Chinese character entries
└── w/        # Chinese word entries
```

### Chinese Dictionary Processing Fix

We've fixed the Chinese dictionary processing to correctly handle MongoDB-style field names, enabling the full integration of Chinese and Japanese dictionaries. The fix addresses:

- For Chinese characters:

  - `_id` → `ID`
  - `char` → `Traditional`
  - `simpVariants` → `Simplified`
  - `gloss` → `Definitions`
  - `strokeCount` → `StrokeCount`

- For Chinese words:
  - `_id` → `ID`
  - `trad` → `Traditional`
  - `simp` → `Simplified`
  - `definitions` or `gloss` → `Definitions`
  - `pinyin` or `pinyinSearchString` → `Pinyin`
  - `statistics.hskLevel` → `HskLevel`

#### ID Handling in Chinese Dictionaries

The Chinese dictionaries use MongoDB ObjectIDs as their primary keys. When these are imported:

1. The MongoDB `_id` field (e.g., "5f523afdde54193ed872ebb6") is preserved as the `ID` field in our structs
2. When these string IDs are converted to numeric values for indexing, they can result in both positive and negative 64-bit integers due to the way two's complement representation works
3. You may see entries with negative ID values in the index files - this is expected and doesn't affect functionality

Example of Chinese word IDs in an index file:

```json
"w": [
  -7375338670634322615,
  -7419189393381929367,
  6323202928105781916,
  6325997886664126603
]
```

The mix of positive and negative IDs is normal and doesn't impact the system's operation. The important thing is that these IDs are unique and consistent throughout the system.

### Removed Legacy Code

We've removed support for the combined-file solution and other legacy code, simplifying the codebase and focusing on the index-based approach.

## Project Structure

- `cmd/kiokun/` - Main application
  - `main.go` - Entry point
  - `internal/` - Internal packages
- `dictionaries/` - Dictionary importers and data structures
  - `common/` - Common interfaces and utilities
  - `jmdict/` - JMdict dictionary importer
  - `jmnedict/` - JMNedict dictionary importer
  - `kanjidic/` - Kanjidic dictionary importer
  - `chinese_chars/` - Chinese character dictionary importer
  - `chinese_words/` - Chinese word dictionary importer
- `processor/` - Dictionary processing logic
  - `index_processor.go` - Index-based processor

## Dictionary Formats

### JMdict (Japanese Words)

```go
type Word struct {
    ID    string       `json:"id"`
    Kanji []KanjiEntry `json:"kanji"`
    Kana  []KanaEntry  `json:"kana"`
    Sense []Sense      `json:"sense"`
}
```

### JMNedict (Japanese Names)

```go
type Name struct {
    ID       string   `json:"id"`
    Kanji    []string `json:"k"`
    Reading  []string `json:"r"`
    Meanings []string `json:"m"`
    Type     []string `json:"type,omitempty"`
}
```

### Kanjidic (Japanese Characters)

```go
type Kanji struct {
    Character string   `json:"c"`
    Meanings  []string `json:"m,omitempty"`
    OnYomi    []string `json:"on,omitempty"`
    KunYomi   []string `json:"kun,omitempty"`
    JLPT      int      `json:"jlpt,omitempty"`
    Grade     int      `json:"grade,omitempty"`
    Stroke    int      `json:"stroke"`
    Frequency int      `json:"freq,omitempty"`
    Radicals  []string `json:"rad,omitempty"`
}
```

### Chinese Characters

```go
type ChineseCharEntry struct {
    ID          string   `json:"id"`
    Traditional string   `json:"traditional"`
    Simplified  string   `json:"simplified"`
    Definitions []string `json:"definitions,omitempty"`
    Pinyin      []string `json:"pinyin,omitempty"`
    StrokeCount int      `json:"strokeCount,omitempty"`
}
```

### Chinese Words

```go
type ChineseWordEntry struct {
    ID          string         `json:"id"`
    Traditional string         `json:"traditional"`
    Simplified  string         `json:"simplified"`
    Pinyin      []string       `json:"pinyin,omitempty"`
    Definitions []string       `json:"definitions,omitempty"`
    HskLevel    int            `json:"hskLevel,omitempty"`
    Frequency   map[string]int `json:"frequency,omitempty"`
}
```

## Usage

### Building the Full Dictionary

To build the complete dictionary with all entries:

```bash
go run cmd/kiokun/main.go
```

This will:

1. Import all dictionaries from the `dictionaries/` directory
2. Process and index entries
3. Write compressed JSON files to the `output/` directory

### Command-Line Options

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
- `--mode <mode>` - Output mode: 'all', 'han-only', 'han-1char', 'han-2char', 'han-3plus', or 'non-han'
- `--test` - Test mode - prioritize entries that have overlap between Chinese and Japanese dictionaries

### Filtering Modes

The `--mode` flag allows filtering entries based on their character composition:

- `all` - Include all entries (default)
- `han-only` - Include only entries with Han characters
- `han-1char` - Include only entries with exactly 1 Han character
- `han-2char` - Include only entries with exactly 2 Han characters
- `han-3plus` - Include only entries with 3 or more Han characters
- `non-han` - Include only entries with at least one non-Han character

Example:

```bash
go run cmd/kiokun/main.go --mode han-1char
```

## Frontend Integration

When using the index mode, the frontend can load dictionary entries as needed:

```javascript
// Example frontend code (using fetch and Promise.all)
async function loadEntry(character) {
  // First, load the index file
  const indexResponse = await fetch(
    `https://cdn.example.com/index/${character}.json.br`
  );
  const index = await indexResponse.json();

  // Prepare promises for exact matches
  const exactPromises = [];

  // Load JMdict exact matches
  if (index.e && index.e.j && index.e.j.length > 0) {
    const jmdictPromises = index.e.j.map((id) =>
      fetch(`https://cdn.example.com/j/${id}.json.br`).then((res) => res.json())
    );
    exactPromises.push(Promise.all(jmdictPromises));
  } else {
    exactPromises.push(Promise.resolve([]));
  }

  // Load other exact matches similarly...

  // Wait for all exact match promises to resolve
  const [
    jmdictExactEntries,
    jmnedictExactEntries,
    kanjidicExactEntries,
    chineseCharsExactEntries,
    chineseWordsExactEntries,
  ] = await Promise.all(exactPromises);

  // Function to load contained-in matches (can be called later)
  async function loadContainedMatches() {
    const containedPromises = [];

    // Load JMdict contained-in matches
    if (index.c && index.c.j && index.c.j.length > 0) {
      const jmdictPromises = index.c.j.map((id) =>
        fetch(`https://cdn.example.com/j/${id}.json.br`).then((res) =>
          res.json()
        )
      );
      containedPromises.push(Promise.all(jmdictPromises));
    } else {
      containedPromises.push(Promise.resolve([]));
    }

    // Load other contained-in matches similarly...

    // Return contained-in matches
    const [
      jmdictContainedEntries,
      jmnedictContainedEntries,
      chineseWordsContainedEntries,
    ] = await Promise.all(containedPromises);

    return {
      jmdict: jmdictContainedEntries,
      jmnedict: jmnedictContainedEntries,
      chineseWords: chineseWordsContainedEntries,
    };
  }
}
```

## GitHub Actions Workflow

The GitHub Actions workflows build and deploy the dictionary files to separate repositories based on the character type:

1. **build-han-1char.yml**: Builds dictionary files for single Han characters
2. **build-han-2char.yml**: Builds dictionary files for two-character Han words
3. **build-han-3plus.yml**: Builds dictionary files for Han words with 3+ characters
4. **build-non-han.yml**: Builds dictionary files for non-Han entries
5. **dictionary-build.yml**: General dictionary build workflow

The workflows need to be updated to use the new index-based approach:

1. Update the output directory paths to match the new structure
2. Add more detailed output reporting
3. Update the Go version
4. Add path triggers for processor and dictionaries

## Testing

Run all tests:

```bash
go test ./...
```

Run specific tests:

```bash
# Test the processor package
go test ./processor

# Test the Chinese-Japanese integration
go test ./processor -run TestChineseJapaneseIntegration

# Test with verbose output
go test -v ./processor
```

## Performance Optimization

For faster development and processing, consider these optimizations:

1. **Use a compiled binary instead of `go run`**:

   ```bash
   cd cmd/kiokun
   go build -o kiokun
   ./kiokun [options]
   ```

2. **Use the `--dev` flag for better I/O performance**:

   ```bash
   ./kiokun --dev --writers 16
   ```

3. **Limit the number of entries for faster testing**:

   ```bash
   ./kiokun --limit 1000 --dev
   ```

## Future Improvements

1. **WebAssembly Support**:

   - Consider using WebAssembly as an alternative to static JSON files for dictionary delivery
   - This could improve download size efficiency and performance

2. **Pagination Support**:

   - Implement pagination for contained-in matches to improve frontend performance
   - This would allow loading a subset of contained-in matches at a time

3. **Improved Testing**:

   - Add more comprehensive tests for the Chinese dictionary processing
   - Add tests for the new index structure and API

4. **Documentation Updates**:
   - Add more detailed examples of how to use the index files on the frontend
   - Add more documentation about the Chinese dictionary structure and processing

## Git Hooks

This repository includes Git hooks to help maintain code quality and prevent common issues.

### Large File Check

The pre-commit hook checks for files that exceed GitHub's size limits:

- Files over 50MB will trigger a warning
- Files over 100MB will block the commit

This helps prevent issues with pushing large files to GitHub, which has a maximum file size limit of 100MB.

### Installation

You can install the hooks in two ways:

#### Option 1: Run the install script

```bash
./install-hooks.sh
```

#### Option 2: Set Git to use the hooks directory

```bash
git config core.hooksPath .githooks
```

### Bypassing Hooks

If you need to bypass the pre-commit hook for a specific commit:

```bash
git commit --no-verify -m "Your commit message"
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
