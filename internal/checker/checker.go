package checker

import (
	"fmt"
	"linkchecker/config"
	"linkchecker/internal/crawler"
	"net/url"
)

func Check(u, root *url.URL, depth, maxDepth int, conf config.Config) {

	Summa := config.Summary{}

	result := crawler.Crawl(u, root, depth, maxDepth, Summa)
	fmt.Println(result.TotalLinks)

}
