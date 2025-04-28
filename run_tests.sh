#!/bin/bash

# run_tests.sh - A simple script to run minification tests

set -e  # Exit on error

echo "Running minification tests..."

# Run processor tests with verbose output
go test -v ./processor

echo "Minification tests completed successfully!"
