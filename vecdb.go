package vecdb

import (
	"fmt"
	"sort"
	"sync"
)

type (
	DocID string
	EmbID string
)

type Doc[T any] struct {
	ID       DocID
	Label    string
	Text     string
	Metadata T
	Ignore   bool
}

type Embedding struct {
	ID     EmbID
	DocID  DocID
	Vector []float64
}

type Result[T any] struct {
	Score float64
	Doc   Doc[T]
}

type Memory[T any] struct {
	Distance   func(a, b []float64) float64
	Embeddings func(text []string) ([][]float64, error)
	docs       map[DocID]Doc[T]
	embeddings map[DocID]Embedding
	cache      Cache[T]
	sync.RWMutex
}

func (m *Memory[T]) Save(docs []Doc[T]) error {
	text := make([]string, len(docs))
	for i, d := range docs {
		text[i] = d.Text
	}

	v, err := m.Embeddings(text)
	if err != nil {
		return fmt.Errorf("embedding: %v", err)
	}


	// save with lock
	m.Lock()
	defer m.Unlock()

	if m.docs == nil {
		m.docs = make(map[DocID]Doc[T])
	}

	if m.embeddings == nil {
		m.embeddings = make(map[DocID]Embedding)
	}

	for i := range v {
		m.docs[docs[i].ID] = Doc[T]{
			ID:       docs[i].ID,
			Label:    docs[i].Label,
			Text:     docs[i].Text,
			Metadata: docs[i].Metadata,
		}

		m.embeddings[docs[i].ID] = Embedding{
			ID:     EmbID(docs[i].ID),
			DocID:  docs[i].ID,
			Vector: v[i],
		}
	}

	return nil
}

func (m *Memory[T]) Search(query string, top int) ([]Result[T], error) {
	if v, ok := m.cache.Get(query); ok {
		return Top(v, top), nil
	}

	vq, err := m.Embeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

	// get with lock
	m.RLock()
	defer m.RUnlock()

	results := make([]Result[T], 0, len(m.docs))
	for _, v := range m.embeddings {
		results = append(results, Result[T]{
			Score: Score(m.Distance(vq[0], v.Vector)),
			Doc:   m.docs[v.DocID],
		})
	}

	m.cache.Put(query, results)
	return Top(results, top), nil
}

func (m *Memory[T]) Modify(query string, modified []Result[T]) {
	m.Lock()
	defer m.Unlock()

	for _, r := range modified {
		if r.Doc.Ignore {
			m.cache.Ignore(query, r.Doc)
		}
	}

	sort.Slice(modified, func(i, j int) bool {
		return modified[i].Score > modified[j].Score
	})

	m.cache.Put(query, modified)
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results[:min(n, len(results))]
}
