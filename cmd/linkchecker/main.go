package main

import (
	"flag"
	"linkchecker/internal/crawler"
	"log"
	"net/url"
	"time"
)

type Config struct {
	URL           string
	MaxDepth      int
	Timeout       time.Duration
	Workers       int
	OutputFormat  string
	Verbose       bool
	SkipSSLVerify bool
}

// visited хранит уже просмотренные URL, чтобы не обходить их повторно.

func main() {
	URL := flag.String("URL", "https://leetcode.com/problemset/", "путь к .md файлу (обязательно)")
	depth := flag.Int("output", 5, "путь для сохранения .html(по умолчанию stdout)")
	flag.Parse()

	conf := Config{
		URL:      *URL,
		MaxDepth: 10,
		Timeout:  300 * time.Second,
		Workers:  5000,
	}

	root, err := url.Parse(conf.URL)
	if err != nil {
		log.Fatalf("bad start URL: %v", err)
	}
	crawler.Crawl(root, root, *depth, conf.MaxDepth)

}
