'use client';

import { extractShardType } from '@/lib/dictionary-utils';

interface JishoEntryCardProps {
  entry: any;
}

/**
 * A component that displays dictionary entries in a 1010-like format
 * - Efficient use of vertical space
 * - Star indicators for common words/readings
 * - Consolidated part of speech labels
 * - Clean badge design
 */
export default function JishoEntryCard({ entry }: JishoEntryCardProps) {
  // Determine the entry type
  const entryType = getEntryType(entry);

  return (
    <div className="border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow mb-4 bg-[#1e1e1e] text-white">
      {/* Display the appropriate content based on entry type */}
      {entryType === 'jmdict' && <JmdictEntry entry={entry} />}
      {entryType === 'jmnedict' && <JmnedictEntry entry={entry} />}
      {entryType === 'kanjidic' && <KanjidicEntry entry={entry} />}
      {entryType === 'chinese_char' && <ChineseCharEntry entry={entry} />}
      {entryType === 'chinese_word' && <ChineseWordEntry entry={entry} />}
      {entryType === 'unknown' && <UnknownEntry entry={entry} />}
    </div>
  );
}

/**
 * Determines the type of dictionary entry
 */
function getEntryType(entry: any): string {
  // Check for lowercase field names (new format)
  if (entry.kanji && entry.kana && entry.sense) return 'jmdict';
  if (entry.k && entry.r && entry.m) return 'jmnedict';
  if (entry.c && entry.on) return 'kanjidic';
  if (entry.traditional && entry.simplified && !entry.traditional_chinese) return 'chinese_char';
  if (entry.traditional_chinese || (entry.traditional && entry.pinyin)) return 'chinese_word';

  // Check for uppercase field names (old format)
  if (entry.Kanji && entry.Kana && entry.Sense) return 'jmdict';
  if (entry.Kanji && entry.Reading && entry.Translation) return 'jmnedict';
  if (entry.Character && entry.Reading && entry.Misc) return 'kanjidic';
  if (entry.Traditional && entry.Simplified && !entry.Components) return 'chinese_char';
  if (entry.Traditional && entry.Pinyin && entry.Components) return 'chinese_word';

  return 'unknown';
}

/**
 * Renders a JMdict word entry in 1010-style format
 */
function JmdictEntry({ entry }: { entry: any }) {
  // Extract ID for display
  const id = entry.id || entry.ID;
  const shardType = id ? extractShardType(id) : null;

  // Get kanji and kana readings
  const kanjiList = entry.kanji || entry.Kanji || [];
  const kanaList = entry.kana || entry.Kana || [];
  const senseList = entry.sense || entry.Sense || [];

  // Group senses by part of speech to avoid repetition
  const sensesByPos: Record<string, any[]> = {};
  senseList.forEach((sense: any) => {
    const partOfSpeech = sense.partOfSpeech || sense.PartOfSpeech || [];
    const posKey = partOfSpeech.length > 0 ? partOfSpeech.join(',') : 'other';

    if (!sensesByPos[posKey]) {
      sensesByPos[posKey] = [];
    }
    sensesByPos[posKey].push(sense);
  });

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-blue-600 text-white text-xs font-medium px-2.5 py-0.5 rounded">
          Japanese Word
        </span>
        {id && (
          <span className="text-gray-400 text-xs">
            ID: {id}
            {shardType !== null && ` (Shard: ${shardType})`}
          </span>
        )}
      </div>

      {/* Main content */}
      <div className="flex flex-col md:flex-row gap-4">
        {/* Left column: Word and definitions */}
        <div className="flex-1">
          {/* Word display */}
          <div className="mb-4">
            {/* Kanji and kana on the same line */}
            <div className="flex flex-wrap items-baseline gap-x-4 gap-y-2">
              {/* Kanji writings with star indicators */}
              {kanjiList.map((k: any, i: number) => (
                <div key={`kanji-${i}`} className="flex items-center">
                  <span className="text-2xl font-bold text-[#3b9cff]">
                    {k.text || k.Text}
                    {(k.common || k.Common) && (
                      <span className="text-[#3b9cff] ml-1">★</span>
                    )}
                  </span>

                  {/* Corresponding kana if available */}
                  {kanaList[i] && (
                    <span className="text-xl text-green-400 ml-2">
                      {kanaList[i].text || kanaList[i].Text}
                      {(kanaList[i].common || kanaList[i].Common) && (
                        <span className="text-green-400 ml-1">★</span>
                      )}
                    </span>
                  )}
                </div>
              ))}

              {/* Additional kana readings that don't match with kanji */}
              {kanaList.slice(kanjiList.length).map((k: any, i: number) => (
                <span key={`kana-extra-${i}`} className="text-xl text-green-400">
                  {k.text || k.Text}
                  {(k.common || k.Common) && (
                    <span className="text-green-400 ml-1">★</span>
                  )}
                </span>
              ))}
            </div>
          </div>

          {/* Definitions grouped by part of speech */}
          <div className="space-y-4">
            {Object.entries(sensesByPos).map(([posKey, senses], groupIndex) => {
              return (
                <div key={`pos-group-${groupIndex}`} className={groupIndex > 0 ? "pt-2 border-t border-gray-700" : ""}>
                  {/* Part of speech badge - shown once per group */}
                  {posKey !== 'other' && (
                    <div className="inline-block bg-gray-700 text-white text-xs px-2 py-0.5 rounded mb-2">
                      {posKey}
                    </div>
                  )}

                  {/* Definitions in this part of speech group */}
                  <div className="space-y-2">
                    {senses.map((sense: any, senseIndex: number) => {
                      const glossList = sense.gloss || sense.Gloss || [];
                      const examples = sense.examples || [];
                      const contexts = sense.context || sense.Context || [];
                      const fields = sense.field || sense.Field || [];
                      const dialects = sense.dialect || sense.Dialect || [];

                      // Calculate the overall index for numbering
                      let overallIndex = 0;
                      Object.entries(sensesByPos).forEach(([, list], idx) => {
                        if (idx < groupIndex) {
                          overallIndex += list.length;
                        }
                      });
                      overallIndex += senseIndex + 1;

                      return (
                        <div key={`sense-${groupIndex}-${senseIndex}`} className="text-white">
                          {/* Definition with inline context badges */}
                          <div className="mb-2">
                            <span className="font-medium text-white">{overallIndex}. </span>

                            {/* Inline context badges */}
                            <span className="inline-flex flex-wrap gap-1 mr-1">
                              {contexts.map((ctx: string, i: number) => (
                                <span key={`ctx-${i}`} className="bg-purple-700 text-white text-xs px-1.5 py-0.5 rounded">
                                  {ctx}
                                </span>
                              ))}
                              {fields.map((field: string, i: number) => (
                                <span key={`field-${i}`} className="bg-green-700 text-white text-xs px-1.5 py-0.5 rounded">
                                  {field}
                                </span>
                              ))}
                              {dialects.map((dialect: string, i: number) => (
                                <span key={`dialect-${i}`} className="bg-yellow-700 text-white text-xs px-1.5 py-0.5 rounded">
                                  {dialect}
                                </span>
                              ))}
                            </span>

                            {/* Definition text */}
                            {glossList.map((g: any, j: number) => (
                              <span key={`gloss-${groupIndex}-${senseIndex}-${j}`}>
                                {j > 0 && ', '}
                                {typeof g === 'string' ? g : g.text || g.Text}
                              </span>
                            ))}
                          </div>

                          {/* Examples */}
                          {examples && examples.length > 0 && (
                            <div className="bg-gray-800 p-2 rounded text-sm mb-2">
                              {examples.map((ex: any, k: number) => {
                                const japSentence = ex.sentences?.find((s: any) => !s.lang || s.lang === '')?.text;
                                const engSentence = ex.sentences?.find((s: any) => s.lang === 'eng')?.text ||
                                  ex.translation || '';

                                return (
                                  <div key={`example-${groupIndex}-${senseIndex}-${k}`} className={k > 0 ? "mt-2 pt-2 border-t border-gray-700" : ""}>
                                    <div className="font-medium">{ex.text}</div>
                                    {japSentence && <div className="text-gray-300">{japSentence}</div>}
                                    {engSentence && <div className="text-gray-400 italic">{engSentence}</div>}
                                  </div>
                                );
                              })}
                            </div>
                          )}
                        </div>
                      );
                    })}
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Right column: Kanji information (if available) */}
        {kanjiList.length > 0 && kanjiList[0].text && (
          <div className="w-full md:w-1/4 bg-gray-800 p-3 rounded-lg">
            <h3 className="text-lg font-semibold mb-2 text-white">Kanji</h3>
            <div className="flex flex-wrap gap-2">
              {Array.from(new Set(kanjiList[0].text.split('')))
                .filter(char => /[\u4e00-\u9faf]/.test(char))
                .map((kanji, i) => (
                  <a
                    key={`kanji-link-${i}`}
                    href={`/word/${kanji}`}
                    className="block text-center p-2 border border-gray-700 rounded bg-gray-900 hover:bg-gray-700 transition-colors"
                  >
                    <div className="text-2xl mb-1 text-white">{kanji}</div>
                  </a>
                ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders a JMnedict name entry in 1010-style format
 */
function JmnedictEntry({ entry }: { entry: any }) {
  // Extract ID for display
  const id = entry.id || entry.ID;
  const shardType = id ? extractShardType(id) : null;

  // Get kanji and reading
  const kanjiList = entry.k || entry.Kanji || [];
  const readingList = entry.r || entry.Reading || [];
  const meaningList = entry.m || entry.Meanings || [];
  const typeList = entry.type || entry.Type || [];

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-purple-600 text-white text-xs font-medium px-2.5 py-0.5 rounded">
          Japanese Name
        </span>
        {id && (
          <span className="text-gray-400 text-xs">
            ID: {id}
            {shardType !== null && ` (Shard: ${shardType})`}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Name display with inline readings */}
        <div className="mb-4">
          <div className="flex flex-wrap items-baseline gap-x-4 gap-y-2">
            {/* Kanji writings with corresponding readings */}
            {kanjiList.map((k: string, i: number) => (
              <div key={`kanji-${i}`} className="flex items-center">
                <span className="text-2xl font-bold text-[#3b9cff]">{k}</span>

                {/* Corresponding reading if available */}
                {readingList[i] && (
                  <span className="text-xl text-green-400 ml-2">{readingList[i]}</span>
                )}
              </div>
            ))}

            {/* Additional readings that don't match with kanji */}
            {readingList.slice(kanjiList.length).map((r: string, i: number) => (
              <span key={`reading-extra-${i}`} className="text-xl text-green-400">{r}</span>
            ))}
          </div>
        </div>

        {/* Meanings with inline type badges */}
        {meaningList.length > 0 && (
          <div className="text-white">
            <div className="inline-block bg-gray-700 text-white text-xs px-2 py-0.5 rounded mb-2">
              name
            </div>
            <div className="space-y-1">
              {meaningList.map((meaning: string, i: number) => (
                <div key={`meaning-${i}`}>
                  <span className="font-medium">{i + 1}. </span>

                  {/* Inline type badges */}
                  {i === 0 && typeList.length > 0 && (
                    <span className="inline-flex flex-wrap gap-1 mr-1">
                      {typeList.map((type: string, j: number) => (
                        <span key={`type-${j}`} className="bg-purple-700 text-white text-xs px-1.5 py-0.5 rounded">
                          {type}
                        </span>
                      ))}
                    </span>
                  )}

                  {meaning}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders a Kanjidic kanji entry in 1010-style format
 */
function KanjidicEntry({ entry }: { entry: any }) {
  // Extract data
  const character = entry.c || entry.Character;
  const id = entry.id || entry.ID;
  const meanings = entry.m || entry.Meanings || [];
  const onYomi = entry.on || (entry.Reading?.OnYomi) || [];
  const kunYomi = entry.kun || (entry.Reading?.KunYomi) || [];
  const grade = entry.grade || (entry.Misc?.Grade);
  const stroke = entry.stroke || (entry.Misc?.StrokeCount);
  const frequency = entry.freq || (entry.Misc?.Frequency);
  const jlpt = entry.jlpt || (entry.Misc?.JLPT);
  const ids = entry.ids || entry.IDS;
  const nanori = entry.nanori || (entry.Reading?.Nanori) || [];

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-green-600 text-white text-xs font-medium px-2.5 py-0.5 rounded">
          Kanji Character
        </span>
        {id && (
          <span className="text-gray-400 text-xs">
            ID: {id}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Character display with inline readings */}
        <div className="flex flex-wrap gap-6 mb-4">
          <div className="flex items-start gap-4">
            <div className="text-6xl font-bold text-[#3b9cff]">{character}</div>

            {/* Readings next to character */}
            <div className="flex flex-col gap-2">
              {onYomi.length > 0 && (
                <div>
                  <div className="text-sm font-medium text-gray-400">On</div>
                  <div className="flex flex-wrap gap-2">
                    {onYomi.map((reading: string, i: number) => (
                      <span key={`on-${i}`} className="text-xl text-green-400">
                        {reading}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {kunYomi.length > 0 && (
                <div>
                  <div className="text-sm font-medium text-gray-400">Kun</div>
                  <div className="flex flex-wrap gap-2">
                    {kunYomi.map((reading: string, i: number) => (
                      <span key={`kun-${i}`} className="text-xl text-green-400">
                        {reading}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="flex flex-col gap-1">
            {grade && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">Grade:</span>
                <span className="text-white">{grade}</span>
              </div>
            )}
            {stroke && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">Strokes:</span>
                <span className="text-white">{stroke}</span>
              </div>
            )}
            {frequency && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">Frequency:</span>
                <span className="text-white">{frequency}</span>
              </div>
            )}
            {jlpt && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">JLPT:</span>
                <span className="text-white">N{jlpt}</span>
              </div>
            )}
          </div>
        </div>

        {/* Nanori readings if available */}
        {nanori.length > 0 && (
          <div className="mb-4">
            <div className="text-sm font-medium text-gray-400 mb-1">Name Readings</div>
            <div className="flex flex-wrap gap-2">
              {nanori.map((reading: string, i: number) => (
                <span key={`nanori-${i}`} className="text-green-400">
                  {reading}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Meanings */}
        {meanings.length > 0 && (
          <div className="mb-4">
            <div className="inline-block bg-gray-700 text-white text-xs px-2 py-0.5 rounded mb-2">
              noun
            </div>
            <div className="text-white space-y-1">
              {meanings.map((meaning: string, i: number) => {
                // Check if meaning has a field/domain in parentheses
                const match = meaning.match(/^(.*?)\s*\((.*?)\)\s*(.*)$/);
                if (match) {
                  // eslint-disable-next-line @typescript-eslint/no-unused-vars
                  const [_, prefix, domain, suffix] = match;
                  return (
                    <div key={`meaning-${i}`}>
                      <span className="font-medium">{i + 1}. </span>
                      <span className="bg-green-700 text-white text-xs px-1.5 py-0.5 rounded mr-1">
                        {domain}
                      </span>
                      {prefix}{suffix}
                    </div>
                  );
                } else {
                  return (
                    <div key={`meaning-${i}`}>
                      <span className="font-medium">{i + 1}. </span>
                      {meaning}
                    </div>
                  );
                }
              })}
            </div>
          </div>
        )}

        {/* IDS data */}
        {ids && (
          <div>
            <div className="text-sm font-medium text-gray-400 mb-1">Character Composition (IDS)</div>
            <div className="text-white">{ids}</div>
            <div className="text-xs text-gray-500 mt-1">
              IDS (Ideographic Description Sequence) shows how the character is composed
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders a Chinese character entry in 1010-style format
 */
function ChineseCharEntry({ entry }: { entry: any }) {
  // Extract data
  const id = entry.id || entry.ID;
  const traditional = entry.traditional || entry.Traditional;
  const simplified = entry.simplified || entry.Simplified;
  const definitions = entry.definitions || entry.Definitions || [];
  const pinyin = entry.pinyin || entry.Pinyin || [];
  const strokeCount = entry.strokeCount || entry.StrokeCount;
  const ids = entry.ids || entry.IDS;
  const frequency = entry.frequency || entry.Frequency;

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-red-600 text-white text-xs font-medium px-2.5 py-0.5 rounded">
          Chinese Word
        </span>
        {id && (
          <span className="text-gray-400 text-xs">
            ID: {id}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Character display with inline pinyin */}
        <div className="flex flex-wrap gap-6 mb-4">
          <div className="flex items-baseline">
            <div>
              <div className="text-sm font-medium text-gray-400">Traditional</div>
              <div className="text-4xl font-bold text-[#3b9cff]">{traditional}</div>
            </div>

            {/* Pinyin next to character */}
            {pinyin.length > 0 && (
              <div className="ml-3 flex flex-col">
                <div className="text-sm font-medium text-gray-400">Pinyin</div>
                <div className="flex flex-wrap gap-2">
                  {pinyin.map((p: string, i: number) => (
                    <span key={`pinyin-${i}`} className="text-xl text-green-400">
                      {p}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>

          {simplified !== traditional && (
            <div className="flex items-baseline">
              <div>
                <div className="text-sm font-medium text-gray-400">Simplified</div>
                <div className="text-4xl font-bold text-[#3b9cff]">{simplified}</div>
              </div>
            </div>
          )}

          <div className="flex flex-col justify-end gap-1">
            {strokeCount && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">Strokes:</span>
                <span className="text-white">{strokeCount}</span>
              </div>
            )}
            {frequency && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">Frequency:</span>
                <span className="text-white">{frequency}</span>
              </div>
            )}
          </div>
        </div>

        {/* Definitions */}
        {definitions.length > 0 && (
          <div className="mb-4">
            <div className="inline-block bg-gray-700 text-white text-xs px-2 py-0.5 rounded mb-2">
              noun
            </div>
            <div className="space-y-1 text-white">
              {definitions.map((def: string, i: number) => {
                // Check if definition has a field/domain in parentheses
                const match = def.match(/^(.*?)\s*\((.*?)\)\s*(.*)$/);
                if (match) {
                  // eslint-disable-next-line @typescript-eslint/no-unused-vars
                  const [_, prefix, domain, suffix] = match;
                  return (
                    <div key={`def-${i}`}>
                      <span className="font-medium">{i + 1}. </span>
                      <span className="bg-green-700 text-white text-xs px-1.5 py-0.5 rounded mr-1">
                        {domain}
                      </span>
                      {prefix}{suffix}
                    </div>
                  );
                } else {
                  return (
                    <div key={`def-${i}`}>
                      <span className="font-medium">{i + 1}. </span>
                      {def}
                    </div>
                  );
                }
              })}
            </div>
          </div>
        )}

        {/* IDS data */}
        {ids && (
          <div>
            <div className="text-sm font-medium text-gray-400 mb-1">Character Composition (IDS)</div>
            <div className="text-white">{ids}</div>
            <div className="text-xs text-gray-500 mt-1">
              IDS (Ideographic Description Sequence) shows how the character is composed
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders a Chinese word entry in 1010-style format
 */
function ChineseWordEntry({ entry }: { entry: any }) {
  // Extract data
  const id = entry.id || entry.ID;
  const traditional = entry.traditional_chinese || entry.traditional || entry.Traditional;
  const simplified = entry.simplified_chinese || entry.simplified || entry.Simplified;
  const definitions = entry.definitions || entry.Definitions || [];
  const pinyin = entry.pinyin || entry.Pinyin || [];
  const hskLevel = entry.hskLevel || entry.HskLevel;
  const frequency = entry.frequency || entry.Frequency;
  const isCommon = entry.common || entry.Common;

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-yellow-600 text-white text-xs font-medium px-2.5 py-0.5 rounded">
          Chinese Word
        </span>
        {id && (
          <span className="text-gray-400 text-xs">
            ID: {id}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Word display with inline pinyin */}
        <div className="flex flex-wrap gap-6 mb-4">
          <div className="flex items-baseline">
            <div>
              <div className="text-sm font-medium text-gray-400">Traditional</div>
              <div className="text-3xl font-bold text-[#3b9cff]">
                {traditional}
                {isCommon && <span className="text-[#3b9cff] ml-1">★</span>}
              </div>
            </div>

            {/* Pinyin next to word */}
            {pinyin.length > 0 && (
              <div className="ml-3 flex flex-col">
                <div className="text-sm font-medium text-gray-400">Pinyin</div>
                <div className="flex flex-wrap gap-2">
                  {pinyin.map((p: string, i: number) => (
                    <span key={`pinyin-${i}`} className="text-xl text-green-400">
                      {p}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>

          {simplified !== traditional && (
            <div className="flex items-baseline">
              <div>
                <div className="text-sm font-medium text-gray-400">Simplified</div>
                <div className="text-3xl font-bold text-[#3b9cff]">
                  {simplified}
                  {isCommon && <span className="text-[#3b9cff] ml-1">★</span>}
                </div>
              </div>
            </div>
          )}

          <div className="flex flex-col justify-end gap-1">
            {hskLevel && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">HSK:</span>
                <span className="text-white">{hskLevel}</span>
              </div>
            )}
            {frequency && (
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-400 mr-2">Frequency:</span>
                <span className="text-white">{frequency}</span>
              </div>
            )}
          </div>
        </div>

        {/* Definitions */}
        {definitions.length > 0 && (
          <div>
            <div className="inline-block bg-gray-700 text-white text-xs px-2 py-0.5 rounded mb-2">
              noun
            </div>
            <div className="space-y-1 text-white">
              {definitions.map((def: string, i: number) => {
                // Check if definition has a field/domain in parentheses
                const match = def.match(/^(.*?)\s*\((.*?)\)\s*(.*)$/);
                if (match) {
                  // eslint-disable-next-line @typescript-eslint/no-unused-vars
                  const [_, prefix, domain, suffix] = match;
                  return (
                    <div key={`def-${i}`}>
                      <span className="font-medium">{i + 1}. </span>
                      <span className="bg-green-700 text-white text-xs px-1.5 py-0.5 rounded mr-1">
                        {domain}
                      </span>
                      {prefix}{suffix}
                    </div>
                  );
                } else {
                  return (
                    <div key={`def-${i}`}>
                      <span className="font-medium">{i + 1}. </span>
                      {def}
                    </div>
                  );
                }
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders an unknown entry type in 1010-style format
 */
function UnknownEntry({ entry }: { entry: any }) {
  return (
    <div>
      <div className="flex justify-between items-start mb-2">
        <span className="bg-gray-600 text-white text-xs font-medium px-2.5 py-0.5 rounded">
          Unknown Entry Type
        </span>
      </div>
      <pre className="bg-gray-800 p-2 rounded overflow-auto text-sm text-gray-300 border border-gray-700">
        {JSON.stringify(entry, null, 2)}
      </pre>
    </div>
  );
}
