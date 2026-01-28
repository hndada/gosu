package piano

import (
	"sort"

	"github.com/hndada/gosu/game"
)

// 2 beats with 150 BPM takes 800ms
const unitDuration = 800

const (
	ScoreScale           = 1_000_000
	MaxScore         int = 1.1 * ScoreScale
	StandardMaxScore int = 1.0 * ScoreScale
)

const baseLevelScale = 0.25

var (
	scoreXs = []float64{
		0.60, 0.70, 0.80, 0.90, 1.00, 1.10}
	decayFactorYs = []float64{
		0.99, 0.98, 0.97, 0.95, 0.90, 0.50}
	levelWeightYs = []float64{
		0.85, 0.90, 0.95, 1.0, 1.1, 1.25}

	decayFactor = game.LinearInterpolate(
		scoreXs, decayFactorYs)
	levelWeight = game.LinearInterpolate(
		scoreXs, levelWeightYs)
)

func (c Chart) StandardLevel() float64 {
	return c.Level(StandardMaxScore)
}

// Each target score gives different level.

// Different BPM make duration of 'diff' different.
// However, it looks fine not to scale each diff based on its duration
// and using the same size of duration on each piece.
// They will be alleviated into diffs.
func (c Chart) Level(score int) float64 {
	factor := decayFactor(float64(score))
	scale := 1 / factor
	scale *= levelWeight(float64(score))
	if len(c.notes) == 0 {
		return 0.0
	}
	var diff float64
	diffs := make([]float64, 0, 200)

	endTime := c.notes[0].Time + unitDuration
	for _, n := range c.notes {
		if n.Time >= endTime {
			diffs = append(diffs, diff)
			diff = 0
			endTime += unitDuration
		}
		diff += n.strain
	}

	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	difficulty := game.WeightedSum(diffs, factor)
	return difficulty * scale
}
