package vecdb

import (
	"fmt"
	"math"
	"sort"
)

type Vector[T any] struct {
	Data     []float64
	Text     string
	Metadata T
}

type Result[T any] struct {
	Similarity float64
	Vector[T]
}

type Memory[T any] struct {
	List       []Vector[T]
	Similarity func(a, b []float64) float64
	Embedding  func(text string) ([]float64, error)
}

func New[T any]() *Memory[T] {
	return &Memory[T]{
		List:       make([]Vector[T], 0),
		Similarity: Cosine,
		Embedding:  Embedding,
	}
}

func (m *Memory[T]) Save(text []string, metadata []T) error {
	for i := range text {
		if err := m.save(text[i], metadata[i]); err != nil {
			return err
		}
	}

	return nil
}

func (m *Memory[T]) save(text string, metadata T) error {
	v, err := m.Embedding(text)
	if err != nil {
		return fmt.Errorf("embedding: %v", err)
	}

	m.List = append(m.List, Vector[T]{
		Data:     v,
		Text:     text,
		Metadata: metadata,
	})

	return nil
}

func (m *Memory[T]) Search(query string, top int) ([]Result[T], error) {
	vq, err := m.Embedding(query)
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

	results := make([]Result[T], len(m.List))
	for i, v := range m.List {
		results[i] = Result[T]{
			Similarity: m.Similarity(vq, v.Data),
			Vector:     v,
		}
	}

	return Top(results, top), nil
}

func Cosine(x, y []float64) float64 {
	xsum, ysum := 0.0, 0.0
	for i := range x {
		xsum += x[i] * x[i]
		ysum += y[i] * y[i]
	}

	xps := math.Sqrt(xsum + 1e-8)
	yps := math.Sqrt(ysum + 1e-8)

	var sum float64
	for i := range x {
		sum += x[i] * y[i]
	}

	return sum / (xps * yps)
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if n > len(results) {
		n = len(results)
	}

	top := make([]Result[T], n)
	for i := 0; i < n; i++ {
		top[i] = results[i]
	}

	return top
}

func Embedding(text string) ([]float64, error) {
	// TODO
	return []float64{}, nil
}