# JMdict Test Data

This directory contains test data and utilities for testing the JMdict package.

## Files

- `test_jmdict.json`: A small sample of JMdict entries for basic testing. Contains examples of different field types.
- `field_finder.go`: A utility to find entries with specific field types in the main JMdict file.
- `field_finder_main.go`: A tool to run the field finder and generate examples.
- `generate_test_data.go`: A tool to create test data with entries containing various field types.

## Test Approach

The tests in this package verify:

1. **Unmarshaling and Marshaling**: Test that JmdictTypes can be correctly unmarshaled from JSON and marshaled back.
2. **Importer Functionality**: Test that the Importer correctly loads entries from JSON files.
3. **Field Type Validation**: Test that various field types in the JMdict entries are correctly parsed.
4. **Dictionary Registration**: Test that the init.go file correctly registers the JMdict dictionary.
5. **Entry Interface**: Test that Word implements the common.Entry interface correctly.

## Running the Tests

From the project root, run:

```bash
go test -v ./dictionaries/jmdict
```

## Finding Examples of Field Types

The `field_finder.go` utility provides a way to scan the main JMdict file to find entries that use specific field types, without loading the entire file into memory. This is useful for creating test data that covers all field types.

To use it:

```bash
cd dictionaries/jmdict/testdata
go run field_finder_main.go
```

## Test Data Generation

The `generate_test_data.go` script is for generating test data with entries containing specific field types. It's not fully implemented yet, but provides a framework for creating more comprehensive test data.

## Adding New Tests

When adding new tests:

1. Use the existing test_jmdict.json file for basic tests
2. If you need examples of specific field types, consider using the field finder
3. Keep tests focused on specific functionality
4. Avoid loading the entire JMdict file into memory
