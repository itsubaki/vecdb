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
	Score float64
	Vector[T]
}

type Memory[T any] struct {
	List       []Vector[T]
	Similarity func(a, b []float64) float64
	Embeddings func(text []string) ([][]float64, error)
}

func (m *Memory[T]) Save(text []string, metadata []T) error {
	v, err := m.Embeddings(text)
	if err != nil {
		return fmt.Errorf("embedding: %v", err)
	}

	for i := range v {
		m.List = append(m.List, Vector[T]{
			Data:     v[i],
			Text:     text[i],
			Metadata: metadata[i],
		})
	}

	return nil
}

func (m *Memory[T]) Search(query string, top int) ([]Result[T], error) {
	vq, err := m.Embeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding: %v", err)
	}

	results := make([]Result[T], len(m.List))
	for i, v := range m.List {
		results[i] = Result[T]{
			Score:  Score(m.Similarity(vq[0], v.Data)),
			Vector: v,
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

	return 1 - (sum / (xps * yps))
}

func Euclid(x, y []float64) float64 {
	var sum float64
	for i := range x {
		sum += math.Pow(x[i]-y[i], 2)
	}

	return math.Sqrt(sum)
}

func Score(v float64) float64 {
	return 1 / (1 + v)
}

func Top[T any](results []Result[T], n int) []Result[T] {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
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
