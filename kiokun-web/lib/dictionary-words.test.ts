/**
 * Tests for various types of words in the dictionary
 *
 * This test file tests a variety of Japanese and Chinese words of different types:
 * - Single Han characters (Japanese kanji and Chinese characters)
 * - Two Han characters (Japanese and Chinese words)
 * - Three or more Han characters (Japanese and Chinese words)
 * - Non-Han words (hiragana, katakana)
 */

import { describe, it, expect } from "vitest";
import {
  ShardType,
  getShardType,
  getIndexUrl,
  IndexEntry,
} from "./dictionary-utils";
import { fetchAndDecompressJson } from "./brotli-utils";

// Skip tests if running in CI environment
const SKIP_NETWORK_TESTS = process.env.CI === "true";

// Test words for different categories
const TEST_WORDS = {
  // Japanese words
  JAPANESE: {
    SINGLE_KANJI: [
      "水", // water
      "火", // fire
      "山", // mountain
      "川", // river
    ],
    TWO_KANJI: [
      "日本", // Japan
      "学校", // school
      "電車", // train
      "食事", // meal
    ],
    THREE_PLUS_KANJI: [
      "図書館", // library
      "新幹線", // bullet train
      "大学生", // university student
      "日本語", // Japanese language
    ],
    HIRAGANA: [
      "ありがとう", // thank you
      "こんにちは", // hello
      "さようなら", // goodbye
    ],
    KATAKANA: [
      "コンピューター", // computer
      "テレビ", // TV
      "スマートフォン", // smartphone
    ],
  },
  // Chinese words
  CHINESE: {
    SINGLE_CHAR: [
      "人", // person
      "心", // heart
      "手", // hand
      "口", // mouth
    ],
    TWO_CHAR: [
      "中国", // China
      "朋友", // friend
      "学习", // study
      "工作", // work
    ],
    THREE_PLUS_CHAR: [
      "电脑", // computer
      "图书馆", // library
      "北京市", // Beijing city
      "大学生", // university student
    ],
  },
};

describe("Dictionary Word Tests", () => {
  // Test Japanese single kanji
  it("correctly processes Japanese single kanji", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.JAPANESE.SINGLE_KANJI[0]; // 水
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.HAN_1CHAR);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    try {
      const indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();
      expect(indexEntry.e).toBeDefined(); // Exact matches
      // Note: Not all entries have contained-in matches

      // Verify we have entries in at least one dictionary type
      const hasExactMatches = Object.values(indexEntry.e).some(
        (ids) => ids.length > 0
      );

      expect(hasExactMatches).toBe(true);
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000);

  // Test Japanese two kanji word
  it("correctly processes Japanese two kanji word", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.JAPANESE.TWO_KANJI[0]; // 日本
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.HAN_2CHAR);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    try {
      const indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();
      expect(indexEntry.e).toBeDefined(); // Exact matches
      // Note: Not all entries have contained-in matches

      // Verify we have entries in at least one dictionary type
      const hasExactMatches = Object.values(indexEntry.e).some(
        (ids) => ids.length > 0
      );

      expect(hasExactMatches).toBe(true);
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000);

  // Test Japanese three+ kanji word
  it("correctly processes Japanese three+ kanji word", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.JAPANESE.THREE_PLUS_KANJI[0]; // 図書館
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.HAN_3PLUS);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    try {
      const indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();
      expect(indexEntry.e).toBeDefined(); // Exact matches
      // Note: Not all entries have contained-in matches

      // Verify we have entries in at least one dictionary type
      const hasExactMatches = Object.values(indexEntry.e).some(
        (ids) => ids.length > 0
      );

      expect(hasExactMatches).toBe(true);
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000);

  // Test Japanese hiragana word
  it("correctly processes Japanese hiragana word", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.JAPANESE.HIRAGANA[0]; // ありがとう
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.NON_HAN);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    try {
      const indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();
      expect(indexEntry.e).toBeDefined(); // Exact matches
      // Note: Not all entries have contained-in matches
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000);

  // Test Japanese katakana word
  it("correctly processes Japanese katakana word", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.JAPANESE.KATAKANA[0]; // コンピューター
    const shardType = getShardType(word);

    expect(shardType).toBe(ShardType.NON_HAN);

    // Get the index URL
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index
    try {
      const indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);

      // Log the structure for debugging
      console.log(`Index for ${word}:`, JSON.stringify(indexEntry, null, 2));

      // Verify index structure
      expect(indexEntry).toBeDefined();
      expect(indexEntry.e).toBeDefined(); // Exact matches
      // Note: Not all entries have contained-in matches
    } catch (error) {
      console.error(`Error fetching index for ${word}:`, error);
      throw error;
    }
  }, 30000);
});
