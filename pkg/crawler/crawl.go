package crawler

import (
	"net/http"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/guillaumerose/sitemap-generator/pkg/types"
	"github.com/sirupsen/logrus"
)

const initialWorkPerWorker = 5

type Request struct {
	URL   string
	Depth int
}

type Crawler struct {
	Spec types.CrawlSpec

	queue *inMemoryQueue

	wg sync.WaitGroup

	visited *treemap.Map
	lock    sync.Mutex

	scraper *scraper
}

func New(spec types.CrawlSpec) *Crawler {
	return &Crawler{
		queue:   newQueue(1_000_000),
		Spec:    spec,
		visited: treemap.NewWithStringComparator(),
		lock:    sync.Mutex{},
		scraper: &scraper{
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
		wg: sync.WaitGroup{},
	}
}

func (c *Crawler) Size() int {
	return c.queue.size()
}

func (c *Crawler) VisitedSize() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.visited.Size()
}

func (c *Crawler) Done() bool {
	s := c.queue.size()
	return s == 0
}

func (c *Crawler) Wait() {
	c.wg.Wait()
}

func (c *Crawler) VisitedURLs() []string {
	c.lock.Lock()
	defer c.lock.Unlock()
	var ans []string
	c.visited.Each(func(link interface{}, valid interface{}) {
		if valid.(bool) {
			ans = append(ans, link.(string))
		}
	})
	return ans
}

func (c *Crawler) process(threadId int, r Request) {
	if c.Spec.MaxDepth >= 0 && r.Depth > c.Spec.MaxDepth {
		return
	}

	if ok := c.checkVisitedAndMark(r.URL); ok {
		return
	}

	logrus.Infof("#%d Visiting %s", threadId, r.URL)
	links, err := c.scraper.scrapeAllLinks(c.Spec.URL + r.URL)
	if err != nil {
		logrus.Warnf("Error while scraping %s: %v", r.URL, err)
		return
	}
	c.markVisited(r.URL)

	for i := range links {
		link := links[i]
		if ok := c.checkVisited(link); !ok {
			c.queue.enqueue(Request{
				URL:   link,
				Depth: r.Depth + 1,
			})
		}
	}
}

func (c *Crawler) worker(threadId int) {
	for c.queue.size() > 0 {
		r := c.queue.pop()
		if r.URL == "" {
			break
		}
		c.process(threadId, r)
	}
	logrus.Debugf("#%d is dead", threadId)
}

func (c *Crawler) Crawl() {
	c.queue.enqueue(Request{
		URL:   "/",
		Depth: 1,
	})
	// Add enough work in the queue before starting workers
	for c.queue.size() < initialWorkPerWorker*c.Spec.Parallelism && c.queue.size() > 0 {
		r := c.queue.pop()
		if r.URL == "" {
			break
		}
		c.process(0, r)
	}
	// Starting workers
	for i := 0; i < c.Spec.Parallelism; i++ {
		c.wg.Add(1)
		go func(threadId int) {
			defer c.wg.Done()
			c.worker(threadId)
		}(i + 1)
	}
}

func (c *Crawler) markVisited(link string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.visited.Put(link, true)
}

func (c *Crawler) checkVisited(link string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
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
