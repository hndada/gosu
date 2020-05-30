package tools

import (
	"math"
	"sort"
)

const InfInt = int(^uint(0) >> 1)

var InfFloat64 = math.Inf(1)

func MaxInt(a, b int) int { return int(math.Max(float64(a), float64(b))) }
func MinInt(a, b int) int { return int(math.Min(float64(a), float64(b))) }
func AbsInt(a int) int    { return int(math.Abs(float64(a))) }
func AvgInt(a []int) float64 {
	total := 0
	for i := 0; i < len(a); i++ {
		total += a[i]
	}
	return float64(total) / float64(len(a))
}
func AvgFloat64(a []float64) float64 {
	total := 0.0
	for i := 0; i < len(a); i++ {
		total += a[i]
	}
	return total / float64(len(a))
}

func Round(v float64, decimal int) float64 {
	scale := math.Pow(10, float64(decimal))
	return math.Round(v*scale) / scale
}

func DecayFactor(decayBase, elapsedTime float64) float64 {
	return math.Pow(decayBase, elapsedTime/1000)
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

func IsIntSameSign(h1, h2 int) bool {
	v := h1 * h2
	switch {
	case v > 0:
		return true
	case v < 0:
		return false
	default: // v==0
		if h1 != 0 || h2 != 0 {
			return false
		} else {
			return true
		}
	}
}
