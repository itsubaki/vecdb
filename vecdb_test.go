package vecdb_test

import (
	"fmt"
	"os"

	"github.com/itsubaki/vecdb"
	"github.com/itsubaki/vecdb/openai"
)

func Example() {
	client := openai.Client{
		Org:     os.Getenv("OPENAI_ORG"),
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		ModelID: openai.TEXT_EMBEDDING_ADA_002,
	}

	type Metadata struct {
		Title   string
		Creator string
	}

	m := vecdb.Memory[Metadata]{
		Similarity: vecdb.Cosine,
		Embedding:  client.Embedding,
	}

	if err := m.Save(
		[]string{
			"1st document is about morning.",
			"2nd document is about night.",
		},
		[]Metadata{
			{Title: "Morning", Creator: "John Doe"},
			{Title: "Night", Creator: "John Doe"},
		},
	); err != nil {
		panic(err)
	}

	top := 3
	results, err := m.Search("Hello", top)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Println(r.Similarity, r.Metadata)
	}

	// Output:
}
