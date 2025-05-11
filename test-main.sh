#!/bin/bash
# Script to test the main program with a small limit

set -e  # Exit on error

echo "=== Testing Main Program ==="

# Run the main program with a small limit
echo "Running the main program..."
go run cmd/kiokun/main.go --limit 10 --silent

# Check the output directories
echo "Checking output directories..."
find output_* -type d -maxdepth 0

echo "=== Test completed successfully ==="
