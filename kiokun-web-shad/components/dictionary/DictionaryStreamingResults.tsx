"use client";

import React, { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2 } from "lucide-react";
import { JishoEntryCard } from "./JishoEntryCard";

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

interface DictionaryStreamingResultsProps {
  word: string;
}

interface StreamingData {
  exactMatches: DictionaryEntriesByType;
  containedMatches: DictionaryEntriesByType;
  containedMatchesPending: boolean;
}

export function DictionaryStreamingResults({ word }: DictionaryStreamingResultsProps) {
  const [data, setData] = useState<StreamingData>({
    exactMatches: { j: [], n: [], d: [], c: [], w: [] },
    containedMatches: { j: [], n: [], d: [], c: [], w: [] },
    containedMatchesPending: true,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const response = await fetch(`/api/lookup-stream?word=${encodeURIComponent(word)}`);

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const reader = response.body?.getReader();
        if (!reader) {
          throw new Error("No response body");
        }

        const decoder = new TextDecoder();
        let buffer = "";

        while (true) {
          const { done, value } = await reader.read();

          if (done) break;

          buffer += decoder.decode(value, { stream: true });

          // Process complete JSON objects
          const lines = buffer.split("\n");
          buffer = lines.pop() || ""; // Keep incomplete line in buffer

          for (const line of lines) {
            if (line.trim()) {
              try {
                const parsed = JSON.parse(line);

                if (parsed.type === "entry") {
                  const { dictType, entry, isExactMatch } = parsed;

                  setData(prev => ({
                    ...prev,
                    [isExactMatch ? "exactMatches" : "containedMatches"]: {
                      ...prev[isExactMatch ? "exactMatches" : "containedMatches"],
                      [dictType]: [
                        ...prev[isExactMatch ? "exactMatches" : "containedMatches"][dictType as keyof DictionaryEntriesByType],
                        entry
                      ]
                    }
                  }));
                } else if (parsed.containedMatches) {
                  // Handle batch contained matches
                  setData(prev => ({
                    ...prev,
                    containedMatches: {
                      j: [...prev.containedMatches.j, ...(parsed.containedMatches.j || [])],
                      n: [...prev.containedMatches.n, ...(parsed.containedMatches.n || [])],
                      d: [...prev.containedMatches.d, ...(parsed.containedMatches.d || [])],
                      c: [...prev.containedMatches.c, ...(parsed.containedMatches.c || [])],
                      w: [...prev.containedMatches.w, ...(parsed.containedMatches.w || [])],
                    },
                    containedMatchesPending: parsed.containedMatchesPending || false,
                  }));
                } else if (parsed.exactMatches) {
                  // Handle batch exact matches
                  setData(prev => ({
                    ...prev,
                    exactMatches: {
                      j: [...prev.exactMatches.j, ...(parsed.exactMatches.j || [])],
                      n: [...prev.exactMatches.n, ...(parsed.exactMatches.n || [])],
                      d: [...prev.exactMatches.d, ...(parsed.exactMatches.d || [])],
                      c: [...prev.exactMatches.c, ...(parsed.exactMatches.c || [])],
                      w: [...prev.exactMatches.w, ...(parsed.exactMatches.w || [])],
                    },
                  }));
                } else if (parsed.type === "complete") {
                  setData(prev => ({
                    ...prev,
                    containedMatchesPending: false,
                  }));
                }
              } catch (parseError) {
                console.error("Error parsing JSON:", parseError, "Line:", line);
              }
            }
          }
        }

        setIsLoading(false);
      } catch (err) {
        console.error("Error fetching streaming data:", err);
        setError(err instanceof Error ? err.message : "An error occurred");
        setIsLoading(false);
      }
    };

    fetchData();
  }, [word]);

  // Count total entries for each type
  const getTotalCount = (entries: DictionaryEntriesByType) => {
    return Object.values(entries).reduce((sum, arr) => sum + arr.length, 0);
  };

  const exactMatchCount = getTotalCount(data.exactMatches);
  const containedMatchCount = getTotalCount(data.containedMatches);

  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-red-600">
            <p>Error loading results: {error}</p>
            <Button
              onClick={() => window.location.reload()}
              variant="outline"
              className="mt-4"
            >
              Retry
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Exact Matches */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              Exact Matches
              {exactMatchCount > 0 && (
                <Badge variant="secondary">{exactMatchCount}</Badge>
              )}
            </CardTitle>
            {isLoading && (
              <Loader2 className="h-4 w-4 animate-spin" />
            )}
          </div>
        </CardHeader>
        <CardContent>
          {exactMatchCount === 0 && !isLoading ? (
            <p className="text-muted-foreground text-center py-8">
              No exact matches found for &ldquo;{word}&rdquo;
            </p>
          ) : (
            <div className="space-y-6">
              {/* Two-column section: Words first, then Characters */}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Chinese Words (Left Column, Position 1) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Chinese Words
                    {data.exactMatches.w.length > 0 && (
                      <Badge variant="outline">{data.exactMatches.w.length}</Badge>
                    )}
                  </h3>
                  {data.exactMatches.w.map((entry, index) => (
                    <JishoEntryCard key={`w-${index}`} entry={entry} type="chinese-word" />
                  ))}
                </div>

                {/* Japanese Words (Right Column, Position 1) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Japanese Words
                    {data.exactMatches.j.length > 0 && (
                      <Badge variant="outline">{data.exactMatches.j.length}</Badge>
                    )}
                  </h3>
                  {data.exactMatches.j.map((entry, index) => (
                    <JishoEntryCard key={`j-${index}`} entry={entry} type="japanese-word" />
                  ))}
                </div>

                {/* Chinese Characters (Left Column, Position 2) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Chinese Characters
                    {data.exactMatches.c.length > 0 && (
                      <Badge variant="outline">{data.exactMatches.c.length}</Badge>
                    )}
                  </h3>
                  {data.exactMatches.c.map((entry, index) => (
                    <JishoEntryCard key={`c-${index}`} entry={entry} type="chinese-char" />
                  ))}
                </div>

                {/* Kanji (Right Column, Position 2) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Kanji
                    {data.exactMatches.d.length > 0 && (
                      <Badge variant="outline">{data.exactMatches.d.length}</Badge>
                    )}
                  </h3>
                  {data.exactMatches.d.map((entry, index) => (
                    <JishoEntryCard key={`d-${index}`} entry={entry} type="kanji" />
                  ))}
                </div>
              </div>

              {/* Japanese Names - Separate four-column section at bottom */}
              {data.exactMatches.n.length > 0 && (
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Japanese Names
                    <Badge variant="outline">{data.exactMatches.n.length}</Badge>
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    {data.exactMatches.n.map((entry, index) => (
                      <JishoEntryCard key={`n-${index}`} entry={entry} type="japanese-name" />
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Contained Matches */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              Words Containing &ldquo;{word}&rdquo;
              {containedMatchCount > 0 && (
                <Badge variant="secondary">{containedMatchCount}</Badge>
              )}
            </CardTitle>
            {data.containedMatchesPending && (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" />
                Loading more...
              </div>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {containedMatchCount === 0 && !data.containedMatchesPending ? (
            <p className="text-muted-foreground text-center py-8">
              No words containing &ldquo;{word}&rdquo; found
            </p>
          ) : (
            <div className="space-y-6">
              {/* Two-column section: Words first, then Characters */}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Chinese Words (Left Column, Position 1) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Chinese Words
                    {data.containedMatches.w.length > 0 && (
                      <Badge variant="outline">{data.containedMatches.w.length}</Badge>
                    )}
                  </h3>
                  {data.containedMatches.w.map((entry, index) => (
                    <JishoEntryCard key={`w-contained-${index}`} entry={entry} type="chinese-word" />
                  ))}
                </div>

                {/* Japanese Words (Right Column, Position 1) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Japanese Words
                    {data.containedMatches.j.length > 0 && (
                      <Badge variant="outline">{data.containedMatches.j.length}</Badge>
                    )}
                  </h3>
                  {data.containedMatches.j.map((entry, index) => (
                    <JishoEntryCard key={`j-contained-${index}`} entry={entry} type="japanese-word" />
                  ))}
                </div>

                {/* Chinese Characters (Left Column, Position 2) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Chinese Characters
                    {data.containedMatches.c.length > 0 && (
                      <Badge variant="outline">{data.containedMatches.c.length}</Badge>
                    )}
                  </h3>
                  {data.containedMatches.c.map((entry, index) => (
                    <JishoEntryCard key={`c-contained-${index}`} entry={entry} type="chinese-char" />
                  ))}
                </div>

                {/* Kanji (Right Column, Position 2) */}
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Kanji
                    {data.containedMatches.d.length > 0 && (
                      <Badge variant="outline">{data.containedMatches.d.length}</Badge>
                    )}
                  </h3>
                  {data.containedMatches.d.map((entry, index) => (
                    <JishoEntryCard key={`d-contained-${index}`} entry={entry} type="kanji" />
                  ))}
                </div>
              </div>

              {/* Japanese Names - Separate four-column section at bottom */}
              {data.containedMatches.n.length > 0 && (
                <div className="space-y-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    Japanese Names
                    <Badge variant="outline">{data.containedMatches.n.length}</Badge>
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    {data.containedMatches.n.map((entry, index) => (
                      <JishoEntryCard key={`n-contained-${index}`} entry={entry} type="japanese-name" />
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
