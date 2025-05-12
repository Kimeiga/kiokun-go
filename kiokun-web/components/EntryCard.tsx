/**
 * Component for displaying a dictionary entry
 */

'use client';

import { extractShardType } from '@/lib/dictionary-utils';

interface EntryCardProps {
  entry: any;
}

/**
 * Renders a JMdict word entry
 */
function renderJMdictEntry(entry: any) {
  return (
    <>
      <div className="mb-2">
        {entry.Kanji && entry.Kanji.length > 0 && (
          <div className="text-xl font-bold">
            {entry.Kanji.map((k: any) => k.Text).join(', ')}
          </div>
        )}
        {entry.Kana && entry.Kana.length > 0 && (
          <div className="text-lg">
            {entry.Kana.map((k: any) => k.Text).join(', ')}
          </div>
        )}
      </div>

      {entry.Sense && entry.Sense.length > 0 && (
        <div className="mt-2">
          <ol className="list-decimal list-inside">
            {entry.Sense.map((sense: any, index: number) => (
              <li key={index} className="mb-1">
                {sense.Gloss && sense.Gloss.length > 0 && (
                  <span>{sense.Gloss.map((g: any) => g.Text).join('; ')}</span>
                )}
                {sense.POS && sense.POS.length > 0 && (
                  <span className="text-gray-500 text-sm ml-2">
                    ({sense.POS.join(', ')})
                  </span>
                )}
              </li>
            ))}
          </ol>
        </div>
      )}
    </>
  );
}

/**
 * Renders a JMnedict name entry
 */
function renderJMnedictEntry(entry: any) {
  return (
    <>
      <div className="mb-2">
        {entry.Kanji && entry.Kanji.length > 0 && (
          <div className="text-xl font-bold">
            {entry.Kanji.join(', ')}
          </div>
        )}
        {entry.Reading && entry.Reading.length > 0 && (
          <div className="text-lg">
            {entry.Reading.join(', ')}
          </div>
        )}
      </div>

      {entry.Translation && entry.Translation.length > 0 && (
        <div className="mt-2">
          <ol className="list-decimal list-inside">
            {entry.Translation.map((trans: any, index: number) => (
              <li key={index} className="mb-1">
                {trans.Translation}
                {trans.Type && trans.Type.length > 0 && (
                  <span className="text-gray-500 text-sm ml-2">
                    ({trans.Type.join(', ')})
                  </span>
                )}
              </li>
            ))}
          </ol>
        </div>
      )}
    </>
  );
}

/**
 * Renders a Kanjidic kanji entry
 */
function renderKanjidicEntry(entry: any) {
  return (
    <>
      <div className="text-4xl font-bold mb-2">{entry.Character}</div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          {entry.Meaning && entry.Meaning.length > 0 && (
            <div className="mb-2">
              <h4 className="font-semibold">Meanings:</h4>
              <ul className="list-disc list-inside">
                {entry.Meaning.map((meaning: string, index: number) => (
                  <li key={index}>{meaning}</li>
                ))}
              </ul>
            </div>
          )}

          {entry.Reading && (
            <div className="mb-2">
              <h4 className="font-semibold">Readings:</h4>
              {entry.Reading.OnYomi && entry.Reading.OnYomi.length > 0 && (
                <div>
                  <span className="font-medium">On:</span> {entry.Reading.OnYomi.join(', ')}
                </div>
              )}
              {entry.Reading.KunYomi && entry.Reading.KunYomi.length > 0 && (
                <div>
                  <span className="font-medium">Kun:</span> {entry.Reading.KunYomi.join(', ')}
                </div>
              )}
            </div>
          )}
        </div>

        <div>
          {entry.Misc && (
            <div className="mb-2">
              <h4 className="font-semibold">Details:</h4>
              {entry.Misc.Grade && (
                <div><span className="font-medium">Grade:</span> {entry.Misc.Grade}</div>
              )}
              {entry.Misc.StrokeCount && (
                <div><span className="font-medium">Strokes:</span> {entry.Misc.StrokeCount}</div>
              )}
              {entry.Misc.Frequency && (
                <div><span className="font-medium">Frequency:</span> {entry.Misc.Frequency}</div>
              )}
              {entry.Misc.JLPT && (
                <div><span className="font-medium">JLPT Level:</span> {entry.Misc.JLPT}</div>
              )}
            </div>
          )}
        </div>
      </div>
    </>
  );
}

/**
 * Renders a Chinese character entry
 */
function renderChineseCharEntry(entry: any) {
  return (
    <>
      <div className="mb-2">
        <div className="text-3xl font-bold">
          {entry.Traditional}
          {entry.Simplified !== entry.Traditional && (
            <span className="ml-2">({entry.Simplified})</span>
          )}
        </div>
        {entry.Pinyin && (
          <div className="text-lg">{entry.Pinyin}</div>
        )}
      </div>

      {entry.Definitions && entry.Definitions.length > 0 && (
        <div className="mt-2">
          <h4 className="font-semibold">Definitions:</h4>
          <ul className="list-disc list-inside">
            {entry.Definitions.map((def: string, index: number) => (
              <li key={index}>{def}</li>
            ))}
          </ul>
        </div>
      )}
    </>
  );
}

/**
 * Renders a Chinese word entry
 */
function renderChineseWordEntry(entry: any) {
  return (
    <>
      <div className="mb-2">
        <div className="text-2xl font-bold">
          {entry.Traditional}
          {entry.Simplified !== entry.Traditional && (
            <span className="ml-2">({entry.Simplified})</span>
          )}
        </div>
        {entry.Pinyin && (
          <div className="text-lg">{entry.Pinyin}</div>
        )}
      </div>

      {entry.Definitions && entry.Definitions.length > 0 && (
        <div className="mt-2">
          <h4 className="font-semibold">Definitions:</h4>
          <ul className="list-disc list-inside">
            {entry.Definitions.map((def: string, index: number) => (
              <li key={index}>{def}</li>
            ))}
          </ul>
        </div>
      )}
    </>
  );
}

export default function EntryCard({ entry }: EntryCardProps) {
  // Determine the entry type and render accordingly
  let content;
  let entryType;

  // Check for lowercase field names first (new format)
  if (entry.kanji && entry.kana) {
    // JMdict word entry
    content = (
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    );
    entryType = 'JMdict Word';
  } else if (entry.k && entry.r) {
    // JMnedict name entry
    content = (
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    );
    entryType = 'JMnedict Name';
  } else if (entry.c && entry.on) {
    // Kanjidic kanji entry
    content = (
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    );
    entryType = 'Kanjidic Character';
  } else if (entry.traditional && entry.simplified) {
    // Chinese character entry
    content = (
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    );
    entryType = 'Chinese Character';
  } else if (entry.traditional_chinese && entry.simplified_chinese) {
    // Chinese word entry
    content = (
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    );
    entryType = 'Chinese Word';
  }
  // Check for uppercase field names (old format)
  else if (entry.Kanji && entry.Kana && entry.Sense) {
    // JMdict word entry
    content = renderJMdictEntry(entry);
    entryType = 'JMdict Word';
  } else if (entry.Kanji && entry.Reading && entry.Translation) {
    // JMnedict name entry
    content = renderJMnedictEntry(entry);
    entryType = 'JMnedict Name';
  } else if (entry.Character && entry.Reading && entry.Misc) {
    // Kanjidic kanji entry
    content = renderKanjidicEntry(entry);
    entryType = 'Kanjidic Character';
  } else if (entry.Traditional && entry.Pinyin && !entry.Components) {
    // Chinese word entry
    content = renderChineseWordEntry(entry);
    entryType = 'Chinese Word';
  } else if (entry.Traditional && entry.Pinyin) {
    // Chinese character entry
    content = renderChineseCharEntry(entry);
    entryType = 'Chinese Character';
  } else {
    // Unknown entry type
    content = (
      <pre className="bg-gray-100 p-2 rounded overflow-auto text-sm text-black">
        {JSON.stringify(entry, null, 2)}
      </pre>
    );
    entryType = 'Unknown';
  }

  // Extract shard type from ID if available
  let shardType = null;
  if (entry.ID) {
    shardType = extractShardType(entry.ID);
  }

  return (
    <div className="border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start mb-2">
        <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">
          {entryType}
        </span>

        {entry.ID && (
          <span className="text-gray-500 text-xs">
            ID: {entry.ID}
            {shardType !== null && ` (Shard: ${shardType})`}
          </span>
        )}
      </div>

      {content}
    </div>
  );
}
