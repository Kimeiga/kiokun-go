// Part of speech mappings from 10ten-ja-reader
// Source: https://github.com/birchill/10ten-ja-reader/blob/main/_locales/en/messages.json

export const partOfSpeechMappings: Record<string, string> = {
  // Adjectives
  'adj-f': 'pre-noun adj.',
  'adj-i': 'i adj.',
  'adj-ix': 'ii/yoi adj.',
  'adj-kari': 'kari adj.',
  'adj-ku': 'ku adj.',
  'adj-na': 'na adj.',
  'adj-nari': 'nari adj.',
  'adj-no': 'no-adj.',
  'adj-pn': 'pre-noun adj.',
  'adj-shiku': 'shiku adj.',
  'adj-t': 'taru adj.',

  // Adverbs
  'adv': 'adverb',
  'adv-to': 'adverb to',

  // Auxiliary
  'aux': 'aux.',
  'aux-adj': 'aux. adj.',
  'aux-v': 'aux. verb',

  // Conjunctions and others
  'conj': 'conj.',
  'cop': 'copula',
  'ctr': 'counter',
  'exp': 'exp.',
  'int': 'int.',

  // Nouns
  'n': 'noun',
  'n-adv': 'adv. noun',
  'n-pr': 'proper noun',
  'n-pref': 'n-pref',
  'n-suf': 'n-suf',
  'n-t': 'n-temp',
  'num': 'numeric',

  // Pronouns and particles
  'pn': 'pronoun',
  'pref': 'prefix',
  'prt': 'particle',
  'suf': 'suffix',
  'unc': '?',

  // Verbs - Basic
  'v-unspec': 'verb',
  'vi': 'intrans.',
  'vt': 'trans.',

  // Ichidan verbs
  'v1': 'Ichidan/ru-verb',
  'v1-s': 'Ichidan/ru-verb (kureru)',

  // Nidan verbs (archaic)
  'v2a-s': '-u Nidan verb',
  'v2b-k': '-bu upper Nidan verb',
  'v2b-s': '-bu lower Nidan verb',
  'v2d-k': '-dzu upper Nidan verb',
  'v2d-s': '-dzu lower Nidan verb',
  'v2g-k': '-gu upper Nidan verb',
  'v2g-s': '-gu lower Nidan verb',
  'v2h-k': '-hu/-fu upper Nidan verb',
  'v2h-s': '-hu/-fu lower Nidan verb',
  'v2k-k': '-ku upper Nidan verb',
  'v2k-s': '-ku lower Nidan verb',
  'v2m-k': '-mu upper Nidan verb',
  'v2m-s': '-mu lower Nidan verb',
  'v2n-s': '-nu Nidan verb',
  'v2r-k': '-ru upper Nidan verb',
  'v2r-s': '-ru lower Nidan verb',
  'v2s-s': '-su Nidan verb',
  'v2t-k': '-tsu upper Nidan verb',
  'v2t-s': '-tsu upper Nidan verb',
  'v2w-s': '-u Nidan verb + we',
  'v2y-k': '-yu upper Nidan verb',
  'v2y-s': '-yu lower Nidan verb',
  'v2z-s': '-zu Nidan verb',

  // Yodan verbs (archaic)
  'v4b': '-bu Yodan verb',
  'v4g': '-gu Yodan verb',
  'v4h': '-hu/-fu Yodan verb',
  'v4k': '-ku Yodan verb',
  'v4m': '-mu Yodan verb',
  'v4n': '-nu Yodan verb',
  'v4r': '-ru Yodan verb',
  'v4s': '-su Yodan verb',
  'v4t': '-tsu Yodan verb',

  // Godan verbs
  'v5aru': '-aru godan verb',
  'v5b': '-bu Godan/u-verb',
  'v5g': '-gu Godan/u-verb',
  'v5k': '-ku Godan/u-verb',
  'v5k-s': 'iku/yuku Godan/u-verb',
  'v5m': '-mu Godan/u-verb',
  'v5n': '-nu Godan/u-verb',
  'v5r': '-ru Godan/u-verb',
  'v5r-i': '-ru Godan/u-verb (irr.)',
  'v5s': '-su Godan/u-verb',
  'v5t': '-tsu Godan/u-verb',
  'v5u': '-u Godan/u-verb',
  'v5u-s': '-u Godan/u-verb (special)',
  'v5uru': '-uru Godan/u-verb',

  // Special verbs
  'vk': 'kuru verb',
  'vn': '-nu irr. verb',
  'vr': '-ru (-ri) irr. verb',
  'vs': '+suru verb',
  'vs-c': '-su(ru) verb',
  'vs-i': '-suru verb',
  'vs-s': '-suru verb (special)',
  'vz': '-zuru Ichidan/ru-verb',
};

/**
 * Maps a part of speech abbreviation to its full form
 * @param pos - The part of speech abbreviation (e.g., "n", "vt", "vi")
 * @returns The full form (e.g., "noun", "trans.", "intrans.") or the original if no mapping exists
 */
export function mapPartOfSpeech(pos: string): string {
  return partOfSpeechMappings[pos] || pos;
}

/**
 * Maps an array of part of speech abbreviations to their full forms
 * @param posArray - Array of part of speech abbreviations
 * @returns Array of full forms
 */
export function mapPartsOfSpeech(posArray: string[]): string[] {
  return posArray.map(mapPartOfSpeech);
}

// Miscellaneous tag mappings from 10ten-ja-reader
export const miscTagMappings: Record<string, string> = {
  'abbr': 'abbreviation',
  'arch': 'archaic',
  'chn': 'children\'s language',
  'col': 'colloquial',
  'derog': 'derogatory',
  'fam': 'familiar language',
  'fem': 'female term or language',
  'gikun': 'gikun (meaning as reading) or jukujikun (special kanji reading)',
  'hon': 'honorific or respectful (sonkeigo) language',
  'hum': 'humble (kenjougo) language',
  'id': 'idiomatic expression',
  'joc': 'jocular, humorous term',
  'male': 'male term or language',
  'male-sl': 'male slang',
  'obs': 'obsolete term',
  'obsc': 'obscure term',
  'on-mim': 'onomatopoeic or mimetic word',
  'poet': 'poetical term',
  'pol': 'polite (teineigo) language',
  'rare': 'rare',
  'sens': 'sensitive',
  'sl': 'slang',
  'uk': 'word usually written using kana alone',
  'vulg': 'vulgar expression or word',
  'yoji': 'yojijukugo',
};

/**
 * Maps a miscellaneous tag abbreviation to its full form
 * @param misc - The misc tag abbreviation (e.g., "abbr", "arch")
 * @returns The full form (e.g., "abbreviation", "archaic") or the original if no mapping exists
 */
export function mapMiscTag(misc: string): string {
  return miscTagMappings[misc] || misc;
}
