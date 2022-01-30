package mania

import (
	"math"
	"sort"
)

const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)

// TODO: Should charts' difficulties be calculated once they are loaded?
func (c *Chart) CalcDifficulty() {
	if len(c.Notes) == 0 {
		return
	}
	sectionCounts := int(c.EndTime()-c.Notes[0].Time) / sectionLength // independent of note offset
	sectionEndTime := sectionLength + c.Notes[0].Time
	var d float64
	ds := make([]float64, 0, sectionCounts)

	{ // LevelNaive
		c.CalcStrain()
		for _, n := range c.Notes {
			for n.Time >= sectionEndTime {
				ds = append(ds, d)
				d = 0.0
				sectionEndTime += sectionLength
			}
			switch n.Type {
			case TypeLNTail:
				d += 0.3 // TEMP
			default:
				d += 1
			}
		}
		if len(ds) != sectionCounts {
			panic("section count mismatch")
		}
		c.LevelNaive = WeightedSum(ds, diffWeightDecay) / 20
	}

	{ // LevelWeighted: ultimate goal
		c.CalcStrain()
		for _, n := range c.Notes {
			for n.Time >= sectionEndTime {
				ds = append(ds, d)
				d = 0.0
				sectionEndTime += sectionLength
			}
			d += n.strain
		}

		if len(ds) != sectionCounts {
			panic("section count mismatch")
		}
		c.LevelWeighted = WeightedSum(ds, diffWeightDecay) / 20
		c.allotScore()
	}
	{ // LevelOsuLegacy
		type noteOsuLegacy struct {
			Note
			strain           float64
			heldUntil        []int64
			individualStrain []float64
		}
		var (
			strainStep          int64   = 400
			weightDecayBase     float64 = 0.9
			individualDecayBase float64 = 0.125
			srScalingFactor     float64 = 0.018
			overallDecayBase    float64 = 0.3
		)

		var notes []noteOsuLegacy = make([]noteOsuLegacy, 0, len(c.Notes))
		var prevNote noteOsuLegacy
		for i, n0 := range c.Notes {
			n := noteOsuLegacy{
				Note:             n0,
				heldUntil:        make([]int64, c.KeyCount),
				individualStrain: make([]float64, c.KeyCount),
			}
			if i == 0 {
				prevNote = n
				notes = append(notes, n)
				continue
			}
			if n0.Type == TypeLNTail {
				continue
			}
			timeElapsed := (n0.Time - prevNote.Time)
			individualDecay := math.Pow(individualDecayBase, float64(timeElapsed)/1000)
			overallDecay := math.Pow(overallDecayBase, float64(timeElapsed)/1000)

			var (
				holdFactor   float64 = 1.0
				holdAddition float64 = 0.0
			)
			for k := 0; k < c.KeyCount; k++ {
				n.heldUntil[k] = prevNote.heldUntil[k]
				if n.Time < n.heldUntil[k] && n.Time2 > n.heldUntil[k] {
					holdAddition = 1
				} else if n.Time2 == n.heldUntil[k] {
					holdAddition = 0
				} else if n.Time2 < n.heldUntil[k] {
					holdFactor = 1.25
				}
				n.individualStrain[k] = prevNote.individualStrain[k] * individualDecay
			}
			n.heldUntil[n.key] = n.Time2
			n.individualStrain[n.key] += 2.0 * holdFactor
			n.strain = prevNote.strain*overallDecay + (1.0+holdAddition)*holdFactor
			prevNote = n
			notes = append(notes, n)
		}

		var (
			strainTable     []float64
			intervalEndTime = strainStep
			maximumStrain   float64
		)
		prevNote = noteOsuLegacy{}
		for _, n := range notes {
			for n.Time > int64(intervalEndTime) {
				strainTable = append(strainTable, maximumStrain)
				if prevNote.heldUntil == nil { // rough way to check whether prevNote is empty
					intervalEndTime += strainStep
					continue
				}
				individualDecay := math.Pow(individualDecayBase, float64(intervalEndTime-prevNote.Time)/1000.0)
				overallDecay := math.Pow(overallDecayBase, float64(intervalEndTime-prevNote.Time)/1000.0)
				maximumStrain = n.individualStrain[n.key]*individualDecay + n.strain*overallDecay

				intervalEndTime += strainStep
			}
			maximumStrain = math.Max(n.individualStrain[n.key]+prevNote.strain, maximumStrain)
			prevNote = n
		}
		c.LevelOsuLegacy = WeightedSum(strainTable, weightDecayBase) * srScalingFactor
	}
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
