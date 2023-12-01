package openai_test

import (
	"fmt"
	"sort"

	"github.com/itsubaki/vecdb/openai"
)

func ExampleClient_Models() {
	c := openai.New("", "")
	models, err := c.Models()
	if err != nil {
		panic(err)
	}

	sort.Slice(models.Data, func(i, j int) bool {
		return models.Data[i].ID < models.Data[j].ID
	})

	for _, m := range models.Data {
		fmt.Println(m)
	}

	// Output:
}
