package cleaner

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	parser "golang.org/x/net/html"
)

var (
	ignoreTags        = []string{"title", "script", "style", "iframe", "frame", "frameset", "noframes", "noembed", "embed", "applet", "object", "base"}
	allowedTags       = []string{"html", "body", "h1", "h2", "h3", "h4", "h5", "h6", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "a", "img", "pre", "code", "blockquote", "table", "tr", "th", "td", "tbody", "thead", "caption"}
	allowedAttributes = []string{"class", "src", "href", "title", "alt", "name"}
)

// CleanHTML cleans the HTML :P
func CleanHTML(f io.Reader) (string, error) {
	tokenizer := parser.NewTokenizer(f)

	buffer := bytes.NewBufferString("")
	ignore := ""

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		switch tokenType {

		case parser.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return cleanEmptyTags(buffer.String()), nil
			}
			return "", err

		case parser.StartTagToken:
			if len(ignore) == 0 && includes(allowedTags, token.Data) {
				token.Attr = cleanAttributes(token.Attr, allowedAttributes)
				buffer.WriteString(token.String())
			} else if includes(ignoreTags, token.Data) {
				ignore = token.Data
			}
		case parser.SelfClosingTagToken:

			if len(ignore) == 0 && includes(allowedTags, token.Data) {
				token.Attr = cleanAttributes(token.Attr, allowedAttributes)
				buffer.WriteString(token.String())
			} else if token.Data == ignore {
				ignore = ""
			}

		case parser.EndTagToken:
			if len(ignore) == 0 && includes(allowedTags, token.Data) {
				token.Attr = []parser.Attribute{}
				buffer.WriteString(token.String())
			} else if token.Data == ignore {
				ignore = ""
			}

		case parser.TextToken:
			if ignore == "" {
				buffer.WriteString(token.String())
			}

		default:
			// Skip
		}
	}

}

func cleanEmptyTags(str string) string {
	str = regexp.MustCompile(`[\xc2\xa0]+`).ReplaceAllLiteralString(str, " ")
	str = regexp.MustCompile(`[\n\s]+`).ReplaceAllLiteralString(str, " ")
	str = regexp.MustCompile(`(<br>\s*)+`).ReplaceAllLiteralString(str, "<br>")
	str = regexp.MustCompile(`<p>(\s*<br>\s*)+`).ReplaceAllLiteralString(str, "<p>")
	str = regexp.MustCompile(`(\s*<br>\s*)+</p>`).ReplaceAllLiteralString(str, "</p>")
	str = strings.Replace(str, "<b><br></b>", "", -1)

	// Remove space-only tags
	str = regexp.MustCompile(`<h.>\s*<\/h.>`).ReplaceAllLiteralString(str, " ")
	str = regexp.MustCompile(`<b>\s*<\/b>`).ReplaceAllLiteralString(str, " ")
	str = regexp.MustCompile(`<i>\s*<\/i>`).ReplaceAllLiteralString(str, " ")
	str = regexp.MustCompile(`<p>\s*<\/p>`).ReplaceAllLiteralString(str, "")

	// Remove paragraphs inside table cells
	str = regexp.MustCompile(`<td>\s*<p>(.*?)<\/p>\s*<\/td>`).ReplaceAllString(str, "<td>$1</td>")

	// Remove plain anchors
	str = regexp.MustCompile(`<a name="[^"]+">(.*?)</a>`).ReplaceAllString(str, "$1")

	str = regexp.MustCompile(`<p>\s*(Tidsskriftet Sakprosa)\s*<\/p>\s*<p>\s*(Bind \d+, Nummer \d+)\s*<\/p>\s*<p>\s*(© \d+)\s*<\/p>`).ReplaceAllString(str, "<p>$1<br>$2<br>$3</p>")

	str = regexp.MustCompile(`\s+`).ReplaceAllLiteralString(str, " ")
	str = strings.Replace(str, "<br><hr>", "<hr>", -1)
	str = strings.Replace(str, "<h2>", "<h3>", -1)
	str = strings.Replace(str, "</h2>", "</h3>", -1)
	return str
}

// cleanAttributes returns an array of attributes after removing malicious ones.
func cleanAttributes(a []parser.Attribute, allowed []string) []parser.Attribute {
	if len(a) == 0 {
		return a
	}

	var cleaned []parser.Attribute
	for _, attr := range a {
		if includes(allowed, attr.Key) {
			val := strings.ToLower(attr.Val)

			// Check for illegal attribute values
			if attr.Key == "class" && val != "abstract" {
				attr.Val = ""
			}

			if attr.Val != "" {
				cleaned = append(cleaned, attr)
			}
		}
	}
	return cleaned
}

// includes checks for inclusion of a string in a []string.
func includes(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}
