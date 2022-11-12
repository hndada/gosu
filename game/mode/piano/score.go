package piano

import (
	"image/color"

	"github.com/hndada/gosu/framework/input"
	"github.com/hndada/gosu/game"
)

var (
	Kool = game.Judgment{Flow: 0.01, Acc: 1, Window: 20}
	Cool = game.Judgment{Flow: 0.01, Acc: 1, Window: 45}
	Good = game.Judgment{Flow: 0.01, Acc: 0.25, Window: 75}
	Bad  = game.Judgment{Flow: 0.01, Acc: 0, Window: 110} // Todo: Flow 0.01 -> 0?
	Miss = game.Judgment{Flow: -1, Acc: 0, Window: 150}
)

var Judgments = []game.Judgment{Kool, Cool, Good, Bad, Miss}
var JudgmentColors = []color.NRGBA{
	game.ColorKool, game.ColorCool, game.ColorGood, game.ColorBad, game.ColorMiss}

func Verdict(noteType int, a input.KeyAction, td int64) game.Judgment {
	if noteType == Tail { // Either Hold or Release when Tail is not scored
		switch {
		case td > Miss.Window:
			if a == input.Release {
				return Miss
			}
		case td < -Miss.Window:
			return Miss
		default: // In range
			if a == input.Release { // a != Hold
				return game.Judge(Judgments, td)
			}
		}
	} else { // Head, Normal
		return game.Verdict(Judgments, a, td)
	}
	return game.Judgment{}
}

// Extra primitive in Piano mode is a count of Kools.
// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *Note, j game.Judgment) {
	if j == Miss {
		s.BreakCombo()
	} else {
		s.AddCombo()
	}
	s.CalcScore(game.Flow, j.Flow, n.Weight())
	s.CalcScore(game.Acc, j.Acc, n.Weight())
	if j.Is(Kool) {
		s.CalcScore(game.Extra, 1, n.Weight())
	} else {
		s.CalcScore(game.Extra, 0, n.Weight())
	}
	for i, j2 := range Judgments {
		if j.Is(j2) {
			s.JudgmentCounts[i]++
			break
		}
	}
	n.Marked = true
	if n.Type == Head && j == Miss {
		s.MarkNote(n.Next, Miss)
	}
	if n.Type != Tail {
		s.Staged[n.Key] = n.Next
	}
}
