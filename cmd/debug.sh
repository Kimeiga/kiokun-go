#!/bin/bash

echo "=== Building and running debug tools ==="

# Clean up previous test output
rm -rf output_test
mkdir -p output_test

# Build and run the debug processor
echo "Building debug processor..."
go build -o debug_processor cmd/debug_processor/main.go
if [ $? -ne 0 ]; then
    echo "Failed to build debug processor"
    exit 1
fi

# Run the debug processor
echo "Running debug processor..."
./debug_processor
if [ $? -ne 0 ]; then
    echo "Failed to run debug processor"
    exit 1
fi

# Build and run the debug analyzer
echo "Building debug analyzer..."
go build -o debug_analyzer cmd/debug_analyzer/main.go
if [ $? -ne 0 ]; then
    echo "Failed to build debug analyzer"
    exit 1
fi

# Run the debug analyzer
echo "Running debug analyzer..."
./debug_analyzer
if [ $? -ne 0 ]; then
    echo "Failed to run debug analyzer"
    exit 1
fi

# Clean up the binaries
rm -f debug_processor debug_analyzer

echo "Debug process completed" 