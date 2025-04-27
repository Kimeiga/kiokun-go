import json
import gzip
import random
import os
from pathlib import Path


def load_gzipped_file(path):
    with gzip.open(path, "rt", encoding="utf-8") as f:
        return json.load(f)


def load_unzipped_file(path):
    with open(path, "rt", encoding="utf-8") as f:
        return json.load(f)


def main():
    # Get list of all gzipped files
    gz_dir = Path("dictionary")
    unzipped_dir = Path("dictionary_unzipped")

    if not gz_dir.exists() or not unzipped_dir.exists():
        print("Error: Both dictionary/ and dictionary_unzipped/ must exist")
        return

    # Get all .json.gz files
    gz_files = list(gz_dir.glob("*.json.gz"))

    # Select 100 random files
    sample_size = min(100, len(gz_files))
    test_files = random.sample(gz_files, sample_size)

    print(f"Testing {sample_size} random files...")

    mismatches = []
    for gz_path in test_files:
        # Get corresponding unzipped path
        unzipped_path = unzipped_dir / gz_path.name[:-3]  # remove .gz

        try:
            gz_data = load_gzipped_file(gz_path)
            unzipped_data = load_unzipped_file(unzipped_path)

            # Deep compare the JSON
            if gz_data != unzipped_data:
                mismatches.append(gz_path.name)
                print(f"\nMismatch in {gz_path.name}:")
                print(f"Gzipped: {gz_data}")
                print(f"Unzipped: {unzipped_data}")
        except Exception as e:
            print(f"\nError processing {gz_path.name}: {e}")
            mismatches.append(gz_path.name)

    # Report results
    if mismatches:
        print(f"\nFound {len(mismatches)} mismatches:")
        for m in mismatches:
            print(f"  - {m}")
    else:
        print("\nAll files match! ✨")


if __name__ == "__main__":
    main()
