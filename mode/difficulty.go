package mode

import (
	"math"
	"sort"
)

func DecayFactor(decayBase float64, time int64) float64 {
	return math.Pow(decayBase, float64(time)/1000)
}

// Maclaurin series
func WeightedSum(series []float64, weightDecay float64) float64 {
	sort.Slice(series, func(i, j int) bool { return series[i] > series[j] })
	sum, weight := 0.0, 1.0
	for _, term := range series {
		sum += weight * term
		weight *= weightDecay
	}
	return sum
}
