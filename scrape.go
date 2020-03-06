package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func scrape(url string) ([]string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}

	defer res.Body.Close()
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			if link := href(n); link != "" {
				if strings.HasPrefix(link, "/") {
					links = append(links, link)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func href(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

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
