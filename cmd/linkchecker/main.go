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
	URL := flag.String("URL", "https://github.com/Konstantin2005/linkchecker/blob/main/go.mod", "путь к .md файлу (обязательно)")
	depth := flag.Int("output", 7, "путь для сохранения .html(по умолчанию stdout)")
	flag.Parse()

	conf := config.Config{
		URL:      *URL,
		MaxDepth: *depth,
		Timeout:  5 * time.Second,
		Workers:  1,
	}

	root, err := url.Parse(conf.URL)
	if err != nil {
		log.Fatalf("bad start URL: %v", err)
	}
	crawler.Check(root, root, 1, *depth, conf)

}
