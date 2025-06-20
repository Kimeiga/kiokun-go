'use client';

import { useState, useEffect } from 'react';
import { getShardType } from '@/lib/dictionary-utils';
import JishoEntryCard from './JishoEntryCard';

interface DictionaryResultsProps {
  word: string;
}

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

interface LookupResponse {
  word: string;
  exactMatches: DictionaryEntriesByType;
  containedMatches: DictionaryEntriesByType;
  error?: string;
}

export default function JishoDictionaryResults({ word }: DictionaryResultsProps) {
  const [data, setData] = useState<LookupResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const shardType = getShardType(word);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);

        // Use relative URL for API calls
        const apiUrl = `/api/lookup?word=${encodeURIComponent(word)}`;

        console.log(`Fetching from URL: ${apiUrl}`);

        const response = await fetch(apiUrl);

        if (!response.ok) {
          throw new Error(`Failed to fetch dictionary data: ${response.status} ${response.statusText}`);
        }

        const result = await response.json();
        setData(result);
      } catch (error) {
        console.error('Error fetching dictionary data:', error);
        setError(error instanceof Error ? error.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    }

    fetchData();
  }, [word]);

  if (loading) {
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

  if (!data) {
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
  const hasAnyExactMatches = Object.values(data.exactMatches).some(entries => entries.length > 0);
  const hasAnyContainedMatches = Object.values(data.containedMatches).some(entries => entries.length > 0);

  // Handle no results
  if (!hasAnyExactMatches && !hasAnyContainedMatches) {
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
  const groupedEntries = {
    exact: {
      // Chinese characters (c)
      chineseChar: data.exactMatches.c || [],
      // Chinese words (w)
      chineseWord: data.exactMatches.w || [],
      // Japanese words (j)
      japaneseWord: data.exactMatches.j || [],
      // Japanese names (n)
      japaneseName: data.exactMatches.n || [],
      // Kanji characters (d)
      kanjiChar: data.exactMatches.d || [],
      // Other entries (should be empty with the new structure)
      other: []
    },
    contained: {
      // Chinese characters (c)
      chineseChar: data.containedMatches.c || [],
      // Chinese words (w)
      chineseWord: data.containedMatches.w || [],
      // Japanese words (j)
      japaneseWord: data.containedMatches.j || [],
      // Japanese names (n)
      japaneseName: data.containedMatches.n || [],
      // Kanji characters (d)
      kanjiChar: data.containedMatches.d || [],
      // Other entries (should be empty with the new structure)
      other: []
    }
  };

  // Check if any group has entries
  const hasExactEntries = Object.values(groupedEntries.exact).some(group => group.length > 0);
  const hasContainedEntries = Object.values(groupedEntries.contained).some(group => group.length > 0);

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
          {groupedEntries.exact.chineseChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.chineseChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-cc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Chinese Word entries */}
          {groupedEntries.exact.chineseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.chineseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-cw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Word entries */}
          {groupedEntries.exact.japaneseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.japaneseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-jw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Kanji Character entries */}
          {groupedEntries.exact.kanjiChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Kanji Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.kanjiChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-kc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Name entries */}
          {groupedEntries.exact.japaneseName.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Names
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.japaneseName.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-jn-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Other entries */}
          {groupedEntries.exact.other.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Other Entries
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.other.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`exact-other-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Contained matches section */}
      {hasContainedEntries && (
        <div>
          <h2 className="text-2xl font-semibold mb-4 text-white">Contained-in Matches</h2>

          {/* Chinese Character entries */}
          {groupedEntries.contained.chineseChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.chineseChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-cc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Chinese Word entries */}
          {groupedEntries.contained.chineseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Chinese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.chineseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-cw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Word entries */}
          {groupedEntries.contained.japaneseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.japaneseWord.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-jw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Kanji Character entries */}
          {groupedEntries.contained.kanjiChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Kanji Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.kanjiChar.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-kc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Japanese Name entries */}
          {groupedEntries.contained.japaneseName.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Japanese Names
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.japaneseName.map((entry: DictionaryEntry, index: number) => (
                  <JishoEntryCard key={`contained-jn-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}

          {/* Other entries */}
          {groupedEntries.contained.other.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
                Other Entries
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.other.map((entry: DictionaryEntry, index: number) => (
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
