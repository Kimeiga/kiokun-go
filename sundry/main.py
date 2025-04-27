import json
import gzip
import os
import shutil
import time
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor
import multiprocessing
from queue import Queue
from threading import Thread, Lock
import argparse


class WordGroup:
    def __init__(self):
        self.word_japanese = []

    def to_dict(self):
        return {"w_j": self.word_japanese}


def show_progress(total_words, processed_count, done_event):
    start = time.time()
    while not done_event.is_set():
        processed = processed_count.value
        elapsed = time.time() - start
        rate = processed / elapsed if elapsed > 0 else 0
        print(
            f"\rProgress: {processed/total_words*100:.1f}% ({processed}/{total_words}) - {rate:.1f} words/sec",
            end="",
            flush=True,
        )
        time.sleep(0.1)

    elapsed = time.time() - start
    rate = total_words / elapsed
    print(
        f"\rCompleted 100% ({total_words}/{total_words}) in {elapsed:.1f}s ({rate:.1f} words/sec)"
    )


def worker(word_queue, word_groups, output_dir, processed_count, groups_lock):
    while True:
        try:
            word = word_queue.get_nowait()
        except:
            break

        # Get filename from kanji or kana
        if word["kanji"]:
            filename = word["kanji"][0]["text"]
        elif word["kana"]:
            filename = word["kana"][0]["text"]
        else:
            filename = word["id"]

        # Add to group with lock
        with groups_lock:
            if filename in word_groups:
                word_groups[filename].word_japanese.append(word)
            else:
                group = WordGroup()
                group.word_japanese.append(word)
                word_groups[filename] = group

        processed_count.value += 1
        word_queue.task_done()


def write_group(args):
    filename, group, output_dir = args
    output_path = output_dir / f"{filename}.json.gz"
    with gzip.open(output_path, "wt", encoding="utf-8") as f:
        json.dump(group.to_dict(), f, ensure_ascii=False)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--input",
        default="jmdict-examples-eng-3.5.0.json",
        help="Input JMDict JSON file",
    )
    parser.add_argument(
        "--workers",
        type=int,
        default=multiprocessing.cpu_count(),
        help="Number of workers",
    )
    parser.add_argument(
        "--unzipped", action="store_true", help="Output uncompressed JSON files"
    )
    parser.add_argument("--silent", action="store_true", help="Disable progress output")
    args = parser.parse_args()

    def logf(format_str, *format_args):
        if not args.silent:
            print(format_str % format_args, end="", flush=True)

    # Determine output directory
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

    # First pass: group words (single-threaded)
    logf("Grouping words...\n")
    word_groups = {}
    processed = 0
    start = time.time()

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
            word_groups[filename].word_japanese.append(word)
        else:
            group = WordGroup()
            group.word_japanese.append(word)
            word_groups[filename] = group

        processed += 1
        if not args.silent and processed % 1000 == 0:
            elapsed = time.time() - start
            rate = processed / elapsed if elapsed > 0 else 0
            logf(
                "\rGrouping: %.1f%% (%d/%d) - %.1f words/sec",
                processed * 100 / total_words,
                processed,
                total_words,
                rate,
            )

    logf("\nGrouping completed in %.2fs\n", time.time() - start)

    # Second pass: write files in parallel
    logf("Writing %d files...\n", len(word_groups))
    processed_count = multiprocessing.Value("i", 0)
    done_event = multiprocessing.Event()

    # Start progress monitoring thread
    progress_thread = Thread(
        target=show_progress, args=(len(word_groups), processed_count, done_event)
    )
    if not args.silent:
        progress_thread.start()

    # Write files using thread pool
    write_args = [(f, g, output_dir, args.unzipped) for f, g in word_groups.items()]
    with ThreadPoolExecutor(max_workers=args.workers) as executor:
        futures = []
        for arg in write_args:
            future = executor.submit(write_group_with_counter, arg, processed_count)
            futures.append(future)

        # Wait for all files to be written
        for future in futures:
            future.result()

    # Signal progress thread to finish
    done_event.set()
    if not args.silent:
        progress_thread.join()


def write_group_with_counter(args, processed_count):
    filename, group, output_dir, unzipped = args

    # Ensure consistent handling of empty arrays
    for word in group.word_japanese:
        for kana in word["kana"]:
            if "appliesToKanji" in kana and not kana["appliesToKanji"]:
                kana["appliesToKanji"] = []

    if unzipped:
        output_path = output_dir / f"{filename}.json"
        with open(output_path, "wt", encoding="utf-8") as f:
            json.dump(group.to_dict(), f, ensure_ascii=False)
    else:
        output_path = output_dir / f"{filename}.json.gz"
        with gzip.open(output_path, "wt", encoding="utf-8") as f:
            json.dump(group.to_dict(), f, ensure_ascii=False)

    with processed_count.get_lock():
        processed_count.value += 1


if __name__ == "__main__":
    main()
