package crawler

import (
	"sync"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/sirupsen/logrus"
)

type inMemoryQueue struct {
	limit int
	list  *doublylinkedlist.List
	lock  *sync.RWMutex
}

func newQueue(limit int) *inMemoryQueue {
	return &inMemoryQueue{
		limit: limit,
		list:  doublylinkedlist.New(),
		lock:  &sync.RWMutex{},
	}
}

func (q *inMemoryQueue) enqueue(r *Request) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.limit > 0 && q.list.Size() > q.limit {
		logrus.Errorf("Full queue, dropping %s", r.URL)
		return
	}
	q.list.Append(r)
}

func (q *inMemoryQueue) pop() *Request {
	q.lock.Lock()
	defer q.lock.Unlock()
	req, ok := q.list.Get(0)
	if !ok {
		return nil
	}
	q.list.Remove(0)
	return req.(*Request)
}

func (q *inMemoryQueue) size() int {
	return q.list.Size()
}
