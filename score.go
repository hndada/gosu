package gosu

import (
	"math"
)

type Judgment struct {
	Flow   float64
	Acc    float64
	Window int64
}

var (
	Kool = Judgment{Flow: 0.01, Acc: 1, Window: 20}
	Cool = Judgment{Flow: 0.01, Acc: 1, Window: 45}
	Good = Judgment{Flow: 0.01, Acc: 0.25, Window: 75}
	Bad  = Judgment{Flow: 0.01, Acc: 0, Window: 110} // Todo: Flow 0.01 -> 0?
	Miss = Judgment{Flow: -1, Acc: 0, Window: 150}
)

var Judgments = []Judgment{Kool, Cool, Good, Bad, Miss}

func Verdict(t NoteType, a KeyAction, td int64) Judgment {
	if t == Tail { // Either Hold or Release when Tail is not scored
		switch {
		case td > Miss.Window:
			if a == Release {
				return Miss
			}
		case td < -Miss.Window:
			return Miss
		default: // In range
			if a == Release { // a != Hold
				return Judge(td)
			}
		}
	} else { // Head, Normal
		switch {
		case td > Miss.Window:
			// Does nothing
		case td < -Miss.Window:
			return Miss
		default: // In range
			if a == Hit {
				return Judge(td)
			}
		}
	}
	return Judgment{}
}
func Judge(td int64) Judgment {
	if td < 0 { // Absolute value
		td *= -1
	}
	for _, j := range Judgments {
		if td <= j.Window {
			return j
		}
	}
	return Judgment{} // Returns None when the input is out of widest range
}

// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *PlayNote, j Judgment) {
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
	s.FlowSum += math.Pow(s.Flow, a) * n.Weight()
	s.AccSum += j.Acc * n.Weight()
	if j.Window == Kool.Window {
		s.KoolSum += n.Weight()
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
	fs = MaxFlowScore * (s.FlowSum / s.MaxNoteWeights)
	as = MaxAccScore * math.Pow(s.AccSum/s.MaxNoteWeights, b)
	ks = MaxKoolRateScore * math.Pow(s.KoolSum/s.MaxNoteWeights, c)
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
	fr := s.MaxNoteWeights - (s.NoteWeights - s.FlowSum)
	ar := s.MaxNoteWeights - (s.NoteWeights - s.AccSum)
	kr := s.MaxNoteWeights - (s.NoteWeights - s.KoolSum)

	fs := MaxFlowScore * (fr / s.MaxNoteWeights)
	as := MaxAccScore * math.Pow(ar/s.MaxNoteWeights, b)
	rs := MaxKoolRateScore * math.Pow(kr/s.MaxNoteWeights, c)
	return math.Ceil(fs + as + rs)
}

// // MarkedNoteCount is for calculating ratio.
// func (s ScenePlay) MarkedNoteCount() int {
// 	sum := 0
// 	for _, c := range s.JudgmentCounts {
// 		sum += c
// 	}
// 	return sum
// }

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }
