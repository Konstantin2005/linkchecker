package crawler

import (
	"io"
	"linkchecker/config"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	StatusNetError  = -1 // любая сетевая / TLS / DNS ошибка
	StatusHTMLError = -2 // не HTML,
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
	if depth > maxDepth || Visited[u.String()] {
		return
	}

	Visited[u.String()] = true // помечаем проверенные сылки
	Sum.TotalLinks++

	start := time.Now()
	resp, err := HttpClient.Get(u.String())
	elapsed := time.Since(start)

	if err != nil { // сетевые, TLS, DNS, time-out ошибки
		AddProblem(Sum, u.String(), depth, root, StatusNetError, err, elapsed)
		return
	}
	status := resp.StatusCode

	if status == 200 {
		Sum.Successful++
	}
	if status >= 300 {
		AddProblem(Sum, u.String(), depth, root, resp.StatusCode, nil, elapsed)
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
		if link.Host == root.Host {
			Crawl(link, root, depth+1, maxDepth, Sum)
		}
	}
	Sum.CheckedLinks = len(Visited)
	return
}

func AddProblem(s *config.Summary, str string, depth int, ref *url.URL,
	code int, err error, dur time.Duration) {
	s.ErrorByType[code]++

	u, _ := url.Parse(str)
	u.Path = path.Dir(u.Path)
	u.RawQuery, u.Fragment = "", ""

	s.ProblemLinks[str] = config.CheckResult{
		URL:          str,
		StatusCode:   code,
		Error:        err,
		Depth:        depth,
		Referrer:     u,
		ResponseTime: dur,
	}
}
