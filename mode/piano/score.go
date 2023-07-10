package piano

import (
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

var (
	Kool = mode.Judgment{Window: 20, Weight: 1}
	Cool = mode.Judgment{Window: 40, Weight: 1}
	Good = mode.Judgment{Window: 80, Weight: 0.5}
	Miss = mode.Judgment{Window: 120, Weight: 0}
)

var Judgments = []mode.Judgment{Kool, Cool, Good, Miss}

func Judge(noteType int, a input.KeyActionType, td int64) mode.Judgment {
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
				return mode.Evaluate(Judgments, td)
			}
		}
	} else { // Head, Normal
		return mode.Judge(Judgments, a, td)
	}
	return mode.Judgment{}
}

// Extra primitive in Piano mode is a count of Kools.
// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *Note, j mode.Judgment) {
	if j == Miss {
		s.BreakCombo()
	} else {
		s.AddCombo()
	}
	s.CalcScore(mode.Flow, j.Flow, n.Weight())
	s.CalcScore(mode.Acc, j.Acc, n.Weight())
	if j.Is(Kool) {
		s.CalcScore(mode.Extra, 1, n.Weight())
	} else {
		s.CalcScore(mode.Extra, 0, n.Weight())
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
