package main

import (
	"flag"
	"os"

	"github.com/guillaumerose/sitemap-generator/pkg/crawler"
	"github.com/sirupsen/logrus"
)

var (
	maxDepth    int
	parallelism int
)

func main() {
	flag.IntVar(&maxDepth, "max", 5, "maximum depth to crawl")
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
	crawler := crawler.New(parallelism)
	crawler.Crawl(url, maxDepth)
	links := crawler.VisitedURLs()
	render(links, os.Stdout)
	return nil
}
