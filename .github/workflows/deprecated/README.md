# Deprecated Workflows

These workflows have been moved here because they were causing conflicts with the new matrix-based build system.

## Issue

Multiple workflows were triggering simultaneously on `push` to `main`:

1. **`build-dictionaries-matrix.yml`** (new matrix approach) ✅
2. **`dictionary-build.yml`** (old coordinator) ❌ 
3. **`build-han-only.yml`** (old individual build) ❌

This caused:
- Resource conflicts during builds
- Concurrent pushes to output repositories  
- Race conditions leading to corrupted/incomplete data
- Missing Chinese word entries (like ID 157102 for 日)

## Solution

Only the matrix workflow (`build-dictionaries-matrix.yml`) now runs on push to main.
The deprecated workflows are preserved here for reference but are inactive.

## Files Moved

- `dictionary-build.yml` - Old coordinator workflow
- `build-han-only.yml` - Old han-only individual build  
- `build-han-1char.yml` - Individual han-1char build (repository_dispatch)
- `build-han-2char.yml` - Individual han-2char build (repository_dispatch)
- `build-han-3plus.yml` - Individual han-3plus build (repository_dispatch)
- `build-non-han.yml` - Individual non-han build (repository_dispatch)

## Date Moved

2025-06-27 - Fixed corruption issue in Chinese word dictionary builds
