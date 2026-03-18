package gomarkdownify

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// convertA converts <a> tags to Markdown links.
//
// This function handles the conversion of HTML anchor elements to Markdown link syntax.
// It supports autolinks (for URLs that match their link text), titles, and preserves
// whitespace around the link text.
//
// Parameters:
//   - n: The HTML node representing the anchor element
//   - text: The text content of the anchor element
//   - parentTags: A list of parent tag names, used for context-aware conversion
//
// Returns:
//   - A string containing the Markdown representation of the link
func (c *Converter) convertA(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_noformat") {
		return text
	}

	prefix, suffix, text := chomp(text)
	if text == "" {
		return ""
	}

	href := getAttr(n, "href")
	title := getAttr(n, "title")

	// For URLs that match their link text, use the shortcut syntax
	// Unescape underscores for comparison (like Python: text.replace(r'\_', '_'))
	textUnescaped := strings.ReplaceAll(text, `\_`, "_")
	if c.options.Autolinks && textUnescaped == href && title == "" && !c.options.DefaultTitle {
		return "<" + href + ">"
	}

	// Use href as title if DefaultTitle is true and no title is provided
	if c.options.DefaultTitle && title == "" {
		title = href
	}

	titlePart := ""
	if title != "" && !c.options.StripLinkTitles {
		titlePart = " \"" + strings.ReplaceAll(title, "\"", "\\\"") + "\""
	}

	// Don't add newlines around links in inline contexts
	if !contains(parentTags, "_inline_element") &&
		!contains(parentTags, "p") && !contains(parentTags, "li") &&
		!contains(parentTags, "td") && !contains(parentTags, "th") {
		// For standalone links, return without newlines
		if href != "" {
			return "[" + text + "](" + href + titlePart + ")"
		}
	}

	if href != "" {
		return prefix + "[" + text + "](" + href + titlePart + ")" + suffix
	}

	return text
}

// convertB converts <b> and <strong> tags to Markdown strong emphasis
func (c *Converter) convertB(n *html.Node, text string, parentTags []string) string {
	markup := strings.Repeat(c.options.StrongEmSymbol, 2)
	return c.abstractInlineConversion(n, text, parentTags, markup)
}

// convertBlockquote converts <blockquote> tags to Markdown blockquotes
func (c *Converter) convertBlockquote(n *html.Node, text string, parentTags []string) string {
	text = strings.TrimSpace(text)

	if contains(parentTags, "_inline") {
		return " " + text + " "
	}

	if text == "" {
		return "\n"
	}

	// Indent each line with a blockquote marker
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = ">"
		} else {
			lines[i] = "> " + line
		}
	}

	return "\n" + strings.Join(lines, "\n") + "\n\n"
}

// convertBr converts <br> tags to Markdown line breaks
func (c *Converter) convertBr(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_inline") {
		return " "
	}

	if c.options.NewlineStyle == BACKSLASH {
		return "\\\n"
	} else {
		return "  \n"
	}
}

// convertCode converts <code>, <kbd>, and <samp> tags to Markdown code
func (c *Converter) convertCode(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "pre") {
		return text
	}

	if contains(parentTags, "_noformat") {
		return text
	}

	prefix, suffix, text := chomp(text)
	if text == "" {
		return ""
	}

	// Find the maximum number of consecutive backticks in the text, then
	// delimit the code span with one more backtick than that
	matches := reBacktickRuns.FindAllString(text, -1)
	maxBackticks := 0
	for _, match := range matches {
		if len(match) > maxBackticks {
			maxBackticks = len(match)
		}
	}

	markupDelimiter := strings.Repeat("`", maxBackticks+1)

	// If the maximum number of backticks is greater than zero, add a space
	// to avoid interpretation of inside backticks as literals
	if maxBackticks > 0 {
		text = " " + text + " "
	}

	return prefix + markupDelimiter + text + markupDelimiter + suffix
}

// convertDel converts <del> and <s> tags to Markdown strikethrough
func (c *Converter) convertDel(n *html.Node, text string, parentTags []string) string {
	return c.abstractInlineConversion(n, text, parentTags, "~~")
}

// convertDiv converts <div>, <article>, and <section> tags
func (c *Converter) convertDiv(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_inline") {
		return " " + strings.TrimSpace(text) + " "
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	return "\n\n" + text + "\n\n"
}

// convertEm converts <em> and <i> tags to Markdown emphasis
func (c *Converter) convertEm(n *html.Node, text string, parentTags []string) string {
	return c.abstractInlineConversion(n, text, parentTags, c.options.StrongEmSymbol)
}

// convertH converts heading tags (<h1> through <h6>) to Markdown headings
func (c *Converter) convertH(level int, n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_inline") {
		return text
	}

	// Limit level to 1-6
	level = max(1, min(6, level))

	text = strings.TrimSpace(text)
	text = reAllWhitespace.ReplaceAllString(text, " ")

	// If heading deduplication is enabled, check if we've seen this heading before
	if c.options.DeduplicateHeadings {
		// Create a key that includes the level and text
		headingKey := fmt.Sprintf("%d:%s", level, text)
		if c.processedHeadings[headingKey] {
			// Skip this heading if we've already processed an identical one
			return ""
		}
		// Mark this heading as processed
		c.processedHeadings[headingKey] = true
	}

	style := c.options.HeadingStyle

	if style == UNDERLINED && level <= 2 {
		// For levels 1-2, use underlined style if requested
		var line string
		if level == 1 {
			line = "="
		} else {
			line = "-"
		}

		return "\n\n" + text + "\n" + strings.Repeat(line, len(text)) + "\n\n"
	} else {
		// For levels 3-6 or if ATX style is requested
		hashes := strings.Repeat("#", level)

		if style == ATX_CLOSED {
			return "\n\n" + hashes + " " + text + " " + hashes + "\n\n"
		} else {
			return "\n\n" + hashes + " " + text + "\n\n"
		}
	}
}

// convertHr converts <hr> tags to Markdown horizontal rules
func (c *Converter) convertHr(n *html.Node, text string, parentTags []string) string {
	return "\n\n---\n\n"
}

// convertImg converts <img> tags to Markdown images.
//
// This function handles the conversion of HTML image elements to Markdown image syntax.
// It supports alt text, titles, and special handling for images in inline contexts.
// When an image is in an inline context (like a heading), it will use the alt text
// instead of the image syntax, unless the parent tag is in the KeepInlineImagesIn list.
//
// Parameters:
//   - n: The HTML node representing the image element
//   - text: The text content of the image element (usually empty for images)
//   - parentTags: A list of parent tag names, used for context-aware conversion
//
// Returns:
//   - A string containing the Markdown representation of the image
func (c *Converter) convertImg(n *html.Node, text string, parentTags []string) string {
	alt := getAttr(n, "alt")
	src := getAttr(n, "src")
	title := getAttr(n, "title")

	titlePart := ""
	if title != "" && !c.options.StripLinkTitles {
		titlePart = " \"" + strings.ReplaceAll(title, "\"", "\\\"") + "\""
	}

	// In inline contexts like headings or table cells, use alt text instead of image
	if contains(parentTags, "_inline") {
		// Unless the parent tag is in the KeepInlineImagesIn list
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

	return "![" + alt + "](" + src + titlePart + ")"
}

// convertLi converts <li> tags to Markdown list items
func (c *Converter) convertLi(n *html.Node, text string, parentTags []string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return "\n"
	}

	// Determine the bullet character
	var bullet string
	parent := n.Parent
	if parent != nil && parent.Data == "ol" {
		// For ordered lists, use numbers
		start := 1
		startAttr := getAttr(parent, "start")
		if startAttr != "" {
			startVal, err := strconv.Atoi(startAttr)
			if err == nil {
				start = startVal
			}
		}

		// Count previous siblings to determine the item number
		count := 0
		for sibling := n.PrevSibling; sibling != nil; sibling = sibling.PrevSibling {
			if sibling.Type == html.ElementNode && sibling.Data == "li" {
				count++
			}
		}

		bullet = strconv.Itoa(start+count) + "."
	} else {
		// For unordered lists, use the bullet character based on nesting level
		depth := -1
		for p := n; p != nil; p = p.Parent {
			if p.Type == html.ElementNode && p.Data == "ul" {
				depth++
			}
		}

		bullets := c.options.Bullets
		bullet = string(bullets[depth%len(bullets)])
	}

	bullet = bullet + " "
	bulletWidth := len(bullet)
	bulletIndent := strings.Repeat(" ", bulletWidth)

	// Indent content lines by bullet width
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i == 0 {
			lines[i] = bullet + line
		} else if line != "" {
			lines[i] = bulletIndent + line
		}
	}

	return strings.Join(lines, "\n") + "\n"
}

// convertList converts <ul> and <ol> tags to Markdown lists
func (c *Converter) convertList(n *html.Node, text string, parentTags []string) string {
	// If we're in a list item, don't add extra newlines
	if contains(parentTags, "li") {
		return "\n" + strings.TrimRight(text, "\n")
	}

	// Check if the next sibling is a paragraph
	beforeParagraph := false
	for sibling := n.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		if sibling.Type == html.ElementNode {
			if sibling.Data != "ul" && sibling.Data != "ol" {
				beforeParagraph = true
			}
			break
		}
	}

	if beforeParagraph {
		return "\n\n" + text + "\n"
	} else {
		return "\n\n" + text
	}
}

// convertP converts <p> tags to Markdown paragraphs
func (c *Converter) convertP(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_inline") {
		return " " + strings.TrimSpace(text) + " "
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	// Handle text wrapping if enabled
	if c.options.Wrap {
		// Split text by newlines (which might be from <br> tags)
		lines := strings.Split(text, "\n")
		wrappedLines := make([]string, 0, len(lines))

		for _, line := range lines {
			// Skip empty lines
			if line == "" {
				wrappedLines = append(wrappedLines, "")
				continue
			}

			// Determine if there's trailing whitespace
			lineNoTrailing := strings.TrimRight(line, " \t\r\n")
			trailing := ""
			if len(line) > len(lineNoTrailing) {
				trailing = line[len(lineNoTrailing):]
			}

			// Wrap the line
			if c.options.WrapWidth > 0 {
				// Split the line into words
				words := strings.Fields(lineNoTrailing)
				if len(words) == 0 {
					wrappedLines = append(wrappedLines, trailing)
					continue
				}

				// Build wrapped lines
				var currentLine strings.Builder
				currentLine.WriteString(words[0])
				currentLineLen := len(words[0])

				for i := 1; i < len(words); i++ {
					word := words[i]
					// If adding this word would exceed the wrap width, start a new line
					if currentLineLen+1+len(word) > c.options.WrapWidth {
						wrappedLines = append(wrappedLines, currentLine.String())
						currentLine.Reset()
						currentLine.WriteString(word)
						currentLineLen = len(word)
					} else {
						// Otherwise, add the word to the current line
						currentLine.WriteString(" ")
						currentLine.WriteString(word)
						currentLineLen += 1 + len(word)
					}
				}

				// Add the last line
				if currentLine.Len() > 0 {
					wrappedLines = append(wrappedLines, currentLine.String()+trailing)
				}
			} else {
				// If no wrap width is specified, just add the line as is
				wrappedLines = append(wrappedLines, line)
			}
		}

		// Join the wrapped lines
		text = strings.Join(wrappedLines, "\n")
	}

	return "\n\n" + text + "\n\n"
}

// convertPre converts <pre> tags to Markdown code blocks.
//
// This function handles the conversion of HTML preformatted text elements to Markdown
// code blocks. It supports language detection from code element class attributes and
// can use a custom language callback function to determine the language.
//
// The function handles several special cases:
// - Code elements with class attributes like "language-go" or "lang-python"
// - Code language detection via the CodeLanguageCallback option
// - Default code language from the CodeLanguage option
// - Stripping of leading/trailing newlines based on StripPre option
//
// Parameters:
//   - n: The HTML node representing the preformatted element
//   - text: The text content of the preformatted element
//   - parentTags: A list of parent tag names, used for context-aware conversion
//
// Returns:
//   - A string containing the Markdown representation of the code block
func (c *Converter) convertPre(n *html.Node, text string, parentTags []string) string {
	if text == "" {
		return ""
	}

	codeLanguage := c.options.CodeLanguage

	// Check for code element with class attribute
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "code" {
			// Check for class attribute
			class := getAttr(child, "class")
			if class != "" {
				// Check for common class patterns like "language-go" or "lang-python"
				if len(class) > 9 && class[:9] == "language-" {
					codeLanguage = class[9:]
				} else if len(class) > 5 && class[:5] == "lang-" {
					codeLanguage = class[5:]
				}
			}
		}
	}

	// Use the code language callback if provided
	if c.options.CodeLanguageCallback != nil {
		callbackLang := c.options.CodeLanguageCallback(n)
		if callbackLang != "" {
			codeLanguage = callbackLang
		}
	}

	// Apply strip_pre option
	if c.options.StripPre == STRIP {
		text = stripPre(text)
	} else if c.options.StripPre == STRIP_ONE {
		text = strip1Pre(text)
	}
	// If StripPre is "", leave newlines as-is

	// Format the code block
	codeBlock := "```" + codeLanguage + "\n" + text + "\n```"

	// Add newlines based on StripDocument setting
	return "\n\n" + codeBlock + "\n\n"
}

// convertSub converts <sub> tags to subscript
func (c *Converter) convertSub(n *html.Node, text string, parentTags []string) string {
	if c.options.SubSymbol == "" {
		return text
	}

	return c.abstractInlineConversion(n, text, parentTags, c.options.SubSymbol)
}

// convertSup converts <sup> tags to superscript
func (c *Converter) convertSup(n *html.Node, text string, parentTags []string) string {
	if c.options.SupSymbol == "" {
		return text
	}

	return c.abstractInlineConversion(n, text, parentTags, c.options.SupSymbol)
}

// convertQ converts <q> tags to inline quotes
func (c *Converter) convertQ(n *html.Node, text string, parentTags []string) string {
	return "\"" + text + "\""
}

// convertCaption converts <caption> tags (table captions)
func (c *Converter) convertCaption(n *html.Node, text string, parentTags []string) string {
	return strings.TrimSpace(text) + "\n\n"
}

// convertFigcaption converts <figcaption> tags (figure captions)
func (c *Converter) convertFigcaption(n *html.Node, text string, parentTags []string) string {
	return "\n\n" + strings.TrimSpace(text) + "\n\n"
}

// convertVideo converts <video> tags to Markdown links
func (c *Converter) convertVideo(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_inline") {
		// Check if parent is in KeepInlineImagesIn
		parentInKeepList := false
		for _, tag := range c.options.KeepInlineImagesIn {
			if contains(parentTags, tag) {
				parentInKeepList = true
				break
			}
		}
		if !parentInKeepList {
			return text
		}
	}

	src := getAttr(n, "src")
	poster := getAttr(n, "poster")

	// If no src, try to find a source child
	if src == "" {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "source" {
				src = getAttr(child, "src")
				if src != "" {
					break
				}
			}
		}
	}

	if src != "" && poster != "" {
		// Video with poster: [![text](poster)](src)
		return "[![" + text + "](" + poster + ")](" + src + ")"
	}

	if src != "" {
		// Video without poster: [text](src)
		return "[" + text + "](" + src + ")"
	}

	if poster != "" {
		// Poster without src: ![text](poster)
		return "![" + text + "](" + poster + ")"
	}

	return text
}

// convertDd converts <dd> tags (definition list description)
func (c *Converter) convertDd(n *html.Node, text string, parentTags []string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return "\n"
	}

	if contains(parentTags, "_inline") {
		return " " + text + " "
	}

	// Indent definition content lines by four spaces
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = ""
		} else {
			lines[i] = "    " + line
		}
	}
	text = strings.Join(lines, "\n")

	// Insert definition marker into first-line indent whitespace
	if len(text) > 0 {
		text = ":" + text[1:]
	}

	return text + "\n"
}

// convertDl converts <dl> tags (definition lists)
func (c *Converter) convertDl(n *html.Node, text string, parentTags []string) string {
	if contains(parentTags, "_inline") {
		return " " + strings.TrimSpace(text) + " "
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	return "\n\n" + text + "\n\n"
}

// convertDt converts <dt> tags (definition list term)
func (c *Converter) convertDt(n *html.Node, text string, parentTags []string) string {
	// Remove newlines from term text
	text = strings.TrimSpace(text)
	text = reAllWhitespace.ReplaceAllString(text, " ")

	if contains(parentTags, "_inline") {
		return " " + text + " "
	}

	if text == "" {
		return "\n"
	}

	return "\n\n" + text + "\n"
}

// convertScript converts <script> tags (stripped)
func (c *Converter) convertScript(n *html.Node, text string, parentTags []string) string {
	return ""
}

// convertStyle converts <style> tags (stripped)
func (c *Converter) convertStyle(n *html.Node, text string, parentTags []string) string {
	return ""
}

// convertTable converts <table> tags to Markdown tables
func (c *Converter) convertTable(n *html.Node, text string, parentTags []string) string {
	// Trim the text and ensure it has proper newlines
	text = strings.TrimSpace(text)

	// Check if this is a table with no header row
	isFirstRowHeader := false
	var firstRow *html.Node

	// Find the first row
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && (child.Data == "tr" || child.Data == "thead") {
			if child.Data == "thead" {
				isFirstRowHeader = true
				break
			}

			// If it's a tr, check if it contains th elements
			if child.Data == "tr" {
				firstRow = child
				for cell := child.FirstChild; cell != nil; cell = cell.NextSibling {
					if cell.Type == html.ElementNode && cell.Data == "th" {
						isFirstRowHeader = true
						break
					}
				}
				break
			}
		}
	}

	// If no header row is found and we need to infer one, we need to add an empty header row
	if !isFirstRowHeader && firstRow != nil && c.options.TableInferHeader {
		// We'll handle this in the convertTr function
	}

	// For tables, we need to ensure they have proper spacing
	return "\n\n" + text + "\n\n"
}

// convertTd converts <td> tags to Markdown table cells
func (c *Converter) convertTd(n *html.Node, text string, parentTags []string) string {
	colspan := 1
	colspanAttr := getAttr(n, "colspan")
	if colspanAttr != "" {
		colspanVal, err := strconv.Atoi(colspanAttr)
		if err == nil && colspanVal > 0 {
			colspan = colspanVal
		}
	}

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")

	if colspan > 1 {
		return " " + text + " |" + strings.Repeat(" |", colspan-1)
	}
	return " " + text + " |"
}

// convertTh converts <th> tags to Markdown table headers
func (c *Converter) convertTh(n *html.Node, text string, parentTags []string) string {
	// Same implementation as convertTd
	return c.convertTd(n, text, parentTags)
}

// convertTr converts <tr> tags to Markdown table rows
func (c *Converter) convertTr(n *html.Node, text string, parentTags []string) string {
	// Count cells and check if they're all th elements
	var cells []*html.Node
	isHeadRow := true
	isFirstRow := true

	// Collect cells and check if this is a header row
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && (child.Data == "td" || child.Data == "th") {
			cells = append(cells, child)
			if child.Data != "th" {
				isHeadRow = false
			}
		}
	}

	// Check if this is the first row in the table
	for sibling := n.PrevSibling; sibling != nil; sibling = sibling.PrevSibling {
		if sibling.Type == html.ElementNode && sibling.Data == "tr" {
			isFirstRow = false
			break
		}
	}

	// Check if we're in a thead
	inThead := false
	for p := n.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && p.Data == "thead" {
			inThead = true
			break
		}
	}

	// Determine if this is a header row
	isHeadRow = isHeadRow || inThead

	// Check if we need to infer a header
	isHeadRowMissing := isFirstRow && !isHeadRow && c.options.TableInferHeader

	// Calculate total colspan
	totalColspan := 0
	for _, cell := range cells {
		colspan := 1
		colspanAttr := getAttr(cell, "colspan")
		if colspanAttr != "" {
			colspanVal, err := strconv.Atoi(colspanAttr)
			if err == nil && colspanVal > 0 {
				colspan = colspanVal
			}
		}
		totalColspan += colspan
	}

	var result strings.Builder

	// Add the row content
	result.WriteString("|")
	result.WriteString(text)
	result.WriteString("\n")

	// If this is a header row or we need to infer a header, add the separator
	if (isHeadRow || isHeadRowMissing) && isFirstRow {
		result.WriteString("| ")
		for i := 0; i < totalColspan; i++ {
			if i > 0 {
				result.WriteString(" | ")
			}
			result.WriteString("---")
		}
		result.WriteString(" |\n")
	}

	return result.String()
}
