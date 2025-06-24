"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
// Simple components since ui components might not be available
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
    className={`px-4 py-2 rounded-md font-medium transition-colors ${
      variant === 'outline'
        ? 'border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
        : 'bg-blue-600 text-white hover:bg-blue-700'
    } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'} ${
      size === 'lg' ? 'px-6 py-3 text-lg' : ''
    } ${className}`}
    {...props}
  >
    {children}
  </button>
);

const ArrowLeft = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
  </svg>
);

const Loader2 = ({ className }: { className?: string }) => (
  <div className={`animate-spin rounded-full border-2 border-current border-t-transparent ${className}`} />
);

interface DictionaryEntriesByType {
  j: any[];
  n: any[];
  d: any[];
  c: any[];
  w: any[];
}

interface PaginationInfo {
  offset: number;
  limit: number;
  hasMore: boolean;
  total: number;
}

interface DictionarySpecificResponse {
  word: string;
  dictType: string;
  exactMatches: DictionaryEntriesByType;
  containedMatches: DictionaryEntriesByType;
  pagination: PaginationInfo;
}

const DICT_TYPE_NAMES = {
  j: "JMdict (Japanese Words)",
  n: "JMNedict (Japanese Names)",
  d: "Kanjidic (Kanji Characters)",
  c: "Chinese Characters",
  w: "Chinese Words"
};

export default function DictionarySpecificPage() {
  const params = useParams();
  const router = useRouter();
  const word = params.word as string;
  const dictType = params.dictType as string;

  const [data, setData] = useState<DictionarySpecificResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async (offset: number = 0, append: boolean = false) => {
    try {
      if (!append) setLoading(true);
      else setLoadingMore(true);

      const response = await fetch(
        `/api/lookup/${encodeURIComponent(word)}/${dictType}?offset=${offset}&limit=50`
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const newData: DictionarySpecificResponse = await response.json();

      if (append && data) {
        // Append new contained matches to existing data
        const updatedData = { ...newData };
        updatedData.containedMatches[dictType as keyof DictionaryEntriesByType] = [
          ...data.containedMatches[dictType as keyof DictionaryEntriesByType],
          ...newData.containedMatches[dictType as keyof DictionaryEntriesByType]
        ];
        setData(updatedData);
      } else {
        setData(newData);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [word, dictType]);

  const handleLoadMore = () => {
    if (data && data.pagination.hasMore) {
      const nextOffset = data.containedMatches[dictType as keyof DictionaryEntriesByType].length;
      fetchData(nextOffset, true);
    }
  };

  const renderEntry = (entry: any, index: number) => {
    // This is a simplified renderer - you might want to use your existing components
    return (
      <div key={index} className="border rounded-lg p-4 mb-4 bg-white shadow-sm">
        <pre className="text-sm overflow-x-auto">
          {JSON.stringify(entry, null, 2)}
        </pre>
      </div>
    );
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center min-h-[400px]">
          <Loader2 className="h-8 w-8 animate-spin" />
          <span className="ml-2">Loading dictionary entries...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Error</h1>
          <p className="text-gray-600 mb-4">{error}</p>
          <Button onClick={() => router.back()}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Go Back
          </Button>
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">No Data Found</h1>
          <Button onClick={() => router.back()}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Go Back
          </Button>
        </div>
      </div>
    );
  }

  const exactEntries = data.exactMatches[dictType as keyof DictionaryEntriesByType];
  const containedEntries = data.containedMatches[dictType as keyof DictionaryEntriesByType];
  const dictTypeName = DICT_TYPE_NAMES[dictType as keyof typeof DICT_TYPE_NAMES] || dictType;

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <Button 
          variant="outline" 
          onClick={() => router.back()}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Search Results
        </Button>
        
        <h1 className="text-3xl font-bold mb-2">
          {dictTypeName} entries for "{word}"
        </h1>
        
        <div className="text-gray-600">
          <p>
            {exactEntries.length} exact match{exactEntries.length !== 1 ? 'es' : ''}, 
            {' '}{containedEntries.length} of {data.pagination.total} contained matches
          </p>
        </div>
      </div>

      {/* Exact Matches */}
      {exactEntries.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Exact Matches</h2>
          <div className="space-y-4">
            {exactEntries.map((entry, index) => renderEntry(entry, index))}
          </div>
        </div>
      )}

      {/* Contained Matches */}
      {containedEntries.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">
            Words Containing "{word}"
          </h2>
          <div className="space-y-4">
            {containedEntries.map((entry, index) => renderEntry(entry, index))}
          </div>
          
          {/* Load More Button */}
          {data.pagination.hasMore && (
            <div className="text-center mt-6">
              <Button 
                onClick={handleLoadMore}
                disabled={loadingMore}
                size="lg"
              >
                {loadingMore ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Loading More...
                  </>
                ) : (
                  `Load More (${data.pagination.total - containedEntries.length} remaining)`
                )}
              </Button>
            </div>
          )}
        </div>
      )}

      {/* No Results */}
      {exactEntries.length === 0 && containedEntries.length === 0 && (
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold mb-2">No Results Found</h2>
          <p className="text-gray-600">
            No entries found in {dictTypeName} for "{word}"
          </p>
        </div>
      )}
    </div>
  );
}
