package vecdb_test

import (
	"fmt"

	"github.com/itsubaki/vecdb"
)

func ExampleCosine() {
	v := []float64{1, 2, 3}
	w := []float64{4, 5, 6}

	fmt.Println(vecdb.Cosine(v, w))

	// Output:
	// 0.02536815421429417
}

func ExampleEuclid() {
	v := []float64{1, 2, 3}
	w := []float64{4, 5, 6}

	fmt.Println(vecdb.Euclid(v, w))

	// Output:
	// 5.196152422706632
}
