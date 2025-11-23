# vecdb

## Example

```go
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

	db := vecdb.Memory[Metadata]{
		Distance:   vecdb.Cosine,
		Embeddings: client.Embeddings,
	}

	if err := db.Save([]vecdb.Doc[Metadata]{
		{
			ID:  "1",
			Text: "1st document is about morning.",
			Metadata: Metadata{
				Title:   "Morning",
				Creator: "John Doe",
			},
		},
		{
			ID:  "2",
			Text: "2nd document is about night.",
			Metadata: Metadata{
				Title:   "Night",
				Creator: "John Doe",
			},
		},
		//
		// and more
		//
	}); err != nil {
		panic(err)
	}

	top := 3
	results, err := db.Search("Night and day", top)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Println(r.Score, r.Doc.Text, r.Doc.Metadata)
	}

	// Output:
}
```
