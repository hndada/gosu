package gosu

import (
	"math"
	"time"

	"github.com/hndada/gosu/input"
)

type Judgment struct {
	Flow   float64
	Acc    float64
	Window int64
	Extra  bool // For distinguishing Big note at Drum mode.
}

// Verdict for normal notes: Note, Head at long note.
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
	if td < 0 { // Absolute value.
		td *= -1
	}
	for _, j := range js {
		if td <= j.Window {
			return j
		}
	}
	return Judgment{} // Returns None when the input is out of widest range
}

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }

// const (
//
//	ScoreFactorFlow = iota
//	ScoreFactorAcc
//	ScoreFactorExtra
//
// )
// const (
//
//	ScoreTotal = iota
//	ScoreFlow
//	ScoreAcc
//	ScoreExtra
//
// )

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

// var Primitives = [3]int{Flow, Acc, Extra}

var DefaultMaxScores = [4]float64{
	7 * 1e5,
	3 * 1e5,
	1 * 1e5,
	11 * 1e5,
}

// Score's fields are temporary.
// Result's fields will be exposed at SceneResult.
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
	// Result
	// Flows          float64 // Sum of aquired Flow.
	// Accs           float64 // Sum of aquired Acc value.
	// NoteWeights    float64 // Works as current max value of note weights.
	// MaxNoteWeights float64 // Works as Upper bound.

	// Extras          float64 // Kool rate in Piano mode.
	// ExtraWeights    float64 // Same with NoteWeights in Piano mode.
	// MaxExtraWeights float64 //  Works as Upper bound.

	// FlowScore  float64
	// AccScore   float64
	// ExtraScore float64
	// TotalScore float64
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

	scoreRate := (s.Primitives[kind] / s.MaxWeights[kind])
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
func (s *Scorer) AddCombo() {
	s.Combo++
	if s.MaxCombo < s.Combo {
		s.MaxCombo = s.Combo
	}
}
func (s *Scorer) BreakCombo() { s.Combo = 0 }
func (s Scorer) NewResult(md5 [16]byte) Result {
	return Result{
		MD5:            md5,
		PlayedTime:     time.Now(),
		ScoreFactors:   s.ScoreFactors,
		Scores:         s.Scores,
		JudgmentCounts: s.JudgmentCounts,
		MaxCombo:       s.MaxCombo,
	}
}

// func (s *Scorer) CalcScores() {
// 	for _, kind := range []int{Flow, Acc, Extra} {
// 		if s.MaxWeights[kind] == 0 {
// 			continue
// 		}
// 		scoreRate := (s.Primitives[kind] / s.MaxWeights[kind])
// 		boundRate := s.MaxWeights[kind] - (s.Weights[kind] - s.Primitives[kind])
// 		if kind != Flow {
// 			scoreRate = math.Pow(scoreRate, s.ScoreFactors[kind])
// 			boundRate = math.Pow(boundRate, s.ScoreFactors[kind])
// 		}
// 		score := s.MaxScores[kind] * scoreRate
// 		bound := s.MaxScores[kind] * boundRate
// 		s.Scores[kind] = score
// 		s.ScoreBounds[kind] = bound
// 		s.Scores[Total] += score
// 		s.ScoreBounds[Total] += bound
// 	}
// 	s.Scores[Total] = math.Ceil(s.Scores[Total])
// 	s.ScoreBounds[Total] = math.Ceil(s.ScoreBounds[Total])
// }
