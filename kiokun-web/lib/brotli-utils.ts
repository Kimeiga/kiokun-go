/**
 * Utilities for handling Brotli-compressed data
 */

import * as brotli from "brotli";

/**
 * Decompresses Brotli-compressed data
 * @param compressedData The Brotli-compressed data as a Buffer
 * @returns The decompressed data as a string
 */
export async function decompressBrotli(
  compressedData: Buffer
): Promise<string> {
  try {
    // Decompress the data
    const decompressedBuffer = brotli.decompress(compressedData);

    if (!decompressedBuffer) {
      throw new Error("Brotli decompression returned null or undefined");
    }

    // Convert the buffer to a string
    return new TextDecoder().decode(decompressedBuffer);
  } catch (error) {
    console.error("Error decompressing Brotli data:", error);
    throw new Error("Failed to decompress Brotli data");
  }
}

/**
 * Fetches and decompresses Brotli-compressed JSON from a URL
 * @param url The URL to fetch the compressed JSON from
 * @returns The parsed JSON object
 */
export async function fetchAndDecompressJson<T>(url: string): Promise<T> {
  try {
    console.log(`Fetching from URL: ${url}`);

    // Fetch the compressed data
    const response = await fetch(url, {
      // Ensure we don't get cached responses during testing
      cache: "no-store",
      // Set longer timeout for potentially large files
      signal: AbortSignal.timeout(30000), // 30 seconds timeout
    });

    if (!response.ok) {
      throw new Error(
        `Failed to fetch data: ${response.status} ${response.statusText}`
      );
    }

    // Get the compressed data as a buffer
    const compressedData = await response.arrayBuffer();
    console.log(
      `Received ${compressedData.byteLength} bytes of compressed data`
    );

    if (compressedData.byteLength === 0) {
      throw new Error("Received empty response");
    }

    // Decompress the data
    const decompressedJson = await decompressBrotli(
      Buffer.from(compressedData)
    );
    console.log(`Decompressed to ${decompressedJson.length} characters`);

    // Parse the JSON
    const result = JSON.parse(decompressedJson) as T;
    return result;
  } catch (error) {
    console.error("Error fetching and decompressing JSON:", error);
    throw error;
  }
}
