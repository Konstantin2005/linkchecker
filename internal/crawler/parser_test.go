package crawler // ← замените на фактическое имя пакета

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

/*** helpers ***/

func mustURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("bad test url %q: %v", raw, err)
	}
	return u
}

func urlsToStrings(in []*url.URL) []string {
	out := make([]string, len(in))
	for i, u := range in {
		out[i] = u.String()
	}
	return out
}

func parseHTML(t *testing.T, s string) *html.Node {
	t.Helper()
	n, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("html.Parse: %v", err)
	}
	return n
}

/*** extractLinks ***/

func TestExtractLinks(t *testing.T) {
	base := mustURL(t, "https://example.com/dir/")

	tests := []struct {
		name string
		html string
		want []string
	}{
		{
			name: "all supported tags + resolution",
			html: `
				<a href="page.html">a</a>
				<img src="/img.png">
				<script src="js/app.js"></script>
				<link href="style.css">
			`,
			want: []string{
				"https://example.com/dir/page.html",
				"https://example.com/img.png",
				"https://example.com/dir/js/app.js",
				"https://example.com/dir/style.css",
			},
		},
		{
			name: "deduplicate keeps first appearance",
			html: `
				<a href="page.html"></a>
				<a href="page.html#section"></a>
				<a href="page.html "></a>
			`,
			want: []string{
				"https://example.com/dir/page.html",
			},
		},
		{
			name: "absolute url is untouched, fragment removed",
			html: `<a href="https://other.com/x#frag"></a>`,
			want: []string{"https://other.com/x"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := parseHTML(t, tc.html)
			got := extractLinks(root, base)
			if !reflect.DeepEqual(urlsToStrings(got), tc.want) {
				t.Fatalf("want %v, got %v", tc.want, urlsToStrings(got))
			}
		})
	}
}

/*** normalizeURL ***/

func TestNormalizeURL(t *testing.T) {
	base := mustURL(t, "https://example.com/dir/")

	tests := []struct {
		raw  string
		ok   bool
		want string
	}{
		{"  page.html ", true, "https://example.com/dir/page.html"},
		{"/x#frag", true, "https://example.com/x"},
		{"https://golang.org", true, "https://golang.org"},
		{"mailto:a@b.c", false, ""},
		{"javascript:alert(1)", false, ""},
		{"::bad::url", false, ""},
		{"   ", false, ""},
	}

	for _, tc := range tests {
		got, ok := normalizeURL(tc.raw, base)
		if ok != tc.ok {
			t.Fatalf("%q: expected ok=%v, got %v", tc.raw, tc.ok, ok)
		}
		if ok && got.String() != tc.want {
			t.Fatalf("%q: expected %q, got %q", tc.raw, tc.want, got)
		}
	}
}

/*** dedup ***/

func TestDedup(t *testing.T) {
	a := mustURL(t, "https://x/a")
	b := mustURL(t, "https://x/b")
	c := mustURL(t, "https://x/a") // дубликат `a`

	in := []*url.URL{a, b, c, a}
	want := []*url.URL{a, b}

	got := dedup(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", urlsToStrings(want), urlsToStrings(got))
	}
}
