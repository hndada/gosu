package gosu

import (
	"math"
)

type Judgment struct {
	Karma  float64
	Acc    float64
	Window int64
}

var (
	Kool = Judgment{Karma: 0.01, Acc: 1, Window: 20}
	Cool = Judgment{Karma: 0.01, Acc: 1, Window: 45}
	Good = Judgment{Karma: 0.01, Acc: 0.25, Window: 75}
	Bad  = Judgment{Karma: 0.01, Acc: 0, Window: 110} // Todo: Karma 0.01 -> 0?
	Miss = Judgment{Karma: -1, Acc: 0, Window: 150}
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

// Todo: no getting karma when hands off the long note
func (s *ScenePlay) MarkNote(n *PlayNote, j Judgment) {
	var a = KarmaScoreFactor
	if j == Miss {
		s.Combo = 0
		s.ComboCountdown = 0
	} else {
		s.Combo++
		s.ComboCountdown = MaxComboCountdown
	}
	s.Karma += j.Karma
	if s.Karma < 0 {
		s.Karma = 0
	} else if s.Karma > 1 {
		s.Karma = 1
	}
	s.KarmaSum += math.Pow(s.Karma, a)
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

// Total score consists of 3 scores: Karma, Acc, and Ratio score.
// Karma score is calculated with sum of Karma.
// Acc score is calculated with sum of Acc of judgments.
// Ratio score is calculated with a ratio of Kool count.
// Karma recovers fast when its value is low, vice versa: math.Pow(x, a); a < 1
// Acc and ratio score increase faster as each parameter approaches to max value: math.Pow(x, b); b > 1
// Karma, Acc, Ratio, Total max score is 700k, 300k, 100k, 1100k each.
const (
	ScoreMaxKarma = 7 * 1e5
	ScoreMaxAcc   = 3 * 1e5
	ScoreMaxRatio = 1 * 1e5
	ScoreMaxTotal = ScoreMaxKarma + ScoreMaxAcc + ScoreMaxRatio
)

// Karma, acc, ratio score in order.
func (s ScenePlay) CalcScore() (ks, as, rs float64) {
	var (
		b = AccScoreFactor
		c = RatioScoreFactor
	)
	nc := float64(len(s.PlayNotes))    // Total note counts
	kc := float64(s.JudgmentCounts[0]) // Kool counts

	ks = ScoreMaxKarma * (s.KarmaSum / nc)
	as = ScoreMaxAcc * math.Pow(s.AccSum/nc, b)
	rs = ScoreMaxRatio * math.Pow(kc/nc, c)
	return
}
func (s ScenePlay) Score() float64 {
	ks, as, rs := s.CalcScore()
	return math.Ceil(ks + as + rs)
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
