package vecdb

import (
	"fmt"
	"sort"
	"sync"
)

type (
	DocID string
	EmbID string
	Label string
)

type Doc[T any] struct {
	ID       DocID
	Label    Label
	Text     string
	Metadata T
}

type Embedding struct {
	ID     EmbID
	DocID  DocID
	Vector []float64
}

type Result[T any] struct {
	Score  float64
	Ignore bool
	Doc    Doc[T]
}

type Memory[T any] struct {
	Distance   func(a, b []float64) float64
	Embeddings func(text []string) ([][]float64, error)
	Ignore     func(doc Doc[T]) bool
	docs       map[DocID]Doc[T]
	embeddings map[DocID]Embedding
	labels     map[Label]map[DocID]Doc[T]
	cache      Cache[T]
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

	if m.embeddings == nil {
		m.embeddings = make(map[DocID]Embedding)
	}

	if m.labels == nil {
		m.labels = make(map[Label]map[DocID]Doc[T])
	}

	for i := range v {
		if m.labels[docs[i].Label] == nil {
			m.labels[docs[i].Label] = make(map[DocID]Doc[T])
		}

		m.labels[docs[i].Label][docs[i].ID] = Doc[T]{
			ID:       docs[i].ID,
			Label:    docs[i].Label,
			Text:     docs[i].Text,
			Metadata: docs[i].Metadata,
		}

		if len(m.labels[docs[i].Label]) > 1 {
			continue
		}

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
	vq, err := m.Embeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

	m.RLock()
	defer m.RUnlock()

	results := make([]Result[T], 0, len(m.docs))
	for _, v := range m.embeddings {
		if m.Ignore != nil && m.Ignore(m.docs[v.DocID]) {
			continue
		}

		results = append(results, Result[T]{
			Score: Score(m.Distance(vq[0], v.Vector)),
			Doc:   m.docs[v.DocID],
		})
	}

	cached, ok := m.cache.Get(query)
	if !ok {
		m.cache.Put(query, results)
		return Top(results, top), nil
	}

	// merge
	var out []Result[T]
	for _, r := range results {
		if result, found := cached[r.Doc.ID]; found {
			if result.Ignore {
				continue
			}

			out = append(out, result)
			continue
		}

		out = append(out, r)
	}

	return Top(out, top), nil
}

func (m *Memory[T]) Dups() map[Label][]Doc[T] {
	m.RLock()
	defer m.RUnlock()

	dups := make(map[Label][]Doc[T])
	for label, docs := range m.labels {
		if len(docs) > 1 {
			dups[label] = make([]Doc[T], 0, len(docs))
			for _, d := range docs {
				dups[label] = append(dups[label], d)
			}

			sort.Slice(dups[label], func(i, j int) bool {
				return dups[label][i].ID < dups[label][j].ID
			})
		}
	}

	return dups
}

func (m *Memory[T]) Remove(docIDs []DocID) {
	m.Lock()
	defer m.Unlock()

	for _, id := range docIDs {
		doc, ok := m.docs[id]
		if !ok {
			continue
		}

		delete(m.labels[doc.Label], id)
		delete(m.docs, id)
		delete(m.embeddings, id)
	}
}

func (m *Memory[T]) Modify(query string, modified []Result[T]) {
	m.Lock()
	defer m.Unlock()

	sort.Slice(modified, func(i, j int) bool {
		return modified[i].Score > modified[j].Score
	})

	m.cache.Put(query, modified)
}

func (m *Memory[T]) Docs() []Doc[T] {
	m.RLock()
	defer m.RUnlock()

	docs := make([]Doc[T], 0, len(m.docs))
	for _, d := range m.docs {
		docs = append(docs, d)
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ID > docs[j].ID
	})

	return docs
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results[:min(n, len(results))]
}
