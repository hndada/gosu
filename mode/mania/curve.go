package mania

import (
	"github.com/hndada/gosu/mode"
	"math"
)

const (
	MaxChordPenalty = -0.3 // must be same or greater than -0.5
	MaxTrillBonus   = 0.08

	MaxJackBonus  = 0.25
	Max2JackBonus = 0.1
	Max2DeltaJack = 120
	maxDeltaJack  = 180

	MinHoldTailStrain  = 0.05
	MaxHoldTailStrain  = 0.2
	ZeroHoldTailStrain = 0.1
)

var (
	maxDeltaChord   int64
	maxDeltaTrill   int64
	curveTrillChord []mode.Segment
	curveJack       []mode.Segment
	curveTail       []mode.Segment
)

func init() {
	curveTrillChord = mode.GetSegments(
		[]float64{
			0,
			float64(good.Window + 30),
			float64(miss.Window + 30)},
		[]float64{
			MaxChordPenalty,
			MaxTrillBonus,
			0})

	curveJack = mode.GetSegments(
		[]float64{
			0,
			Max2DeltaJack,
			maxDeltaJack},
		[]float64{
			MaxJackBonus,
			Max2JackBonus,
			0})

	curveTail = mode.GetSegments(
		[]float64{
			0,
			float64(kool.Window),
			float64(bad.Window)},
		[]float64{
			ZeroHoldTailStrain,
			MinHoldTailStrain,
			MaxHoldTailStrain})

	xValues := mode.SolveX(curveTrillChord, 0)
	if len(xValues) != 2 {
		panic("incorrect numbers of xValues")
	}
	maxDeltaChord = int64(math.Round(xValues[0]))
	maxDeltaTrill = int64(math.Round(xValues[1]))

	// maxDeltaJack = int(math.Round(tools.SolveX(beatmap.Curves["Jack"], 0)[0]))
}
