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

export default function DictionaryResults({ word }: DictionaryResultsProps) {
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
      <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
        <h2 className="text-red-800 text-lg font-semibold mb-2">Error</h2>
        <p className="text-red-700">{error}</p>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4 mb-6">
        <h2 className="text-yellow-800 text-lg font-semibold mb-2">No Data</h2>
        <p className="text-yellow-700">
          No data was returned for "{word}".
        </p>
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
              <JishoEntryCard key={`exact-${index}`} entry={entry} />
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
              <JishoEntryCard key={`contained-${index}`} entry={entry} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
