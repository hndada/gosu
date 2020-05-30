package mania

import (
	"log"

	"github.com/hndada/gosu/game/tools"
)

// minimize parameters so that no need to do extra process: learning
const (
	MaxChordPenalty      = -0.465 // ok
	MaxTrillBonus        = 0.08
	TrillChordXOffset300 = -10
	TrillChordXOffset200 = 0
	TrillChordXOffset0   = 0

	MaxJackBonus  = 0.05
	Max2JackBonus = 0.025
	JackXOffset1  = 20
	JackXOffset2  = 40

	MinHoldTailStrain  = 0.05
	MaxHoldTailStrain  = 0.2
	ZeroHoldTailStrain = 0.1
)

func (beatmap *ManiaBeatmap) SetCurves() {
	if MaxChordPenalty < -0.5 {
		log.Fatal("Chord penalty should not be lower than -0.5")
	}
	beatmap.Curves = make(map[string][]tools.Segment)
	beatmap.Curves["TrillChord"] = tools.GetSegments(
		[]float64{
			0,
			beatmap.HitWindows["300"] + TrillChordXOffset300,
			beatmap.HitWindows["200"] + TrillChordXOffset200,
			beatmap.HitWindows["0"] + TrillChordXOffset0},
		[]float64{
			MaxChordPenalty,
			0,
			MaxTrillBonus,
			0})
	beatmap.Curves["Jack"] = tools.GetSegments(
		[]float64{
			0,
			beatmap.HitWindows["100"] + JackXOffset1,
			beatmap.HitWindows["0"] + JackXOffset2},
		[]float64{
			MaxJackBonus,
			Max2JackBonus,
			0})
	beatmap.Curves["HoldTail"] = tools.GetSegments(
		[]float64{
			0,
			beatmap.HitWindows["320"],
			beatmap.HitWindows["100"]},
		[]float64{
			ZeroHoldTailStrain,
			MinHoldTailStrain,
			MaxHoldTailStrain})
}
