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

    // Check if the URL is for a Chinese word dictionary (w/ folder)
    const isChineseWordDict = url.includes("/w/");

    // Fetch the compressed data with custom headers
    const response = await fetch(url, {
      // Ensure we don't get cached responses during testing
      cache: "no-store",
      // Set longer timeout for potentially large files
      signal: AbortSignal.timeout(30000), // 30 seconds timeout
      headers: {
        // Add a user agent to avoid being blocked
        "User-Agent": "Kiokun-Dictionary/1.0",
        // Add a referrer to indicate the source
        Referer: "https://kiokun.com/",
      },
    });

    // Log response status and headers for debugging
    console.log(`Response status: ${response.status} ${response.statusText}`);
    console.log(
      `Response headers:`,
      Object.fromEntries(response.headers.entries())
    );

    // If it's a 403 error and it's a Chinese word dictionary, try GitHub raw content
    if (!response.ok && response.status === 403 && isChineseWordDict) {
      console.log(
        `jsDelivr returned 403 for Chinese word dictionary. Package likely exceeds 50MB limit.`
      );

      // For Chinese word dictionaries that exceed jsDelivr's size limit,
      // we'll return an empty object instead of trying to fetch from GitHub
      // since the repository structure might not match our expectations
      console.log(
        `Returning empty object for oversized Chinese word dictionary entry`
      );

      // Return an empty object that matches the expected structure
      return {} as T;
    } else if (!response.ok) {
      if (response.status === 403) {
        console.error(`403 Forbidden error for URL: ${url}`);
        console.error(
          `This might be due to rate limiting by jsDelivr or missing files.`
        );
        console.error(
          `Headers:`,
          Object.fromEntries(response.headers.entries())
        );
      }

      throw new Error(
        `Failed to fetch data: ${response.status} ${response.statusText}`
      );
    }

    // Get the compressed data as a buffer
    const compressedData = await response.arrayBuffer();
    console.log(
      `Received ${compressedData.byteLength} bytes of compressed data from ${url}`
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
    console.error(`Error fetching and decompressing JSON from ${url}:`, error);

    // Check if it's a network error
    if (error instanceof TypeError && error.message.includes("fetch")) {
      console.error(
        `Network error when fetching ${url}. This might be a connectivity issue.`
      );
    }

    throw error;
  }
}
