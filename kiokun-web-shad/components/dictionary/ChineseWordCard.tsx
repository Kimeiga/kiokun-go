"use client";

import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface ChineseWordEntry {
  id: string;
  traditional: string;
  simplified?: string;
  pinyin?: string[];
  definitions: string[];
  hskLevel?: number;
}

interface ChineseWordCardProps {
  entry: ChineseWordEntry;
  isExactMatch?: boolean;
}

export function ChineseWordCard({ entry, isExactMatch = false }: ChineseWordCardProps) {
  return (
    <Card className="w-full">
      <CardHeader className="pb-0">
        <CardTitle className="text-xl font-bold flex items-center gap-2">
          <span className="flex items-center gap-1">
            <Link
              href={`/word/${encodeURIComponent(entry.traditional)}`}
              className="hover:text-primary transition-colors"
            >
              {entry.traditional}
            </Link>
          </span>
          {entry.simplified && entry.simplified !== entry.traditional && (
            <span className="text-lg text-muted-foreground">
              ({entry.simplified})
            </span>
          )}
          {isExactMatch && (
            <Badge variant="default" className="ml-2">Exact Match</Badge>
          )}
        </CardTitle>

        {entry.pinyin && entry.pinyin.length > 0 && (
          <div className="text-sm text-muted-foreground mt-1">
            <span className="font-medium">Pinyin:</span> {entry.pinyin.join(", ")}
          </div>
        )}
      </CardHeader>

      <CardContent className="pt-2 pb-2">
        <div className="space-y-3">
          {/* Definitions */}
          {entry.definitions && entry.definitions.length > 0 && (
            <div className="space-y-1">
              {entry.definitions.map((definition, index) => (
                <div key={index} className="flex items-start gap-2">
                  {/* Show number only if multiple definitions */}
                  {entry.definitions.length > 1 && (
                    <span className="text-sm font-medium min-w-[1.5rem]">
                      {index + 1}.
                    </span>
                  )}
                  <span className="text-sm">{definition}</span>
                </div>
              ))}
            </div>
          )}

          {/* Additional info */}
          {entry.hskLevel && (
            <div className="flex flex-wrap gap-2 text-xs">
              <Badge variant="outline" className="text-xs">
                HSK {entry.hskLevel}
              </Badge>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
