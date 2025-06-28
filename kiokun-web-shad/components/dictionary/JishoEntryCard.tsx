"use client";

import { ChineseCharacterCard } from "./ChineseCharacterCard";
import { ChineseWordCard } from "./ChineseWordCard";
import { JapaneseWordCard } from "./JapaneseWordCard";
import { KanjiCard } from "./KanjiCard";
import { JapaneseNameCard } from "./JapaneseNameCard";

interface JishoEntryCardProps {
  entry: any;
  type: string;
  isExactMatch?: boolean;
}

/**
 * Main entry card component that routes to the appropriate specialized card
 * based on the entry type
 */
export function JishoEntryCard({ entry, type, isExactMatch = false }: JishoEntryCardProps) {
  switch (type) {
    case 'japanese-word':
      return <JapaneseWordCard entry={entry} isExactMatch={isExactMatch} />;
    case 'japanese-name':
      return <JapaneseNameCard entry={entry} isExactMatch={isExactMatch} />;
    case 'kanji':
      return <KanjiCard entry={entry} isExactMatch={isExactMatch} />;
    case 'chinese-char':
      return <ChineseCharacterCard entry={entry} isExactMatch={isExactMatch} />;
    case 'chinese-word':
      return <ChineseWordCard entry={entry} isExactMatch={isExactMatch} />;
    default:
      return <UnknownEntryCard entry={entry} type={type} isExactMatch={isExactMatch} />;
  }
}

/**
 * Fallback component for unknown entry types
 */
function UnknownEntryCard({ entry, type, isExactMatch }: { entry: any; type: string; isExactMatch?: boolean }) {
  return (
    <div className="p-4 border rounded-lg bg-muted">
      <div className="flex items-center justify-between mb-2">
        <h3 className="font-semibold">Unknown Entry Type: {type}</h3>
        {isExactMatch && (
          <span className="text-xs bg-primary text-primary-foreground px-2 py-1 rounded">
            Exact Match
          </span>
        )}
      </div>
      <pre className="text-xs text-muted-foreground overflow-auto">
        {JSON.stringify(entry, null, 2)}
      </pre>
    </div>
  );
}

// Export the main component as default for backward compatibility
export default JishoEntryCard;
