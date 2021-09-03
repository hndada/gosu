package mania

import (
	"sort"
)

const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)

// Chart가 Load됨과 동시에 계산되어야 할까?
func (c *Chart) CalcDifficulty() {
	if len(c.Notes) == 0 {
		return
	}
	c.CalcStrain()
	sectionCounts := int(c.EndTime()-c.Notes[0].Time) / sectionLength // independent of note offset
	sectionEndTime := sectionLength + c.Notes[0].Time

	var d float64
	ds := make([]float64, 0, sectionCounts)
	for _, n := range c.Notes {
		for n.Time >= sectionEndTime {
			ds = append(ds, d)
			d = 0.0
			sectionEndTime += sectionLength
		}
		d += n.strain // + n.stamina // n.read (->SV) someday
	}

	if len(ds) != sectionCounts {
		// fmt.Println(len(ds), sectionCounts)
		panic("section count mismatch")
	}
	c.Level = WeightedSum(ds, diffWeightDecay) / 20
	c.allotScore()
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
