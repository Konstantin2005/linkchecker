package main

import (
	"flag"
	"linkchecker/config"
	"linkchecker/internal/crawler"
	"log"
	"net/url"
	"time"
)

// visited хранит уже просмотренные URL, чтобы не обходить их повторно.

func main() {
	URL := flag.String("URL", "https://www.google.com/", "путь к .md файлу (обязательно)")
	depth := flag.Int("output", 10, "путь для сохранения .html(по умолчанию stdout)")
	flag.Parse()

	conf := config.Config{
		URL:      *URL,
		MaxDepth: 10,
		Timeout:  300 * time.Second,
		Workers:  5000,
	}

	root, err := url.Parse(conf.URL)
	if err != nil {
		log.Fatalf("bad start URL: %v", err)
	}
	crawler.Check(root, root, 1, *depth, conf)

}
