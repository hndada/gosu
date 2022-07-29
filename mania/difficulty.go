package mania

import (
	"fmt"
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
		// There's some deviation with an amount of around 0.01 ~ 0.2
		type noteOsuLegacy struct {
			Note

			holdEndtimes      []int64
			individualStrains []float64

			individualStrain float64
			overallStrain    float64
		}
		const (
			individualDecayBase float64 = 0.125
			overallDecayBase    float64 = 0.3
		)

		var notes []noteOsuLegacy = make([]noteOsuLegacy, 0, len(c.Notes))
		var prevNote noteOsuLegacy
		for i, n0 := range c.Notes {
			n := noteOsuLegacy{
				Note:              n0,
				holdEndtimes:      make([]int64, c.KeyCount),
				individualStrains: make([]float64, c.KeyCount),

				overallStrain: 1,
			}
			if i == 0 {
				prevNote = n
				notes = append(notes, n)
				continue
			}
			if n.Type == TypeLNTail {
				continue
			}
			if n.Time > n.Time2 {
				fmt.Printf("%+v\n", n)
				panic("LN tail reaches")
			}
			timeElapsed := float64(n.Time-prevNote.Time) / 1000 // second
			indDecay := math.Pow(individualDecayBase, timeElapsed)
			ovDecay := math.Pow(overallDecayBase, timeElapsed)
			var (
				holdFactor   float64 = 1.0
				holdAddition float64 = 0.0
			)
			for k := 0; k < c.KeyCount; k++ {
				n.holdEndtimes[k] = prevNote.holdEndtimes[k]
				if definitelyBigger(n.holdEndtimes[k], n.Time, 1) && definitelyBigger(n.Time2, n.holdEndtimes[k], 1) {
					holdAddition = 1
				} else if almostEquals(n.Time2, n.holdEndtimes[k], 1) {
					holdAddition = 0
				} else if definitelyBigger(n.holdEndtimes[k], n.Time2, 1) {
					holdFactor = 1.25
				}
				n.individualStrains[k] = prevNote.individualStrains[k] * indDecay
			}
			n.holdEndtimes[n.Key] = n.Time2
			n.individualStrains[n.Key] += 2.0 * holdFactor
			n.individualStrain = n.individualStrains[n.Key]
			n.overallStrain = prevNote.overallStrain*ovDecay + (1+holdAddition)*holdFactor
			prevNote = n
			notes = append(notes, n)
		}

		const (
			strainStep      int64   = 400
			weightDecayBase float64 = 0.9
			srScalingFactor float64 = 0.018
		)
		var (
			strainTable     []float64
			intervalEndTime = strainStep
			maxStrain       float64
		)
		prevNote = noteOsuLegacy{}
		// StrainValueAt: CurrentStrain + individualStrain + overallStrain - CurrentStrain
		for _, n := range notes {
			for n.Time > int64(intervalEndTime) {
				strainTable = append(strainTable, maxStrain)
				if prevNote.holdEndtimes == nil { // rough way to check whether prevNote is empty
					maxStrain = 0
				} else {
					deltaTime := float64(intervalEndTime-prevNote.Time) / 1000
					// should be prev note
					maxStrain = prevNote.individualStrains[n.Key]*math.Pow(individualDecayBase, deltaTime) + prevNote.overallStrain*math.Pow(overallDecayBase, deltaTime)
				}
				intervalEndTime += strainStep
			}
			maxStrain = math.Max(n.individualStrains[n.Key]+n.overallStrain, maxStrain)
			prevNote = n
		}
		c.LevelOsuLegacy = WeightedSum(strainTable, weightDecayBase) * srScalingFactor
	}
	c.Level = c.LevelOsuLegacy
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

func definitelyBigger(v1, v2, e int64) bool { return v1-e > v2 }
func almostEquals(v1, v2, e int64) bool     { return math.Abs(float64(v1)-float64(v2)) < float64(e) }
