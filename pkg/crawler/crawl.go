package crawler

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Crawler struct {
	visited map[string]bool
	lock    sync.Mutex

	waitChan chan bool
	wg       sync.WaitGroup

	scraper *scraper
}

func New(parallelism int) *Crawler {
	return &Crawler{
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
func (c *Crawler) VisitedURLs() []string {
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

func (c *Crawler) markVisited(link string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.visited[link] = true
}

func (c *Crawler) markError(link string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.visited[link] = false
}

func (c *Crawler) isVisited(link string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.visited[link]
	return ok
}

func (c *Crawler) Crawl(base string, maxDepth int) {
	c.doCrawl(base, "/", maxDepth)
}

func (c *Crawler) doCrawl(base, current string, depth int) {
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
	links, err := c.scraper.scrapeAllLinks(url)
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
				c.doCrawl(base, link, depth-1)
			}()
		}
	}
}
