package vecdb_test

import (
	"fmt"

	"github.com/itsubaki/vecdb"
)

func ExampleCache() {
	s := &vecdb.Cache[string]{}
	s.Put("foo", []vecdb.Result[string]{
		{
			Score: 1.2,
			Doc: vecdb.Doc[string]{
				ID:   "1",
				Text: "bar",
			},
		},
	})

	v, ok := s.Get("foo")
	if !ok {
		panic("no such entity")
	}
	for _, r := range v {
		fmt.Printf("%.1f %s\n", r.Score, r.Doc.Text)
	}

	// Output:
	// 1.2 bar
}
