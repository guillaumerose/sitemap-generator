package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		logrus.Fatal("url is mandatory")
	}
	if err := run(flag.Arg(0)); err != nil {
		logrus.Fatal(err)
	}
}

func run(base string) error {
	links, err := scrape(base)
	if err != nil {
		return err
	}
	for _, link := range links {
		fmt.Println(link)
	}
	return nil
}
