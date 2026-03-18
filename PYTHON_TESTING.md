# Python Markdownify Testing Setup

This document explains how to set up and use the Python markdownify reference implementation for testing and comparison with the Go implementation.

## Setup

### Using uv (Recommended)

```bash
# Create virtual environment (already done)
uv venv refcode/python-markdownify/.venv

# Install dependencies
uv pip install six beautifulsoup4 pytest --python refcode/python-markdownify/.venv/bin/python
```

### Activation

```bash
# Activate the virtual environment
source refcode/python-markdownify/.venv/bin/activate

# Or use directly
refcode/python-markdownify/.venv/bin/python
```

## Running Python Tests

```bash
# Run all tests
cd refcode/python-markdownify
.venv/bin/python -m pytest tests/

# Run specific test file
.venv/bin/python -m pytest tests/test_basic.py -v

# Run with verbose output
.venv/bin/python -m pytest tests/ -v
```

## Comparing Go vs Python Output

### Quick Comparison Script

Create a test script to compare outputs:

```python
import markdownify

# Test case
html = '<p>Hello</p>'
result = markdownify.markdownify(html)
print(repr(result))

# With options
result = markdownify.markdownify(
    '<h1>Hello</h1>',
    heading_style='atx',
    strip_document='strip'
)
print(repr(result))
```

### Go Equivalent

```go
import gomarkdownify

// Test case
result, _ := gomarkdownify.Convert("<p>Hello</p>")
fmt.Printf("%q\n", result)

// With options
opts := gomarkdownify.DefaultOptions()
opts.HeadingStyle = gomarkdownify.ATX
opts.StripDocument = gomarkdownify.STRIP
result, _ := gomarkdownify.Convert("<h1>Hello</h1>", opts)
fmt.Printf("%q\n", result)
```

## Key Differences Between Python and Go

### Default StripDocument Option

- **Python**: Default is `strip_document='strip'` (removes both leading and trailing newlines)
- **Go**: Default is `StripDocument=LSTRIP` (removes only leading newlines)

**Impact**: When comparing outputs, ensure both use the same `strip_document` setting.

### Example Comparison

| HTML | Python (default) | Go (default LSTRIP) | Go (STRIP) |
|------|------------------|---------------------|------------|
| `<p>Hello</p>` | `'Hello'` | `"Hello\n\n"` | `"Hello"` |
| `<h1>Hello</h1>` | `'Hello\n====='` | `"Hello\n=====\n\n"` | `"Hello\n====="` |

## Testing Guidelines

1. **Always use matching strip_document settings** when comparing
2. **Test with both default and custom options**
3. **Verify edge cases** like:
   - Empty content
   - Nested elements
   - Code blocks with backticks
   - Tables with captions
   - Definition lists
   - Video elements

## Running Comparison Tests

See `comparison_tests.py` for automated comparison between Go and Python implementations.
