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
	docs       map[DocID]Doc[T]
	embeddings map[DocID]Embedding
	labels     map[Label]Doc[T]
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
		m.labels = make(map[Label]Doc[T])
	}

	for i := range v {
		m.Add(docs[i], v[i])
	}

	return nil
}

func (m *Memory[T]) Add(doc Doc[T], embed []float64) error {
	if v, ok := m.labels[doc.Label]; ok {
		isDup, err := m.IsDuplicated(doc, v)
		if err != nil {
			return fmt.Errorf("is duplicated: %v", err)
		}

		if isDup {
			return nil
		}
	}

	m.labels[doc.Label] = doc

	m.docs[doc.ID] = Doc[T]{
		ID:       doc.ID,
		Label:    doc.Label,
		Text:     doc.Text,
		Metadata: doc.Metadata,
	}

	m.embeddings[doc.ID] = Embedding{
		ID:     EmbID(doc.ID),
		DocID:  doc.ID,
		Vector: embed,
	}

	return nil
}

func (m *Memory[T]) IsDuplicated(latest, old Doc[T]) (bool, error) {
	// TODO: AI task
	return latest.Text == old.Text, nil
}

func (m *Memory[T]) Search(query string, top int) ([]Result[T], error) {
	if v, ok := m.cache.Get(query); ok {
		return Top(v, top), nil
	}

	vq, err := m.Embeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

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
		if r.Ignore {
			m.cache.Ignore(query, r.Doc)
		}
	}

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
