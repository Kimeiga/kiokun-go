'use client';

import { useState, useEffect } from 'react';
import { getShardType } from '@/lib/dictionary-utils';
import JishoEntryCard from './JishoEntryCard';

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

interface DictionaryResultsProps {
  word: string;
}

interface InitialLookupResponse {
  word: string;
  exactMatches: DictionaryEntriesByType;
  containedMatchesPending: boolean;
}

interface ContainedMatchesResponse {
  containedMatches: DictionaryEntriesByType;
  containedMatchesPending: boolean;
}

export default function JishoDictionaryStreamingResults({ word }: DictionaryResultsProps) {
  const [exactMatches, setExactMatches] = useState<DictionaryEntriesByType | null>(null);
  const [containedMatches, setContainedMatches] = useState<DictionaryEntriesByType | null>(null);
  const [containedMatchesPending, setContainedMatchesPending] = useState<boolean>(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const shardType = getShardType(word);

  useEffect(() => {
    async function fetchStreamingData() {
      try {
        setLoading(true);
        setExactMatches(null);
        setContainedMatches(null);
        setContainedMatchesPending(false);
        setError(null);

        // Use relative URL for API calls
        const apiUrl = `/api/lookup-stream?word=${encodeURIComponent(word)}`;
        console.log(`Fetching from streaming URL: ${apiUrl}`);

        const response = await fetch(apiUrl);

        if (!response.ok) {
          throw new Error(`Failed to fetch dictionary data: ${response.status} ${response.statusText}`);
        }

        if (!response.body) {
          throw new Error('ReadableStream not supported in this browser.');
        }

        // Get a reader from the response body
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        // Read the stream
        while (true) {
          const { done, value } = await reader.read();

          if (done) {
            // Process any remaining buffer content
            if (buffer.trim()) {
              try {
                const result = JSON.parse(buffer);
                console.log('Final buffer parse result:', result);

                if ('exactMatches' in result) {
                  setExactMatches(result.exactMatches);
                  setContainedMatchesPending(result.containedMatchesPending);
                  setLoading(false);
                }

                if ('containedMatches' in result) {
                  setContainedMatches(result.containedMatches);
                  setContainedMatchesPending(result.containedMatchesPending);
                }

                if ('error' in result) {
                  setError(result.error);
                  setLoading(false);
                }
              } catch (error) {
                console.error('Error parsing final buffer:', error, 'Buffer:', buffer);
              }
            }
            break;
          }

          // Decode the chunk and add it to our buffer
          const chunk = decoder.decode(value, { stream: true });
          buffer += chunk;
          console.log('Received chunk:', chunk);
          console.log('Current buffer:', buffer);

          // Try to parse complete JSON objects from the buffer
          // Split by potential JSON object boundaries and try to parse each
          const jsonObjects = [];
          let currentJson = '';
          let braceCount = 0;
          let inString = false;
          let escaped = false;

          for (let i = 0; i < buffer.length; i++) {
            const char = buffer[i];
            currentJson += char;

            if (!inString) {
              if (char === '{') {
                braceCount++;
              } else if (char === '}') {
                braceCount--;
                if (braceCount === 0 && currentJson.trim()) {
                  // We have a complete JSON object
                  jsonObjects.push(currentJson.trim());
                  currentJson = '';
                }
              } else if (char === '"') {
                inString = true;
              }
            } else {
              if (escaped) {
                escaped = false;
              } else if (char === '\\') {
                escaped = true;
              } else if (char === '"') {
                inString = false;
              }
            }
          }

          // Process each complete JSON object
          for (const jsonStr of jsonObjects) {
            try {
              const result = JSON.parse(jsonStr);
              console.log('Parsed JSON object:', result);

              // Process the result based on its content
              if ('exactMatches' in result) {
                setExactMatches(result.exactMatches);
                setContainedMatchesPending(result.containedMatchesPending);
                setLoading(false);
                console.log('Set exact matches and loading to false');
              }

              if ('containedMatches' in result) {
                setContainedMatches(result.containedMatches);
                setContainedMatchesPending(result.containedMatchesPending);
                console.log('Set contained matches');
              }

              if ('error' in result) {
                setError(result.error);
                setLoading(false);
                console.log('Set error and loading to false');
              }
            } catch (parseError) {
              console.error('Error parsing JSON object:', parseError, 'JSON:', jsonStr);
            }
          }

          // Update buffer to keep any incomplete JSON
          buffer = currentJson;
        }
      } catch (error) {
        console.error('Error fetching streaming dictionary data:', error);
        setError(error instanceof Error ? error.message : 'Unknown error');
        setLoading(false);
      }
    }

    fetchStreamingData();
  }, [word]);

  if (loading && !exactMatches) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-400"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-900/30 border border-red-800 rounded-lg p-4 mb-6">
        <h2 className="text-red-400 text-lg font-semibold mb-2">Error</h2>
        <p className="text-red-300">{error}</p>
      </div>
    );
  }

  if (!exactMatches) {
    return (
      <div className="bg-yellow-900/30 border border-yellow-800 rounded-lg p-4 mb-6">
        <h2 className="text-yellow-400 text-lg font-semibold mb-2">No Data</h2>
        <p className="text-yellow-300">
          No data was returned for &ldquo;{word}&rdquo;.
        </p>
      </div>
    );
  }

  // Check if there are any entries
  const hasAnyExactMatches = Object.values(exactMatches).some(entries => entries.length > 0);
  const hasAnyContainedMatches = containedMatches ?
    Object.values(containedMatches).some(entries => entries.length > 0) :
    false;

  // Handle no results
  if (!hasAnyExactMatches && !hasAnyContainedMatches && !containedMatchesPending) {
    return (
      <div className="bg-yellow-900/30 border border-yellow-800 rounded-lg p-4 mb-6">
        <h2 className="text-yellow-400 text-lg font-semibold mb-2">No Results Found</h2>
        <p className="text-yellow-300">
          No dictionary entries were found for &ldquo;{word}&rdquo;.
        </p>
        <p className="text-yellow-300 mt-2">
          Shard type: {shardType} (
          {shardType === 0 ? 'Non-Han' :
            shardType === 1 ? 'Han 1 Character' :
              shardType === 2 ? 'Han 2 Characters' :
                'Han 3+ Characters'})
        </p>
      </div>
    );
  }

  // The entries are already organized by dictionary type in the API response
  const groupedExactEntries = {
    // Chinese characters (c)
    chineseChar: exactMatches.c || [],
    // Chinese words (w)
    chineseWord: exactMatches.w || [],
    // Japanese words (j)
    japaneseWord: exactMatches.j || [],
    // Japanese names (n)
    japaneseName: exactMatches.n || [],
    // Kanji characters (d)
    kanjiChar: exactMatches.d || [],
    // Other entries (should be empty with the new structure)
    other: []
  };

  // Only process contained matches if they're available
  const groupedContainedEntries = containedMatches ? {
    // Chinese characters (c)
    chineseChar: containedMatches.c || [],
    // Chinese words (w)
    chineseWord: containedMatches.w || [],
    // Japanese words (j)
    japaneseWord: containedMatches.j || [],
    // Japanese names (n)
    japaneseName: containedMatches.n || [],
    // Kanji characters (d)
    kanjiChar: containedMatches.d || [],
    // Other entries (should be empty with the new structure)
    other: []
  } : {
    chineseChar: [],
    chineseWord: [],
    japaneseWord: [],
    japaneseName: [],
    kanjiChar: [],
    other: []
  };

  // Check if any group has entries
  const hasExactEntries = Object.values(groupedExactEntries).some(group => group.length > 0);
  const hasContainedEntries = Object.values(groupedContainedEntries).some(group => group.length > 0);

  return (
    <div>
      {/* Display shard type info */}
      <div className="bg-blue-900/30 border border-blue-800 rounded-lg p-4 mb-6">
        <p className="text-blue-400">
          Shard type: {shardType} (
          {shardType === 0 ? 'Non-Han' :
            shardType === 1 ? 'Han 1 Character' :
              shardType === 2 ? 'Han 2 Characters' :
                'Han 3+ Characters'})
        </p>
      </div>

      {/* Exact matches section */}
      {hasExactEntries && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4 text-white">Exact Matches</h2>

          {/* Chinese Character entries */}
          {groupedExactEntries.chineseChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedExactEntries.chineseChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-cc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Chinese Word entries */}
          {groupedExactEntries.chineseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedExactEntries.chineseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-cw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Word entries */}
          {groupedExactEntries.japaneseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedExactEntries.japaneseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-jw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Kanji Character entries */}
          {groupedExactEntries.kanjiChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Kanji Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedExactEntries.kanjiChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-kc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Name entries */}
          {groupedExactEntries.japaneseName.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Names
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedExactEntries.japaneseName.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-jn-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Other entries */}
          {groupedExactEntries.other.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Other Entries
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedExactEntries.other.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-other-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Contained matches section */}
      {containedMatchesPending && !containedMatches && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4 text-white">Contained-in Matches</h2>
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-400"></div>
            <span className="ml-3 text-gray-300">Loading contained matches...</span>
          </div>
        </div>
      )}

      {hasContainedEntries && (
        <div>
          <h2 className="text-2xl font-semibold mb-4 text-white">Contained-in Matches</h2>

          {/* Chinese Character entries */}
          {groupedContainedEntries.chineseChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedContainedEntries.chineseChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-cc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Chinese Word entries */}
          {groupedContainedEntries.chineseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedContainedEntries.chineseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-cw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Word entries */}
          {groupedContainedEntries.japaneseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedContainedEntries.japaneseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-jw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Kanji Character entries */}
          {groupedContainedEntries.kanjiChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Kanji Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedContainedEntries.kanjiChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-kc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Name entries */}
          {groupedContainedEntries.japaneseName.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Names
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedContainedEntries.japaneseName.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-jn-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Other entries */}
          {groupedContainedEntries.other.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Other Entries
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedContainedEntries.other.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-other-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

