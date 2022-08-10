package gosu

import (
	"math"
	"sort"
)

var FingerMap = map[int][]int{
	0:  {},
	1:  {0},
	2:  {1, 1},
	3:  {1, 0, 1},
	4:  {2, 1, 1, 2},
	5:  {2, 1, 0, 1, 2},
	6:  {3, 2, 1, 1, 2, 3},
	7:  {3, 2, 1, 0, 1, 2, 3},
	8:  {4, 3, 2, 1, 1, 2, 3, 4},
	9:  {4, 3, 2, 1, 0, 1, 2, 3, 4},
	10: {4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
}

func init() {
	for k := 2; k <= 8; k++ {
		FingerMap[k|LeftScratch] = append([]int{FingerMap[k-1][0] + 1}, FingerMap[k-1]...)
		FingerMap[k|RightScratch] = append(FingerMap[k-1], FingerMap[k-1][k-2]+1)
	}
}

// Todo: find the best SliceDuration value
const (
	SliceDuration = 800
	DecayFactor   = 0.95
	LevelPower    = 1.15
	LevelScale    = 0.02
)

// Todo: Mods as input parameter
func (c Chart) Level() float64 {
	// Todo: new mods-applied chart here
	ds := c.Difficulties()
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })
	sum, weight := 0.0, 1.0
	for _, term := range ds {
		sum += weight * term
		weight *= DecayFactor
	}
	return math.Pow(sum, LevelPower) * LevelScale
}

// Mods may change the duration of chart.
// Todo: implement actual calculating chart difficulties
func (c Chart) Difficulties() []float64 {
	if len(c.Notes) == 0 {
		return make([]float64, 0)
	}
	endTime := c.Notes[len(c.Notes)-1].Time
	ds := make([]float64, 0, endTime/SliceDuration+1)
	t := c.Notes[0].Time
	var d float64
	for _, n := range c.Notes {
		for n.Time > t+SliceDuration {
			ds = append(ds, d)
			d = 0
			t += SliceDuration
		}
		switch n.Type {
		case Tail:
			d += 0.15
		default:
			d += 1
		}
	}
	return ds
}
