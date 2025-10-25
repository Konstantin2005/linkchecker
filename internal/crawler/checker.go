package crawler

import (
	"fmt"
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

	fmt.Println(Sum.TotalLinks)
	fmt.Println(Sum.ErrorByType)
	f := output.NewFormatter()
	f.PrintSummary(Sum)
}
