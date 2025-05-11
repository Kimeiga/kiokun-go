#!/bin/bash

# Build the binary
echo "Building binary..."
cd cmd/kiokun
go build -o kiokun

# Run with a small dataset for benchmarking
echo "Running benchmark with limit=10000..."
time ./kiokun --limit 10000 --dev

# Clean up
rm kiokun
