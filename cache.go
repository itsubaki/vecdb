package vecdb

import "sync"

type Query string

type Cache[T any] struct {
	m map[Query]map[DocID]Result[T]
	sync.RWMutex
}

func (c *Cache[T]) Put(query string, results []Result[T]) {
	c.Lock()
	defer c.Unlock()

	if c.m == nil {
		c.m = make(map[Query]map[DocID]Result[T])
	}

	q := Query(query)
	if _, ok := c.m[q]; !ok {
		c.m[q] = make(map[DocID]Result[T])
	}

	for _, r := range results {
		c.m[q][r.Doc.ID] = r
	}
}

func (c *Cache[T]) Get(query string) (map[DocID]Result[T], bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.m[Query(query)]
	return v, ok
}
