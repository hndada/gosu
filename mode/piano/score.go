package piano

import (
	"image/color"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

var (
	Kool = gosu.Judgment{Flow: 0.01, Acc: 1, Window: 20}
	Cool = gosu.Judgment{Flow: 0.01, Acc: 1, Window: 45}
	Good = gosu.Judgment{Flow: 0.01, Acc: 0.25, Window: 75}
	Bad  = gosu.Judgment{Flow: 0.01, Acc: 0, Window: 110} // Todo: Flow 0.01 -> 0?
	Miss = gosu.Judgment{Flow: -1, Acc: 0, Window: 150}
)

var Judgments = []gosu.Judgment{Kool, Cool, Good, Bad, Miss}
var JudgmentColors = []color.NRGBA{
	gosu.ColorKool, gosu.ColorCool, gosu.ColorGood, gosu.ColorBad, gosu.ColorMiss}

func Verdict(noteType int, a input.KeyAction, td int64) gosu.Judgment {
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
				return gosu.Judge(Judgments, td)
			}
		}
	} else { // Head, Normal
		return gosu.Verdict(Judgments, a, td)
	}
	return gosu.Judgment{}
}

// Extra primitive in Piano mode is a count of Kools.
// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *Note, j gosu.Judgment) {
	if j == Miss {
		s.BreakCombo()
	} else {
		s.AddCombo()
	}
	s.CalcScore(gosu.Flow, j.Flow, n.Weight())
	s.CalcScore(gosu.Acc, j.Acc, n.Weight())
	if j.Window == Kool.Window {
		s.CalcScore(gosu.Extra, 1, n.Weight())
	}
	for i, judgment := range Judgments {
		if judgment.Window == j.Window {
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

// const (
// 	MaxFlowScore  = 7 * 1e5
// 	MaxAccScore   = 3 * 1e5
// 	MaxExtraScore = 1 * 1e5
// 	MaxScore      = MaxFlowScore + MaxAccScore + MaxExtraScore // 1.1m
// )

// // CalcScores and CalcTotalScore are separated for convenient debugging.
// func (s ScenePlay) CalcTotalScore() float64 {
// 	fs, as, es := s.CalcScores()
// 	return math.Ceil(fs + as + es)
// }

// // Flow, acc, extra score in order.
// func (s ScenePlay) CalcScores() (fs, as, es float64) {
// 	var (
// 		b = AccScoreFactor
// 		c = ExtraScoreFactor
// 	)
// 	if s.MaxNoteWeights == 0 {
// 		fs, as = 0, 0 // No score when no notes.
// 	} else {
// 		fs = MaxFlowScore * (s.Flows / s.MaxNoteWeights)
// 		as = MaxAccScore * math.Pow(s.Accs/s.MaxNoteWeights, b)
// 	}
// 	if s.MaxExtraWeights == 0 {
// 		es = 0
// 	} else {
// 		es = MaxExtraScore * math.Pow(s.Extras/s.MaxExtraWeights, c)
// 	}
// 	return
// }
// func (s ScenePlay) ScoreBound() float64 {
// 	if s.MaxNoteWeights == 0 {
// 		return 0 // No score when no notes.
// 	}
// 	var (
// 		b = AccScoreFactor
// 		c = ExtraScoreFactor
// 	)
// 	fr := s.MaxNoteWeights - (s.NoteWeights - s.Flows)
// 	ar := s.MaxNoteWeights - (s.NoteWeights - s.Accs)
// 	er := s.MaxExtraWeights - (s.NoteWeights - s.Extras)

// 	fs := MaxFlowScore * (fr / s.MaxNoteWeights)
// 	as := MaxAccScore * math.Pow(ar/s.MaxNoteWeights, b)
// 	es := MaxExtraScore * math.Pow(er/s.MaxExtraWeights, c)
// 	return math.Ceil(fs + as + es)
// }
