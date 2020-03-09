package main

import (
	"flag"
	"os"

	"github.com/guillaumerose/sitemap-generator/pkg/crawler"
	"github.com/guillaumerose/sitemap-generator/pkg/render"
	"github.com/guillaumerose/sitemap-generator/pkg/types"
	"github.com/sirupsen/logrus"
)

var (
	maxDepth    int
	parallelism int
)

func main() {
	flag.IntVar(&maxDepth, "d", 5, "maximum depth to crawl (-1 is unlimited depth)")
	flag.IntVar(&parallelism, "p", 2, "maximum number of concurrent requests")
	flag.Parse()
	if flag.NArg() != 1 {
		logrus.Fatal("url is mandatory")
	}
	if err := run(flag.Arg(0), maxDepth, parallelism); err != nil {
		logrus.Fatal(err)
	}
}

func run(url string, maxDepth, parallelism int) error {
	crawler := crawler.New(types.CrawlSpec{
		URL:         url,
		MaxDepth:    maxDepth,
		Parallelism: parallelism,
	})
	crawler.Crawl()
	crawler.Wait()
	links := crawler.VisitedURLs()
	logrus.Infof("Found %d URLs", len(links))
	render.AsTree(links, os.Stdout)
	return nil
}
