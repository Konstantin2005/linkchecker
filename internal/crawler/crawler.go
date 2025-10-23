package crawler

import (
	"fmt"
	"io"
	"linkchecker/config"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var visited = make(map[string]bool)

// crawl обходит страницу u, извлекает внутренние ссылки и рекурсивно обходит их.

func Crawl(u, root *url.URL, depth, maxDepth int, Sum config.Summary) config.Summary {
	if depth > maxDepth {
		return Sum
	}
	if visited[u.String()] {
		return Sum
	}
	visited[u.String()] = true

	fmt.Printf("%s (depth %d)\n", u, depth)

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("GET %s: %v", u, err)
		return Sum
	}
	fmt.Println(resp.Status)
	Sum.TotalLinks++
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// интересуют только HTML
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		return Sum
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("parse %s: %v", u, err)
		return Sum
	}

	for _, link := range extractLinks(doc, u) { // u как base для относительных
		// внутренние = тот же host

		Crawl(link, root, depth+1, maxDepth, Sum)
	}
	return Sum
}
