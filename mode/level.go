package mode

import (
	"math"
	"sort"
)

// Todo: find the best SliceDuration value
const (
	SliceDuration = 800
	DecayFactor   = 0.95
	LevelPower    = 1.15
	LevelScale    = 0.02
)

// Todo: Mods as input parameter
// Input is Difficulties.
func Level(ds []float64) float64 {
	// Todo: new mods-applied chart here
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })
	sum, weight := 0.0, 1.0
	for _, term := range ds {
		sum += weight * term
		weight *= DecayFactor
	}
	return math.Pow(sum, LevelPower) * LevelScale
}
