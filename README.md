# Kiokun-Go

A multilingual dictionary processor for Japanese and Chinese dictionary sources (JMdict, JMNedict, Kanjidic, Chinese characters, and Chinese words). It creates an optimized index-based structure for efficient lookup and retrieval.

This project uses a sharded architecture to handle the large number of dictionary files. The dictionary is split into four shards: Han-1 character, Han-2 character, Han-3+ character, and Non-Han.

## Recent Updates

### Sharded Repository Architecture

To handle the large number of dictionary files, we've implemented a sharded repository architecture:

1. **Shard Types**:

   - `han-1char`: Entries with exactly 1 Han character
   - `han-2char`: Entries with exactly 2 Han characters
   - `han-3plus`: Entries with 3 or more Han characters
   - `non-han`: Entries with at least one non-Han character

2. **Repository Structure**:

   - Each shard type has its own dedicated repository
   - The main repository coordinates builds and triggers workflows for each shard
   - Each shard repository follows the same internal structure (index/, j/, n/, d/, c/, w/)

3. **Numeric ID System**:
   - All dictionaries now use numeric IDs for maximum compression
   - IDs include a shard type prefix (0-3) to identify which shard they belong to
   - Example: `11000001` for a single Han character JMdict entry with original ID "1000001"
   - The first digit indicates the shard: 0=non-han, 1=han-1char, 2=han-2char, 3=han-3plus

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

### Unified Numeric ID System

All dictionaries now use a unified numeric ID system:

1. **JMdict and JMNedict**: Already use numeric IDs (e.g., "1000001")
2. **Kanjidic**: Now uses sequential numeric IDs (e.g., "1", "2", "3", ...)
3. **Chinese Chars**: Now uses sequential numeric IDs (e.g., "1", "2", "3", ...)
4. **Chinese Words**: Now uses sequential numeric IDs (e.g., "1", "2", "3", ...)

#### Shard Type Prefixes

Each ID is prefixed with a single digit indicating its shard type:

- **0**: Non-Han entries (contains at least one non-Han character)
- **1**: Han-1char entries (exactly 1 Han character)
- **2**: Han-2char entries (exactly 2 Han characters)
- **3**: Han-3plus entries (3 or more Han characters)

Examples:

- `01000001`: A non-Han JMdict entry with original ID "1000001"
- `11000002`: A single Han character JMdict entry with original ID "1000002"
- `21000003`: A two Han character JMdict entry with original ID "1000003"
- `31000004`: A three+ Han character JMdict entry with original ID "1000004"

#### Character-to-ID Mappings

For dictionaries that previously used non-numeric IDs (Kanjidic, Chinese Chars, Chinese Words), we maintain mappings from characters/words to their assigned numeric IDs. These mappings ensure that:

1. The same character or word always gets the same numeric ID
2. Entries can be looked up by their original characters or words
3. The system remains backward compatible with existing code

#### Benefits of the Numeric ID System

1. **Smaller JSON Size**: Numeric IDs take less space in JSON files
2. **Efficient Indexing**: Numeric IDs can be efficiently indexed and looked up
3. **Shard Identification**: The shard type prefix allows quick identification of which shard an entry belongs to
4. **Consistent Format**: All dictionaries now use the same ID format

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
    NumericID string   `json:"id"`  // Added numeric ID
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

When using the sharded architecture, the frontend needs to determine which repository to query based on the search term:

```javascript
// Configuration for dictionary repositories
const config = {
  // Repositories for each shard
  repositories: {
    nonHan: "japanese-dict-non-han",
    han1Char: "japanese-dict-han-1char",
    han2Char: "japanese-dict-han-2char",
    han3Plus: "japanese-dict-han-3plus",
  },
  cdnBase: "https://cdn.jsdelivr.net/gh/",
  branch: "main",
  fileExtension: ".json.br",
};

// Check if a string contains only Han characters
function isHanOnly(word) {
  for (const char of word) {
    const code = char.codePointAt(0);
    // Check if character is not in the Han Unicode block
    if (!(code >= 0x4e00 && code <= 0x9fff)) {
      return false;
    }
  }
  return true;
}

// Get the repository for a search term
function getRepositoryForWord(word) {
  // First check if it contains only Han characters
  if (isHanOnly(word)) {
    // Split based on character length
    const charCount = word.length;
    if (charCount === 1) {
      return config.repositories.han1Char;
    } else if (charCount === 2) {
      return config.repositories.han2Char;
    } else {
      // 3 or more characters
      return config.repositories.han3Plus;
    }
  } else {
    // Contains at least one non-Han character
    return config.repositories.nonHan;
  }
}

// Example frontend code (using fetch and Promise.all)
async function loadEntry(word) {
  // Determine which repository to query
  const repo = getRepositoryForWord(word);

  // First, load the index file
  const indexResponse = await fetch(
    `${config.cdnBase}${repo}@${config.branch}/index/${word}${config.fileExtension}`
  );
  const index = await indexResponse.json();

  // Prepare promises for exact matches
  const exactPromises = [];

  // Load JMdict exact matches
  if (index.e && index.e.j && index.e.j.length > 0) {
    const jmdictPromises = index.e.j.map((id) =>
      fetch(
        `${config.cdnBase}${repo}@${config.branch}/j/${id}${config.fileExtension}`
      ).then((res) => res.json())
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
        fetch(
          `${config.cdnBase}${repo}@${config.branch}/j/${id}${config.fileExtension}`
        ).then((res) => res.json())
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

### Workflow Architecture

1. **dictionary-build.yml**: Coordinator workflow that:

   - Downloads and prepares dictionary source files
   - Triggers individual shard build workflows

2. **Shard-specific workflows**:
   - **build-han-1char.yml**: Builds dictionary files for single Han characters
   - **build-han-2char.yml**: Builds dictionary files for two-character Han words
   - **build-han-3plus.yml**: Builds dictionary files for Han words with 3+ characters
   - **build-non-han.yml**: Builds dictionary files for non-Han entries

### Repository Structure

Each shard has its own dedicated repository:

- `japanese-dict-han-1char`: For single Han character entries
- `japanese-dict-han-2char`: For two-character Han entries
- `japanese-dict-han-3plus`: For Han entries with 3+ characters
- `japanese-dict-non-han`: For non-Han entries

### Required Secrets

To enable the workflows to push to separate repositories, you need to set up these secrets:

- `DICTIONARY_DEPLOY_TOKEN`: A personal access token with write access to the shard repositories
- `REPO_ACCESS_TOKEN`: A token with permission to trigger repository dispatch events

### Local Testing

You can test the sharded build process locally using the provided script:

```bash
./build-sharded-dictionary.sh
```

This script:

1. Downloads necessary dictionary files
2. Builds each shard sequentially
3. Reports the number of files generated for each shard

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
