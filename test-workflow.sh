#!/bin/bash
# Script to test the GitHub Actions workflow locally using act

set -e  # Exit on error

echo "=== Testing GitHub Actions Workflow ==="

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "Error: act is not installed. Please install it first."
    echo "See https://github.com/nektos/act for installation instructions."
    exit 1
fi

# Run the workflow
echo "Running the workflow..."
act -j build -W .github/workflows/build-han-1char.yml

echo "=== Test completed successfully ==="
