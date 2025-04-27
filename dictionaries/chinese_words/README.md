# Chinese Word Dictionary

This package provides functionality for importing, processing, and accessing Chinese word data.

## Overview

The Chinese Word Dictionary integrates with the common dictionary framework and provides:

1. **Automatic Import**: The dictionary is automatically registered and loaded when the application starts
2. **Simple Interface**: Words are accessible through the common dictionary interface
3. **Cross-Language Support**: Designed to work alongside other language dictionaries

## Usage

The dictionary is automatically loaded on application start via the `init()` function, which:

1. Locates the dictionary file in the `dictionaries/chinese_words/source` directory
2. Registers the dictionary with the common framework
3. Makes the entries available through the common dictionary API

### Accessing Dictionary Entries

Once loaded, you can access the Chinese words through the common dictionary interface:

```go
import (
    "kiokun-go/dictionaries/common"
    _ "kiokun-go/dictionaries/chinese_words" // Import for side effects (auto-registration)
)

func main() {
    // Get all entries
    entries, err := common.ImportAllDictionaries()
    if err != nil {
        log.Fatalf("Failed to import dictionaries: %v", err)
    }

    // Process entries
    for _, entry := range entries {
        if word, ok := entry.(chinese_words.ChineseWordEntry); ok {
            fmt.Printf("Traditional: %s\n", word.Traditional)
            fmt.Printf("Simplified: %s\n", word.Simplified)
            fmt.Printf("Definitions: %v\n", word.Definitions)
        }
    }
}
```

## Data Structure

Each word entry contains:

- **ID**: The unique identifier for the word
- **Traditional**: The traditional form of the word
- **Simplified**: The simplified form of the word
- **Definitions**: A list of definitions and glosses
- **Pinyin**: A list of Pinyin pronunciations
- **HskLevel**: The HSK proficiency level of the word
- **Frequency**: Frequency statistics for various contexts

## Integration with Character Dictionary

The word dictionary is designed to work alongside the character dictionary:

1. Words can be decomposed into individual characters
2. Each character can be looked up in the character dictionary
3. This allows for creating rich language learning tools

## Example Applications

- **Word decomposition**: Break down words into their component characters
- **Character frequency analysis**: Determine which characters appear most frequently in words
- **Cross-referencing**: Look up related words containing a specific character
- **Frequency-based learning**: Order words by frequency for more efficient learning
