package lv

import (
	"github.com/hndada/gosu/internal/tools"
	"github.com/hndada/gosu/mode/mania"
	"math"
)

const (
	MaxChordPenalty = -0.3
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
	maxDeltaChord, maxDeltaTrill          int
	curveTrillChord, curveJack, curveTail []tools.Segment
)

func init() {
	// if MaxChordPenalty < -0.5 {
	// 	panic("Chord penalty should not be lower than -0.5")
	// }

	curveTrillChord = tools.GetSegments(
		[]float64{
			0,
			float64(mania.Judgements[2].Window+30),
			float64(mania.Judgements[4].Window+30)},
		[]float64{
			MaxChordPenalty,
			MaxTrillBonus,
			0})

	curveJack = tools.GetSegments(
		[]float64{
			0,
			Max2DeltaJack,
			maxDeltaJack},
		[]float64{
			MaxJackBonus,
			Max2JackBonus,
			0})

	curveTail = tools.GetSegments(
		[]float64{
			0,
			float64(mania.Judgements[0].Window),
			float64(mania.Judgements[3].Window)},
		[]float64{
			ZeroHoldTailStrain,
			MinHoldTailStrain,
			MaxHoldTailStrain})

	xValues := tools.SolveX(curveTrillChord, 0)
	if len(xValues) != 2 {
		panic("incorrect numbers of xValues")
	}
	maxDeltaChord = int(math.Round(xValues[0]))
	maxDeltaTrill = int(math.Round(xValues[1]))

	// maxDeltaJack = int(math.Round(tools.SolveX(beatmap.Curves["Jack"], 0)[0]))
}
