/**
 * Integration test for dictionary lookup functionality
 *
 * This test verifies the entire dictionary lookup flow:
 * 1. Determining the correct shard based on Han character count
 * 2. Fetching the index file from the appropriate repository
 * 3. Following links to get dictionary entries
 * 4. Verifying both exact matches (E) and contained-in matches (C)
 */

import { describe, it, expect } from "vitest";
import {
  getShardType,
  getIndexUrl,
  getDictionaryEntryUrl,
  ShardType,
  DictionaryType,
  IndexEntry,
  extractShardType,
} from "./dictionary-utils";
import { fetchAndDecompressJson } from "./brotli-utils";

// Test words for different shard types
const TEST_WORDS = {
  NON_HAN: "コンピューター", // Non-Han word (computer in katakana)
  HAN_1CHAR: "水", // Single Han character (water)
  HAN_2CHAR: "日本", // Two Han characters (Japan)
  HAN_3PLUS: "図書館", // Three Han characters (library)
};

// Skip tests if running in CI environment
const SKIP_NETWORK_TESTS = process.env.CI === "true";

describe("Dictionary Lookup Integration", () => {
  // Test shard type determination
  it("correctly determines shard type based on Han character count", () => {
    expect(getShardType(TEST_WORDS.NON_HAN)).toBe(ShardType.NON_HAN);
    expect(getShardType(TEST_WORDS.HAN_1CHAR)).toBe(ShardType.HAN_1CHAR);
    expect(getShardType(TEST_WORDS.HAN_2CHAR)).toBe(ShardType.HAN_2CHAR);
    expect(getShardType(TEST_WORDS.HAN_3PLUS)).toBe(ShardType.HAN_3PLUS);
  });

  // Test URL generation
  it("generates correct index URLs", () => {
    const nonHanUrl = getIndexUrl(TEST_WORDS.NON_HAN, ShardType.NON_HAN);
    const han1CharUrl = getIndexUrl(TEST_WORDS.HAN_1CHAR, ShardType.HAN_1CHAR);
    const han2CharUrl = getIndexUrl(TEST_WORDS.HAN_2CHAR, ShardType.HAN_2CHAR);
    const han3PlusUrl = getIndexUrl(TEST_WORDS.HAN_3PLUS, ShardType.HAN_3PLUS);

    expect(nonHanUrl).toContain(
      "japanese-dict-non-han/index/コンピューター.json.br"
    );
    expect(han1CharUrl).toContain("japanese-dict-han-1char/index/水.json.br");
    expect(han2CharUrl).toContain("japanese-dict-han-2char/index/日本.json.br");
    expect(han3PlusUrl).toContain(
      "japanese-dict-han-3plus/index/図書館.json.br"
    );
  });

  // Test dictionary entry URL generation
  it("generates correct dictionary entry URLs", () => {
    const jmdictUrl = getDictionaryEntryUrl(
      "1123456",
      DictionaryType.JMDICT,
      ShardType.HAN_1CHAR
    );
    const jmnedictUrl = getDictionaryEntryUrl(
      "1789012",
      DictionaryType.JMNEDICT,
      ShardType.HAN_2CHAR
    );
    const kanjidicUrl = getDictionaryEntryUrl(
      "1345678",
      DictionaryType.KANJIDIC,
      ShardType.HAN_1CHAR
    );
    const chineseCharsUrl = getDictionaryEntryUrl(
      "1901234",
      DictionaryType.CHINESE_CHARS,
      ShardType.HAN_1CHAR
    );
    const chineseWordsUrl = getDictionaryEntryUrl(
      "2567890",
      DictionaryType.CHINESE_WORDS,
      ShardType.HAN_3PLUS
    );

    expect(jmdictUrl).toContain("japanese-dict-han-1char/j/1123456.json.br");
    expect(jmnedictUrl).toContain("japanese-dict-han-2char/n/1789012.json.br");
    expect(kanjidicUrl).toContain("japanese-dict-han-1char/d/1345678.json.br");
    expect(chineseCharsUrl).toContain(
      "japanese-dict-han-1char/c/1901234.json.br"
    );
    expect(chineseWordsUrl).toContain(
      "japanese-dict-han-3plus/w/2567890.json.br"
    );
  });

  // Test extracting shard type from ID
  it("correctly extracts shard type from sharded ID", () => {
    expect(extractShardType("0123456")).toBe(ShardType.NON_HAN);
    expect(extractShardType("1123456")).toBe(ShardType.HAN_1CHAR);
    expect(extractShardType("2123456")).toBe(ShardType.HAN_2CHAR);
    expect(extractShardType("3123456")).toBe(ShardType.HAN_3PLUS);
  });

  // Integration test for single Han character (水)
  it("fetches and processes index for single Han character", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.HAN_1CHAR; // 水 (water)
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.HAN_1CHAR);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    let indexEntry: IndexEntry;
    try {
      indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the actual data
      console.log("Index entry:", JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify we have entries in at least one dictionary type
      const hasExactMatches = Object.values(indexEntry.e).some(
        (ids) => ids.length > 0
      );
      const hasContainedMatches = Object.values(indexEntry.c).some(
        (ids) => ids.length > 0
      );

      expect(hasExactMatches || hasContainedMatches).toBe(true);

      // Fetch one entry from each dictionary type that has entries
      for (const [dictType, ids] of Object.entries(indexEntry.e)) {
        if (ids.length > 0) {
          const id = ids[0].toString();
          const shardedId = id.startsWith(shardType.toString())
            ? id
            : `${shardType}${id}`;
          const entryUrl = getDictionaryEntryUrl(
            shardedId,
            dictType as DictionaryType,
            extractShardType(shardedId)
          );

          try {
            const entry = await fetchAndDecompressJson(entryUrl);
            console.log(`Entry for ${dictType} ${shardedId}:`, entry);
            expect(entry).toBeDefined();
          } catch (error) {
            console.error(
              `Error fetching entry for ${dictType} ${shardedId}:`,
              error
            );
          }
        }
      }
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000); // Increase timeout to 30 seconds for network requests

  // Integration test for two Han characters (日本)
  it("fetches and processes index for two Han characters", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.HAN_2CHAR; // 日本 (Japan)
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.HAN_2CHAR);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    let indexEntry: IndexEntry;
    try {
      indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the raw data
      console.log("Raw index data:", JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify we have entries in at least one dictionary type
      const hasExactMatches = indexEntry.e
        ? Object.values(indexEntry.e).some((ids) => ids.length > 0)
        : false;
      const hasContainedMatches = indexEntry.c
        ? Object.values(indexEntry.c).some((ids) => ids.length > 0)
        : false;

      expect(hasExactMatches || hasContainedMatches).toBe(true);

      // Fetch one entry from each dictionary type that has entries
      if (indexEntry.e) {
        for (const [dictType, ids] of Object.entries(indexEntry.e)) {
          if (ids.length > 0) {
            const id = ids[0].toString();
            const shardedId = id.startsWith(shardType.toString())
              ? id
              : `${shardType}${id}`;
            const entryUrl = getDictionaryEntryUrl(
              shardedId,
              dictType as DictionaryType,
              extractShardType(shardedId)
            );

            try {
              const entry = await fetchAndDecompressJson(entryUrl);
              console.log(`Entry for ${dictType} ${shardedId}:`, entry);
              expect(entry).toBeDefined();
            } catch (error) {
              console.error(
                `Error fetching entry for ${dictType} ${shardedId}:`,
                error
              );
            }
          }
        }
      }
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000); // Increase timeout to 30 seconds for network requests
});
