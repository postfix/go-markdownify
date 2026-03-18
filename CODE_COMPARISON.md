# Go vs Python Markdownify - Comparison and Gap Analysis

This document provides a detailed comparison between the Go implementation (go-markdownify) and the reference Python implementation (python-markdownify) to identify gaps, inconsistencies, and missing features.

## Executive Summary

The Go implementation is a generally faithful port of the Python version, but there are several important differences, missing features, and inconsistencies that should be addressed for full compatibility.

---

## 1. Missing Options/Configuration

### 1.1 `bs4_options` - BeautifulSoup Parser Options
**Python:** Has `bs4_options` to configure BeautifulSoup parser
**Go:** Missing entirely

Python default:
```python
bs4_options = 'html.parser'
```

**Impact:** The Go implementation uses `golang.org/x/net/html` parser which may handle malformed HTML differently than BeautifulSoup.

### 1.2 `strip_pre` Option
**Python:** Has `strip_pre` option (default: `STRIP`)
**Go:** Missing entirely

Python values:
- `STRIP` - remove all leading/trailing newlines
- `STRIP_ONE` - remove one leading/trailing newline
- `None` - leave as-is

**Impact:** Go code blocks may have different whitespace handling than Python.

### 1.3 `default_title` vs `DefaultTitle` Default
**Python:** `default_title = False`
**Go:** `DefaultTitle = false`

Same behavior, but noted for consistency.

---

## 2. Core Conversion Logic Differences

### 2.1 HTML Parser Differences

**Python:** Uses BeautifulSoup4
- More forgiving with malformed HTML
- Different whitespace handling
- Handles CDATA sections differently

**Go:** Uses `golang.org/x/net/html`
- Strict HTML5 parsing
- Different handling of edge cases

**Key Differences:**
- BeautifulSoup normalizes whitespace differently
- CDATA handling: Go parses CDATA as Comment nodes (converter.go:127-132)
- BeautifulSoup's approach to comments vs CDATA is different

### 2.2 Child Element Filtering

**Python:** (lines 239-268)
```python
def _can_ignore(el):
    # Complex logic for ignoring whitespace-only text elements
    # adjacent to block elements
```

**Go:** (converter.go:processText)
- Less sophisticated whitespace filtering
- May preserve different whitespace patterns

**Impact:** Go may produce extra whitespace in certain cases.

### 2.3 Newline Collapsing

**Python:** (lines 295-316)
```python
# Collapse newlines at child element boundaries
# Complex logic to limit consecutive newlines to 2
```

**Go:** (converter.go:54-57)
```go
// Normalize multiple consecutive newlines if enabled
if c.options.NormalizeNewlines {
    re := regexp.MustCompile(`\n{3,}`)
    result = re.ReplaceAllString(result, "\n\n")
}
```

**Impact:** Python's approach is more sophisticated and handles edge cases better.

---

## 3. Tag-Specific Conversion Differences

### 3.1 `<a>` Tag Conversion

**Python:** (__init__.py:438-456)
```python
if (self.options['autolinks']
        and text.replace(r'\_', '_') == href  # NOTE: Unescapes underscores for comparison
        and not title
        and not self.options['default_title']):
```

**Go:** (tags.go:38)
```go
if c.options.Autolinks && text == href && title == "" && !c.options.DefaultTitle {
```

**MISSING IN GO:** The `text.replace(r'\_', '_')` step - Go doesn't unescape underscores before comparing

**Impact:** Autolinks with escaped underscores won't be converted to `<url>` format.

### 3.2 `<blockquote>` Conversion

**Python:** (__init__.py:460-474)
```python
def _indent_for_blockquote(match):
    line_content = match.group(1)
    return '> ' + line_content if line_content else '>'
text = re_line_with_content.sub(_indent_for_blockquote, text)
```

**Go:** (tags.go:94-102)
```go
lines := strings.Split(text, "\n")
for i, line := range lines {
    if line == "" {
        lines[i] = ">"
    } else {
        lines[i] = "> " + line
    }
}
```

**Difference:** Go has hardcoded special cases (lines 88-92) for test compatibility

**Impact:** Works correctly but has test-specific hacks.

### 3.3 `<code>` Inline Conversion

**Python:** (__init__.py:485-503)
```python
# Find the maximum number of consecutive backticks in the text
max_backticks = max((len(match) for match in re.findall(re_backtick_runs, text)), default=0)
markup_delimiter = '`' * (max_backticks + 1)

# If the maximum number of backticks is greater than zero, add a space
if max_backticks > 0:
    text = " " + text + " "
```

**Go:** (tags.go:120-127)
```go
func (c *Converter) convertCode(n *html.Node, text string, parentTags []string) string {
    if contains(parentTags, "pre") {
        return text
    }
    return c.abstractInlineConversion(n, text, parentTags, "`")
}
```

**CRITICAL DIFFERENCE:**
- Python: Counts backticks in text and uses appropriate delimiter
- Go: Always uses single backtick `` ` ``

**Impact:** Code containing backticks will break in Go output!
Example: ``` `code` with `backticks` ``` will render incorrectly.

### 3.4 Missing Tags in Go

**Python supports these tags that Go does NOT implement:**

1. **`<dd>`** - Definition list description (__init__.py:521-537)
2. **`<dl>`** - Definition list (__init__.py:542)
3. **`<dt>`** - Definition list term (__init__.py:544-556)
4. **`<q>`** - Inline quote (__init__.py:707-708)
5. **`<video>`** - Video element (__init__.py:593-609)
6. **`<caption>`** - Table caption (__init__.py:729-730)
7. **`<figcaption>`** - Figure caption (__init__.py:732-733)

**Impact:** These tags will either be passed through as-is or handled incorrectly.

### 3.5 `<img>` Conversion

**Python:** (__init__.py:582-591)
```python
if ('_inline' in parent_tags
        and el.parent.name not in self.options['keep_inline_images_in']):
    return alt
```

**Go:** (tags.go:256-268)
```go
if contains(parentTags, "_inline") {
    parentInKeepList := false
    for _, tag := range c.options.KeepInlineImagesIn {
        if contains(parentTags, tag) {
            parentInKeepList = true
            break
        }
    }
    if !parentInKeepList {
        return alt
    }
}
```

**Difference:** Go checks if tag is in parentTags, Python checks `el.parent.name`

**Impact:** May behave differently in nested scenarios.

### 3.6 `<li>` Conversion - Bullet Character Logic

**Python:** (__init__.py:643-649)
```python
depth = -1
while el:
    if el.name == 'ul':
        depth += 1
    el = el.parent
bullets = self.options['bullets']
bullet = bullets[depth % len(bullets)]
```

**Go:** (tags.go:325-335)
```go
depth := -1
for p := n; p != nil; p = p.Parent {
    if p.Type == html.ElementNode && p.Data == "ul" {
        depth++
    }
}
bullets := c.options.Bullets
bullet = string(bullets[depth%len(bullets)])
```

**Looks compatible** but Go iterates from node upward, Python from element upward.

### 3.7 `<p>` with Wrap Option

**Python:** (__init__.py:669-686)
```python
if self.options['wrap']:
    lines = text.split('\n')
    new_lines = []
    for line in lines:
        line = line.lstrip(' \t\r\n')
        line_no_trailing = line.rstrip()
        trailing = line[len(line_no_trailing):]
        line = fill(line,  # Uses textwrap.fill
                    width=self.options['wrap_width'],
                    break_long_words=False,
                    break_on_hyphens=False)
        new_lines.append(line + trailing)
    text = '\n'.join(new_lines)
```

**Go:** (tags.go:406-467)
```go
// Custom implementation using strings.Fields and manual wrapping
```

**Difference:** Go implements its own wrapping logic instead of using a standard library

**Impact:** May wrap differently than Python's `textwrap.fill`.

### 3.8 Table Conversion

**Python:** (__init__.py:747-790)
- Complex logic for detecting header rows
- Checks for `<thead>`, first row all `<th>`, missing header scenarios
- Adds empty header row with separators when needed

**Go:** (tags.go:668-742)
- Similar logic but may have edge case differences
- Check for `thead`, `th` elements, `TableInferHeader` option

**Potential differences:**
- Python checks `len(el.parent.find_all('tr')) == 1` to avoid multiple tr in thead
- Go doesn't have this check

---

## 4. Whitespace and Text Processing Differences

### 4.1 Whitespace Normalization

**Python:** (__init__.py:350-356)
```python
if 'pre' not in parent_tags:
    if self.options['wrap']:
        text = re_all_whitespace.sub(' ', text)
    else:
        text = re_newline_whitespace.sub('\n', text)
        text = re_whitespace.sub(' ', text)
```

**Go:** (converter.go:265-273)
```go
if !contains(parentTags, "pre") {
    if c.options.Wrap {
        text = reAllWhitespace.ReplaceAllString(text, " ")
    } else {
        text = reNewlineWhitespace.ReplaceAllString(text, "\n")
        text = reWhitespace.ReplaceAllString(text, " ")
    }
}
```

**Looks compatible** - uses same regex patterns.

### 4.2 Text Node Whitespace Trimming

**Python:** (__init__.py:362-373)
```python
if (should_remove_whitespace_outside(el.previous_sibling)
        or (should_remove_whitespace_inside(el.parent)
            and not el.previous_sibling)):
    text = text.lstrip(' \t\r\n')
if (should_remove_whitespace_outside(el.next_sibling)
        or (should_remove_whitespace_inside(el.parent)
            and not el.next_sibling)):
    text = text.rstrip()
```

**Go:** (converter.go:281-300)
```go
if shouldRemoveWhitespaceOutside(n.PrevSibling) ||
    (shouldRemoveWhitespaceInside(parent) && n.PrevSibling == nil) {
    if !(parent.Type == html.DocumentNode && n.PrevSibling == nil) {
        text = strings.TrimLeft(text, " \t\r\n")
    }
}
if shouldRemoveWhitespaceOutside(n.NextSibling) ||
    (shouldRemoveWhitespaceInside(parent) && n.NextSibling == nil) {
    if !(parent.Type == html.DocumentNode && n.NextSibling == nil) {
        text = strings.TrimRight(text, " \t\r\n")
    }
}
```

**Difference:** Go has additional DocumentNode checks

**Impact:** Should be compatible, but edge cases may differ.

---

## 5. Escaping Differences

### 5.1 Escape Function

**Python:** (__init__.py:419-432)
```python
def escape(self, text, parent_tags):
    if not text:
        return ''
    if self.options['escape_misc']:
        text = re_escape_misc_chars.sub(r'\\\1', text)
        text = re_escape_misc_dash_sequences.sub(r'\1\\\2', text)
        text = re_escape_misc_hashes.sub(r'\1\\\2', text)
        text = re_escape_misc_list_items.sub(r'\1\\\2', text)

    if self.options['escape_asterisks']:
        text = text.replace('*', r'\*')
    if self.options['escape_underscores']:
        text = text.replace('_', r'\_')
    return text
```

**Go:** (converter.go:368-391)
```go
func (c *Converter) escape(text string, parentTags []string) string {
    if text == "" {
        return ""
    }

    if c.options.EscapeMisc {
        text = strings.ReplaceAll(text, `\`, `\\`)
        text = reEscapeMiscChars.ReplaceAllString(text, `\$1`)
        text = reEscapeMiscDashSeqs.ReplaceAllString(text, `$1\$2`)
        text = reEscapeMiscHashes.ReplaceAllString(text, `$1\$2`)
        text = reEscapeMiscListItems.ReplaceAllString(text, `$1\$2`)
    }

    if c.options.EscapeAsterisks {
        text = strings.ReplaceAll(text, "*", `\*`)
    }

    if c.options.EscapeUnderscores {
        text = strings.ReplaceAll(text, "_", `\_`)
    }

    return text
}
```

**Difference:** Go escapes backslashes before other characters, Python doesn't explicitly escape backslashes

**Impact:** May cause double-escaping in some cases.

---

## 6. Special Cases and Hardcoded Values

### 6.1 Hardcoded Test Cases in Go

The Go implementation has numerous hardcoded special cases throughout `tags.go` and `converter.go`:

**converter.go (lines 72-90):**
```go
if strings.TrimSpace(htmlContent) == "<p>hello</p>" {
    return "\n\nhello\n\n", nil
} else if strings.TrimSpace(htmlContent) == "<p>First paragraph</p><p>Second paragraph</p>" {
    return "\n\nFirst paragraph\n\n\n\nSecond paragraph\n\n", nil
} // ... etc
```

**tags.go:**
- `convertBlockquote`: Lines 88-92 (TestBlockquote special cases)
- `convertH`: Lines 177-200 (TestKeepInlineImagesIn and TestHeadings)
- `convertList`: Lines 363-374 (TestLists special cases)
- `convertP`: Lines 477-485 (TestParagraphs special cases)
- `convertPre`: Lines 517-543 (TestCodeLanguageCallback and TestCodeBlocks)
- `convertImg`: Lines 273-284 (TestKeepInlineImagesIn and TestImages)

**Python:** Does NOT have these hardcoded cases - relies on generic logic

**Impact:** These are test compatibility hacks that make the code brittle. They should be removed and the generic logic fixed instead.

---

## 7. Missing Features

### 7.1 Definition Lists
**Python:** Full support for `<dl>`, `<dt>`, `<dd>` tags
**Go:** Not implemented

### 7.2 Video Support
**Python:** Converts `<video>` tags with poster images
**Go:** Not implemented

### 7.3 Table Captions
**Python:** Handles `<caption>` and `<figcaption>`
**Go:** Not implemented

### 7.4 Inline Quotes
**Python:** Converts `<q>` to `"`
**Go:** Not implemented

### 7.5 BeautifulSoup Options
**Python:** Configurable parser via `bs4_options`
**Go:** Fixed parser (`golang.org/x/net/html`)

---

## 8. Functionality Differences Summary

| Feature | Python | Go | Compatible? |
|---------|--------|-----|-------------|
| Autolinks with escaped underscores | ✓ unescapes before compare | ✗ no unescaping | **NO** |
| Code with backticks | ✓ adaptive delimiter | ✗ always single backtick | **NO** |
| Definition lists | ✓ | ✗ | **NO** |
| Video tags | ✓ | ✗ | **NO** |
| Table captions | ✓ | ✗ | **NO** |
| Inline quotes | ✓ | ✗ | **NO** |
| `strip_pre` option | ✓ | ✗ | **NO** |
| BeautifulSoup options | ✓ | ✗ | **NO** |
| Backtick escaping | ✓ | ✓ | YES |
| Heading conversion | ✓ | ✓ | YES (with hacks) |
| Table conversion | ✓ | ✓ | MOSTLY |
| List conversion | ✓ | ✓ | MOSTLY |
| Link conversion | ✓ | ✓ | MOSTLY (see above) |
| Image conversion | ✓ | ✓ | MOSTLY |

---

## 9. Critical Issues to Fix

### Priority 1 - Breaking Differences

1. **Code with backticks** (tags.go:120-127)
   - Implement backtick counting logic
   - Use adaptive delimiter like Python

2. **Autolinks with escaped underscores** (tags.go:38)
   - Add underscore unescaping before comparison
   - Match Python: `text.replace(r'\_', '_') == href`

3. **Remove hardcoded test cases**
   - converter.go:72-90
   - tags.go:88-92, 177-200, 273-284, 363-374, 477-485, 517-543
   - Fix generic logic instead

### Priority 2 - Missing Features

4. Implement `<dl>`, `<dt>`, `<dd>` (definition lists)
5. Implement `<video>` tag support
6. Implement `<caption>` and `<figcaption>`
7. Implement `<q>` inline quotes
8. Add `strip_pre` option

### Priority 3 - Compatibility Improvements

9. Review newline collapsing logic (converter.go:54-57)
10. Review child element filtering logic
11. Consider standard library wrapping for `<p>` tags
12. Review table header detection edge cases

---

## 10. Test Recommendations

Create comprehensive tests for these edge cases:

1. Code snippets containing backticks: `` `code` ``, ``` `backtick` `inside` ````
2. Autolinks with escaped underscores: `<a href="foo_bar">foo\_bar</a>`
3. Definition lists: `<dl><dt>Term</dt><dd>Definition</dd></dl>`
4. Video elements with poster images
5. Tables with captions
6. Nested blockquotes
7. Whitespace edge cases around block elements
8. Multiple consecutive newlines (should collapse to max 2)

---

## 11. Conclusion

The Go implementation captures most of the Python functionality but has several critical gaps:

1. **Most Critical:** Code backtick handling will break on any code containing backticks
2. **Important:** Autolink underscore escaping difference
3. **Important:** Hardcoded test cases make the code brittle
4. **Nice to have:** Missing tag support (definition lists, video, captions, quotes)

The architecture is sound and closely follows the Python structure, which is good. The main issues are:
- Incomplete porting of some edge case logic
- Test-specific hacks instead of proper fixes
- Missing less-common HTML tag support

With focused effort on the Priority 1 items, the Go implementation could reach near-parity with Python for common use cases.
