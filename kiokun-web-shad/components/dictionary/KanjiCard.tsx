"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface KanjiEntry {
  id: string;
  c: string;
  on?: string[];
  kun?: string[];
  meanings?: string[];
  grade?: number;
  strokes?: number;
  frequency?: number;
  jlpt?: string;
}

interface KanjiCardProps {
  entry: KanjiEntry;
  isExactMatch?: boolean;
}

export function KanjiCard({ entry, isExactMatch = false }: KanjiCardProps) {
  return (
    <Card className="w-full">
      <CardHeader className="pb-0">
        <CardTitle className="text-3xl font-bold flex items-center gap-2">
          {entry.c}
          {isExactMatch && (
            <Badge variant="default" className="ml-2">Exact Match</Badge>
          )}
        </CardTitle>
      </CardHeader>

      <CardContent className="pt-2 pb-2">
        <div className="space-y-3">
          {/* Readings */}
          {entry.on && entry.on.length > 0 && (
            <div>
              <h4 className="font-medium text-sm mb-1">On reading:</h4>
              <p className="text-sm text-muted-foreground">{entry.on.join(", ")}</p>
            </div>
          )}

          {entry.kun && entry.kun.length > 0 && (
            <div>
              <h4 className="font-medium text-sm mb-1">Kun reading:</h4>
              <p className="text-sm text-muted-foreground">{entry.kun.join(", ")}</p>
            </div>
          )}

          {/* Meanings */}
          {entry.meanings && entry.meanings.length > 0 && (
            <div>
              <h4 className="font-medium text-sm mb-1">Meanings:</h4>
              <p className="text-sm text-muted-foreground">{entry.meanings.join(", ")}</p>
            </div>
          )}

          {/* Additional info */}
          <div className="flex flex-wrap gap-2 text-xs">
            {entry.grade && (
              <Badge variant="outline">
                Grade {entry.grade}
              </Badge>
            )}
            {entry.strokes && (
              <Badge variant="outline">
                {entry.strokes} strokes
              </Badge>
            )}
            {entry.frequency && (
              <Badge variant="outline">
                Freq. {entry.frequency}
              </Badge>
            )}
            {entry.jlpt && (
              <Badge variant="outline">
                JLPT {entry.jlpt}
              </Badge>
            )}

          </div>
        </div>
      </CardContent>
    </Card>
  );
}
