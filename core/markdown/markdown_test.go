package markdown

import (
	"strings"
	"testing"
)

func TestRenderDisablesRawHTML(t *testing.T) {
	html, err := Render("# title\n\n<script>alert(1)</script>\n\n[bad](javascript:alert(1))")
	if err != nil {
		t.Fatalf("Render err: %v", err)
	}
	if !strings.Contains(html, "<h1>title</h1>") {
		t.Fatalf("heading was not rendered: %q", html)
	}
	if strings.Contains(html, "<script") || strings.Contains(html, "href=\"javascript:") {
		t.Fatalf("unsafe HTML rendered: %q", html)
	}
}
