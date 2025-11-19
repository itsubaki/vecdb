package vecdb

import "sync"

type Cache[T any] struct {
	m map[string][]Result[T]
	sync.RWMutex
}

func (c *Cache[T]) Put(id string, results []Result[T]) {
	c.Lock()
	defer c.Unlock()

	if c.m == nil {
		c.m = make(map[string][]Result[T])
	}

	c.m[id] = results
}

func (c *Cache[T]) Get(id string) ([]Result[T], bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.m[id]
	if !ok {
		return nil, false
	}

	return v, true
}
