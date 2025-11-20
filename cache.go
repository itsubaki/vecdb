package vecdb

import "sync"

type Cache[T any] struct {
	m      map[string][]Result[T]
	ignore map[string]map[DocID]struct{}
	sync.RWMutex
}

func (c *Cache[T]) Ignore(query string, id DocID) {
	if c.ignore[query] == nil {
		c.ignore[query] = make(map[DocID]struct{})
	}

	c.ignore[query][id] = struct{}{}
}

func (c *Cache[T]) Put(query string, results []Result[T]) {
	c.Lock()
	defer c.Unlock()

	if c.m == nil {
		c.m = make(map[string][]Result[T])
	}

	c.m[query] = results
}

func (c *Cache[T]) Get(query string) ([]Result[T], bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.m[query]
	if !ok {
		return nil, false
	}

	ig, ok := c.ignore[query]
	if !ok {
		return v, true
	}

	var results []Result[T]
	for _, r := range v {
		if _, ok := ig[r.Doc.ID]; ok {
			continue
		}

		results = append(results, r)
	}

	return results, true
}
