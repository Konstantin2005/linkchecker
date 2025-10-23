package crawler

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// extractLinks проходит весь DOM и собирает ссылки из a/img/link/script.

func extractLinks(n *html.Node, base *url.URL) []*url.URL {
	var out []*url.URL
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		if node.Type == html.ElementNode {
			var attrKey string
			switch node.DataAtom.String() {
			case "a":
				attrKey = "href"
			case "img", "script":
				attrKey = "src"
			case "link":
				attrKey = "href"
			}
			if attrKey != "" {
				for _, a := range node.Attr {
					if a.Key == attrKey && a.Val != "" {
						if link, ok := normalizeURL(a.Val, base); ok {
							out = append(out, link)
						}
					}
				}
			}

		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(n)
	return dedup(out)
}

// normalizeURL приводит строку к абсолютному URL относительно base.
func normalizeURL(raw string, base *url.URL) (*url.URL, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.HasPrefix(trimmed, "mailto:") || strings.HasPrefix(trimmed, "javascript:") {
		return nil, false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, false
	}
	resolved := base.ResolveReference(parsed)
	// уберём фрагмент, чтобы одна и та же страница с #id не считалась разной
	resolved.Fragment = ""
	return resolved, true
}

// dedup убирает дубликаты, сохраняя порядок.
func dedup(in []*url.URL) []*url.URL {
	seen := make(map[string]struct{}, len(in))
	var out []*url.URL
	for _, u := range in {
		s := u.String()
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, u)
		}
	}
	return out
}
