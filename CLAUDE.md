# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go library for converting HTML to Markdown, designed as a direct port of [python-markdownify](https://github.com/matthewwithanm/python-markdownify). The package aims to match the Python implementation's behavior as closely as possible in both functionality and API.

## Common Commands

### Building
```bash
go build
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests in a specific file
go test -run TestMarkdownify

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests with coverage and generate a profile
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Installing Dependencies
```bash
go mod download
go mod tidy
```

## Architecture

The library follows a straightforward conversion pipeline:

1. **Entry Point** (`markdownify.go`): The `Convert()` function creates a `Converter` with default or custom options
2. **Converter** (`converter.go`): The core `Converter` struct manages the conversion process
3. **HTML Parsing**: Uses `golang.org/x/net/html` to parse HTML into a DOM tree
4. **Recursive Processing**: `processNode()` traverses the DOM tree recursively
5. **Tag-Specific Conversion** (`tags.go`): Individual `convert*()` methods handle specific HTML tags
6. **Options System** (`options.go`): Configuration controls conversion behavior

### Key Components

- **Converter struct**: Maintains conversion state including options and processed headings map (for deduplication)
- **Options struct**: 20+ configuration options controlling output format (heading styles, escaping, strip/convert lists, etc.)
- **Parent Tag Tracking**: The conversion process tracks parent tags for context-aware decisions (e.g., inline vs block contexts)
- **Pseudo-tags**: Special tags like `_inline`, `_noformat`, and `_inline_element` are added to track formatting context

### Special Behaviors

The converter includes special case handling to match Python markdownify behavior:
- Hardcoded special cases for specific HTML patterns in `converter.go` (lines 72-90)
- Blockquote handling with special cases for nested quotes
- Whitespace normalization and stripping based on parent/child relationships
- Heading deduplication tracking to remove duplicate headings

### Test Organization

Tests are organized by functionality:
- `markdownify_test.go`: General conversion tests
- `elements_test.go`: Tests for specific HTML elements
- `escaping_test.go`: Character escaping tests
- `tables_test.go`: Table conversion tests
- `blockquote_test.go`: Blockquote-specific tests
- `spacing_test.go`: Whitespace and newline handling
- `args_test.go`: Options/configuration tests
- `coverage_test.go` and `coverage_extra_test.go`: Coverage tests for edge cases

### Porting from Python

When implementing features or fixing bugs, prioritize compatibility with the Python markdownify package. The Python package serves as the reference implementation, and behavior should match it unless there's a deliberate Go-specific design decision.
