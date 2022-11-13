package chart

import (
	"math"
	"sort"
)

// Todo: find the best SliceDuration value
const (
	DecayFactor = 0.95
	LevelPower  = 1.15
	LevelScale  = 0.02
)

// No need to define interface{} called ChartAnalyzer:
// https://go.dev/play/p/PtgBkwKZFhP
func Level(c interface{ Difficulties() []float64 }) (float64, [3]float64) {
	ds := c.Difficulties()
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })
	sum, weight := 0.0, 1.0
	for _, term := range ds {
		sum += weight * term
		weight *= DecayFactor
	}
	level := math.Pow(sum, LevelPower) * LevelScale
	// Todo: Variate factors based on difficulty-skewed charts
	return level, [3]float64{0.5, 5, 2}
}
