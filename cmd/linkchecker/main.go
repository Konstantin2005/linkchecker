package main

import (
	"flag"
	"linkchecker/config"
	"linkchecker/internal/checker"
	"log"
	"net/url"
	"time"
)

// visited хранит уже просмотренные URL, чтобы не обходить их повторно.

func main() {
	URL := flag.String("URL", "https://practicum.yandex.ru/profile/go-developer-basic/", "путь к .md файлу (обязательно)")
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
	checker.Check(root, root, 0, *depth, conf)

}
