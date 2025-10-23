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

var Sum config.Summary

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
	Sum.TotalLinks++

	if err != nil {
		log.Printf("GET %s: %v", u, err)
		return
	}
	fmt.Println(resp.Status)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

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
	return
}
