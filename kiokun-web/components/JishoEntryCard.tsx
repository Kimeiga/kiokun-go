'use client';

import { extractShardType } from '@/lib/dictionary-utils';

interface JishoEntryCardProps {
  entry: any;
}

/**
 * A component that displays dictionary entries in a jisho.org-like format
 */
export default function JishoEntryCard({ entry }: JishoEntryCardProps) {
  // Determine the entry type
  const entryType = getEntryType(entry);

  return (
    <div className="border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow mb-4">
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
  if (entry.traditional_chinese && entry.simplified_chinese) return 'chinese_word';

  // Check for uppercase field names (old format)
  if (entry.Kanji && entry.Kana && entry.Sense) return 'jmdict';
  if (entry.Kanji && entry.Reading && entry.Translation) return 'jmnedict';
  if (entry.Character && entry.Reading && entry.Misc) return 'kanjidic';
  if (entry.Traditional && entry.Simplified && !entry.Components) return 'chinese_char';

  return 'unknown';
}

/**
 * Renders a JMdict word entry
 */
function JmdictEntry({ entry }: { entry: any }) {
  // Extract ID for display
  const id = entry.id || entry.ID;
  const shardType = id ? extractShardType(id) : null;

  // Get kanji and kana readings
  const kanjiList = entry.kanji || entry.Kanji || [];
  const kanaList = entry.kana || entry.Kana || [];
  const senseList = entry.sense || entry.Sense || [];

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">
          Japanese Word
        </span>
        {id && (
          <span className="text-gray-500 text-xs">
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
            {/* Kanji writings */}
            <div className="flex flex-wrap gap-2 mb-2">
              {kanjiList.map((k: any, i: number) => (
                <span key={`kanji-${i}`} className="text-2xl font-bold">
                  {k.text || k.Text}
                  {(k.common || k.Common) && (
                    <span className="text-xs text-green-600 ml-1 align-top">common</span>
                  )}
                </span>
              ))}
            </div>

            {/* Kana readings */}
            <div className="flex flex-wrap gap-2 text-gray-600">
              {kanaList.map((k: any, i: number) => (
                <span key={`kana-${i}`} className="text-xl">
                  {k.text || k.Text}
                  {(k.common || k.Common) && (
                    <span className="text-xs text-green-600 ml-1 align-top">common</span>
                  )}
                </span>
              ))}
            </div>
          </div>

          {/* Definitions */}
          <div className="space-y-4">
            {senseList.map((sense: any, i: number) => {
              const partOfSpeech = sense.partOfSpeech || sense.PartOfSpeech || [];
              const glossList = sense.gloss || sense.Gloss || [];
              const examples = sense.examples || [];

              return (
                <div key={`sense-${i}`} className={i > 0 ? "pt-2 border-t border-gray-100" : ""}>
                  {/* Part of speech */}
                  {partOfSpeech.length > 0 && (
                    <div className="text-sm text-gray-500 mb-1">
                      {partOfSpeech.join(', ')}
                    </div>
                  )}

                  {/* Definition */}
                  <div className="text-gray-800 mb-2">
                    <span className="font-medium text-gray-700">{i + 1}. </span>
                    {glossList.map((g: any, j: number) => (
                      <span key={`gloss-${i}-${j}`}>
                        {j > 0 && ', '}
                        {typeof g === 'string' ? g : g.text || g.Text}
                      </span>
                    ))}
                  </div>

                  {/* Examples */}
                  {examples && examples.length > 0 && (
                    <div className="bg-gray-50 p-2 rounded text-sm">
                      {examples.map((ex: any, k: number) => {
                        const japSentence = ex.sentences?.find((s: any) => !s.lang || s.lang === '')?.text;
                        const engSentence = ex.sentences?.find((s: any) => s.lang === 'eng')?.text ||
                          ex.translation || '';

                        return (
                          <div key={`example-${i}-${k}`} className={k > 0 ? "mt-2 pt-2 border-t border-gray-100" : ""}>
                            <div className="font-medium">{ex.text}</div>
                            {japSentence && <div className="text-gray-700">{japSentence}</div>}
                            {engSentence && <div className="text-gray-600 italic">{engSentence}</div>}
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

        {/* Right column: Kanji information (if available) */}
        {kanjiList.length > 0 && kanjiList[0].text && (
          <div className="w-full md:w-1/4 bg-gray-50 p-3 rounded-lg">
            <h3 className="text-lg font-semibold mb-2">Kanji</h3>
            <div className="flex flex-wrap gap-2">
              {Array.from(new Set(kanjiList[0].text.split(''))).filter(char => /[\u4e00-\u9faf]/.test(char)).map((kanji, i) => (
                <a
                  key={`kanji-link-${i}`}
                  href={`/word/${kanji}`}
                  className="block text-center p-2 border border-gray-200 rounded bg-white hover:bg-blue-50"
                >
                  <div className="text-2xl mb-1">{kanji}</div>
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
 * Renders a JMnedict name entry
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
        <span className="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded">
          Japanese Name
        </span>
        {id && (
          <span className="text-gray-500 text-xs">
            ID: {id}
            {shardType !== null && ` (Shard: ${shardType})`}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Name display */}
        <div className="mb-4">
          {/* Kanji writings */}
          {kanjiList.length > 0 && (
            <div className="flex flex-wrap gap-2 mb-2">
              {kanjiList.map((k: string, i: number) => (
                <span key={`kanji-${i}`} className="text-2xl font-bold">{k}</span>
              ))}
            </div>
          )}

          {/* Readings */}
          {readingList.length > 0 && (
            <div className="flex flex-wrap gap-2 text-gray-600">
              {readingList.map((r: string, i: number) => (
                <span key={`reading-${i}`} className="text-xl">{r}</span>
              ))}
            </div>
          )}
        </div>

        {/* Name type */}
        {typeList.length > 0 && (
          <div className="mb-2">
            <span className="text-sm font-medium text-gray-500">Type: </span>
            <span className="text-gray-700">{typeList.join(', ')}</span>
          </div>
        )}

        {/* Meanings */}
        {meaningList.length > 0 && (
          <div>
            <span className="text-sm font-medium text-gray-500">Meaning: </span>
            <span className="text-gray-800">{meaningList.join(', ')}</span>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders a Kanjidic kanji entry
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

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded">
          Kanji Character
        </span>
        {id && (
          <span className="text-gray-500 text-xs">
            ID: {id}
          </span>
        )}
      </div>

      {/* Main content */}
      <div className="flex flex-col md:flex-row gap-4">
        {/* Left column: Character and readings */}
        <div className="flex-1">
          {/* Character display */}
          <div className="flex items-center gap-6 mb-4">
            <div className="text-6xl font-bold">{character}</div>
            <div>
              {grade && (
                <div className="mb-1">
                  <span className="text-sm font-medium text-gray-500">Grade: </span>
                  <span className="text-gray-800">{grade}</span>
                </div>
              )}
              {stroke && (
                <div className="mb-1">
                  <span className="text-sm font-medium text-gray-500">Strokes: </span>
                  <span className="text-gray-800">{stroke}</span>
                </div>
              )}
              {frequency && (
                <div className="mb-1">
                  <span className="text-sm font-medium text-gray-500">Frequency: </span>
                  <span className="text-gray-800">{frequency}</span>
                </div>
              )}
              {jlpt && (
                <div>
                  <span className="text-sm font-medium text-gray-500">JLPT Level: </span>
                  <span className="text-gray-800">N{jlpt}</span>
                </div>
              )}
            </div>
          </div>

          {/* Readings */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            {onYomi.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-gray-500 mb-1">On Readings</h3>
                <div className="text-gray-800">{onYomi.join(', ')}</div>
              </div>
            )}

            {kunYomi.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-gray-500 mb-1">Kun Readings</h3>
                <div className="text-gray-800">{kunYomi.join(', ')}</div>
              </div>
            )}
          </div>

          {/* Meanings */}
          {meanings.length > 0 && (
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-500 mb-1">Meanings</h3>
              <div className="text-gray-800">{meanings.join(', ')}</div>
            </div>
          )}

          {/* IDS data */}
          {ids && (
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Character Composition (IDS)</h3>
              <div className="text-gray-800">{ids}</div>
              <div className="text-xs text-gray-500 mt-1">
                IDS (Ideographic Description Sequence) shows how the character is composed
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

/**
 * Renders a Chinese character entry
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

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded">
          Chinese Character
        </span>
        {id && (
          <span className="text-gray-500 text-xs">
            ID: {id}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Character display */}
        <div className="flex gap-6 mb-4">
          <div>
            <div className="text-sm font-medium text-gray-500">Traditional</div>
            <div className="text-4xl font-bold">{traditional}</div>
          </div>

          {simplified !== traditional && (
            <div>
              <div className="text-sm font-medium text-gray-500">Simplified</div>
              <div className="text-4xl font-bold">{simplified}</div>
            </div>
          )}

          {strokeCount && (
            <div className="self-end">
              <div className="text-sm font-medium text-gray-500">Strokes</div>
              <div className="text-xl">{strokeCount}</div>
            </div>
          )}
        </div>

        {/* Pinyin */}
        {pinyin.length > 0 && (
          <div className="mb-4">
            <div className="text-sm font-medium text-gray-500 mb-1">Pinyin</div>
            <div className="text-xl text-gray-800">{pinyin.join(', ')}</div>
          </div>
        )}

        {/* Definitions */}
        {definitions.length > 0 && (
          <div className="mb-4">
            <div className="text-sm font-medium text-gray-500 mb-1">Definitions</div>
            <ul className="list-disc list-inside text-gray-800">
              {definitions.map((def: string, i: number) => (
                <li key={`def-${i}`}>{def}</li>
              ))}
            </ul>
          </div>
        )}

        {/* IDS data */}
        {ids && (
          <div>
            <div className="text-sm font-medium text-gray-500 mb-1">Character Composition (IDS)</div>
            <div className="text-gray-800">{ids}</div>
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
 * Renders a Chinese word entry
 */
function ChineseWordEntry({ entry }: { entry: any }) {
  // Extract data
  const id = entry.id || entry.ID;
  const traditional = entry.traditional_chinese || entry.traditional || entry.Traditional;
  const simplified = entry.simplified_chinese || entry.simplified || entry.Simplified;
  const definitions = entry.definitions || entry.Definitions || [];
  const pinyin = entry.pinyin || entry.Pinyin || [];
  const hskLevel = entry.hskLevel || entry.HskLevel;

  return (
    <div>
      {/* Header with dictionary type */}
      <div className="flex justify-between items-start mb-2">
        <span className="bg-yellow-100 text-yellow-800 text-xs font-medium px-2.5 py-0.5 rounded">
          Chinese Word
        </span>
        {id && (
          <span className="text-gray-500 text-xs">
            ID: {id}
          </span>
        )}
      </div>

      {/* Main content */}
      <div>
        {/* Word display */}
        <div className="flex gap-6 mb-4">
          <div>
            <div className="text-sm font-medium text-gray-500">Traditional</div>
            <div className="text-3xl font-bold">{traditional}</div>
          </div>

          {simplified !== traditional && (
            <div>
              <div className="text-sm font-medium text-gray-500">Simplified</div>
              <div className="text-3xl font-bold">{simplified}</div>
            </div>
          )}

          {hskLevel && (
            <div className="self-end">
              <div className="text-sm font-medium text-gray-500">HSK Level</div>
              <div className="text-xl">{hskLevel}</div>
            </div>
          )}
        </div>

        {/* Pinyin */}
        {pinyin.length > 0 && (
          <div className="mb-4">
            <div className="text-sm font-medium text-gray-500 mb-1">Pinyin</div>
            <div className="text-xl text-gray-800">{pinyin.join(', ')}</div>
          </div>
        )}

        {/* Definitions */}
        {definitions.length > 0 && (
          <div>
            <div className="text-sm font-medium text-gray-500 mb-1">Definitions</div>
            <ul className="list-disc list-inside text-gray-800">
              {definitions.map((def: string, i: number) => (
                <li key={`def-${i}`}>{def}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Renders an unknown entry type
 */
function UnknownEntry({ entry }: { entry: any }) {
  return (
    <div>
      <div className="flex justify-between items-start mb-2">
        <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded">
          Unknown Entry Type
        </span>
      </div>
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    </div>
  );
}
