# vecdb

## Example

```go

func Example() {
	type Document struct {
		Title   string
		Creator string
	}

	m := vecdb.New[Document]()
	if err := m.Save(
		[]string{
			"1st document is about morning.",
			"2nd document is about night.",
		},
		[]Document{
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
```
