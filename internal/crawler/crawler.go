package crawler

import (
	"io"
	"linkchecker/config"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var Visited = make(map[string]bool)

var HttpClient *http.Client

func init() {
	HttpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   10,
			MaxConnsPerHost:       15,
		},
		Timeout: 60 * time.Second,
	}
}

// crawl обходит страницу u, извлекает внутренние ссылки и рекурсивно обходит их.

func Crawl(u, root *url.URL, depth, maxDepth int, Sum *config.Summary) {
	if depth > maxDepth {
		return
	}
	if Visited[u.String()] {
		return
	}
	Visited[u.String()] = true

	resp, err := HttpClient.Get(u.String())
	Sum.ErrorByType[resp.StatusCode]++

	if resp.StatusCode > 299 {
		Sum.ProblemLinks[u.String()] = config.CheckResult{
			StatusCode:   resp.StatusCode,
			Error:        err,
			Workers:      1,
			Depth:        depth,
			Referrer:     root,
			ResponseTime: time.Duration(resp.ContentLength),
		}
	}
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
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
		Crawl(link, root, depth+1, maxDepth, Sum)
	}
	Sum.CheckedLinks = len(Visited)
	return
}
