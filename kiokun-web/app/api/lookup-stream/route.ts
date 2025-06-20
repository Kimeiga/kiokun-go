/**
 * Streaming API route handler for dictionary lookups
 */

import { NextRequest } from "next/server";
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

// Dictionary entry type
type DictionaryEntry = Record<string, unknown>;

// Dictionary entries by type
interface DictionaryEntriesByType {
  j: DictionaryEntry[]; // JMdict (Japanese words)
  n: DictionaryEntry[]; // JMnedict (Japanese names)
  d: DictionaryEntry[]; // Kanjidic (Kanji characters)
  c: DictionaryEntry[]; // Chinese characters
  w: DictionaryEntry[]; // Chinese words
}

// Initial response structure (for streaming)
interface InitialLookupResponse {
  word: string;
  exactMatches: DictionaryEntriesByType;
  containedMatchesPending: boolean;
}

// Contained matches response structure (for streaming)
interface ContainedMatchesResponse {
  containedMatches: DictionaryEntriesByType;
  containedMatchesPending: boolean;
}

/**
 * Fetches dictionary entries based on IDs and dictionary type
 */
async function fetchDictionaryEntries(
  dictType: string,
  ids: number[],
  shardType: ShardType
): Promise<Record<string, unknown>[]> {
  try {
    console.log(
      `Fetching ${ids.length} entries for dictionary type: ${dictType}, shard type: ${shardType}`
    );

    // Add delay between requests to avoid rate limiting
    const delay = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));

    // Process entries sequentially with delay to avoid rate limiting
    const entries: Record<string, unknown>[] = [];

    // Check if this is a Chinese word dictionary (which might exceed jsDelivr's 50MB limit)
    const isChineseWordDict = dictType === DictionaryType.CHINESE_WORDS;
    if (isChineseWordDict) {
      console.log(
        `Processing Chinese word dictionary entries - these may exceed jsDelivr's 50MB limit`
      );
    }

    for (const id of ids) {
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

      console.log(
        `Processing entry ID: ${shardedId}, Dictionary Type: ${dictType}, URL: ${url}`
      );

      try {
        // Add a small delay between requests to avoid rate limiting
        if (entries.length > 0) {
          await delay(100); // 100ms delay between requests
        }

        // Our enhanced fetchAndDecompressJson will return an empty object
        // if jsDelivr returns a 403 error for Chinese word dictionaries
        const entry = await fetchAndDecompressJson<Record<string, unknown>>(
          url
        );

        // Only add non-empty entries
        if (entry && Object.keys(entry).length > 0) {
          entries.push(entry);
        } else if (isChineseWordDict && Object.keys(entry).length === 0) {
          console.log(
            `Skipping empty Chinese word dictionary entry for ${shardedId}`
          );
        }
      } catch (error) {
        console.error(
          `Error fetching entry ${shardedId} (Dict: ${dictType}):`,
          error
        );

        // Only retry if it's not a 403 error for Chinese word dictionaries
        // (since those are already handled in fetchAndDecompressJson)
        if (
          !(
            error instanceof Error &&
            error.message.includes("403") &&
            isChineseWordDict
          )
        ) {
          // Try again with a longer delay
          try {
            console.log(`Retrying ${url} after delay...`);
            await delay(1000); // 1 second delay before retry
            const entry = await fetchAndDecompressJson<Record<string, unknown>>(
              url
            );

            // Only add non-empty entries
            if (entry && Object.keys(entry).length > 0) {
              console.log(`Retry successful for ${url}`);
              entries.push(entry);
            } else if (isChineseWordDict && Object.keys(entry).length === 0) {
              console.log(
                `Skipping empty Chinese word dictionary entry on retry for ${shardedId}`
              );
            }
          } catch (retryError) {
            console.error(`Retry failed for ${url}:`, retryError);
          }
        }
      }
    }

    console.log(
      `Successfully fetched ${entries.length}/${ids.length} entries for dictionary type: ${dictType}`
    );
    return entries;
  } catch (error) {
    console.error(`Error fetching ${dictType} entries:`, error);
    return [];
  }
}

/**
 * GET handler for the streaming lookup API
 */
export async function GET(request: NextRequest) {
  // Get the word from the query parameters
  const searchParams = request.nextUrl.searchParams;
  const word = searchParams.get("word");

  // Create a new ReadableStream
  const stream = new ReadableStream({
    async start(controller) {
      try {
        // Validate the word parameter
        if (!word) {
          const errorResponse = JSON.stringify({
            error: "Word parameter is required",
          });
          controller.enqueue(new TextEncoder().encode(errorResponse));
          controller.close();
          return;
        }

        // Determine the shard type based on the word
        const shardType = getShardType(word);

        // Get the URL for the index file
        const indexUrl = getIndexUrl(word, shardType);

        // Fetch and decompress the index file
        let indexEntry: IndexEntry;
        try {
          indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (error) {
          // If the index file doesn't exist, return an empty response
          const emptyResponse: DictionaryEntriesByType = {
            j: [],
            n: [],
            d: [],
            c: [],
            w: [],
          };

          // Send initial response with empty results
          const initialResponse: InitialLookupResponse = {
            word,
            exactMatches: emptyResponse,
            containedMatchesPending: false,
          };

          controller.enqueue(
            new TextEncoder().encode(JSON.stringify(initialResponse))
          );
          controller.close();
          return;
        }

        // Initialize dictionary entries by type
        const exactMatches: DictionaryEntriesByType = {
          j: [], // JMdict (Japanese words)
          n: [], // JMnedict (Japanese names)
          d: [], // Kanjidic (Kanji characters)
          c: [], // Chinese characters
          w: [], // Chinese words
        };

        // Process exact matches first
        if (indexEntry.e) {
          for (const [dictType, ids] of Object.entries(indexEntry.e)) {
            // Validate that dictType is one of the supported types
            if (dictType in exactMatches) {
              const entries = await fetchDictionaryEntries(
                dictType,
                ids,
                shardType
              );
              // Add entries to the appropriate dictionary type array
              exactMatches[dictType as keyof DictionaryEntriesByType] = entries;
            } else {
              console.warn(
                `Unknown dictionary type in exact matches: ${dictType}`
              );
            }
          }
        }

        // Send initial response with exact matches
        const initialResponse: InitialLookupResponse = {
          word,
          exactMatches,
          containedMatchesPending:
            indexEntry.c && Object.keys(indexEntry.c).length > 0,
        };

        controller.enqueue(
          new TextEncoder().encode(JSON.stringify(initialResponse))
        );

        // Process contained-in matches
        const containedMatches: DictionaryEntriesByType = {
          j: [], // JMdict (Japanese words)
          n: [], // JMnedict (Japanese names)
          d: [], // Kanjidic (Kanji characters)
          c: [], // Chinese characters
          w: [], // Chinese words
        };

        if (indexEntry.c) {
          for (const [dictType, ids] of Object.entries(indexEntry.c)) {
            // Validate that dictType is one of the supported types
            if (dictType in containedMatches) {
              const entries = await fetchDictionaryEntries(
                dictType,
                ids,
                shardType
              );
              // Add entries to the appropriate dictionary type array
              containedMatches[dictType as keyof DictionaryEntriesByType] =
                entries;
            } else {
              console.warn(
                `Unknown dictionary type in contained matches: ${dictType}`
              );
            }
          }
        }

        // Send contained matches response
        const containedResponse: ContainedMatchesResponse = {
          containedMatches,
          containedMatchesPending: false,
        };

        controller.enqueue(
          new TextEncoder().encode(JSON.stringify(containedResponse))
        );
        controller.close();
      } catch (error) {
        console.error("Error processing streaming lookup request:", error);

        // Create empty response with the new structure
        const emptyResponse: DictionaryEntriesByType = {
          j: [],
          n: [],
          d: [],
          c: [],
          w: [],
        };

        // Send error response
        const errorResponse = {
          word,
          exactMatches: emptyResponse,
          containedMatches: emptyResponse,
          error: "Failed to process lookup request",
          containedMatchesPending: false,
        };

        controller.enqueue(
          new TextEncoder().encode(JSON.stringify(errorResponse))
        );
        controller.close();
      }
    },
  });

  // Return the stream as a response
  return new Response(stream, {
    headers: {
      "Content-Type": "application/json",
      "Transfer-Encoding": "chunked",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
    },
  });
}
