package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var linksExtractor = regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)

func scrape(base string) ([]string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", base, nil)
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	matches := linksExtractor.FindAllStringSubmatch(string(body), -1)
	var links []string
	for _, match := range matches {
		link := match[1]
		if strings.HasPrefix(link, "/") {
			links = append(links, link)
		}
	}
	return links, nil
}
