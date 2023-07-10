package piano

import (
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

var (
	Kool  = mode.Judgment{Window: 20, Weight: 1}
	Cool  = mode.Judgment{Window: 40, Weight: 1}
	Good  = mode.Judgment{Window: 80, Weight: 0.5}
	Miss  = mode.Judgment{Window: 120, Weight: 0}
	blank = mode.Judgment{}
)

var Judgments = []mode.Judgment{Kool, Cool, Good, Miss}

const (
	MaxFlow = 50
	MaxAcc  = 20
)

// Score consists of three parts: Flow, Acc, and Extra.
// Ratios of Flow and Acc to their max values are multiplied to unit scores.
// Flow drops to zero when Miss, and recovers when Kool, Cool, and Good.
// Acc drops to zero when Miss or Good, and recovers when Kool and Cool.
// Extra will be simply added to the score when hit Kool.
const (
	Flow = iota
	Acc
	Extra
)

type Scorer struct {
	Mods      interface{}
	Judgments []mode.Judgment // May change by mods

	Combo          int
	Score          float64
	UnitScores     [3]float64
	JudgmentCounts []int

	Flow float64
	Acc  float64

	notes  []*Note
	Staged []*Note

	// Todo: FlowPoint

	worstJudgment mode.Judgment
	isHits        []bool
}

func NewScorer(c *Chart) Scorer {
	unit := 1e6 / float64(len(c.Notes))
	unitScores := [3]float64{unit * 0.7, unit * 0.3, unit * 0.1}
	js := Judgments

	return Scorer{
		Mods:      c.Mods,
		Judgments: js,

		UnitScores:     unitScores,
		JudgmentCounts: make([]int, len(Judgments)),

		Flow: MaxFlow,
		Acc:  MaxAcc,

		notes:  c.Notes,
		Staged: newStaged(c),
	}
}
func newStaged(c *Chart) []*Note {
	staged := make([]*Note, c.KeyCount)
	for k := range staged {
		for _, n := range c.Notes {
			if k == n.Key {
				staged[n.Key] = n
				break
			}
		}
	}
	return staged
}

func (s *Scorer) Check(ka input.KeyboardAction) {
	s.worstJudgment = blank
	s.isHits = make([]bool, len(s.Staged)) // for drawing hit lighting
	for k, n := range s.Staged {
		if n == nil {
			continue
		}

		timeError := n.Time - ka.Time

		// Flush marked tail notes.
		if n.Marked {
			if n.Type != Tail {
				panic("marked yet remained note is not Tail")
			}
			// Keep Tail staged until near ends.
			if timeError < Miss.Window {
				s.Staged[k] = n.Next
			}
			// continue // I think no continue is right.
		}

		j := Judge(n.Type, timeError, ka.Action[k])
		if j != blank { // Comparison between two structs is possible.
			s.Mark(n, j)
			if s.worstJudgment.Window < j.Window {
				s.worstJudgment = j
			}
			if !j.Is(Miss) { // && n.Type != Head
				s.isHits[k] = true
			}
			// Todo: Add time error meter mark
			// Todo: Use different color for error meter of Tail
		}
	}
}

func Judge(noteType int, timeError int32, a input.KeyActionType) mode.Judgment {
	switch noteType {
	case Normal, Head:
		return mode.Judge(Judgments, timeError, a)
	case Tail:
		return judgeTail(timeError, a)
	}
	return blank
}

// Either Hold or Release when Tail is not scored
func judgeTail(timeError int32, a input.KeyActionType) mode.Judgment {
	switch {
	case timeError > Miss.Window:
		if a == input.Release {
			return Miss
		}
	case timeError < -Miss.Window:
		return Miss
	default: // In range
		if a == input.Release { // a != Hold
			j := mode.Evaluate(Judgments, timeError)
			if j.Is(Cool) { // Cool at Tail goes Kool
				j = Kool
			}
			return j
		}
	}
	return blank
}

// Todo: no getting Flow when hands off the long note
func (s *Scorer) Mark(n *Note, j mode.Judgment) {
	if j == Miss {
		s.Combo = 0
		s.Flow = 0
	} else { // Kool, Cool, Good
		s.Combo++
		s.Flow++
		if s.Flow > MaxFlow {
			s.Flow = MaxFlow
		}

		if j.Is(Good) {
			s.Acc = 0
		} else {
			s.Acc++
			if s.Acc > MaxAcc {
				s.Acc = MaxAcc
			}
		}
	}

	flow := s.UnitScores[Flow] * (s.Flow / MaxFlow)
	acc := s.UnitScores[Acc] * (s.Acc / MaxAcc)
	var extra float64
	if j.Is(Kool) {
		extra = s.UnitScores[Extra]
	}

	s.Score += j.Weight * (flow + acc + extra)
	s.addJugdmentCount(j)
	n.Marked = true

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && j.Is(Miss) {
		s.Mark(n.Next, Miss)
	}

	// Tail is flushed at Check().
	if n.Type != Tail {
		s.Staged[n.Key] = n.Next
	}
}

func (s *Scorer) addJugdmentCount(j mode.Judgment) {
	for i, j2 := range Judgments {
		if j.Is(j2) {
			s.JudgmentCounts[i]++
			break
		}
	}
}
