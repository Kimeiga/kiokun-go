import { NextRequest, NextResponse } from "next/server";
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

// Dictionary entry types
interface DictionaryEntriesByType {
  j: any[]; // JMdict (Japanese words)
  n: any[]; // JMnedict (Japanese names)
  d: any[]; // Kanjidic (Kanji characters)
  c: any[]; // Chinese characters
  w: any[]; // Chinese words
}



async function fetchDictionaryEntries(
  dictType: string,
  ids: number[],
  shardType: ShardType
): Promise<any[]> {
  const entries: any[] = [];
  
  for (const id of ids) {
    try {
      const entryUrl = getDictionaryEntryUrl(dictType, id, shardType);
      const entry = await fetchAndDecompressJson(entryUrl);
      entries.push(entry);
    } catch (error) {
      console.warn(`Failed to fetch entry ${id} from ${dictType}:`, error);
    }
  }
  
  return entries;
}

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ word: string; dictType: string }> }
) {
  const { word, dictType } = await params;
  const { searchParams } = new URL(request.url);
  const offset = parseInt(searchParams.get("offset") || "0");
  const limit = parseInt(searchParams.get("limit") || "50");

  // Validate dictionary type
  const validDictTypes = ["j", "n", "d", "c", "w"];
  if (!validDictTypes.includes(dictType)) {
    return NextResponse.json(
      { error: "Invalid dictionary type" },
      { status: 400 }
    );
  }

  try {
    // Determine the shard type based on the word
    const shardType = getShardType(word);
    console.log(`Dictionary-specific lookup for word: ${word}, dictType: ${dictType}, shardType: ${shardType}`);

    // Initialize response structure
    const exactMatches: DictionaryEntriesByType = {
      j: [], n: [], d: [], c: [], w: [],
    };
    const containedMatches: DictionaryEntriesByType = {
      j: [], n: [], d: [], c: [], w: [],
    };

    // Get the URL for the index file
    const indexUrl = getIndexUrl(word, shardType);
    console.log(`Trying to fetch index from: ${indexUrl}`);

    // Fetch and decompress the index file for exact matches
    let indexEntry: IndexEntry;
    try {
      indexEntry = await fetchAndDecompressJson<IndexEntry>(indexUrl);
      console.log(`Found index entry:`, indexEntry);
    } catch (error) {
      console.log(`Index file not found in primary shard, error:`, error);
      // If the index file doesn't exist, check other shards for single characters
      const isSingleCharacter = word.length === 1 && isHanCharacter(word);
      
      if (isSingleCharacter) {
        // Search across all shards for single characters
        const allShardTypes = [ShardType.HAN_1CHAR, ShardType.HAN_2CHAR, ShardType.HAN_3PLUS, ShardType.NON_HAN];
        
        for (const searchShardType of allShardTypes) {
          try {
            const searchIndexUrl = getIndexUrl(word, searchShardType);
            const searchIndexEntry = await fetchAndDecompressJson<IndexEntry>(searchIndexUrl);
            
            // Process exact matches
            if (searchIndexEntry.e && searchIndexEntry.e[dictType]) {
              const entries = await fetchDictionaryEntries(
                dictType,
                searchIndexEntry.e[dictType],
                searchShardType
              );
              exactMatches[dictType as keyof DictionaryEntriesByType].push(...entries);
            }
            
            // Process contained matches with pagination
            if (searchIndexEntry.c && searchIndexEntry.c[dictType]) {
              const allIds = searchIndexEntry.c[dictType];
              const paginatedIds = allIds.slice(offset, offset + limit);
              
              const entries = await fetchDictionaryEntries(
                dictType,
                paginatedIds,
                searchShardType
              );
              containedMatches[dictType as keyof DictionaryEntriesByType].push(...entries);
            }
          } catch (shardError) {
            // Continue to next shard if this one doesn't exist
            continue;
          }
        }
      }
      
      return NextResponse.json({
        word,
        dictType,
        exactMatches,
        containedMatches,
        pagination: {
          offset,
          limit,
          hasMore: false, // We'll calculate this properly below
          total: containedMatches[dictType as keyof DictionaryEntriesByType].length
        }
      });
    }

    // Process exact matches from primary shard
    if (indexEntry.e && indexEntry.e[dictType]) {
      const entries = await fetchDictionaryEntries(
        dictType,
        indexEntry.e[dictType],
        shardType
      );
      exactMatches[dictType as keyof DictionaryEntriesByType] = entries;
    }

    // Process contained matches with pagination
    let totalContainedIds: number[] = [];
    
    // For single characters, collect from all shards
    const isSingleCharacter = word.length === 1 && isHanCharacter(word);
    
    if (isSingleCharacter) {
      const allShardTypes = [ShardType.HAN_1CHAR, ShardType.HAN_2CHAR, ShardType.HAN_3PLUS, ShardType.NON_HAN];
      
      for (const searchShardType of allShardTypes) {
        try {
          const searchIndexUrl = getIndexUrl(word, searchShardType);
          const searchIndexEntry = await fetchAndDecompressJson<IndexEntry>(searchIndexUrl);
          
          if (searchIndexEntry.c && searchIndexEntry.c[dictType]) {
            totalContainedIds.push(...searchIndexEntry.c[dictType]);
          }
        } catch (error) {
          // Continue to next shard
          continue;
        }
      }
    } else {
      // For multi-character words, only use primary shard
      if (indexEntry.c && indexEntry.c[dictType]) {
        totalContainedIds = indexEntry.c[dictType];
      }
    }

    // Apply pagination to contained matches and fetch from appropriate shards
    const paginatedIds = totalContainedIds.slice(offset, offset + limit);

    if (paginatedIds.length > 0) {
      // For single characters, we need to fetch from all shards where the IDs exist
      if (isSingleCharacter) {
        const allShardTypes = [ShardType.HAN_1CHAR, ShardType.HAN_2CHAR, ShardType.HAN_3PLUS, ShardType.NON_HAN];

        for (const searchShardType of allShardTypes) {
          try {
            const searchIndexUrl = getIndexUrl(word, searchShardType);
            const searchIndexEntry = await fetchAndDecompressJson<IndexEntry>(searchIndexUrl);

            if (searchIndexEntry.c && searchIndexEntry.c[dictType]) {
              // Find which of our paginated IDs exist in this shard
              const shardIds = searchIndexEntry.c[dictType];
              const idsToFetch = paginatedIds.filter(id => shardIds.includes(id));

              if (idsToFetch.length > 0) {
                const entries = await fetchDictionaryEntries(
                  dictType,
                  idsToFetch,
                  searchShardType
                );
                containedMatches[dictType as keyof DictionaryEntriesByType].push(...entries);
              }
            }
          } catch (error) {
            // Continue to next shard
            continue;
          }
        }
      } else {
        // For multi-character words, use primary shard
        const entries = await fetchDictionaryEntries(
          dictType,
          paginatedIds,
          shardType
        );
        containedMatches[dictType as keyof DictionaryEntriesByType] = entries;
      }
    }

    const hasMore = offset + limit < totalContainedIds.length;

    return NextResponse.json({
      word,
      dictType,
      exactMatches,
      containedMatches,
      pagination: {
        offset,
        limit,
        hasMore,
        total: totalContainedIds.length
      }
    });

  } catch (error) {
    console.error("Error in dictionary-specific lookup:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    );
  }
}
