package crawler

import (
	"linkchecker/config"
	"linkchecker/internal/output"
	"net/url"
)

func NewSummary() *config.Summary {
	return &config.Summary{
		ErrorByType: make(map[int]int),
	}
}

func Check(u, root *url.URL, depth, maxDepth int, conf config.Config) {

	Sum := NewSummary()

	Crawl(u, root, depth, maxDepth, Sum)

	f := output.NewFormatter()

	f.PrintResult(conf)
	f.PrintSummary(Sum)
}
