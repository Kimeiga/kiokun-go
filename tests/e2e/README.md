# End-to-End Tests for Kiokun

This directory contains end-to-end tests for the Kiokun dictionary system.

## Character Contained Matches Test

The `character_contained_matches_test.go` verifies that the contained matches feature works correctly for any given character.

### Usage

```bash
# Run all E2E tests
go test ./tests/e2e/ -v

# Run only character contained matches tests
go test ./tests/e2e/ -v -run TestCharacterContainedMatches

# Run test for specific character
go test ./tests/e2e/ -v -run TestCharacterContainedMatches/Character_日
```

### What it does

1. **Creates test data** with the specified character and words containing it
2. **Processes entries** through the sharded index processor
3. **Verifies sharding** - checks that data is distributed correctly across shards
4. **Tests API logic** - simulates the frontend API calls to verify contained matches work
5. **Cleans up** - removes test output files

### Test Structure

- `character_contained_matches_test.go` - Complete E2E test using Go testing framework
- Uses root `go.mod` for dependency management
- Standard Go test conventions with `testing.T`

### Expected Results

For any character (e.g., 日, 人, 水):
- ✅ Character itself should be in the appropriate shard
- ✅ Words containing the character should be in their respective shards
- ✅ API should find ALL contained matches across all shards
- ✅ Total contained matches should meet minimum expected count

### Benefits of This Approach

- **Standard Go testing**: Uses `go test` command and `testing.T`
- **Single dependency management**: Uses root `go.mod` instead of separate module
- **Better IDE integration**: Standard test structure works with all Go IDEs
- **CI/CD friendly**: Easy to integrate into automated testing pipelines
- **Maintainable**: Follows Go community best practices
- ✅ Words containing the character should be in their respective shards
- ✅ API should find ALL contained matches across all shards
- ✅ Total contained matches should match expected count
