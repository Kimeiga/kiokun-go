/**
 * Service for fetching dictionary entries from multiple repositories
 */

/**
 * Check if a string contains only Han characters
 */
export function isHanOnly(word: string): boolean {
  // Unicode range for Han characters (CJK Unified Ideographs)
  for (const char of word) {
    const code = char.codePointAt(0)!;
    // Check if character is not in the Han Unicode block
    if (!(code >= 0x4e00 && code <= 0x9fff)) {
      return false;
    }
  }
  return true;
}

/**
 * Base configuration for dictionary repositories
 */
const config = {
  // Han character repositories by length
  han1CharRepo: "your-username/japanese-dict-han-1char",
  han2CharRepo: "your-username/japanese-dict-han-2char",
  han3PlusRepo: "your-username/japanese-dict-han-3plus",
  nonHanRepo: "your-username/japanese-dict-non-han",
  cdnBase: "https://cdn.jsdelivr.net/gh/",
  branch: "main",
  fileExtension: ".json.br",
};

/**
 * Get the URL for a dictionary word file
 */
export function getDictionaryFileUrl(word: string): string {
  let repoName: string;

  // First check if it contains only Han characters
  if (isHanOnly(word)) {
    // Split based on character length
    const charCount = word.length;
    if (charCount === 1) {
      repoName = config.han1CharRepo;
    } else if (charCount === 2) {
      repoName = config.han2CharRepo;
    } else {
      // 3 or more characters
      repoName = config.han3PlusRepo;
    }
  } else {
    // Contains at least one non-Han character
    repoName = config.nonHanRepo;
  }

  return `${config.cdnBase}${repoName}@${config.branch}/${word}${config.fileExtension}`;
}

/**
 * Fetch a dictionary entry
 */
export async function fetchDictionaryEntry(word: string): Promise<any> {
  const url = getDictionaryFileUrl(word);

  try {
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`Failed to fetch word: ${response.status}`);
    }

    // For client-side, you may need to handle brotli decompression
    // This depends on your frontend framework and environment
    // Some environments handle this automatically

    const data = await response.json();
    return data;
  } catch (error) {
    console.error(`Error fetching dictionary entry for "${word}":`, error);
    throw error;
  }
}
