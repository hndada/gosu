package taiko

import (
	"github.com/hndada/gosu/game"
)

// minimize parameters so that no need to do extra process: learning
const (
	MaxTrillBonus      = 0.1
	MaxChordPenalty    = -0.3
	TrillBonusXOffset1 = 30

	MaxJackBonus      = 0.5
	Max2JackBonus     = 0.3
	JackBonusXOffset2 = 0
)

func (beatmap *TaikoBeatmap) SetCurves() {
	beatmap.Curves = make(map[string]game.Segments)
	beatmap.Curves["Trill"] = game.NewSegments(
		[]float64{0, beatmap.HitWindows["100"] + TrillBonusXOffset1, beatmap.HitWindows["0"]},
		[]float64{MaxChordPenalty, MaxTrillBonus, 0})
	beatmap.Curves["Jack"] = game.NewSegments(
		[]float64{0, beatmap.HitWindows["100"], beatmap.HitWindows["0"] + JackBonusXOffset2},
		[]float64{MaxJackBonus, Max2JackBonus, 0})
}
