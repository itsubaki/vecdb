package vecdb

import (
	"fmt"
	"sort"
	"sync"
)

type DocID string

type Doc[T any] struct {
	ID        DocID
	Text      string
	Metadata  T
	Embedding []float64
}

type Result[T any] struct {
	Score float64
	Doc   Doc[T]
}

type Memory[T any] struct {
	Distance   func(a, b []float64) float64
	Embeddings func(text []string) ([][]float64, error)
	docs       map[DocID]Doc[T]
	sync.RWMutex
}

func (m *Memory[T]) Save(docs []Doc[T]) error {
	// embeddings
	text := make([]string, len(docs))
	for i, d := range docs {
		text[i] = d.Text
	}

	v, err := m.Embeddings(text)
	if err != nil {
		return fmt.Errorf("embedding: %v", err)
	}

	// save
	m.Lock()
	defer m.Unlock()

	if m.docs == nil {
		m.docs = make(map[DocID]Doc[T])
	}

	for i := range v {
		m.docs[docs[i].ID] = Doc[T]{
			ID:        docs[i].ID,
			Text:      docs[i].Text,
			Metadata:  docs[i].Metadata,
			Embedding: v[i],
		}
	}

	return nil
}

func (m *Memory[T]) Search(query string, top int) ([]Result[T], error) {
	vq, err := m.Embeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

	m.RLock()
	defer m.RUnlock()

	results := make([]Result[T], 0, len(m.docs))
	for _, v := range m.docs {
		results = append(results, Result[T]{
			Score: Score(m.Distance(vq[0], v.Embedding)),
			Doc:   v,
		})
	}

	return Top(results, top), nil
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results[:min(n, len(results))]
}
