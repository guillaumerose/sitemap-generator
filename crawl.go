package main

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type crawler struct {
	visited map[string]bool
	lock    sync.Mutex

	waitChan chan bool
	wg       sync.WaitGroup

	scraper *scraper
}

func newCrawler(parallelism int) *crawler {
	return &crawler{
		visited:  make(map[string]bool),
		lock:     sync.Mutex{},
		waitChan: make(chan bool, parallelism),
		wg:       sync.WaitGroup{},
		scraper: &scraper{
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
	}
}
func (c *crawler) visitedURLs() []string {
	c.wg.Wait()
	var ans []string
	for link, valid := range c.visited {
		if valid {
			ans = append(ans, link)
		}
	}
	sort.Strings(ans)
	return ans
}

func (c *crawler) markVisited(link string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.visited[link] = true
}

func (c *crawler) markError(link string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.visited[link] = false
}

func (c *crawler) isVisited(link string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.visited[link]
	return ok
}

func (c *crawler) crawl(base, current string, depth int) {
	c.waitChan <- true
	defer func() {
		<-c.waitChan
	}()

	if depth <= 0 {
		return
	}

	if ok := c.isVisited(current); ok {
		return
	}
	c.markError(current)

	url := base + current
	logrus.Infof("Visiting %s", url)
	links, err := c.scraper.scrape(url)
	if err != nil {
		logrus.Warnf("Error while scraping %s: %v", url, err)
		return
	}
	c.markVisited(current)

	for i := range links {
		link := links[i]
		if ok := c.isVisited(link); !ok {
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.crawl(base, link, depth-1)
			}()
		}
	}
}
