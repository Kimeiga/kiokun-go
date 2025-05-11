package processor

import (
	"fmt"
	"strconv"
	"unicode"

	"kiokun-go/dictionaries/chinese_chars"
	"kiokun-go/dictionaries/chinese_words"
	"kiokun-go/dictionaries/common"
	"kiokun-go/dictionaries/jmdict"
	"kiokun-go/dictionaries/jmnedict"
	"kiokun-go/dictionaries/kanjidic"
)

// ShardType identifies which shard an entry belongs to
type ShardType int

const (
	ShardNonHan   ShardType = 0 // Contains at least one non-Han character
	ShardHan1Char ShardType = 1 // Exactly 1 Han character
	ShardHan2Char ShardType = 2 // Exactly 2 Han characters
	ShardHan3Plus ShardType = 3 // 3 or more Han characters
)

// GetShardType determines which shard an entry belongs to
func GetShardType(entry common.Entry) ShardType {
	// Get primary text for the entry
	var primaryText string
	switch e := entry.(type) {
	case jmdict.Word:
		if len(e.Kanji) > 0 {
			primaryText = e.Kanji[0].Text
		} else if len(e.Kana) > 0 {
			primaryText = e.Kana[0].Text
		} else {
			primaryText = e.ID
		}
	case jmnedict.Name:
		if len(e.Kanji) > 0 {
			primaryText = e.Kanji[0]
		} else if len(e.Reading) > 0 {
			primaryText = e.Reading[0]
		} else {
			primaryText = e.ID
		}
	case kanjidic.Kanji:
		primaryText = e.Character
	case chinese_chars.ChineseCharEntry:
		primaryText = e.Traditional
	case chinese_words.ChineseWordEntry:
		primaryText = e.Traditional
	default:
		primaryText = entry.GetID()
	}

	// Check if it contains only Han characters
	isHan := isHanOnly(primaryText)
	charCount := len([]rune(primaryText))

	if !isHan {
		return ShardNonHan
	} else if charCount == 1 {
		return ShardHan1Char
	} else if charCount == 2 {
		return ShardHan2Char
	} else {
		return ShardHan3Plus
	}
}

// isHanOnly checks if a string contains only Han characters
func isHanOnly(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

// GetShardedID returns an ID that includes the shard information
func GetShardedID(entry common.Entry) string {
	originalID := entry.GetID()
	shardType := GetShardType(entry)

	// Create the sharded ID by prepending the shard type
	return fmt.Sprintf("%d%s", shardType, originalID)
}

// ExtractOriginalID extracts the original ID from a sharded ID
func ExtractOriginalID(shardedID string) string {
	// Remove the first digit (shard type)
	if len(shardedID) > 1 {
		return shardedID[1:]
	}
	return shardedID
}

// ExtractShardType extracts the shard type from a sharded ID
func ExtractShardType(shardedID string) (ShardType, error) {
	if len(shardedID) == 0 {
		return ShardNonHan, fmt.Errorf("empty sharded ID")
	}

	// Get the first digit
	shardTypeStr := string(shardedID[0])
	shardTypeInt, err := strconv.Atoi(shardTypeStr)
	if err != nil {
		return ShardNonHan, fmt.Errorf("invalid shard type: %s", shardTypeStr)
	}

	return ShardType(shardTypeInt), nil
}

// GetOutputDirForShard returns the output directory name for a shard type
func GetOutputDirForShard(baseDir string, shardType ShardType) string {
	switch shardType {
	case ShardNonHan:
		return fmt.Sprintf("%s_non_han", baseDir)
	case ShardHan1Char:
		return fmt.Sprintf("%s_han_1char", baseDir)
	case ShardHan2Char:
		return fmt.Sprintf("%s_han_2char", baseDir)
	case ShardHan3Plus:
		return fmt.Sprintf("%s_han_3plus", baseDir)
	default:
		return baseDir
	}
}
