'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { getShardType } from '@/lib/dictionary-utils';
import JishoEntryCard from './JishoEntryCard';
// Simple button component since ui/button doesn't exist
const Button = ({ children, onClick, disabled, variant, size, className = "", ...props }: {
  children: React.ReactNode;
  onClick?: () => void;
  disabled?: boolean;
  variant?: string;
  size?: string;
  className?: string;
  [key: string]: any;
}) => (
  <button
    onClick={onClick}
    disabled={disabled}
    className={`px-4 py-2 rounded-md font-medium transition-colors ${variant === 'outline'
      ? 'border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
      : 'bg-blue-600 text-white hover:bg-blue-700'
      } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'} ${size === 'sm' ? 'px-3 py-1 text-sm' : ''
      } ${className}`}
    {...props}
  >
    {children}
  </button>
);
// Simple icon components since lucide-react might not be available
const Loader2 = ({ className }: { className?: string }) => (
  <div className={`animate-spin rounded-full border-2 border-current border-t-transparent ${className}`} />
);

const ExternalLink = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
  </svg>
);

const Plus = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
  </svg>
);

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

// Dictionary section component with Load More functionality
interface DictionarySectionProps {
  title: string;
  dictType: 'j' | 'n' | 'd' | 'c' | 'w';
  entries: DictionaryEntry[];
  word: string;
  isContained?: boolean;
}

function DictionarySection({ title, dictType, entries, word, isContained = false }: DictionarySectionProps) {
  const router = useRouter();
  const [additionalEntries, setAdditionalEntries] = useState<DictionaryEntry[]>([]);
  const [loadingMore, setLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const allEntries = [...entries, ...additionalEntries];

  const handleLoadMore = async () => {
    setLoadingMore(true);
    try {
      const offset = allEntries.length;
      const response = await fetch(
        `/api/lookup/${encodeURIComponent(word)}/${dictType}?offset=${offset}&limit=10`
      );

      if (response.ok) {
        const data = await response.json();
        const newEntries = isContained
          ? data.containedMatches[dictType] || []
          : data.exactMatches[dictType] || [];

        setAdditionalEntries(prev => [...prev, ...newEntries]);
        setHasMore(data.pagination?.hasMore || false);
      }
    } catch (error) {
      console.error('Error loading more entries:', error);
    } finally {
      setLoadingMore(false);
    }
  };

  const handleSeeMore = () => {
    const url = `/lookup/${encodeURIComponent(word)}/${dictType}`;
    console.log('Navigating to:', url);
    console.log('Router:', router);

    // Use window.location for more reliable navigation
    window.location.href = url;
  };

  if (allEntries.length === 0) return null;

  return (
    <div className="mb-6">
      <h3 className="text-xl font-medium mb-3 border-b border-gray-700 pb-2 text-gray-300">
        {title}
      </h3>
      <div className="grid grid-cols-1 gap-4">
        {allEntries.map((entry: DictionaryEntry, index: number) => (
          <JishoEntryCard key={`${dictType}-${index}`} entry={entry} />
        ))}
      </div>

      {/* Action buttons */}
      {isContained && (
        <div className="flex gap-2 mt-4 justify-center">
          {hasMore && (
            <Button
              variant="outline"
              onClick={handleLoadMore}
              disabled={loadingMore}
              size="sm"
            >
              {loadingMore ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Loading...
                </>
              ) : (
                <>
                  <Plus className="h-4 w-4 mr-2" />
                  Load More
                </>
              )}
            </Button>
          )}

          <Button
            variant="outline"
            onClick={handleSeeMore}
            size="sm"
          >
            <ExternalLink className="h-4 w-4 mr-2" />
            See More in New Page
          </Button>
        </div>
      )}
    </div>
  );
}

function JishoDictionaryStreamingResults({ word }: DictionaryResultsProps) {
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

        // Initialize empty collections for real-time streaming
        const streamingExactMatches: DictionaryEntriesByType = {
          j: [], n: [], d: [], c: [], w: []
        };
        const streamingContainedMatches: DictionaryEntriesByType = {
          j: [], n: [], d: [], c: [], w: []
        };

        // Set initial empty state to show UI immediately
        setExactMatches(streamingExactMatches);
        setContainedMatches(streamingContainedMatches);
        setLoading(false); // Show UI immediately, entries will stream in

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
                }

                if ('containedMatches' in result) {
                  setContainedMatches(result.containedMatches);
                  setContainedMatchesPending(result.containedMatchesPending);
                }

                if ('error' in result) {
                  setError(result.error);
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

          // Try to parse complete JSON objects from the buffer (line-delimited JSON)
          const lines = buffer.split('\n');

          // Process all complete lines (except the last one which might be incomplete)
          for (let i = 0; i < lines.length - 1; i++) {
            const line = lines[i].trim();
            if (!line) continue;

            try {
              // Additional validation: check if line looks like valid JSON
              if (!line.startsWith('{') || !line.endsWith('}')) {
                console.warn('Skipping invalid JSON line:', line.substring(0, 100) + '...');
                continue;
              }

              const result = JSON.parse(line);
              console.log('📦 Parsed streaming response:', result);

              // Handle individual streaming entries
              if (result.type === 'entry') {
                const { dictType, entry, isExactMatch } = result;
                console.log(`🎯 Streaming ${isExactMatch ? 'exact' : 'contained'} match: ${dictType}`);

                if (isExactMatch) {
                  // Add to exact matches immediately
                  setExactMatches(prev => {
                    if (!prev) return prev;
                    const updated = { ...prev };
                    if (dictType in updated) {
                      updated[dictType as keyof DictionaryEntriesByType] = [
                        ...updated[dictType as keyof DictionaryEntriesByType],
                        entry
                      ];
                    }
                    return updated;
                  });
                } else {
                  // Add to contained matches immediately
                  setContainedMatches(prev => {
                    if (!prev) return prev;
                    const updated = { ...prev };
                    if (dictType in updated) {
                      updated[dictType as keyof DictionaryEntriesByType] = [
                        ...updated[dictType as keyof DictionaryEntriesByType],
                        entry
                      ];
                    }
                    return updated;
                  });
                }
              }
              // Handle legacy bulk responses (for compatibility)
              else if ('exactMatches' in result) {
                setExactMatches(result.exactMatches);
                setContainedMatchesPending(result.containedMatchesPending);
                console.log('📋 Set bulk exact matches');
              }
              else if ('containedMatches' in result) {
                setContainedMatches(result.containedMatches);
                setContainedMatchesPending(result.containedMatchesPending);
                console.log('📋 Set bulk contained matches');
              }
              else if ('error' in result) {
                setError(result.error);
                console.log('❌ Set error');
              }
            } catch (parseError) {
              // Only log parsing errors for lines that look like they should be JSON
              if (line.startsWith('{')) {
                const errorMessage = parseError instanceof Error ? parseError.message : 'Unknown error';
                console.warn('Failed to parse JSON line (might be incomplete):', errorMessage, 'Line preview:', line.substring(0, 100) + '...');
              }
            }
          }

          // Keep the last (potentially incomplete) line in the buffer
          buffer = lines[lines.length - 1];
        }
      } catch (error) {
        console.error('Error fetching streaming dictionary data:', error);
        const errorMessage = error instanceof Error ? error.message : 'Unknown error';
        setError(`Failed to load dictionary data: ${errorMessage}`);
        setLoading(false);
      }
    }

    fetchStreamingData();
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
      {containedMatchesPending && !hasContainedEntries && (
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

          <DictionarySection
            title="Chinese Characters"
            dictType="c"
            entries={groupedContainedEntries.chineseChar}
            word={word}
            isContained={true}
          />

          <DictionarySection
            title="Chinese Words"
            dictType="w"
            entries={groupedContainedEntries.chineseWord}
            word={word}
            isContained={true}
          />

          <DictionarySection
            title="Japanese Words"
            dictType="j"
            entries={groupedContainedEntries.japaneseWord}
            word={word}
            isContained={true}
          />

          <DictionarySection
            title="Kanji Characters"
            dictType="d"
            entries={groupedContainedEntries.kanjiChar}
            word={word}
            isContained={true}
          />

          <DictionarySection
            title="Japanese Names"
            dictType="n"
            entries={groupedContainedEntries.japaneseName}
            word={word}
            isContained={true}
          />
        </div>
      )}
    </div>
  );
}

export default JishoDictionaryStreamingResults;