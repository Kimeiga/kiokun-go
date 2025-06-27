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
  isHanCharacter,
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

interface StreamingEntryResponse {
  type: 'entry';
  dictType: string;
  entry: Record<string, unknown>;
  isExactMatch: boolean;
}

/**
 * Fetches dictionary entries in parallel and streams results as they come back
 */
async function fetchDictionaryEntriesStreaming(
  dictType: string,
  ids: number[],
  shardType: ShardType,
  onEntryReady: (dictType: string, entry: Record<string, unknown>) => void
): Promise<Record<string, unknown>[]> {
  try {
    console.log(
      `🚀 Starting parallel fetch of ${ids.length} entries for dictionary type: ${dictType}, shard type: ${shardType}`
    );

    const startTime = Date.now();

    // Create parallel promises for all entries - NO MORE DELAYS!
    const promises = ids.map(async (id, index) => {
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
        console.log(
          `📥 Fetching entry ${index + 1}/${ids.length}: ${shardedId} (${dictType})`
        );

        const entry = await fetchAndDecompressJson<Record<string, unknown>>(url);

        // Stream the entry immediately when it's ready
        if (entry && Object.keys(entry).length > 0) {
          console.log(`✅ Entry ready: ${shardedId} (${dictType}) - streaming immediately`);
          onEntryReady(dictType, entry);
          return entry;
        } else {
          console.log(`⚠️ Empty entry skipped: ${shardedId} (${dictType})`);
          return null;
        }
      } catch (error) {
        console.error(`❌ Failed to fetch entry ${shardedId} (${dictType}):`, error);
        return null;
      }
    });

    // Wait for all requests to complete
    const results = await Promise.all(promises);
    const validEntries = results.filter(entry => entry !== null);

    console.log(
      `🎉 Completed ${dictType}: ${validEntries.length}/${ids.length} entries in ${Date.now() - startTime}ms`
    );

    return validEntries;
  } catch (error) {
    console.error(`💥 Error in fetchDictionaryEntriesStreaming for ${dictType}:`, error);
    return [];
  }
}

/**
 * GET handler for the streaming lookup API
 */
export async function GET(request: NextRequest) {
  const startTime = Date.now();
  // Get the word from the query parameters
  const searchParams = request.nextUrl.searchParams;
  const word = searchParams.get("word");

  console.log(`🚀 Starting lookup-stream for word: ${word}`);

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
        console.log(`📊 Shard type determined: ${shardType} (${Date.now() - startTime}ms)`);

        // Get the URL for the index file
        const indexUrl = getIndexUrl(word, shardType);

        // Fetch and decompress the index file
        let indexEntry: IndexEntry;
        try {
          const indexStartTime = Date.now();
          indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);
          console.log(`📥 Index fetched in ${Date.now() - indexStartTime}ms`);
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

        // Process exact matches with real-time streaming
        if (indexEntry.e) {
          console.log(`🎯 Processing exact matches for ${Object.keys(indexEntry.e).length} dictionary types`);

          // Create streaming callback for exact matches
          const streamExactMatch = (dictType: string, entry: Record<string, unknown>) => {
            // Add to exact matches collection
            if (dictType in exactMatches) {
              exactMatches[dictType as keyof DictionaryEntriesByType].push(entry);
            }

            // Stream the entry immediately
            const streamResponse: StreamingEntryResponse = {
              type: 'entry',
              dictType,
              entry,
              isExactMatch: true
            };

            controller.enqueue(
              new TextEncoder().encode(JSON.stringify(streamResponse) + '\n')
            );
            console.log(`📤 Streamed exact match: ${dictType}`);
          };

          // Start parallel fetching for all exact match dictionary types
          const exactMatchPromises = Object.entries(indexEntry.e).map(async ([dictType, ids]) => {
            if (dictType in exactMatches) {
              return fetchDictionaryEntriesStreaming(
                dictType,
                ids,
                shardType,
                streamExactMatch
              );
            }
            return [];
          });

          // Wait for all exact matches to complete
          await Promise.all(exactMatchPromises);
          console.log(`✅ All exact matches completed`);
        }

        // Check if we need to search across multiple shards for contained matches
        const isSingleCharacter = word.length === 1 && isHanCharacter(word);
        let hasContainedMatches = indexEntry.c && Object.keys(indexEntry.c).length > 0;

        // For single characters, we might have contained matches in other shards
        if (isSingleCharacter && !hasContainedMatches) {
          // Check if there might be contained matches in other shards
          const allShardTypes = [ShardType.HAN_1CHAR, ShardType.HAN_2CHAR, ShardType.HAN_3PLUS, ShardType.NON_HAN];
          for (const searchShardType of allShardTypes) {
            if (searchShardType !== shardType) {
              try {
                const searchIndexUrl = getIndexUrl(word, searchShardType);
                const searchIndexEntry = await fetchAndDecompressJson<IndexEntry>(searchIndexUrl);
                if (searchIndexEntry.c && Object.keys(searchIndexEntry.c).length > 0) {
                  hasContainedMatches = true;
                  break;
                }
              } catch (error) {
                // Index file doesn't exist in this shard, continue
              }
            }
          }
        }

        // Send initial response with exact matches
        const initialResponse: InitialLookupResponse = {
          word,
          exactMatches,
          containedMatchesPending: hasContainedMatches,
        };

        controller.enqueue(
          new TextEncoder().encode(JSON.stringify(initialResponse))
        );

        // Process contained-in matches
        console.log(`🔍 Starting contained matches processing (${Date.now() - startTime}ms total)`);
        const containedMatches: DictionaryEntriesByType = {
          j: [], // JMdict (Japanese words)
          n: [], // JMnedict (Japanese names)
          d: [], // Kanjidic (Kanji characters)
          c: [], // Chinese characters
          w: [], // Chinese words
        };

        if (isSingleCharacter) {
          // Search for contained matches across all shards with round-robin pagination
          const allShardTypes = [ShardType.HAN_1CHAR, ShardType.HAN_2CHAR, ShardType.HAN_3PLUS, ShardType.NON_HAN];

          // Collect all contained match IDs first
          const allContainedIds: { [dictType: string]: { ids: number[], shardType: ShardType }[] } = {};

          for (const searchShardType of allShardTypes) {
            try {
              const searchIndexUrl = getIndexUrl(word, searchShardType);
              const searchIndexEntry = await fetchAndDecompressJson<IndexEntry>(searchIndexUrl);

              if (searchIndexEntry.c) {
                for (const [dictType, ids] of Object.entries(searchIndexEntry.c)) {
                  if (dictType in containedMatches) {
                    if (!allContainedIds[dictType]) {
                      allContainedIds[dictType] = [];
                    }
                    allContainedIds[dictType].push({ ids, shardType: searchShardType });
                  }
                }
              }
            } catch (error) {
              // Index file doesn't exist in this shard, continue to next shard
              console.log(`No index file for ${word} in shard ${searchShardType}`);
            }
          }

          // Round-robin pagination: take 1 from each dict type, repeat up to 15 times (reduced for faster initial load)
          const maxRounds = 15;
          const dictTypes = Object.keys(allContainedIds);

          for (let round = 0; round < maxRounds; round++) {
            let addedInThisRound = false;

            for (const dictType of dictTypes) {
              const shardGroups = allContainedIds[dictType];
              if (!shardGroups) continue;

              // Find the next available ID across all shards for this dict type
              let totalProcessed = containedMatches[dictType as keyof DictionaryEntriesByType].length;
              let currentIndex = totalProcessed;

              // Find which shard and index within that shard
              let cumulativeCount = 0;
              for (const shardGroup of shardGroups) {
                if (currentIndex < cumulativeCount + shardGroup.ids.length) {
                  const indexInShard = currentIndex - cumulativeCount;
                  const idToFetch = shardGroup.ids[indexInShard];

                  try {
                    // Create streaming callback for contained matches
                    const streamContainedMatch = (dictType: string, entry: Record<string, unknown>) => {
                      // Add to contained matches collection
                      if (dictType in containedMatches) {
                        containedMatches[dictType as keyof DictionaryEntriesByType].push(entry);
                      }

                      // Stream the entry immediately
                      const streamResponse: StreamingEntryResponse = {
                        type: 'entry',
                        dictType,
                        entry,
                        isExactMatch: false
                      };

                      controller.enqueue(
                        new TextEncoder().encode(JSON.stringify(streamResponse) + '\n')
                      );
                      console.log(`📤 Streamed contained match: ${dictType}`);
                    };

                    const entries = await fetchDictionaryEntriesStreaming(
                      dictType,
                      [idToFetch],
                      shardGroup.shardType,
                      streamContainedMatch
                    );
                    addedInThisRound = entries.length > 0;
                  } catch (error) {
                    console.warn(`Error fetching entry ${idToFetch} from ${dictType}:`, error);
                  }
                  break;
                }
                cumulativeCount += shardGroup.ids.length;
              }
            }

            // If no entries were added in this round, we've exhausted all available entries
            if (!addedInThisRound) {
              break;
            }
          }
        } else {
          // For multi-character words, only search in the primary shard
          if (indexEntry.c) {
            for (const [dictType, ids] of Object.entries(indexEntry.c)) {
              // Validate that dictType is one of the supported types
              if (dictType in containedMatches) {
                // Limit to first 15 entries for multi-character words (reduced for faster load)
                const limitedIds = ids.slice(0, 15);

                // Create streaming callback for multi-char contained matches
                const streamMultiCharMatch = (dictType: string, entry: Record<string, unknown>) => {
                  // Add to contained matches collection
                  if (dictType in containedMatches) {
                    containedMatches[dictType as keyof DictionaryEntriesByType].push(entry);
                  }

                  // Stream the entry immediately
                  const streamResponse: StreamingEntryResponse = {
                    type: 'entry',
                    dictType,
                    entry,
                    isExactMatch: false
                  };

                  controller.enqueue(
                    new TextEncoder().encode(JSON.stringify(streamResponse) + '\n')
                  );
                  console.log(`📤 Streamed multi-char contained match: ${dictType}`);
                };

                await fetchDictionaryEntriesStreaming(
                  dictType,
                  limitedIds,
                  shardType,
                  streamMultiCharMatch
                );
              } else {
                console.warn(
                  `Unknown dictionary type in contained matches: ${dictType}`
                );
              }
            }
          }
        }

        // Send contained matches response
        console.log(`✅ Lookup-stream completed in ${Date.now() - startTime}ms total`);
        const containedResponse: ContainedMatchesResponse = {
          containedMatches,
          containedMatchesPending: false,
        };

        controller.enqueue(
          new TextEncoder().encode('\n' + JSON.stringify(containedResponse))
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
