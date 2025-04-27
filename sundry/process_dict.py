import json
import gzip
import os
import shutil
import time
from pathlib import Path


def main():
    # Setup
    output_dir = Path("dictionary")
    if output_dir.exists():
        shutil.rmtree(output_dir)
    output_dir.mkdir()

    # Read input file
    print("Reading input file...")
    with open("jmdict-eng-3.5.0.json") as f:
        dict_data = json.load(f)

    words = dict_data["words"]
    total_words = len(words)
    print(f"Processing {total_words} words...")

    # Group words
    word_groups = {}
    start = time.time()
    processed = 0

    for word in words:
        # Get filename from kanji or kana
        if word["kanji"]:
            filename = word["kanji"][0]["text"]
        elif word["kana"]:
            filename = word["kana"][0]["text"]
        else:
            filename = word["id"]

        # Add to group
        if filename in word_groups:
            word_groups[filename]["w_j"].append(word)
        else:
            word_groups[filename] = {"w_j": [word]}

        processed += 1
        if processed % 1000 == 0:
            elapsed = time.time() - start
            rate = processed / elapsed
            print(
                f"\rProgress: {processed/total_words*100:.1f}% ({processed}/{total_words}) - {rate:.1f} words/sec",
                end="",
            )

    # Write files
    for filename, group in word_groups.items():
        output_path = output_dir / f"{filename}.json.gz"
        with gzip.open(output_path, "wt", encoding="utf-8") as f:
            json.dump(group, f, ensure_ascii=False)

    elapsed = time.time() - start
    rate = total_words / elapsed
    print(
        f"\nCompleted 100% ({total_words}/{total_words}) in {elapsed:.1f}s ({rate:.1f} words/sec)"
    )


if __name__ == "__main__":
    main()
