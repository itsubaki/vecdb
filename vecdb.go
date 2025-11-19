package vecdb

import (
	"fmt"
	"sort"
)

type Doc[T any] struct {
	Text     string
	Vector   []float64
	Metadata T
}

type Result[T any] struct {
	Score float64
	Doc   Doc[T]
}

type Memory[T any] struct {
	List       []Doc[T]
	Distance   func(a, b []float64) float64
	Embeddings func(text []string) ([][]float64, error)
	Cache      Cache[T]
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

	for i := range v {
		m.List = append(m.List, Doc[T]{
			Text:     docs[i].Text,
			Metadata: docs[i].Metadata,
			Vector:   v[i],
		})
	}

	return nil
}

func (m *Memory[T]) Search(query string, top int) ([]Result[T], error) {
	if v, ok := m.Cache.Get(query); ok {
		return v, nil
	}

	vq, err := m.Embeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

	results := make([]Result[T], len(m.List))
	for i, doc := range m.List {
		results[i] = Result[T]{
			Score: Score(m.Distance(vq[0], doc.Vector)),
			Doc:   doc,
		}
	}

	m.Cache.Put(query, results)
	return Top(results, top), nil
}

func (m *Memory[T]) Rerank(id string, old, latest []Result[T]) {
	m.Cache.Put(id, latest)

	// TODO: logging
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	n = min(n, len(results))
	top := make([]Result[T], n)
	for i := range n {
		top[i] = results[i]
	}

	return top
}
