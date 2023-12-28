package game

func LinearInterpolate(xs, ys []float64) func(float64) float64 {
	return func(x float64) float64 {
		// No out of index panic.
		// https://go.dev/play/p/PzSCpSce1to
		for i, b := range xs[:len(xs)-1] {
			if x < b {
				x0 := xs[i]
				x1 := xs[i+1]
				y0 := ys[i]
				y1 := ys[i+1]
				return y0 + (y1-y0)*(x-x0)/(x1-x0)
			}
		}
		return ys[len(ys)-1]
	}
}

func WeightedSum(vs []float64, decayFactor float64) float64 {
	sum, weight := 0.0, 1.0
	for _, term := range vs {
		sum += weight * term
		weight *= decayFactor
	}
	return sum
}
