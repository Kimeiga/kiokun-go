"use client";

import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { mapPartsOfSpeech, mapMiscTag } from "@/lib/pos-mappings";

interface JapaneseWordEntry {
  id: string;
  kanji?: Array<{ text: string; common?: boolean }>;
  kana?: Array<{ text: string; common?: boolean; appliesToKanji?: string[] }>;
  sense?: Array<{
    partOfSpeech?: string[];
    field?: string[];        // Domain/field tags like "sumo", "comp", "med"
    misc?: string[];         // Miscellaneous tags like "abbr", "arch", "vulg"
    dialect?: string[];      // Dialect tags like "ksb", "osb"
    info?: string[];         // Additional information
    gloss?: Array<{ text: string; lang?: string }>;
    examples?: Array<{
      text: string;
      sentences: Array<{ text: string; lang?: string }>;
    }>;
  }>;
}

interface JapaneseWordCardProps {
  entry: JapaneseWordEntry;
  isExactMatch?: boolean;
}

export function JapaneseWordCard({ entry, isExactMatch = false }: JapaneseWordCardProps) {
  const kanjiText = entry.kanji?.[0]?.text || "";
  const kanaText = entry.kana?.[0]?.text || "";
  const isCommon = entry.kanji?.[0]?.common || entry.kana?.[0]?.common;

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
              {isCommon && <span className="text-yellow-500">★</span>}
            </span>
          )}
          {kanaText && (
            <span className="text-lg text-muted-foreground">
              {kanaText}
              {isCommon && !kanjiText && <span className="text-yellow-500">★</span>}
            </span>
          )}
          {isExactMatch && (
            <Badge variant="default" className="ml-2">Exact Match</Badge>
          )}
        </CardTitle>
      </CardHeader>

      <CardContent className="pt-2 pb-2">
        <div className="space-y-2">
          {/* Process senses into part-of-speech groups */}
          {(() => {
            if (!entry.sense) return null;

            // Step 1: Group consecutive senses by part of speech
            const posGroups: Array<{
              partOfSpeechKey: string;
              allPartsOfSpeech: string[];
              senseDefinitions: Array<{ text: string; inlineTags: Array<{ text: string; type: 'field' | 'misc' | 'dialect' | 'info' }> }>;
            }> = [];

            entry.sense.forEach(sense => {
              if (!sense.gloss || sense.gloss.length === 0) return;

              // Use all parts of speech, joined with space for grouping key
              const rawPartsOfSpeech = sense.partOfSpeech || ['unknown'];
              const posKey = rawPartsOfSpeech.join(' ');
              const allPartsOfSpeech = mapPartsOfSpeech(rawPartsOfSpeech);
              const senseDefinitionString = sense.gloss.map(gloss => gloss.text).join('; ');

              // Collect inline tags with their types for styling
              const inlineTags: Array<{ text: string; type: 'field' | 'misc' | 'dialect' | 'info' }> = [
                ...(sense.field || []).map(tag => ({ text: tag, type: 'field' as const })),
                ...(sense.misc || []).map(tag => ({ text: mapMiscTag(tag), type: 'misc' as const })),
                ...(sense.dialect || []).map(tag => ({ text: tag, type: 'dialect' as const })),
                ...(sense.info || []).map(tag => ({ text: tag, type: 'info' as const }))
              ];

              // Check if we can merge with the previous group (same parts of speech)
              const lastGroup = posGroups[posGroups.length - 1];
              if (lastGroup && lastGroup.partOfSpeechKey === posKey) {
                lastGroup.senseDefinitions.push({ text: senseDefinitionString, inlineTags });
              } else {
                posGroups.push({
                  partOfSpeechKey: posKey,
                  allPartsOfSpeech: allPartsOfSpeech,
                  senseDefinitions: [{ text: senseDefinitionString, inlineTags }]
                });
              }
            });

            // Step 2: Render each part-of-speech group
            return posGroups.map((group, groupIndex) => {
              if (group.senseDefinitions.length === 1) {
                // Single sense in group: Part of speech inline
                const senseData = group.senseDefinitions[0];
                return (
                  <div key={groupIndex} className="flex items-start gap-2">
                    {/* Show group number only if multiple groups */}
                    {posGroups.length > 1 && (
                      <span className="text-sm font-medium min-w-[1.5rem]">
                        {groupIndex + 1}.
                      </span>
                    )}
                    <div className="flex items-center gap-2 flex-1 flex-wrap">
                      {/* All parts of speech badges inline */}
                      {group.allPartsOfSpeech.filter(pos => pos !== 'unknown').map((pos, posIndex) => (
                        <Badge key={posIndex} variant="outline" className="text-xs">
                          {pos}
                        </Badge>
                      ))}
                      {/* Inline tags with different colors by type */}
                      {senseData.inlineTags.map((tag, tagIndex) => {
                        const colorClass = tag.type === 'misc'
                          ? 'bg-blue-100 text-blue-800'
                          : 'bg-green-100 text-green-800';
                        return (
                          <Badge key={tagIndex} variant="secondary" className={`text-xs ${colorClass}`}>
                            {tag.text}
                          </Badge>
                        );
                      })}
                      <span className="text-sm">{senseData.text}</span>
                    </div>
                  </div>
                );
              } else {
                // Multiple senses in group: Part of speech as block, then ordered list
                return (
                  <div key={groupIndex} className="space-y-1">
                    <div className="flex items-start gap-2">
                      {/* Show group number only if multiple groups */}
                      {posGroups.length > 1 && (
                        <span className="text-sm font-medium min-w-[1.5rem]">
                          {groupIndex + 1}.
                        </span>
                      )}
                      <div className="flex flex-wrap gap-1">
                        {/* All parts of speech badges as block */}
                        {group.allPartsOfSpeech.filter(pos => pos !== 'unknown').map((pos, posIndex) => (
                          <Badge key={posIndex} variant="outline" className="text-xs">
                            {pos}
                          </Badge>
                        ))}
                      </div>
                    </div>
                    {/* Ordered list of definitions */}
                    <div className="space-y-1 ml-6">
                      {group.senseDefinitions.map((senseData, defIndex) => (
                        <div key={defIndex} className="flex items-start gap-2">
                          <span className="text-sm font-medium min-w-[1.5rem]">
                            {defIndex + 1}.
                          </span>
                          <div className="flex items-center gap-2 flex-1 flex-wrap">
                            {/* Inline tags for this specific sense with different colors */}
                            {senseData.inlineTags.map((tag, tagIndex) => {
                              const colorClass = tag.type === 'misc'
                                ? 'bg-blue-100 text-blue-800'
                                : 'bg-green-100 text-green-800';
                              return (
                                <Badge key={tagIndex} variant="secondary" className={`text-xs ${colorClass}`}>
                                  {tag.text}
                                </Badge>
                              );
                            })}
                            <span className="text-sm">{senseData.text}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              }
            });
          })()}


        </div>
      </CardContent>
    </Card>
  );
}
