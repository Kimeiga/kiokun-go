/**
 * API route handler for dictionary lookups
 */

import { NextRequest, NextResponse } from "next/server";
import {
  getShardType,
  getIndexUrl,
  getDictionaryEntryUrl,
  IndexEntry,
  DictionaryType,
  ShardType,
  extractShardType,
} from "@/lib/dictionary-utils";
import { fetchAndDecompressJson } from "@/lib/brotli-utils";

// Response structure
interface LookupResponse {
  word: string;
  exactMatches: any[];
  containedMatches: any[];
  error?: string;
}

/**
 * Fetches dictionary entries based on IDs and dictionary type
 */
async function fetchDictionaryEntries(
  dictType: string,
  ids: number[],
  shardType: ShardType
): Promise<any[]> {
  try {
    // Fetch entries in parallel
    const entryPromises = ids.map(async (id) => {
      // Convert ID to string and prepend shard type if not already included
      const idStr = id.toString();
      const shardedId = idStr.startsWith(shardType.toString())
        ? idStr
        : `${shardType}${idStr}`;

      const url = getDictionaryEntryUrl(
        shardedId,
        dictType as DictionaryType,
        extractShardType(shardedId)
      );

      try {
        return await fetchAndDecompressJson(url);
      } catch (error) {
        console.error(`Error fetching entry ${shardedId}:`, error);
        return null;
      }
    });

    // Wait for all entries to be fetched
    const entries = await Promise.all(entryPromises);

    // Filter out null entries (failed fetches)
    return entries.filter((entry) => entry !== null);
  } catch (error) {
    console.error(`Error fetching ${dictType} entries:`, error);
    return [];
  }
}

/**
 * GET handler for the lookup API
 */
export async function GET(request: NextRequest): Promise<NextResponse> {
  // Get the word from the query parameters
  const searchParams = request.nextUrl.searchParams;
  const word = searchParams.get("word");

  // Validate the word parameter
  if (!word) {
    return NextResponse.json(
      { error: "Word parameter is required" },
      { status: 400 }
    );
  }

  try {
    // Determine the shard type based on the word
    const shardType = getShardType(word);

    // Get the URL for the index file
    const indexUrl = getIndexUrl(word, shardType);

    // Fetch and decompress the index file
    let indexEntry: IndexEntry;
    try {
      indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);
    } catch (error) {
      // If the index file doesn't exist, return an empty response
      return NextResponse.json({
        word,
        exactMatches: [],
        containedMatches: [],
      });
    }

    // Initialize arrays for exact and contained matches
    const exactMatches: any[] = [];
    const containedMatches: any[] = [];

    // Process exact matches
    if (indexEntry.e) {
      for (const [dictType, ids] of Object.entries(indexEntry.e)) {
        const entries = await fetchDictionaryEntries(dictType, ids, shardType);
        exactMatches.push(...entries);
      }
    }

    // Process contained-in matches
    if (indexEntry.c) {
      for (const [dictType, ids] of Object.entries(indexEntry.c)) {
        const entries = await fetchDictionaryEntries(dictType, ids, shardType);
        containedMatches.push(...entries);
      }
    }

    // Return the response
    return NextResponse.json({
      word,
      exactMatches,
      containedMatches,
    });
  } catch (error) {
    console.error("Error processing lookup request:", error);

    return NextResponse.json(
      {
        word,
        exactMatches: [],
        containedMatches: [],
        error: "Failed to process lookup request",
      },
      { status: 500 }
    );
  }
}
