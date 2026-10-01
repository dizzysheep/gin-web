package article

import "testing"

func TestBuildSlug(t *testing.T) {
	cases := map[string]string{
		"Hello, Go World!": "hello-go-world",
		"  custom_slug  ":  "custom-slug",
	}
	for input, want := range cases {
		if got := buildSlug(input, "ignored"); got != want {
			t.Fatalf("buildSlug(%q) = %q, want %q", input, got, want)
		}
	}
	if got := buildSlug("", "中文标题"); got == "" {
		t.Fatal("non-ASCII title should receive a deterministic fallback slug")
	}
}
