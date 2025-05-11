# Dictionary Build Workflow Documentation

## Overview

This document explains the GitHub Actions workflows used to build the dictionary files for the Kiokun project. The dictionary is split into four separate repositories to handle the large file sizes:

1. **Non-Han Dictionary**: Contains non-Han characters
2. **Han-1Char Dictionary**: Contains single Han characters
3. **Han-2Char Dictionary**: Contains two-character Han combinations
4. **Han-3Plus Dictionary**: Contains three or more Han character combinations

## Recent Changes

The workflows have been modified to address build failures:

1. **Sequential Builds**: The main workflow now triggers builds sequentially instead of concurrently to avoid resource contention.
2. **Timeouts**: Each build job now has a 60-minute timeout to prevent indefinite runs.
3. **Concurrency Controls**: Each workflow has concurrency settings to ensure only one instance runs at a time.

## Workflow Files

### Main Workflow: `.github/workflows/dictionary-build.yml`

This workflow is the entry point that triggers the individual build workflows in sequence:

1. First, it triggers the Non-Han build and waits for it to complete
2. Then, it triggers the Han-1Char build and waits for it to complete
3. Next, it triggers the Han-2Char build and waits for it to complete
4. Finally, it triggers the Han-3Plus build

### Individual Build Workflows

Each dictionary type has its own workflow file:

- `.github/workflows/build-non-han.yml`
- `.github/workflows/build-han-1char.yml`
- `.github/workflows/build-han-2char.yml`
- `.github/workflows/build-han-3plus.yml`

These workflows:
1. Check out the code
2. Download the Chinese dictionary files if needed
3. Run the main program to generate the dictionary files
4. Push the generated files to their respective repositories

## Troubleshooting

If builds are still failing:

1. **Memory Issues**: Consider optimizing the code to use less memory or process data in smaller chunks
2. **Timeout Issues**: If 60 minutes is not enough, increase the timeout or break the processing into smaller parts
3. **Repository Secrets**: Ensure the `DICTIONARY_DEPLOY_TOKEN` secret is properly set up with write access to all repositories

## Manual Triggering

You can manually trigger the builds:

1. Go to the Actions tab in GitHub
2. Select the "Dictionary Build" workflow
3. Click "Run workflow"

Alternatively, you can trigger individual builds by selecting their specific workflow.
