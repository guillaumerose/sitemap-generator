package crawler

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/sirupsen/logrus"
)

type Crawler struct {
	visited *treemap.Map
	lock    sync.RWMutex

	waitChan chan bool
	wg       sync.WaitGroup
	count    uint64

	scraper *scraper
}

func New(parallelism int) *Crawler {
	return &Crawler{
		visited:  treemap.NewWithStringComparator(),
		lock:     sync.RWMutex{},
		waitChan: make(chan bool, parallelism),
		wg:       sync.WaitGroup{},
		scraper: &scraper{
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
	}
}

func (c *Crawler) Size() int {
	return c.visited.Size()
}

func (c *Crawler) Done() bool {
	count := atomic.LoadUint64(&c.count)
	return count == 0
}

func (c *Crawler) Wait() {
	c.wg.Wait()
}

func (c *Crawler) VisitedURLs() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var ans []string
	c.visited.Each(func(link interface{}, valid interface{}) {
		if valid.(bool) {
			ans = append(ans, link.(string))
		}
	})
	return ans
}

func (c *Crawler) markVisited(link string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.visited.Put(link, true)
}

func (c *Crawler) checkVisited(link string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.visited.Get(link)
	return ok
}

func (c *Crawler) checkVisitedAndMark(link string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.visited.Get(link)
	if !ok {
		c.visited.Put(link, false)
	}
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

	if ok := c.checkVisitedAndMark(current); ok {
		return
	}

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
		if ok := c.checkVisited(link); !ok {
			c.wg.Add(1)
			atomic.AddUint64(&c.count, 1)
			go func() {
				defer atomic.AddUint64(&c.count, ^uint64(0))
				defer c.wg.Done()
				c.doCrawl(base, link, depth-1)
			}()
		}
	}
}
