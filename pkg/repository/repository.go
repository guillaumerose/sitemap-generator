package repository

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/guillaumerose/sitemap-generator/pkg/crawler"
	"github.com/guillaumerose/sitemap-generator/pkg/types"
)

type Repository interface {
	Create(req *types.Crawl) (*types.Crawl, error)
	Get(id string) (*types.Crawl, error)
	GetLinks(id string) ([]string, error)
}

type InMemoryRepository struct {
	lock   sync.RWMutex
	crawls map[string]*crawler.Crawler
	seq    int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		lock:   sync.RWMutex{},
		crawls: make(map[string]*crawler.Crawler),
		seq:    1,
	}
}

func (r *InMemoryRepository) Create(req *types.Crawl) (*types.Crawl, error) {
	if req.Spec.URL == "" {
		return nil, errors.New("spec.url is mandatory")
	}
	if req.Spec.MaxDepth == 0 {
		req.Spec.MaxDepth = 3
	}
	if req.Spec.Parallelism == 0 {
		req.Spec.Parallelism = 2
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	crawler := crawler.New(req.Spec)
	crawler.Crawl()

	id := strconv.Itoa(r.seq)
	r.seq++
	r.crawls[id] = crawler

	return &types.Crawl{
		ID:   id,
		Spec: req.Spec,
	}, nil
}

func (r *InMemoryRepository) Get(id string) (*types.Crawl, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if crawl, ok := r.crawls[id]; ok {
		return &types.Crawl{
			ID:   id,
			Spec: crawl.Spec,
			Status: types.CrawlStatus{
				Done: crawl.Done(),
				Size: crawl.Size(),
			},
		}, nil
	}
	return nil, fmt.Errorf("cannot find crawl with id %s", id)
}

func (r *InMemoryRepository) GetLinks(id string) ([]string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if crawl, ok := r.crawls[id]; ok {
		return crawl.VisitedURLs(), nil
	}
	return nil, fmt.Errorf("cannot find crawl with id %s", id)
}
