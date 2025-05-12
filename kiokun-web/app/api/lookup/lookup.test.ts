/**
 * Integration test for the dictionary lookup API
 *
 * This test verifies the entire API flow with real dictionary data:
 * 1. Making requests to the API with different words
 * 2. Verifying the response structure
 * 3. Checking that exact matches (E) and contained-in matches (C) are correctly processed
 */

import { describe, it, expect, beforeAll } from "vitest";
import { NextRequest } from "next/server";
import { GET } from "./route";
import { ShardType, getShardType } from "@/lib/dictionary-utils";

// Test words for different shard types
const TEST_WORDS = {
  NON_HAN: "コンピューター", // Non-Han word (computer in katakana)
  HAN_1CHAR: "水", // Single Han character (water)
  HAN_2CHAR: "日本", // Two Han characters (Japan)
  HAN_3PLUS: "図書館", // Three Han characters (library)
};

// Skip tests if running in CI environment
const SKIP_NETWORK_TESTS = process.env.CI === "true";

describe("Dictionary Lookup API Integration", () => {
  // Test with a single Han character (水)
  it("returns correct results for a single Han character", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.HAN_1CHAR; // 水 (water)
    const request = new NextRequest(
      `http://localhost:3000/api/lookup?word=${encodeURIComponent(word)}`
    );

    // Call the API handler
    const response = await GET(request);
    expect(response.status).toBe(200);

    // Parse the response
    const data = await response.json();

    // Verify the response structure
    expect(data).toBeDefined();
    expect(data.word).toBe(word);
    expect(data.exactMatches).toBeDefined();
    expect(data.containedMatches).toBeDefined();

    // Log the response for debugging
    console.log(
      `API response for ${word}:`,
      JSON.stringify(
        {
          word: data.word,
          exactMatchCount: data.exactMatches.length,
          containedMatchCount: data.containedMatches.length,
          exactMatchTypes: data.exactMatches.map((entry: any) => {
            if (entry.Kanji && entry.Kana && entry.Sense) return "JMdict";
            if (entry.Kanji && entry.Reading && entry.Translation)
              return "JMnedict";
            if (entry.Character && entry.Reading && entry.Misc)
              return "Kanjidic";
            if (entry.Traditional && entry.Pinyin && !entry.Components)
              return "ChineseWord";
            if (entry.Traditional && entry.Pinyin) return "ChineseChar";
            return "Unknown";
          }),
          containedMatchTypes: data.containedMatches.map((entry: any) => {
            if (entry.Kanji && entry.Kana && entry.Sense) return "JMdict";
            if (entry.Kanji && entry.Reading && entry.Translation)
              return "JMnedict";
            if (entry.Character && entry.Reading && entry.Misc)
              return "Kanjidic";
            if (entry.Traditional && entry.Pinyin && !entry.Components)
              return "ChineseWord";
            if (entry.Traditional && entry.Pinyin) return "ChineseChar";
            return "Unknown";
          }),
        },
        null,
        2
      )
    );

    // Verify we have at least some results
    // For 水 (water), we should have at least some exact matches
    expect(data.exactMatches.length).toBeGreaterThan(0);

    // Check that the shard type is correct
    const shardType = getShardType(word);
    expect(shardType).toBe(ShardType.HAN_1CHAR);

    // Verify the structure of at least one exact match
    if (data.exactMatches.length > 0) {
      const entry = data.exactMatches[0];
      expect(entry).toBeDefined();

      // The entry should have some properties depending on its type
      if (entry.Kanji && entry.Kana && entry.Sense) {
        // JMdict entry
        expect(entry.Kanji).toBeInstanceOf(Array);
        expect(entry.Kana).toBeInstanceOf(Array);
        expect(entry.Sense).toBeInstanceOf(Array);
      } else if (entry.Kanji && entry.Reading && entry.Translation) {
        // JMnedict entry
        expect(entry.Kanji).toBeInstanceOf(Array);
        expect(entry.Reading).toBeInstanceOf(Array);
        expect(entry.Translation).toBeInstanceOf(Array);
      } else if (entry.Character && entry.Reading && entry.Misc) {
        // Kanjidic entry
        expect(entry.Character).toBeDefined();
        expect(entry.Reading).toBeDefined();
        expect(entry.Misc).toBeDefined();
      } else if (entry.Traditional && entry.Pinyin) {
        // Chinese character or word entry
        expect(entry.Traditional).toBeDefined();
        expect(entry.Pinyin).toBeDefined();
        if (entry.Definitions) {
          expect(entry.Definitions).toBeInstanceOf(Array);
        }
      }
    }
  }, 30000); // Increase timeout to 30 seconds for network requests

  // Test with two Han characters (日本)
  it("returns correct results for two Han characters", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.HAN_2CHAR; // 日本 (Japan)
    const request = new NextRequest(
      `http://localhost:3000/api/lookup?word=${encodeURIComponent(word)}`
    );

    // Call the API handler
    const response = await GET(request);
    expect(response.status).toBe(200);

    // Parse the response
    const data = await response.json();

    // Verify the response structure
    expect(data).toBeDefined();
    expect(data.word).toBe(word);
    expect(data.exactMatches).toBeDefined();
    expect(data.containedMatches).toBeDefined();

    // Log the response for debugging
    console.log(
      `API response for ${word}:`,
      JSON.stringify(
        {
          word: data.word,
          exactMatchCount: data.exactMatches.length,
          containedMatchCount: data.containedMatches.length,
          exactMatchTypes: data.exactMatches.map((entry: any) => {
            if (entry.Kanji && entry.Kana && entry.Sense) return "JMdict";
            if (entry.Kanji && entry.Reading && entry.Translation)
              return "JMnedict";
            if (entry.Character && entry.Reading && entry.Misc)
              return "Kanjidic";
            if (entry.Traditional && entry.Pinyin && !entry.Components)
              return "ChineseWord";
            if (entry.Traditional && entry.Pinyin) return "ChineseChar";
            return "Unknown";
          }),
          containedMatchTypes: data.containedMatches.map((entry: any) => {
            if (entry.Kanji && entry.Kana && entry.Sense) return "JMdict";
            if (entry.Kanji && entry.Reading && entry.Translation)
              return "JMnedict";
            if (entry.Character && entry.Reading && entry.Misc)
              return "Kanjidic";
            if (entry.Traditional && entry.Pinyin && !entry.Components)
              return "ChineseWord";
            if (entry.Traditional && entry.Pinyin) return "ChineseChar";
            return "Unknown";
          }),
        },
        null,
        2
      )
    );

    // Verify we have at least some results
    // For 日本 (Japan), we should have at least some exact matches
    expect(data.exactMatches.length).toBeGreaterThan(0);

    // Check that the shard type is correct
    const shardType = getShardType(word);
    expect(shardType).toBe(ShardType.HAN_2CHAR);
  }, 30000); // Increase timeout to 30 seconds for network requests

  // Test with a non-Han word (computer)
  it("returns correct results for a non-Han word", async () => {
    if (SKIP_NETWORK_TESTS) {
      console.log("Skipping network test in CI environment");
      return;
    }

    const word = TEST_WORDS.NON_HAN; // computer
    const request = new NextRequest(
      `http://localhost:3000/api/lookup?word=${encodeURIComponent(word)}`
    );

    // Call the API handler
    const response = await GET(request);
    expect(response.status).toBe(200);

    // Parse the response
    const data = await response.json();

    // Verify the response structure
    expect(data).toBeDefined();
    expect(data.word).toBe(word);
    expect(data.exactMatches).toBeDefined();
    expect(data.containedMatches).toBeDefined();

    // Log the response for debugging
    console.log(
      `API response for ${word}:`,
      JSON.stringify(
        {
          word: data.word,
          exactMatchCount: data.exactMatches.length,
          containedMatchCount: data.containedMatches.length,
        },
        null,
        2
      )
    );

    // Check that the shard type is correct
    const shardType = getShardType(word);
    expect(shardType).toBe(ShardType.NON_HAN);
  }, 30000); // Increase timeout to 30 seconds for network requests
});
