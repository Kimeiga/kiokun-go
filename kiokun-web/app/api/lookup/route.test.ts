/**
 * Tests for the dictionary lookup API
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { GET } from './route';
import { NextRequest } from 'next/server';
import * as dictionaryUtils from '@/lib/dictionary-utils';
import * as brotliUtils from '@/lib/brotli-utils';

// Mock the dictionary utilities
vi.mock('@/lib/dictionary-utils', () => ({
  getShardType: vi.fn(),
  getIndexUrl: vi.fn(),
  getDictionaryEntryUrl: vi.fn(),
  extractShardType: vi.fn(),
  ShardType: {
    NON_HAN: 0,
    HAN_1CHAR: 1,
    HAN_2CHAR: 2,
    HAN_3PLUS: 3,
  },
  DictionaryType: {
    JMDICT: 'j',
    JMNEDICT: 'n',
    KANJIDIC: 'd',
    CHINESE_CHARS: 'c',
    CHINESE_WORDS: 'w',
  },
}));

// Mock the Brotli utilities
vi.mock('@/lib/brotli-utils', () => ({
  fetchAndDecompressJson: vi.fn(),
}));

describe('Dictionary Lookup API', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('should return 400 if word parameter is missing', async () => {
    // Create a mock request without a word parameter
    const request = new NextRequest('http://localhost:3000/api/lookup');
    
    // Call the API handler
    const response = await GET(request);
    
    // Check the response
    expect(response.status).toBe(400);
    const data = await response.json();
    expect(data.error).toBe('Word parameter is required');
  });

  it('should return empty arrays if index file does not exist', async () => {
    // Create a mock request with a word parameter
    const request = new NextRequest('http://localhost:3000/api/lookup?word=test');
    
    // Mock the getShardType function
    vi.mocked(dictionaryUtils.getShardType).mockReturnValue(dictionaryUtils.ShardType.NON_HAN);
    
    // Mock the getIndexUrl function
    vi.mocked(dictionaryUtils.getIndexUrl).mockReturnValue('https://example.com/index/test.json.br');
    
    // Mock the fetchAndDecompressJson function to throw an error
    vi.mocked(brotliUtils.fetchAndDecompressJson).mockRejectedValue(new Error('Not found'));
    
    // Call the API handler
    const response = await GET(request);
    
    // Check the response
    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data.word).toBe('test');
    expect(data.exactMatches).toEqual([]);
    expect(data.containedMatches).toEqual([]);
  });

  it('should return exact and contained matches', async () => {
    // Create a mock request with a word parameter
    const request = new NextRequest('http://localhost:3000/api/lookup?word=test');
    
    // Mock the getShardType function
    vi.mocked(dictionaryUtils.getShardType).mockReturnValue(dictionaryUtils.ShardType.NON_HAN);
    
    // Mock the getIndexUrl function
    vi.mocked(dictionaryUtils.getIndexUrl).mockReturnValue('https://example.com/index/test.json.br');
    
    // Mock the extractShardType function
    vi.mocked(dictionaryUtils.extractShardType).mockReturnValue(dictionaryUtils.ShardType.NON_HAN);
    
    // Mock the getDictionaryEntryUrl function
    vi.mocked(dictionaryUtils.getDictionaryEntryUrl).mockImplementation(
      (id, dictType, shardType) => `https://example.com/${dictType}/${id}.json.br`
    );
    
    // Mock the fetchAndDecompressJson function for the index
    vi.mocked(brotliUtils.fetchAndDecompressJson).mockImplementation(async (url) => {
      if (url === 'https://example.com/index/test.json.br') {
        return {
          E: {
            j: [1, 2],
            n: [3],
          },
          C: {
            c: [4],
          },
        };
      } else if (url === 'https://example.com/j/01.json.br') {
        return { id: '01', type: 'jmdict', text: 'Test 1' };
      } else if (url === 'https://example.com/j/02.json.br') {
        return { id: '02', type: 'jmdict', text: 'Test 2' };
      } else if (url === 'https://example.com/n/03.json.br') {
        return { id: '03', type: 'jmnedict', text: 'Test 3' };
      } else if (url === 'https://example.com/c/04.json.br') {
        return { id: '04', type: 'chinese_chars', text: 'Test 4' };
      }
      throw new Error(`Unexpected URL: ${url}`);
    });
    
    // Call the API handler
    const response = await GET(request);
    
    // Check the response
    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data.word).toBe('test');
    expect(data.exactMatches).toHaveLength(3);
    expect(data.containedMatches).toHaveLength(1);
    
    // Check the exact matches
    expect(data.exactMatches[0]).toEqual({ id: '01', type: 'jmdict', text: 'Test 1' });
    expect(data.exactMatches[1]).toEqual({ id: '02', type: 'jmdict', text: 'Test 2' });
    expect(data.exactMatches[2]).toEqual({ id: '03', type: 'jmnedict', text: 'Test 3' });
    
    // Check the contained matches
    expect(data.containedMatches[0]).toEqual({ id: '04', type: 'chinese_chars', text: 'Test 4' });
  });
});
