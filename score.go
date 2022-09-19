package gosu

import (
	"math"

	"github.com/hndada/gosu/input"
)

type Judgment struct {
	Flow   float64
	Acc    float64
	Window int64
	// Extra  bool // For distinguishing Big note at Drum mode.
}

// Is returns whether j and j2 are equal by its window size.
func (j Judgment) Is(j2 Judgment) bool { return j.Window == j2.Window }

// Valid returns whether j is not a blank judgment by its window size.
func (j Judgment) Valid() bool { return j.Window != 0 }

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }
// Verdict for normal notes, e.g., Note, Head at Piano mode.
func Verdict(js []Judgment, a input.KeyAction, td int64) Judgment {
	Miss := js[len(js)-1]
	switch {
	case td > Miss.Window:
		// Does nothing.
	case td < -Miss.Window:
		return Miss
	default: // In range
		if a == input.Hit {
			return Judge(js, td)
		}
	}
	return Judgment{}
}

func Judge(js []Judgment, td int64) Judgment {
	if td < 0 {
		td *= -1
	}
	for _, j := range js {
		if td <= j.Window {
			return j
		}
	}
	return Judgment{} // Returns None when the input is out of widest range
}

// Total score consists of 3 scores: Flow, Acc, and Kool rate score.
// Flow score is calculated with sum of Flow. Flow once named as Karma.
// Acc score is calculated with sum of Acc of judgments.
// Extra score is calculated with a sum of Extra primitive.
// Flow recovers fast when its value is low, vice versa: math.Pow(x, a); a < 1
// Acc and Extra score increase faster as each parameter approaches to max value: math.Pow(x, b); b > 1
const (
	Flow = iota
	Acc
	Extra
	Total
)

var DefaultMaxScores = [4]float64{
	7 * 1e5,
	3 * 1e5,
	1 * 1e5,
	11 * 1e5,
}

type Scorer struct {
	Flow       float64
	Combo      int
	Primitives [3]float64 // Sum of aquired primitive.
	Ratios     [3]float64
	Weights    [3]float64 // Works as current max value of note weights.
	MaxWeights [3]float64 // Works as Upper bound.

	ScoreFactors   [3]float64
	Scores         [4]float64
	ScoreBounds    [4]float64
	MaxScores      [4]float64
	JudgmentCounts []int
	MaxCombo       int
}

func NewScorer(scoreFactors [3]float64) Scorer {
	return Scorer{
		Flow:         1,
		Ratios:       [3]float64{1, 1, 1},
		ScoreFactors: scoreFactors,
		ScoreBounds:  DefaultMaxScores,
		MaxScores:    DefaultMaxScores,
		// JudgmentCounts: []int{},
	}
}

func (s *Scorer) SetMaxScores(maxScores [4]float64) {
	s.ScoreBounds = maxScores
	s.MaxScores = maxScores
}
func (s *Scorer) AddCombo() {
	s.Combo++
	if s.MaxCombo < s.Combo {
		s.MaxCombo = s.Combo
	}
}
func (s *Scorer) BreakCombo() { s.Combo = 0 }

// s.Primitives[Flow]+=math.Pow(s.Flow, a) * n.Weight()
func (s *Scorer) CalcScore(kind int, value, weight float64) {
	if kind == Flow {
		s.Flow += value * weight
		if s.Flow < 0 {
			s.Flow = 0
		} else if s.Flow > 1 {
			s.Flow = 1
		}
		s.Primitives[kind] += math.Pow(s.Flow, s.ScoreFactors[kind])
	} else {
		s.Primitives[kind] += value * weight
	}
	s.Weights[kind] += weight
	s.Ratios[kind] = s.Primitives[kind] / s.Weights[kind]

	scoreRate := s.Primitives[kind] / s.MaxWeights[kind]
	boundRate := 1 - (s.Weights[kind]-s.Primitives[kind])/s.MaxWeights[kind]
	if kind != Flow {
		scoreRate = math.Pow(scoreRate, s.ScoreFactors[kind])
		boundRate = math.Pow(boundRate, s.ScoreFactors[kind])
	}
	s.Scores[kind] = s.MaxScores[kind] * scoreRate
	s.ScoreBounds[kind] = s.MaxScores[kind] * boundRate
	s.Scores[Total] = math.Floor(Sum(s.Scores[:Total]) + 0.1)
	s.ScoreBounds[Total] = math.Floor(Sum(s.ScoreBounds[:Total]) + 0.1)
}
func Sum(vs []float64) (sum float64) {
	for _, v := range vs {
		sum += v
	}
	return
}
