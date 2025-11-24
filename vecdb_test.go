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
		Ignore: func(doc vecdb.Doc[Metadata]) bool {
			return doc.Label == "allnight"
		},
	}

	if err := db.Save([]vecdb.Doc[Metadata]{
		{
			ID:    "1",
			Label: "morning",
			Text:  "1st document is about morning.",
			Metadata: Metadata{
				Title:   "Morning",
				Creator: "John Doe",
			},
		},
		{
			ID:    "2",
			Label: "night",
			Text:  "2nd document is about night.",
			Metadata: Metadata{
				Title:   "Night",
				Creator: "John Doe",
			},
		},
		{
			ID:    "3",
			Label: "midnight",
			Text:  "3rd document is about midnight",
			Metadata: Metadata{
				Title:   "Midnight",
				Creator: "John Doe",
			},
		},
		{
			ID:    "4",
			Label: "daybreak",
			Text:  "4th document is about daybreak",
			Metadata: Metadata{
				Title:   "Daybreak",
				Creator: "John Doe",
			},
		},
		{
			ID:    "5",
			Label: "morning", // duplicated label
			Text:  "1st document is about morning.",
			Metadata: Metadata{
				Title:   "Morning",
				Creator: "John Doe",
			},
		},
		{
			ID:    "6",
			Label: "allnight", // this will be ignored
			Text:  "1st document is about allnight.",
			Metadata: Metadata{
				Title:   "Allnight",
				Creator: "John Doe",
			},
		},
	}); err != nil {
		panic(err)
	}

	dup := db.Dups()
	for label, docs := range dup {
		for _, doc := range docs {
			fmt.Printf("label: %q, doc: %v\n", label, doc)
		}
	}

	// TODO: AI task to decide which one to remove, update or noop
	db.Remove([]vecdb.DocID{"5"})
	fmt.Println("removed docID 5")
	fmt.Println("-")

	// search
	query, top := "Night and day", 5
	results, err := db.Search(query, top)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Printf("%.4f, %v %q %+v\n", r.Score, r.Doc.ID, r.Doc.Text, r.Doc.Metadata)
	}

	// Ignore the 3rd document
	results[2].Ignore = true
	// Rescore the 1st document
	results[0].Score = 0.1

	// modify
	db.Modify(query, results)
	fmt.Println("ignored 3rd document")
	fmt.Println("rescored 1st document")
	fmt.Println("-")

	results, err = db.Search(query, top)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Printf("%.4f, %v %q %+v\n", r.Score, r.Doc.ID, r.Doc.Text, r.Doc.Metadata)
	}

	// Output:
}
