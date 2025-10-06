package main

import "math"

func SimpleDecay(maxX, maxY, power float64) func(x float64) float64 {
	return func(x float64) float64 {
		base := (maxY - x/maxX)
		if base < 0 {
			return 0
		}
		return math.Pow(base, power)
	}
}

func LinearFunc(xs, ys []float64) func(x float64) float64 {
	n := len(xs)
	if n == 0 || len(ys) != n {
		panic("xs and ys must have same non-zero length")
	}

	return func(x float64) float64 {
		// before first point
		if x <= xs[0] {
			return ys[0]
		}
		// after last point
		if x >= xs[n-1] {
			return ys[n-1]
		}

		// find interval
		for i := 0; i < n-1; i++ {
			if x >= xs[i] && x <= xs[i+1] {
				x0, y0 := xs[i], ys[i]
				x1, y1 := xs[i+1], ys[i+1]
				t := (x - x0) / (x1 - x0)
				return y0 + t*(y1-y0)
			}
		}
		// fallback (shouldn’t reach)
		return ys[n-1]
	}
}
