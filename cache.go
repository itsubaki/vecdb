package vecdb

import "sync"

type Query string

type Cache[T any] struct {
	m      map[Query][]Result[T]
	ignore map[Query]map[Text]struct{}
	sync.RWMutex
}

func (c *Cache[T]) Ignore(query string, doc Doc[T]) {
	c.Lock()
	defer c.Unlock()

	if c.ignore == nil {
		c.ignore = make(map[Query]map[Text]struct{})
	}

	q := Query(query)
	if c.ignore[q] == nil {
		c.ignore[q] = make(map[Text]struct{})
	}

	c.ignore[q][doc.Text] = struct{}{}
}

func (c *Cache[T]) Put(query string, results []Result[T]) {
	c.Lock()
	defer c.Unlock()

	if c.m == nil {
		c.m = make(map[Query][]Result[T])
	}

	c.m[Query(query)] = results
}

func (c *Cache[T]) Get(query string) ([]Result[T], bool) {
	c.RLock()
	defer c.RUnlock()

	q := Query(query)
	v, ok := c.m[q]
	if !ok {
		return nil, false
	}

	ig, ok := c.ignore[q]
	if !ok {
		return v, true
	}

	var results []Result[T]
	for _, r := range v {
		if _, ok := ig[r.Doc.Text]; ok {
			continue
		}

		results = append(results, r)
	}

	return results, true
}
