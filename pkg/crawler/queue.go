package crawler

import (
	"sync"

	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/sirupsen/logrus"
)

type inMemoryQueue struct {
	limit int
	list  *linkedhashset.Set
	lock  *sync.Mutex
}

func newQueue(limit int) *inMemoryQueue {
	return &inMemoryQueue{
		limit: limit,
		list:  linkedhashset.New(),
		lock:  &sync.Mutex{},
	}
}

func (q *inMemoryQueue) enqueue(r Request) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.limit > 0 && q.list.Size() > q.limit {
		logrus.Errorf("Full queue, dropping %s", r.URL)
		return
	}
	q.list.Add(r)
}

func (q *inMemoryQueue) pop() Request {
	q.lock.Lock()
	defer q.lock.Unlock()
	it := q.list.Iterator()
	if !it.Next() {
		return Request{}
	}
	req := it.Value()
	q.list.Remove(req)
	return req.(Request)
}

func (q *inMemoryQueue) size() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.list.Size()
}
