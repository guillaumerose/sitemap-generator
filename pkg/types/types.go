package types

type Crawl struct {
	ID     string      `json:"id"`
	Spec   CrawlSpec   `json:"spec"`
	Status CrawlStatus `json:"status"`
}

type CrawlSpec struct {
	URL         string `json:"url"`
	MaxDepth    int    `json:"max_depth"`
	Parallelism int    `json:"parallelism"`
}

type CrawlStatus struct {
	Done bool `json:"done"`
	Size int  `json:"size"`
}
