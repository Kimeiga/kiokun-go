import json
import gzip
import os
import shutil
import time
from pathlib import Path
import argparse


class WordGroup:
    def __init__(self):
        self.word_japanese = []

    def to_dict(self):
        return {"w_j": self.word_japanese}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--input", default="jmdict-eng-3.5.0.json", help="Input JMDict JSON file"
    )
    parser.add_argument(
        "--unzipped", action="store_true", help="Output uncompressed JSON files"
    )
    parser.add_argument("--silent", action="store_true", help="Disable progress output")
    args = parser.parse_args()

    def logf(format_str, *format_args):
        if not args.silent:
            print(format_str % format_args, end="", flush=True)

    logf("Starting initialization...\n")

    # Determine output directory based on unzipped flag
    output_dir = Path("dictionary_unzipped" if args.unzipped else "dictionary")
    if output_dir.exists():
        shutil.rmtree(output_dir)
    output_dir.mkdir()

    # Read input file
    logf("Reading input file...\n")
    with open(args.input) as f:
        dict_data = json.load(f)

    words = dict_data["words"]
    total_words = len(words)
    logf(f"Processing {total_words} words...\n")

    # Create a map to group words by filename
    word_groups = {}
    start = time.time()
    processed = 0

    # First pass: group all words
    for word in words:
        # Get filename from kanji or kana
        if word["kanji"]:
            filename = word["kanji"][0]["text"]
        elif word["kana"]:
            filename = word["kana"][0]["text"]
        else:
            filename = word["id"]

        # Add word to appropriate group
        if filename in word_groups:
            word_groups[filename].word_japanese.append(word)
        else:
            group = WordGroup()
            group.word_japanese.append(word)
            word_groups[filename] = group

        processed += 1
        if processed % 1000 == 0:
            elapsed = time.time() - start
            rate = processed / elapsed if elapsed > 0 else 0
            logf(
                f"\rProgress: {processed/total_words*100:.1f}% ({processed}/{total_words}) - {rate:.1f} words/sec",
            )

    # Second pass: write all groups to files
    for filename, group in word_groups.items():
        output_path = output_dir / f"{filename}.json.gz"
        with gzip.open(output_path, "wt", encoding="utf-8") as f:
            json.dump(group.to_dict(), f, ensure_ascii=False)

    elapsed = time.time() - start
    rate = total_words / elapsed
    logf(
        f"\nCompleted 100% ({total_words}/{total_words}) in {elapsed:.1f}s ({rate:.1f} words/sec)\n"
    )


if __name__ == "__main__":
    main()
