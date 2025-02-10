import { error } from "@sveltejs/kit";
import pako from "pako";

export async function load({ params, fetch }) {
  try {
    const response = await fetch(`/dictionary/${params.word}.json.gz`);

    if (!response.ok) {
      throw error(404, "Word not found");
    }

    const buffer = await response.arrayBuffer();
    const decompressed = pako.inflate(new Uint8Array(buffer), { to: "string" });
    const data = JSON.parse(decompressed);

    return {
      entries: data,
    };
  } catch (e) {
    throw error(404, "Word not found");
  }
}
