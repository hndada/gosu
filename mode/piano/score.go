package piano

import (
	"image/color"
	"math"

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

func Verdict(t NoteType, a input.KeyAction, td int64) gosu.Judgment {
	if t == Tail { // Either Hold or Release when Tail is not scored
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

// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *PlayNote, j gosu.Judgment) {
	var a = FlowScoreFactor
	if j == Miss {
		s.Combo = 0
	} else {
		s.Combo++
	}
	s.Flow += j.Flow * n.Weight()
	if s.Flow < 0 {
		s.Flow = 0
	} else if s.Flow > 1 {
		s.Flow = 1
	}
	s.Flows += math.Pow(s.Flow, a) * n.Weight()
	s.Accs += j.Acc * n.Weight()
	if j.Window == Kool.Window {
		s.Extras += n.Weight()
	}
	for i, jk := range Judgments {
		if jk.Window == j.Window {
			s.JudgmentCounts[i]++
			break
		}
		if i == 4 {
			panic("no reach")
		}
	}
	s.NoteWeights += n.Weight()
	n.Marked = true
	if n.Type == Head && j == Miss {
		s.MarkNote(n.Next, Miss)
	}
	if n.Type != Tail {
		s.StagedNotes[n.Key] = n.Next
	}
}

// Total score consists of 3 scores: Flow, Acc, and Kool rate score.
// Flow score is calculated with sum of Flow. Flow once named as Karma.
// Acc score is calculated with sum of Acc of judgments.
// Kool rate score is calculated with a rate of Kool counts.
// Flow recovers fast when its value is low, vice versa: math.Pow(x, a); a < 1
// Acc and Kool rate score increase faster as each parameter approaches to max value: math.Pow(x, b); b > 1
const (
	MaxFlowScore  = 7 * 1e5
	MaxAccScore   = 3 * 1e5
	MaxExtraScore = 1 * 1e5
	MaxScore      = MaxFlowScore + MaxAccScore + MaxExtraScore // 1.1m
)

func (s ScenePlay) Score() float64 {
	fs, as, es := s.CalcScore()
	return math.Ceil(fs + as + es)
}

// Flow, acc, kool rate score in order.
func (s ScenePlay) CalcScore() (fs, as, es float64) {
	if s.MaxNoteWeights == 0 {
		return 0, 0, 0 // No score when no notes.
	}
	var (
		b = AccScoreFactor
		c = KoolRateScoreFactor
	)
	fs = MaxFlowScore * (s.Flows / s.MaxNoteWeights)
	as = MaxAccScore * math.Pow(s.Accs/s.MaxNoteWeights, b)
	es = MaxExtraScore * math.Pow(s.Extras/s.MaxNoteWeights, c)
	return
}
func (s ScenePlay) ScoreBound() float64 {
	if s.MaxNoteWeights == 0 {
		return 0 // No score when no notes.
	}
	var (
		b = AccScoreFactor
		c = KoolRateScoreFactor
	)
	fr := s.MaxNoteWeights - (s.NoteWeights - s.Flows)
	ar := s.MaxNoteWeights - (s.NoteWeights - s.Accs)
	er := s.MaxNoteWeights - (s.NoteWeights - s.Extras)

	fs := MaxFlowScore * (fr / s.MaxNoteWeights)
	as := MaxAccScore * math.Pow(ar/s.MaxNoteWeights, b)
	es := MaxExtraScore * math.Pow(er/s.MaxNoteWeights, c)
	return math.Ceil(fs + as + es)
}
