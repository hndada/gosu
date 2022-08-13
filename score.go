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
	s.Flow += j.Flow
	if s.Flow < 0 {
		s.Flow = 0
	} else if s.Flow > 1 {
		s.Flow = 1
	}
	s.FlowSum += math.Pow(s.Flow, a)
	s.AccSum += j.Acc
	for i, jk := range Judgments {
		if jk.Window == j.Window {
			s.JudgmentCounts[i]++
			break
		}
		if i == 4 {
			panic("no reach")
		}
	}
	n.Marked = true
	if n.Type == Head && j == Miss {
		s.MarkNote(n.Next, Miss)
	}
	if n.Type != Tail {
		s.StagedNotes[n.Key] = n.Next
	}
}

// Total score consists of 3 scores: Flow, Acc, and Ratio score.
// Flow score is calculated with sum of Flow. Flow once named as Karma.
// Acc score is calculated with sum of Acc of judgments.
// Ratio score is calculated with a ratio of Kool count.
// Flow recovers fast when its value is low, vice versa: math.Pow(x, a); a < 1
// Acc and ratio score increase faster as each parameter approaches to max value: math.Pow(x, b); b > 1
// Flow, Acc, Ratio, Total max score is 700k, 300k, 100k, 1100k each.
const (
	ScoreMaxFlow  = 7 * 1e5
	ScoreMaxAcc   = 3 * 1e5
	ScoreMaxRatio = 1 * 1e5
	ScoreMaxTotal = ScoreMaxFlow + ScoreMaxAcc + ScoreMaxRatio
)

// Flow, acc, ratio score in order.
func (s ScenePlay) CalcScore() (fs, as, rs float64) {
	var (
		b = AccScoreFactor
		c = RatioScoreFactor
	)
	nc := float64(len(s.PlayNotes))    // Total note counts
	kc := float64(s.JudgmentCounts[0]) // Kool counts

	fs = ScoreMaxFlow * (s.FlowSum / nc)
	as = ScoreMaxAcc * math.Pow(s.AccSum/nc, b)
	rs = ScoreMaxRatio * math.Pow(kc/nc, c)
	return
}
func (s ScenePlay) Score() float64 {
	fs, as, rs := s.CalcScore()
	return math.Ceil(fs + as + rs)
}

// MarkedNoteCount is for calculating ratio.
func (s ScenePlay) MarkedNoteCount() int {
	sum := 0
	for _, c := range s.JudgmentCounts {
		sum += c
	}
	return sum
}

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }
