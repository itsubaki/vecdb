package vecdb

import (
	"fmt"
	"sort"
)

type Text string

type Doc[T any] struct {
	Text     Text
	Vector   []float64
	Metadata T
	Ignore   bool
}

type Result[T any] struct {
	Score float64
	Doc   Doc[T]
}

type Memory[T any] struct {
	Distance   func(a, b []float64) float64
	Embeddings func(text []string) ([][]float64, error)
	docs       map[Text]Doc[T]
	cache      Cache[T]
}

func (m *Memory[T]) Save(docs []Doc[T]) error {
	text := make([]string, len(docs))
	for i, d := range docs {
		text[i] = string(d.Text)
	}

	v, err := m.Embeddings(text)
	if err != nil {
		return fmt.Errorf("embedding: %v", err)
	}

	for i := range v {
		m.docs[docs[i].Text] = Doc[T]{
			Text:     docs[i].Text,
			Metadata: docs[i].Metadata,
			Vector:   v[i],
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

	results := make([]Result[T], 0, len(m.docs))
	for _, doc := range m.docs {
		results = append(results, Result[T]{
			Score: Score(m.Distance(vq[0], doc.Vector)),
			Doc:   doc,
		})
	}

	m.cache.Put(query, results)
	return Top(results, top), nil
}

func (m *Memory[T]) Modify(query string, modified Result[T]) {
	if modified.Doc.Ignore {
		m.cache.Ignore(query, modified.Doc)
	}

	// TODO
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	return results[:min(n, len(results))]
}
