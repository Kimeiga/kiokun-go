"use client";

import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface JapaneseNameEntry {
  id: string;
  k?: string[];
  r?: string[];
  m?: string[];
  type?: string[];
}

interface JapaneseNameCardProps {
  entry: JapaneseNameEntry;
  isExactMatch?: boolean;
}

export function JapaneseNameCard({ entry, isExactMatch = false }: JapaneseNameCardProps) {
  const kanjiText = entry.k?.[0] || "";
  const readingText = entry.r?.[0] || "";
  const meaningText = entry.m?.[0] || "";

  return (
    <Card className="w-full">
      <CardHeader className="pb-0">
        <CardTitle className="text-xl font-bold flex items-center gap-2">
          {kanjiText && (
            <span className="flex items-center gap-1">
              <Link
                href={`/word/${encodeURIComponent(kanjiText)}`}
                className="hover:text-primary transition-colors"
              >
                {kanjiText}
              </Link>
            </span>
          )}
          {readingText && (
            <span className="text-lg text-muted-foreground">
              ({readingText})
            </span>
          )}
          {isExactMatch && (
            <Badge variant="default" className="ml-2">Exact Match</Badge>
          )}
        </CardTitle>

        {meaningText && (
          <div className="text-sm text-muted-foreground mt-1">
            <span className="font-medium">Meaning:</span> {meaningText}
          </div>
        )}
      </CardHeader>

      <CardContent className="pt-2 pb-2">
        <div className="space-y-3">
          {/* Name types */}
          {entry.type && entry.type.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {entry.type.map((type, index) => (
                <Badge key={index} variant="outline" className="text-xs capitalize">
                  {type}
                </Badge>
              ))}
            </div>
          )}

          {/* Additional readings/meanings */}
          {entry.r && entry.r.length > 1 && (
            <div>
              <h4 className="font-medium text-sm mb-1">Other readings:</h4>
              <p className="text-sm text-muted-foreground">
                {entry.r.slice(1).join(", ")}
              </p>
            </div>
          )}

          {entry.m && entry.m.length > 1 && (
            <div>
              <h4 className="font-medium text-sm mb-1">Other meanings:</h4>
              <p className="text-sm text-muted-foreground">
                {entry.m.slice(1).join(", ")}
              </p>
            </div>
          )}


        </div>
      </CardContent>
    </Card>
  );
}
