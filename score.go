package main

import (
	"math"
)

func Verdict(n *PlayNote, td int64, a KeyAction) Judgment {
	if n.Type == Tail { // Either Hold or Release when Tail is not scored
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

func (s *ScenePlay) Score(n *PlayNote, j Judgment) {
	if j == Miss {
		s.Combo = 0
	} else {
		s.Combo++
	}

	s.Karma += j.Karma
	if s.Karma < 0 {
		s.Karma = 0
	} else if s.Karma > 1 {
		s.Karma = 1
	}
	s.KarmaSum += s.Karma

	for i, jk := range Judgments {
		if jk.Window == j.Window {
			s.JudgmentCounts[i]++
			break
		}
		if i == 4 {
			panic("no reach")
		}
	}

	n.Scored = true
	if n.Type == Head && j == Miss {
		s.Score(n.Next, Miss)
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
func (s ScenePlay) CurrentScore() int {
	const (
		a = 0.5
		b = 5
		c = 2
	)
	nc := float64(len(s.PlayNotes))    // Total note counts
	kc := float64(s.JudgmentCounts[0]) // Kool counts
	var accSum float64
	for j, c := range s.JudgmentCounts {
		accSum += Judgments[j].Acc * float64(c)
	}
	ks := 7 * 1e5 * (s.KarmaSum / nc)
	as := 3 * 1e5 * math.Pow(accSum/nc, b)
	rs := 1 * 1e5 * math.Pow(kc/nc, c)
	return int(math.Ceil(ks + as + rs))
}
