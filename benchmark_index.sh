#!/bin/bash

# Build the binary
echo "Building binary..."
cd cmd/kiokun
go build -o kiokun

# Run with index mode
echo "Running benchmark with index mode (limit=10000)..."
time ./kiokun --limit 10000 --dev --index=true

# Run with combined mode
echo "Running benchmark with combined mode (limit=10000)..."
time ./kiokun --limit 10000 --dev --index=false

# Clean up
rm kiokun
