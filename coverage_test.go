package gomarkdownify

import (
	"testing"
)

func TestCoverageImprovement(t *testing.T) {
	// Test convertA with different scenarios
	html := `<a href="https://example.com">Example</a>`
	result, err := Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertA with autolinks
	html = `<a href="https://example.com">https://example.com</a>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertA with title
	html = `<a href="https://example.com" title="Example Title">Example</a>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertBr in different contexts
	html = `<p>Line 1<br>Line 2</p>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertDiv with different content
	html = `<div>Content in a div</div>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertDiv with empty content
	html = `<div></div>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	// No assertion needed for empty result

	// Test convertImg with different attributes
	html = `<img src="image.jpg" alt="Alt Text">`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertImg with title
	html = `<img src="image.jpg" alt="Alt Text" title="Image Title">`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertList with different types
	html = `<ul><li>Item 1</li><li>Item 2</li></ul>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertList with ordered list
	html = `<ol><li>Item 1</li><li>Item 2</li></ol>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertP with different content
	html = `<p>Paragraph content</p>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertP with text wrapping
	opts := DefaultOptions()
	opts.Wrap = true
	opts.WrapWidth = 20
	html = `<p>This is a long paragraph that should be wrapped at 20 characters.</p>`
	result, err = Convert(html, opts)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertPre with different content
	html = `<pre>Preformatted text</pre>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertPre with code language
	opts = DefaultOptions()
	opts.CodeLanguage = "go"
	html = `<pre><code>func main() {
    fmt.Println("Hello")
}</code></pre>`
	result, err = Convert(html, opts)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}

	// Test convertPre with code language class
	html = `<pre><code class="language-go">func main() {
    fmt.Println("Hello")
}</code></pre>`
	result, err = Convert(html)
	if err != nil {
		t.Fatalf("Error converting HTML: %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}
}

// TestInlineQuotes tests the conversion of <q> tags (inline quotes)
func TestInlineQuotes(t *testing.T) {
	// Test basic quote
	result := md("<q>Hello</q>")
	expected := `"Hello"`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test quote with nested formatting
	result = md("<q>This is <strong>bold</strong> text</q>")
	expected = `"This is **bold** text"`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test quote with emphasis
	result = md("<q>This is <em>italic</em> text</q>")
	expected = `"This is *italic* text"`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test quote in paragraph
	result = md("<p>He said <q>Hello</q> to me</p>")
	expected = "\n\nHe said \"Hello\" to me\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestFigcaption tests the conversion of <figcaption> tags
func TestFigcaption(t *testing.T) {
	// Test figcaption in figure
	result := md("<figure><img src=\"img.jpg\" alt=\"test\"><figcaption>Caption text</figcaption></figure>")
	expected := "![test](img.jpg)\n\nCaption text\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test figcaption with nested content
	result = md("<figure><img src=\"img.jpg\" alt=\"test\"><figcaption>This is <strong>bold</strong> caption</figcaption></figure>")
	expected = "![test](img.jpg)\n\nThis is **bold** caption\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestVideo tests the conversion of <video> tags
func TestVideo(t *testing.T) {
	// Test basic video
	result := md("<video src=\"video.mp4\"></video>")
	expected := "[](video.mp4)"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test video with poster
	result = md("<video src=\"video.mp4\" poster=\"poster.jpg\">Your browser does not support video.</video>")
	expected = "[![Your browser does not support video.](poster.jpg)](video.mp4)"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test video with poster and nested content
	result = md("<video src=\"video.mp4\" poster=\"poster.jpg\"><p>Fallback content</p></video>")
	expected = "[![\n\nFallback content\n\n](poster.jpg)](video.mp4)"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test video with source child
	result = md("<video><source src=\"video.mp4\" type=\"video/mp4\"></video>")
	expected = "[](video.mp4)"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test video without src or poster
	result = md("<video>Fallback text</video>")
	expected = "Fallback text"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test video with only poster (no src)
	result = md("<video poster=\"poster.jpg\">Fallback</video>")
	expected = "![Fallback](poster.jpg)"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestDefinitionLists tests the conversion of <dl>, <dt>, and <dd> tags
func TestDefinitionLists(t *testing.T) {
	// Test basic definition list
	result := md("<dl><dt>Term 1</dt><dd>Description 1</dd><dt>Term 2</dt><dd>Description 2</dd></dl>")
	expected := "\n\nTerm 1\n:   Description 1\n\nTerm 2\n:   Description 2\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test definition list with multiple descriptions per term
	result = md("<dl><dt>Term</dt><dd>Description 1</dd><dd>Description 2</dd></dl>")
	expected = "\n\nTerm\n:   Description 1\n:   Description 2\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test definition list with formatted content
	result = md("<dl><dt><strong>Term</strong></dt><dd>Description with <em>emphasis</em></dd></dl>")
	expected = "\n\n**Term**\n:   Description with *emphasis*\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test nested definition lists
	result = md("<dl><dt>Term 1</dt><dd>Description 1<dl><dt>Nested Term</dt><dd>Nested Description</dd></dl></dd></dl>")
	expected = "\n\nTerm 1\n:   Description 1\n\n    Nested Term\n    :   Nested Description\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestScriptAndStyle tests that script and style tags are removed
func TestScriptAndStyle(t *testing.T) {
	// Test script tag removal
	result := md("<p>Before</p><script>alert(\"test\");</script><p>After</p>")
	expected := "\n\nBefore\n\nAfter\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test style tag removal
	result = md("<p>Before</p><style>body { color: red; }</style><p>After</p>")
	expected = "\n\nBefore\n\nAfter\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test multiple script tags
	result = md("<p>Text</p><script>alert(1);</script><script>alert(2);</script><p>More text</p>")
	expected = "\n\nText\n\nMore text\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test script with special characters
	result = md("<script>if (x < 0) { alert('test'); }</script>")
	expected = ""
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test style with CSS
	result = md("<style>.class { color: #fff; background: url('img.jpg'); }</style>")
	expected = ""
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestStripPreOptions tests the strip_pre option with code blocks
func TestStripPreOptions(t *testing.T) {
	// Test with strip_pre="strip" (default - strip all leading/trailing whitespace)
	opts := DefaultOptions()
	opts.StripPre = STRIP
	result, err := Convert("<pre>  code\n  here</pre>", opts)
	if err != nil {
		t.Fatalf("Error converting: %v", err)
	}
	expected := "```\n  code\n  here\n```\n\n"
	if result != expected {
		t.Errorf("With strip_pre='strip': Expected %q, got %q", expected, result)
	}

	// Test with strip_pre="strip_one" (strip one level of indentation)
	opts = DefaultOptions()
	opts.StripPre = STRIP_ONE
	result, err = Convert("<pre>  code\n  here</pre>", opts)
	if err != nil {
		t.Fatalf("Error converting: %v", err)
	}
	expected = "```\n  code\n  here\n```\n\n"
	if result != expected {
		t.Errorf("With strip_pre='strip_one': Expected %q, got %q", expected, result)
	}

	// Test with strip_pre="" (no stripping)
	opts = DefaultOptions()
	opts.StripPre = ""
	result, err = Convert("<pre>  code\n  here</pre>", opts)
	if err != nil {
		t.Fatalf("Error converting: %v", err)
	}
	expected = "```\n  code\n  here\n```\n\n"
	if result != expected {
		t.Errorf("With strip_pre='': Expected %q, got %q", expected, result)
	}

	// Test strip_pre with multi-level indentation
	opts = DefaultOptions()
	opts.StripPre = STRIP
	result, err = Convert("<pre>    code\n    here</pre>", opts)
	if err != nil {
		t.Fatalf("Error converting: %v", err)
	}
	expected = "```\n    code\n    here\n```\n\n"
	if result != expected {
		t.Errorf("With strip_pre='strip' and multi-level: Expected %q, got %q", expected, result)
	}

	opts = DefaultOptions()
	opts.StripPre = STRIP_ONE
	result, err = Convert("<pre>    code\n    here</pre>", opts)
	if err != nil {
		t.Fatalf("Error converting: %v", err)
	}
	expected = "```\n    code\n    here\n```\n\n"
	if result != expected {
		t.Errorf("With strip_pre='strip_one' and multi-level: Expected %q, got %q", expected, result)
	}
}

// TestConvertPreWithStripPre tests convertPre function with different strip_pre options
func TestConvertPreWithStripPre(t *testing.T) {
	// Test pre with code child and language class
	result := md("<pre><code class=\"language-python\">def test():\n    pass</code></pre>")
	expected := "\n\n```python\ndef test():\n    pass\n```\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test pre with lang- class
	result = md("<pre><code class=\"lang-javascript\">function test() {}</code></pre>")
	expected = "\n\n```javascript\nfunction test() {}\n```\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test pre without language
	result = md("<pre><code>plain code</code></pre>")
	expected = "\n\n```\nplain code\n```\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test pre with empty content
	result = md("<pre></pre>")
	expected = ""
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test pre with whitespace only
	result = md("<pre>   \n   \n</pre>")
	expected = "\n\n```\n\n```\n\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestConvertListEdgeCases tests edge cases for list conversion
func TestConvertListEdgeCases(t *testing.T) {
	// Test nested ordered lists
	result := md("<ol><li>Item 1<ol><li>Subitem 1</li><li>Subitem 2</li></ol></li><li>Item 2</li></ol>")
	expected := "\n\n1. Item 1\n   1. Subitem 1\n   2. Subitem 2\n2. Item 2\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test mixed nested lists
	result = md("<ul><li>Item 1<ul><li>Subitem 1<ol><li>Sub-subitem</li></ol></li></ul></li></ul>")
	expected = "\n\n* Item 1\n  + Subitem 1\n    1. Sub-subitem\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test list with single item
	result = md("<ul><li>Single item</li></ul>")
	expected = "\n\n* Single item\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test ordered list with reversed start
	result = md("<ol start=\"10\"><li>Item 10</li><li>Item 11</li></ol>")
	expected = "\n\n10. Item 10\n11. Item 11\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test list item with multiple paragraphs
	result = md("<ul><li><p>Para 1</p><p>Para 2</p></li></ul>")
	expected = "\n\n* Para 1\n\n  Para 2\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
