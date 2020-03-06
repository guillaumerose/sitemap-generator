package main

import (
	"sort"

	"github.com/sirupsen/logrus"
)

func crawl(base string, maxURLs int) ([]string, error) {
	visited := make(map[string]struct{})
	queue := []string{"/"}
	for len(queue) > 0 && len(visited) < maxURLs {
		current := queue[0]
		queue = queue[1:]

		if _, ok := visited[current]; ok {
			continue
		}
		visited[current] = struct{}{}

		url := base + current
		logrus.Infof("Visiting %s", url)
		links, err := scrape(url)
		if err != nil {
			logrus.Warnf("Error while scraping %s: %v", url, err)
			continue
		}

		for _, link := range links {
			if _, ok := visited[link]; !ok {
				queue = append(queue, link)
			}
		}
	}

	var links []string
	for link := range visited {
		links = append(links, link)
	}
	sort.Strings(links)
	return links, nil
}
