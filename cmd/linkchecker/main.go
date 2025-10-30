package main

import (
	"flag"
	"linkchecker/config"
	"linkchecker/internal/crawler"
	"linkchecker/internal/output"
	"linkchecker/pkg"
	"log"
	"net/url"
	"time"
)

func NewSummary() *config.Summary {
	return &config.Summary{
		ErrorByType: make(map[int]int),
	}
}

// visited хранит уже просмотренные URL, чтобы не обходить их повторно.

func main() {
	done := make(chan bool)
	start := time.Now()

	URL := flag.String("URL", "https://github.com/Konstantin2005/linkchecker/blob/main/go.mod", "Сслыка для обхода")
	depth := flag.Int("depth", 2, "Глубина обхода")
	timeout := flag.Int("timeout", 10, "Таймаут запроса в секундах")
	workers := flag.Int("workers", 1, "Количесвто Горутин")
	OutputFormat := flag.String("output", "json", "Как выводить")

	help := flag.Bool("help", false, "показать справку")
	flag.Parse()

	if *help || *URL == "" {
		flag.Usage()
		return
	}

	conf := config.Config{
		URL:          *URL,
		MaxDepth:     *depth,
		Timeout:      time.Duration(*timeout),
		Workers:      *workers,
		OutputFormat: *OutputFormat,
	}

	go pkg.Loading(done)
	root, err := url.Parse(conf.URL)
	if err != nil {
		log.Fatalf("bad start URL: %v", err)
	}

	Sum := NewSummary()

	crawler.Crawl(root, root, 1, *depth, Sum)
	done <- true
	Sum.Duration = time.Since(start)

	output.MainFormate(conf, Sum)

}
