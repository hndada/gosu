package mania

import (
	"math"

	"github.com/hndada/gosu/game"
)

const (
	MaxChordPenalty = -0.3 // must be same or greater than -0.5
	MaxTrillBonus   = 0.08

	MaxJackBonus  = 0.2
	Max2JackBonus = 0.05
	Max2DeltaJack = 120
	maxDeltaJack  = 180

	MinHoldTailStrain  = 0   // 0.05
	MaxHoldTailStrain  = 0.3 // 0.2
	ZeroHoldTailStrain = 0   // 0.1
)

var (
	maxDeltaChord   int64
	maxDeltaTrill   int64
	curveTrillChord []game.Segment
	curveJack       []game.Segment
	curveTail       []game.Segment
)

func init() {
	curveTrillChord = game.GetSegments(
		[]float64{
			0,
			float64(good.Window + 30),
			float64(miss.Window + 30)},
		[]float64{
			MaxChordPenalty,
			MaxTrillBonus,
			0})

	curveJack = game.GetSegments(
		[]float64{
			0,
			Max2DeltaJack,
			maxDeltaJack},
		[]float64{
			MaxJackBonus,
			Max2JackBonus,
			0})

	curveTail = game.GetSegments(
		[]float64{
			0,
			float64(kool.Window),
			float64(bad.Window) + 50},
		[]float64{
			ZeroHoldTailStrain,
			MinHoldTailStrain,
			MaxHoldTailStrain})

	xValues := game.SolveX(curveTrillChord, 0)
	if len(xValues) != 2 {
		panic("incorrect numbers of xValues")
	}
	maxDeltaChord = int64(math.Round(xValues[0]))
	maxDeltaTrill = int64(math.Round(xValues[1]))

	// maxDeltaJack = int(math.Round(tools.SolveX(beatmap.Curves["Jack"], 0)[0]))
}
