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
				Text: "bar",
			},
		},
	})

	v, ok := s.Get("foo")
	if !ok {
		panic("no such entity")
	}

	fmt.Println(v[0].Score, v[0].Doc.Text)

	// Output:
	// 1.2 bar
}
