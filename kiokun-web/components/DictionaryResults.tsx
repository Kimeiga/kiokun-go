/**
 * Component for displaying dictionary lookup results
 */

import { getShardType } from '@/lib/dictionary-utils';
import EntryCard from './EntryCard';

interface DictionaryResultsProps {
  word: string;
}

interface LookupResponse {
  word: string;
  exactMatches: any[];
  containedMatches: any[];
  error?: string;
}

async function fetchDictionaryData(word: string): Promise<LookupResponse> {
  try {
    // Use relative URL for API calls
    const apiUrl = `/api/lookup?word=${encodeURIComponent(word)}`;

    console.log(`Fetching from URL: ${apiUrl}`);

    const response = await fetch(apiUrl, {
      // This ensures the request is made on the server during SSR
      cache: 'no-store',
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch dictionary data: ${response.status} ${response.statusText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Error fetching dictionary data:', error);
    return {
      word,
      exactMatches: [],
      containedMatches: [],
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}

export default async function DictionaryResults({ word }: DictionaryResultsProps) {
  const data = await fetchDictionaryData(word);
  const shardType = getShardType(word);

  // Handle errors
  if (data.error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
        <h2 className="text-red-800 text-lg font-semibold mb-2">Error</h2>
        <p className="text-red-700">{data.error}</p>
      </div>
    );
  }

  // Handle no results
  if (data.exactMatches.length === 0 && data.containedMatches.length === 0) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4 mb-6">
        <h2 className="text-yellow-800 text-lg font-semibold mb-2">No Results Found</h2>
        <p className="text-yellow-700">
          No dictionary entries were found for "{word}".
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

  return (
    <div>
      {/* Display shard type info */}
      <div className="bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
        <p className="text-blue-700">
          Shard type: {shardType} (
          {shardType === 0 ? 'Non-Han' :
            shardType === 1 ? 'Han 1 Character' :
              shardType === 2 ? 'Han 2 Characters' :
                'Han 3+ Characters'})
        </p>
      </div>

      {/* Exact matches section */}
      {data.exactMatches.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Exact Matches</h2>
          <div className="grid grid-cols-1 gap-4">
            {data.exactMatches.map((entry, index) => (
              <EntryCard key={`exact-${index}`} entry={entry} />
            ))}
          </div>
        </div>
      )}

      {/* Contained matches section */}
      {data.containedMatches.length > 0 && (
        <div>
          <h2 className="text-2xl font-semibold mb-4">Contained-in Matches</h2>
          <div className="grid grid-cols-1 gap-4">
            {data.containedMatches.map((entry, index) => (
              <EntryCard key={`contained-${index}`} entry={entry} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
