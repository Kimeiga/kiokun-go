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

// Dictionary entry type - removed unused type

// Dictionary entries by type - removed unused interface

/**
 * Fetches dictionary entries and streams them as they become available
 */
async function fetchDictionaryEntriesStreaming(
  dictType: string,
  ids: number[],
  shardType: ShardType,
  writer: WritableStreamDefaultWriter<Uint8Array>
): Promise<void> {
  try {
    console.log(`Fetching ${ids.length} entries for dictionary type: ${dictType}, shard type: ${shardType}`);

    // Process entries in parallel
    const promises = ids.map(async (id) => {
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
        const entry = await fetchAndDecompressJson<Record<string, unknown>>(url);
        
        if (entry && Object.keys(entry).length > 0) {
          // Stream the entry immediately
          const streamData = {
            type: 'entry',
            dictType,
            entry,
            isExactMatch: true // We'll determine this in the caller
          };
          
          const chunk = JSON.stringify(streamData) + '\n';
          await writer.write(new TextEncoder().encode(chunk));
          
          return entry;
        }
        return null;
      } catch (error) {
        console.error(`Error fetching entry ${shardedId}:`, error);
        return null;
      }
    });

    await Promise.all(promises);
  } catch (error) {
    console.error(`Error fetching ${dictType} entries:`, error);
  }
}

/**
 * GET handler for the streaming lookup API
 */
export async function GET(request: NextRequest): Promise<Response> {
  const searchParams = request.nextUrl.searchParams;
  const word = searchParams.get("word");

  if (!word) {
    return new Response(
      JSON.stringify({ error: "Word parameter is required" }),
      { status: 400, headers: { "Content-Type": "application/json" } }
    );
  }

  // Create a readable stream
  const stream = new ReadableStream({
    async start(controller) {
      const writer = controller;
      
      try {
        // Determine the shard type based on the word
        const shardType = getShardType(word);
        const indexUrl = getIndexUrl(word, shardType);

        // Fetch the index file for exact matches
        let indexEntry: IndexEntry;
        try {
          indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);
        } catch {
          // If no index file, send empty response and close
          const emptyResponse = {
            exactMatches: { j: [], n: [], d: [], c: [], w: [] },
            containedMatches: { j: [], n: [], d: [], c: [], w: [] },
            containedMatchesPending: false
          };
          
          const chunk = JSON.stringify(emptyResponse) + '\n';
          writer.enqueue(new TextEncoder().encode(chunk));
          writer.close();
          return;
        }

        // Send initial exact matches
        if (indexEntry.e) {
          for (const [dictType, ids] of Object.entries(indexEntry.e)) {
            if (['j', 'n', 'd', 'c', 'w'].includes(dictType)) {
              // Limit to first 10 entries for performance
              const limitedIds = ids.slice(0, 10);
              
              // Create a custom writer that marks entries as exact matches
              const exactMatchWriter = {
                write: async (chunk: Uint8Array) => {
                  const text = new TextDecoder().decode(chunk);
                  const data = JSON.parse(text.trim());
                  data.isExactMatch = true;
                  const newChunk = JSON.stringify(data) + '\n';
                  writer.enqueue(new TextEncoder().encode(newChunk));
                }
              };
              
              await fetchDictionaryEntriesStreaming(dictType, limitedIds, shardType, exactMatchWriter as unknown as WritableStreamDefaultWriter<Uint8Array>);
            }
          }
        }

        // Send contained matches for single characters
        const isSingleCharacter = word.length === 1 && isHanCharacter(word);

        if (isSingleCharacter) {
          // Search for contained matches across all shards
          const allShardTypes = [ShardType.HAN_1CHAR, ShardType.HAN_2CHAR, ShardType.HAN_3PLUS, ShardType.NON_HAN];

          // Collect all contained match IDs first
          const allContainedIds: { [dictType: string]: { ids: number[], shardType: ShardType }[] } = {};

          for (const searchShardType of allShardTypes) {
            try {
              const searchIndexUrl = getIndexUrl(word, searchShardType);
              const searchIndexEntry = await fetchAndDecompressJson<IndexEntry>(searchIndexUrl);

              if (searchIndexEntry.c) {
                for (const [dictType, ids] of Object.entries(searchIndexEntry.c)) {
                  if (['j', 'n', 'd', 'c', 'w'].includes(dictType)) {
                    if (!allContainedIds[dictType]) {
                      allContainedIds[dictType] = [];
                    }
                    allContainedIds[dictType].push({ ids, shardType: searchShardType });
                  }
                }
              }
            } catch {
              // Index file doesn't exist in this shard, continue to next shard
              console.log(`No index file for ${word} in shard ${searchShardType}`);
            }
          }

          // Round-robin streaming: take a few from each dict type, repeat up to 10 times
          const maxRounds = 10;
          const dictTypes = Object.keys(allContainedIds);

          for (let round = 0; round < maxRounds; round++) {
            let addedInThisRound = false;

            for (const dictType of dictTypes) {
              const shardGroups = allContainedIds[dictType];
              if (!shardGroups) continue;

              // Find the next available ID across all shards for this dict type
              const totalProcessed = round * dictTypes.length + dictTypes.indexOf(dictType);
              const currentIndex = totalProcessed;

              // Find which shard and index within that shard
              let cumulativeCount = 0;
              for (const shardGroup of shardGroups) {
                if (currentIndex < cumulativeCount + shardGroup.ids.length) {
                  const indexInShard = currentIndex - cumulativeCount;
                  const idToFetch = shardGroup.ids[indexInShard];

                  try {
                    // Create a custom writer that marks entries as contained matches
                    const containedMatchWriter = {
                      write: async (chunk: Uint8Array) => {
                        const text = new TextDecoder().decode(chunk);
                        const data = JSON.parse(text.trim());
                        data.isExactMatch = false;
                        const newChunk = JSON.stringify(data) + '\n';
                        writer.enqueue(new TextEncoder().encode(newChunk));
                      }
                    };

                    await fetchDictionaryEntriesStreaming(dictType, [idToFetch], shardGroup.shardType, containedMatchWriter as unknown as WritableStreamDefaultWriter<Uint8Array>);
                    addedInThisRound = true;
                  } catch {
                    console.warn(`Error fetching entry ${idToFetch} from ${dictType}`);
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
        }

        // Send completion signal
        const completionSignal = {
          type: 'complete',
          containedMatchesPending: false
        };
        
        const chunk = JSON.stringify(completionSignal) + '\n';
        writer.enqueue(new TextEncoder().encode(chunk));
        
      } catch (error) {
        console.error("Error in streaming lookup:", error);
        
        const errorResponse = {
          error: "Failed to process lookup request"
        };
        
        const chunk = JSON.stringify(errorResponse) + '\n';
        writer.enqueue(new TextEncoder().encode(chunk));
      } finally {
        writer.close();
      }
    }
  });

  return new Response(stream, {
    headers: {
      "Content-Type": "application/json",
      "Cache-Control": "no-cache",
      "Connection": "keep-alive",
    },
  });
}
