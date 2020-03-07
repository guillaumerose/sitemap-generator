package main

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"golang.org/x/net/html"
)

type scraper struct {
	client *http.Client
}

func (s *scraper) scrape(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}
	if !strings.HasPrefix(res.Header.Get("content-type"), "text/html") {
		return []string{}, nil
	}

	defer res.Body.Close()
	return links(res.Body)
}

func links(body io.ReadCloser) ([]string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			if link := href(n); link != "" {
				if strings.HasPrefix(link, "/") && !strings.HasPrefix(link, "//") {
					withoutDash := strings.SplitN(link, "#", 2)[0]
					withoutQueryParams := strings.SplitN(withoutDash, "?", 2)[0]
					withoutExtraSlash := path.Clean(withoutQueryParams)
					links = append(links, withoutExtraSlash)
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
