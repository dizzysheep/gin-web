// Package markdown renders the canonical HTML representation stored for an article.
package markdown

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
)

var renderer = goldmark.New()

// Render converts Markdown to HTML. Goldmark's safe defaults omit raw HTML and
// dangerous link destinations, so content_html never trusts client-provided HTML.
func Render(source string) (string, error) {
	var output bytes.Buffer
	if err := renderer.Convert([]byte(source), &output); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	return output.String(), nil
}
