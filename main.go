package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

var maxURLs int

func main() {
	flag.IntVar(&maxURLs, "max", 20, "maximum number of URLs in sitemap")
	flag.Parse()
	if flag.NArg() != 1 {
		logrus.Fatal("url is mandatory")
	}
	if err := run(flag.Arg(0), maxURLs); err != nil {
		logrus.Fatal(err)
	}
}

func run(base string, maxURLs int) error {
	links, err := crawl(base, maxURLs)
	if err != nil {
		return err
	}
	for _, link := range links {
		fmt.Println(link)
	}
	return nil
}
