package vecdb_test

import (
	"fmt"

	"github.com/itsubaki/vecdb"
)

func Example() {
	embeddings := func(text []string) ([][]float64, error) {
		var emb [][]float64
		for i := range text {
			emb = append(emb, []float64{
				float64(i+len(text)) + 0,
				float64(i+len(text)) + 1,
				float64(i+len(text)) + 2,
				float64(i+len(text)) + 3,
			})
		}

		return emb, nil
	}

	type Metadata struct {
		Title   string
		Creator string
	}

	db := vecdb.Memory[Metadata]{
		Distance:   vecdb.Cosine,
		Embeddings: embeddings,
	}

	if err := db.Save([]vecdb.Doc[Metadata]{
		{
			ID:   "1",
			Text: "1st document is about morning.",
			Metadata: Metadata{
				Title:   "Morning",
				Creator: "John Doe",
			},
		},
		{
			ID:   "2",
			Text: "2nd document is about night.",
			Metadata: Metadata{
				Title:   "Night",
				Creator: "John Doe",
			},
		},
		{
			ID:   "3",
			Text: "3rd document is about midnight",
			Metadata: Metadata{
				Title:   "Midnight",
				Creator: "John Doe",
			},
		},
		{
			ID:   "4",
			Text: "4th document is about daybreak",
			Metadata: Metadata{
				Title:   "Daybreak",
				Creator: "John Doe",
			},
		},
	}); err != nil {
		panic(err)
	}

	query, top := "Night and day", 3
	results, err := db.Search(query, top)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Printf("%.4f, %q %+v\n", r.Score, r.Doc.Text, r.Doc.Metadata)
	}

	// Ignore the 3rd document
	results[2].Doc.Ignore = true
	db.Modify(query, results)
	fmt.Println("ignored the 3rd document")

	results, err = db.Search(query, top)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Printf("%.4f, %q, %+v\n", r.Score, r.Doc.Text, r.Doc.Metadata)
	}

	// Output:
}
