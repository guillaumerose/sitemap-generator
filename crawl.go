package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

func crawl(base string, maxURLs int) ([]string, error) {
	scraper := scraper{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	visited := make(map[string]bool)
	queue := []string{"/"}
	for len(queue) > 0 && len(visited) < maxURLs {
		current := queue[0]
		queue = queue[1:]

		if _, ok := visited[current]; ok {
			continue
		}

		url := base + current
		logrus.Infof("Visiting %s", url)
		links, err := scraper.scrape(url)
		if err != nil {
			visited[current] = false
			logrus.Warnf("Error while scraping %s: %v", url, err)
			continue
		}
		visited[current] = true

		for _, link := range links {
			if _, ok := visited[link]; !ok {
				queue = append(queue, link)
			}
		}
	}

	var links []string
	for link, valid := range visited {
		if valid {
			links = append(links, link)
		}
	}
	sort.Strings(links)
	return links, nil
}
