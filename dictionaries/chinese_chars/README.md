# Chinese Character Dictionary

This package provides functionality for importing, processing, and accessing Chinese character data.

## Overview

The Chinese Character Dictionary integrates with the common dictionary framework and provides:

1. **Automatic Import**: The dictionary is automatically registered and loaded when the application starts
2. **Simple Interface**: Characters are accessible through the common dictionary interface
3. **Cross-Language Support**: Designed to work alongside other language dictionaries

## Usage

The dictionary is automatically loaded on application start via the `init()` function, which:

1. Locates the dictionary file in the `dictionaries/chinese_characters/source` directory
2. Registers the dictionary with the common framework
3. Makes the entries available through the common dictionary API

### Accessing Dictionary Entries

Once loaded, you can access the Chinese characters through the common dictionary interface:

```go
import (
    "kiokun-go/dictionaries/common"
    _ "kiokun-go/dictionaries/chinese_chars" // Import for side effects (auto-registration)
)

func main() {
    // Get all entries
    entries, err := common.ImportAllDictionaries()
    if err != nil {
        log.Fatalf("Failed to import dictionaries: %v", err)
    }

    // Process entries
    for _, entry := range entries {
        if char, ok := entry.(chinese_chars.ChineseCharEntry); ok {
            fmt.Printf("Traditional: %s\n", char.Traditional)
            fmt.Printf("Simplified: %s\n", char.Simplified)
        }
    }
}
```

## Data Structure

Each character entry contains:

- **ID**: The unique identifier for the character
- **Traditional**: The traditional form of the character
- **Simplified**: The simplified form of the character
- **Definitions**: A list of definitions and glosses
- **Pinyin**: A list of Pinyin pronunciations
- **StrokeCount**: The number of strokes in the character

## Integration with Other Dictionaries

This dictionary is designed to be integrated with other language dictionaries, particularly Japanese, to create comprehensive cross-language character mappings.

### Example Mapping Strategy

Traditional Chinese characters can be used as an index to map between different character variants:

```
圖.json:
c_c: 圖 <- Chinese character dictionary entry
c_j: 図 <- Japanese equivalent

图.json:
c_c: 图 <- Simplified Chinese
c_j: 図 <- Japanese equivalent

図.json:
c_c: 圖 <- Traditional Chinese equivalent
c_j: 図 <- Japanese character
```

This allows for seamless lookups across writing systems.
