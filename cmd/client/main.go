package main

import (
	"flag"
	"os"
	"time"

	"github.com/guillaumerose/sitemap-generator/pkg/client"
	"github.com/guillaumerose/sitemap-generator/pkg/render"
	"github.com/guillaumerose/sitemap-generator/pkg/types"
	"github.com/sirupsen/logrus"
)

var (
	maxDepth    int
	parallelism int
	crawlerURL  string
)

func main() {
	flag.StringVar(&crawlerURL, "s", "http://127.0.0.1:8080", "crawler URL")
	flag.IntVar(&maxDepth, "d", 5, "maximum depth to crawl")
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
	client := client.New(crawlerURL)
	crawl, err := client.CreateCrawl(&types.Crawl{
		Spec: types.CrawlSpec{
			URL:         url,
			MaxDepth:    maxDepth,
			Parallelism: parallelism,
		},
	})
	if err != nil {
		return err
	}
	logrus.Infof("Crawling %s (parallelism: %d, maxDepth: %d)", url, parallelism, maxDepth)
	for {
		crawl, err = client.GetCrawl(crawl.ID)
		if err != nil {
			return err
		}
		if crawl.Status.Done {
			break
		}
		logrus.Infof("%d URLs found", crawl.Status.Size)
		time.Sleep(time.Second)
	}
	logrus.Infof("Finished! %d URLs found", crawl.Status.Size)
	links, err := client.GetCrawlLinks(crawl.ID)
	if err != nil {
		return err
	}
	render.AsTree(links, os.Stdout)
	return nil
}
