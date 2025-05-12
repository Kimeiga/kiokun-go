# IDS Dictionary

This package provides support for Ideographic Description Sequences (IDS) in the Kiokun Dictionary application. IDS data is sourced from the [CHISE IDS database](https://gitlab.chise.org/CHISE/ids), which provides detailed information about the composition of Han characters.

## Data Format

The IDS data is stored in text files with the following format:

```
<CODEPOINT><TAB><CHARACTER><TAB><IDS>(<TAB>@apparent=<IDS>)
```

Where:
- `<CODEPOINT>` is the Unicode code point (e.g., U+4E00)
- `<CHARACTER>` is the character itself
- `<IDS>` is the Ideographic Description Sequence
- `@apparent=<IDS>` is an optional field for the apparent structure

## Files

- `IDS-UCS-Basic.txt`: CJK Unified Ideographs (U+4E00 - U+9FA5)
- `IDS-UCS-Ext-A.txt`: CJK Unified Ideographs Extension A (U+3400 - U+4DB5)

## Usage

The IDS data is automatically loaded when the application starts. It can be accessed through the common dictionary interface:

```go
import (
    "kiokun-go/dictionaries/common"
    "kiokun-go/dictionaries/ids"
    _ "kiokun-go/dictionaries/ids" // Import for side effects (auto-registration)
)

func main() {
    // Get all entries
    entries, err := common.ImportAllDictionaries()
    if err != nil {
        log.Fatalf("Failed to import dictionaries: %v", err)
    }

    // Process entries
    for _, entry := range entries {
        if idsEntry, ok := entry.(ids.IDSEntry); ok {
            fmt.Printf("Character: %s\n", idsEntry.Character)
            fmt.Printf("IDS: %s\n", idsEntry.IDS)
            if idsEntry.ApparentIDS != "" {
                fmt.Printf("Apparent IDS: %s\n", idsEntry.ApparentIDS)
            }
        }
    }
}
```

## Credits

The IDS data is sourced from the [CHISE IDS database](https://gitlab.chise.org/CHISE/ids), which is maintained by MORIOKA Tomohiko and contributors. The data is licensed under the GNU General Public License.
