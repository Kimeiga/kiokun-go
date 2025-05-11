#!/bin/bash
# Script to install Git hooks

# Create the hooks directory if it doesn't exist
mkdir -p .git/hooks

# Copy the pre-commit hook
cp .githooks/pre-commit .git/hooks/

# Make it executable
chmod +x .git/hooks/pre-commit

echo "Git hooks installed successfully!"
