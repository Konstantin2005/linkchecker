package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CheckResult struct {
	URL          string
	StatusCode   int
	Error        string
	Depth        int
	Referrer     string
	ResponseTime time.Duration
}

type Summary struct {
	TotalLinks   int
	CheckedLinks int
	Successful   int
	Errors       int
	ErrorByType  map[int]int
	Duration     time.Duration
}

var visited = make(map[string]bool)

// crawl обходит страницу u, извлекает внутренние ссылки и рекурсивно обходит их.
func Crawl(u, root *url.URL, depth, maxDepth int) {
	if depth > maxDepth {
		return
	}
	if visited[u.String()] {
		return
	}
	visited[u.String()] = true

	fmt.Printf("%s (depth %d)\n", u, depth)

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("GET %s: %v", u, err)
		return
	}
	defer resp.Body.Close()

	// интересуют только HTML
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("parse %s: %v", u, err)
		return
	}

	for _, link := range extractLinks(doc, u) { // u как base для относительных
		// внутренние = тот же host
		if link.Host != root.Host {
			continue
		}
		Crawl(link, root, depth+1, maxDepth)
	}
}
