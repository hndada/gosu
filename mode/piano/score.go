package piano

import (
	"math"

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

func Verdict(t NoteType, a input.KeyAction, td int64) mode.Judgment {
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
				return mode.Judge(Judgments, td)
			}
		}
	} else { // Head, Normal
		mode.Verdict(Judgments, a, td)
	}
	return mode.Judgment{}
}

// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *PlayNote, j mode.Judgment) {
	var a = FlowScoreFactor
	if j == Miss {
		s.Combo = 0
		s.ComboCountdown = 0
	} else {
		s.Combo++
		s.ComboCountdown = MaxComboCountdown
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
	MaxFlowScore     = 7 * 1e5
	MaxAccScore      = 3 * 1e5
	MaxKoolRateScore = 1 * 1e5
	MaxScore         = MaxFlowScore + MaxAccScore + MaxKoolRateScore // 1.1m
)

func (s ScenePlay) Score() float64 {
	fs, as, rs := s.CalcScore()
	return math.Ceil(fs + as + rs)
}

// Flow, acc, kool rate score in order.
func (s ScenePlay) CalcScore() (fs, as, ks float64) {
	var (
		b = AccScoreFactor
		c = KoolRateScoreFactor
	)
	fs = MaxFlowScore * (s.Flows / s.MaxNoteWeights)
	as = MaxAccScore * math.Pow(s.Accs/s.MaxNoteWeights, b)
	ks = MaxKoolRateScore * math.Pow(s.Extras/s.MaxNoteWeights, c)
	return
}
func (s ScenePlay) ScoreBound() float64 {
	var (
		b = AccScoreFactor
		c = KoolRateScoreFactor
	)
	// tnc := float64(len(s.PlayNotes)) // Total note counts
	// nc := float64(s.MarkedNoteCount())
	// s.CurrentMaxNoteWeights - s.NoteWeights
	// kc := float64(s.JudgmentCounts[0]) // Kool counts
	fr := s.MaxNoteWeights - (s.NoteWeights - s.Flows)
	ar := s.MaxNoteWeights - (s.NoteWeights - s.Accs)
	kr := s.MaxNoteWeights - (s.NoteWeights - s.Extras)

	fs := MaxFlowScore * (fr / s.MaxNoteWeights)
	as := MaxAccScore * math.Pow(ar/s.MaxNoteWeights, b)
	rs := MaxKoolRateScore * math.Pow(kr/s.MaxNoteWeights, c)
	return math.Ceil(fs + as + rs)
}
