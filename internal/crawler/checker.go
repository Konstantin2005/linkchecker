package crawler

import (
	"fmt"
	"linkchecker/config"
	"net/url"
)

func Check(u, root *url.URL, depth, maxDepth int, conf config.Config) {

	Crawl(u, root, depth, maxDepth)
	fmt.Println(Sum.TotalLinks)

}
