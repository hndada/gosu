package piano

import (
	"image/color"

	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

var (
	Kool = mode.Judgment{Flow: 0.01, Acc: 1, Window: 20}
	Cool = mode.Judgment{Flow: 0.01, Acc: 1, Window: 45}
	Good = mode.Judgment{Flow: 0.01, Acc: 0.25, Window: 75}
	Bad  = mode.Judgment{Flow: 0.01, Acc: 0, Window: 110} // Todo: Flow 0.01 -> 0?
	Miss = mode.Judgment{Flow: -1, Acc: 0, Window: 150}
)

var Judgments = []mode.Judgment{Kool, Cool, Good, Bad, Miss}
var JudgmentColors = []color.NRGBA{
	mode.ColorKool, mode.ColorCool, mode.ColorGood, mode.ColorBad, mode.ColorMiss}

func Verdict(noteType int, a input.KeyAction, td int64) mode.Judgment {
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
				return mode.Judge(Judgments, td)
			}
		}
	} else { // Head, Normal
		return mode.Verdict(Judgments, a, td)
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
	// s.NoteCount++
}

// func (s ScenePlay) LinearScore() float64 {
// 	return s.ScoreBounds[mode.Total] * float64(s.NoteCount) / float64(s.MaxNoteCount)
// }
