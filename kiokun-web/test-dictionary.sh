#!/bin/bash

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
  echo "Installing dependencies..."
  npm install
fi

# Run the dictionary tests
echo "Running dictionary tests..."
npx vitest run lib/dictionary.test.ts app/api/lookup/lookup.test.ts
