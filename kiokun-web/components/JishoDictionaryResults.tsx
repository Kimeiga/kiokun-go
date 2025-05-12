'use client';

import { useState, useEffect } from 'react';
import { getShardType } from '@/lib/dictionary-utils';
import JishoEntryCard from './JishoEntryCard';

interface DictionaryResultsProps {
  word: string;
}

interface LookupResponse {
  word: string;
  exactMatches: any[];
  containedMatches: any[];
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
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }
  
  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
        <h2 className="text-red-800 text-lg font-semibold mb-2">Error</h2>
        <p className="text-red-700">{error}</p>
      </div>
    );
  }
  
  if (!data) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
        <h2 className="text-yellow-800 text-lg font-semibold mb-2">No Data</h2>
        <p className="text-yellow-700">
          No data was returned for &ldquo;{word}&rdquo;.
        </p>
      </div>
    );
  }
  
  // Handle no results
  if (data.exactMatches.length === 0 && data.containedMatches.length === 0) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
        <h2 className="text-yellow-800 text-lg font-semibold mb-2">No Results Found</h2>
        <p className="text-yellow-700">
          No dictionary entries were found for &ldquo;{word}&rdquo;.
        </p>
        <p className="text-yellow-700 mt-2">
          Shard type: {shardType} (
          {shardType === 0 ? 'Non-Han' : 
           shardType === 1 ? 'Han 1 Character' : 
           shardType === 2 ? 'Han 2 Characters' : 
           'Han 3+ Characters'})
        </p>
      </div>
    );
  }
  
  // Group entries by dictionary type
  const groupedEntries = {
    exact: {
      chineseChar: data.exactMatches.filter(entry => 
        (entry.traditional && entry.simplified && !entry.traditional_chinese) || 
        (entry.Traditional && entry.Simplified && !entry.Components)
      ),
      chineseWord: data.exactMatches.filter(entry => 
        entry.traditional_chinese || 
        (entry.Traditional && entry.Pinyin && entry.Components)
      ),
      japaneseWord: data.exactMatches.filter(entry => 
        (entry.kanji && entry.kana) || 
        (entry.Kanji && entry.Kana)
      ),
      japaneseName: data.exactMatches.filter(entry => 
        (entry.k && entry.r) || 
        (entry.Kanji && entry.Reading && entry.Translation)
      ),
      kanjiChar: data.exactMatches.filter(entry => 
        (entry.c && entry.on) || 
        (entry.Character && entry.Reading && entry.Misc)
      ),
      other: data.exactMatches.filter(entry => 
        !((entry.traditional && entry.simplified && !entry.traditional_chinese) || 
          (entry.Traditional && entry.Simplified && !entry.Components)) &&
        !(entry.traditional_chinese || 
          (entry.Traditional && entry.Pinyin && entry.Components)) &&
        !((entry.kanji && entry.kana) || 
          (entry.Kanji && entry.Kana)) &&
        !((entry.k && entry.r) || 
          (entry.Kanji && entry.Reading && entry.Translation)) &&
        !((entry.c && entry.on) || 
          (entry.Character && entry.Reading && entry.Misc))
      )
    },
    contained: {
      chineseChar: data.containedMatches.filter(entry => 
        (entry.traditional && entry.simplified && !entry.traditional_chinese) || 
        (entry.Traditional && entry.Simplified && !entry.Components)
      ),
      chineseWord: data.containedMatches.filter(entry => 
        entry.traditional_chinese || 
        (entry.Traditional && entry.Pinyin && entry.Components)
      ),
      japaneseWord: data.containedMatches.filter(entry => 
        (entry.kanji && entry.kana) || 
        (entry.Kanji && entry.Kana)
      ),
      japaneseName: data.containedMatches.filter(entry => 
        (entry.k && entry.r) || 
        (entry.Kanji && entry.Reading && entry.Translation)
      ),
      kanjiChar: data.containedMatches.filter(entry => 
        (entry.c && entry.on) || 
        (entry.Character && entry.Reading && entry.Misc)
      ),
      other: data.containedMatches.filter(entry => 
        !((entry.traditional && entry.simplified && !entry.traditional_chinese) || 
          (entry.Traditional && entry.Simplified && !entry.Components)) &&
        !(entry.traditional_chinese || 
          (entry.Traditional && entry.Pinyin && entry.Components)) &&
        !((entry.kanji && entry.kana) || 
          (entry.Kanji && entry.Kana)) &&
        !((entry.k && entry.r) || 
          (entry.Kanji && entry.Reading && entry.Translation)) &&
        !((entry.c && entry.on) || 
          (entry.Character && entry.Reading && entry.Misc))
      )
    }
  };
  
  // Check if any group has entries
  const hasExactEntries = Object.values(groupedEntries.exact).some(group => group.length > 0);
  const hasContainedEntries = Object.values(groupedEntries.contained).some(group => group.length > 0);
  
  return (
    <div>
      {/* Display shard type info */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
        <p className="text-blue-700">
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
          <h2 className="text-2xl font-semibold mb-4">Exact Matches</h2>
          
          {/* Chinese Character entries */}
          {groupedEntries.exact.chineseChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Chinese Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.chineseChar.map((entry, index) => (
                  <JishoEntryCard key={`exact-cc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Chinese Word entries */}
          {groupedEntries.exact.chineseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Chinese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.chineseWord.map((entry, index) => (
                  <JishoEntryCard key={`exact-cw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Japanese Word entries */}
          {groupedEntries.exact.japaneseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Japanese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.japaneseWord.map((entry, index) => (
                  <JishoEntryCard key={`exact-jw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Kanji Character entries */}
          {groupedEntries.exact.kanjiChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Kanji Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.kanjiChar.map((entry, index) => (
                  <JishoEntryCard key={`exact-kc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Japanese Name entries */}
          {groupedEntries.exact.japaneseName.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Japanese Names
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.japaneseName.map((entry, index) => (
                  <JishoEntryCard key={`exact-jn-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Other entries */}
          {groupedEntries.exact.other.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Other Entries
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.exact.other.map((entry, index) => (
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
          <h2 className="text-2xl font-semibold mb-4">Contained-in Matches</h2>
          
          {/* Chinese Character entries */}
          {groupedEntries.contained.chineseChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Chinese Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.chineseChar.map((entry, index) => (
                  <JishoEntryCard key={`contained-cc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Chinese Word entries */}
          {groupedEntries.contained.chineseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Chinese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.chineseWord.map((entry, index) => (
                  <JishoEntryCard key={`contained-cw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Japanese Word entries */}
          {groupedEntries.contained.japaneseWord.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Japanese Words
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.japaneseWord.map((entry, index) => (
                  <JishoEntryCard key={`contained-jw-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Kanji Character entries */}
          {groupedEntries.contained.kanjiChar.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Kanji Characters
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.kanjiChar.map((entry, index) => (
                  <JishoEntryCard key={`contained-kc-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Japanese Name entries */}
          {groupedEntries.contained.japaneseName.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Japanese Names
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.japaneseName.map((entry, index) => (
                  <JishoEntryCard key={`contained-jn-${index}`} entry={entry} />
                ))}
              </div>
            </div>
          )}
          
          {/* Other entries */}
          {groupedEntries.contained.other.length > 0 && (
            <div className="mb-6">
              <h3 className="text-xl font-medium mb-3 border-b border-gray-200 pb-2">
                Other Entries
              </h3>
              <div className="grid grid-cols-1 gap-4">
                {groupedEntries.contained.other.map((entry, index) => (
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
