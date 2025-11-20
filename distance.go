package vecdb

import "math"

func Cosine(x, y []float64) float64 {
	var xsum, ysum float64
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
