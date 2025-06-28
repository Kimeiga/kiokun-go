/**
 * Utility functions for dictionary lookups
 */

// Base URL for jsDelivr CDN
const BASE_URL = "https://cdn.jsdelivr.net/gh/Kimeiga";

// Repository names for different shard types
const REPOS = {
  NON_HAN: "japanese-dict-non-han",
  HAN_1CHAR: "japanese-dict-han-1char",
  HAN_2CHAR: "japanese-dict-han-2char",
  HAN_3PLUS: "japanese-dict-han-3plus",
};

// Shard types
export enum ShardType {
  NON_HAN = 0,
  HAN_1CHAR = 1,
  HAN_2CHAR = 2,
  HAN_3PLUS = 3,
}

// Dictionary types
export enum DictionaryType {
  JMDICT = "j",
  JMNEDICT = "n",
  KANJIDIC = "d",
  CHINESE_CHARS = "c",
  CHINESE_WORDS = "w",
}

// Index entry structure
export interface IndexEntry {
  e: Record<string, number[]>; // Exact matches (lowercase in the actual data)
  c: Record<string, number[]>; // Contained-in matches (lowercase in the actual data)

  // For backward compatibility with our code
  E?: Record<string, number[]>; // Exact matches (uppercase in our code)
  C?: Record<string, number[]>; // Contained-in matches (uppercase in our code)
}

// Dictionary entry interface
export interface DictionaryEntry {
  id: string;
  [key: string]: unknown;
}

/**
 * Checks if a character is a Han (Chinese/Japanese) character
 */
export function isHanCharacter(char: string): boolean {
  const code = char.codePointAt(0);
  if (!code) return false;

  // CJK Unified Ideographs range
  return (
    (code >= 0x4e00 && code <= 0x9fff) || // CJK Unified Ideographs
    (code >= 0x3400 && code <= 0x4dbf) || // CJK Unified Ideographs Extension A
    (code >= 0x20000 && code <= 0x2a6df) || // CJK Unified Ideographs Extension B
    (code >= 0x2a700 && code <= 0x2b73f) || // CJK Unified Ideographs Extension C
    (code >= 0x2b740 && code <= 0x2b81f) || // CJK Unified Ideographs Extension D
    (code >= 0x2b820 && code <= 0x2ceaf) || // CJK Unified Ideographs Extension E
    (code >= 0x2ceb0 && code <= 0x2ebef) || // CJK Unified Ideographs Extension F
    (code >= 0x30000 && code <= 0x3134f) // CJK Unified Ideographs Extension G
  );
}

/**
 * Determines the shard type for a word based on Han character count
 */
export function getShardType(word: string): ShardType {
  // Count Han characters
  let hanCount = 0;
  for (const char of word) {
    if (isHanCharacter(char)) {
      hanCount++;
    }
  }

  // Determine shard type based on Han character count
  if (hanCount === 0) {
    return ShardType.NON_HAN;
  } else if (hanCount === 1) {
    return ShardType.HAN_1CHAR;
  } else if (hanCount === 2) {
    return ShardType.HAN_2CHAR;
  } else {
    return ShardType.HAN_3PLUS;
  }
}

/**
 * Gets the repository name for a shard type
 */
export function getRepoForShardType(shardType: ShardType): string {
  switch (shardType) {
    case ShardType.NON_HAN:
      return REPOS.NON_HAN;
    case ShardType.HAN_1CHAR:
      return REPOS.HAN_1CHAR;
    case ShardType.HAN_2CHAR:
      return REPOS.HAN_2CHAR;
    case ShardType.HAN_3PLUS:
      return REPOS.HAN_3PLUS;
    default:
      return REPOS.NON_HAN;
  }
}

/**
 * Builds the URL for an index file
 */
export function getIndexUrl(word: string, shardType: ShardType): string {
  const repo = getRepoForShardType(shardType);
  return `${BASE_URL}/${repo}/index/${word}.json.br`;
}

/**
 * Builds the URL for a dictionary entry
 */
export function getDictionaryEntryUrl(
  id: string,
  dictType: DictionaryType,
  shardType: ShardType
): string {
  const repo = getRepoForShardType(shardType);
  // The ID already includes the shard type as the first digit
  return `${BASE_URL}/${repo}/${dictType}/${id}.json.br`;
}

/**
 * Extracts the shard type from a sharded ID
 */
export function extractShardType(shardedId: string): ShardType {
  if (shardedId.length === 0) {
    return ShardType.NON_HAN;
  }

  // Get the first digit
  const shardTypeStr = shardedId[0];
  const shardTypeInt = parseInt(shardTypeStr, 10);

  if (isNaN(shardTypeInt)) {
    return ShardType.NON_HAN;
  }

  return shardTypeInt as ShardType;
}
