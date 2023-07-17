package piano

import (
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

const (
	Kool = iota
	Cool
	Good
	Miss
)

func DefaultJudgments() []mode.Judgment {
	return []mode.Judgment{
		{Window: 20, Weight: 1},
		{Window: 40, Weight: 1},
		{Window: 80, Weight: 0.5},
		{Window: 120, Weight: 0},
	}
}

const (
	maxFlow = 50
	maxAcc  = 20
)

// Todo: FlowPoint
type Scorer struct {
	stagedNotes []*Note // ScenePlay has same slice.
	flow        float64
	acc         float64
	judgments   []mode.Judgment // May change by mods.
	unitScores  [3]float64

	Combo          int
	Score          float64
	JudgmentCounts []int
}

// It is separated from ScenePlay because it can be used for score simulation.
func (s ScenePlay) newScorer() Scorer {
	unit := 1e6 / float64(len(s.Notes))
	js := DefaultJudgments()
	unitScores := [3]float64{unit * 0.7, unit * 0.3, unit * 0.1}

	return Scorer{
		stagedNotes: s.stagedNotes,
		flow:        maxFlow,
		acc:         maxAcc,
		judgments:   js,
		unitScores:  unitScores,

		Combo: 0,
		// Accumulating floating-point numbers may result in imprecise values.
		// To ensure that the maximum score is attainable,
		// we initialize the score with a small value in advance.
		Score:          0.01,
		JudgmentCounts: make([]int, len(js)),
	}
}

func (s Scorer) kool() mode.Judgment { return s.judgments[Kool] }
func (s Scorer) cool() mode.Judgment { return s.judgments[Cool] }
func (s Scorer) good() mode.Judgment { return s.judgments[Good] }
func (s Scorer) miss() mode.Judgment { return s.judgments[Miss] }

func (s *Scorer) flushStagedNotes(now int32) (missed bool) { // return: for draw
	for k, n := range s.stagedNotes {
		for ; n != nil; n = n.Next {
			if e := n.Time - now; e >= -s.miss().Window {
				break
			}

			// Tail note may remain in staged even if it is missed.
			if !n.Marked {
				s.mark(n, s.miss())
				missed = true
			} else {
				if n.Type != Tail {
					panic("remained marked note is not Tail")
				}
				// Other notes go flushed at mark().
				s.stagedNotes[k] = n.Next
			}
		}
	}
	return missed
}

func (s *Scorer) tryJudge(ka input.KeyboardAction) []mode.Judgment {
	js := make([]mode.Judgment, len(s.stagedNotes)) // draw
	for k, n := range s.stagedNotes {
		if n == nil || n.Marked {
			continue
		}
		e := n.Time - ka.Time
		j := s.judge(n.Type, e, ka.KeyActions[k])
		if !j.IsBlank() {
			s.mark(n, j)
		}
		js[k] = j
	}
	return js
}

func (s Scorer) judge(noteType int, e int32, a input.KeyActionType) mode.Judgment {
	switch noteType {
	case Normal, Head:
		return mode.Judge(s.judgments, e, a)
	case Tail:
		return s.judgeTail(e, a)
	default:
		panic("invalid note type")
	}
}

// Either Hold or Release when Tail is not scored
func (s Scorer) judgeTail(e int32, a input.KeyActionType) mode.Judgment {
	switch {
	case e > s.miss().Window:
		if a == input.Release {
			return s.miss()
		}
	case e < -s.miss().Window:
		return s.miss()
	default: // In range
		if a == input.Release { // a != Hold
			j := mode.Evaluate(s.judgments, e)
			if j.Is(s.cool()) { // Cool at Tail goes Kool
				j = s.kool()
			}
			return j
		}
	}
	return mode.Judgment{}
}

// Todo: no getting Flow when hands off the long note
func (s *Scorer) mark(n *Note, j mode.Judgment) {
	// Score consists of three parts: Flow, Acc, and Extra.
	// Ratios of Flow and Acc to their max values are multiplied to unit scores.
	// Flow drops to zero when Miss, and recovers when Kool, Cool, and Good.
	// Acc drops to zero when Miss or Good, and recovers when Kool and Cool.
	// Extra will be simply added to the score when hit Kool.
	const (
		flow = iota
		acc
		extra
	)

	if j == s.miss() {
		s.Combo = 0
		s.flow = 0
	} else { // Kool, Cool, Good
		s.Combo++
		s.flow++
		if s.flow > maxFlow {
			s.flow = maxFlow
		}

		if j.Is(s.good()) {
			s.acc = 0
		} else {
			s.acc++
			if s.acc > maxAcc {
				s.acc = maxAcc
			}
		}
	}

	flowScore := s.unitScores[flow] * (s.flow / maxFlow)
	accScore := s.unitScores[acc] * (s.acc / maxAcc)
	var extraScore float64
	if j.Is(s.kool()) {
		extraScore = s.unitScores[extra]
	}

	s.Score += j.Weight * (flowScore + accScore + extraScore)
	s.addJugdmentCount(j)
	n.Marked = true

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && j.Is(s.miss()) {
		s.mark(n.Next, s.miss())
	}

	// Tail is flushed separately at flushStagedNotes().
	if n.Type != Tail {
		s.stagedNotes[n.Key] = n.Next
	}
}

func (s *Scorer) addJugdmentCount(j mode.Judgment) {
	for i, j2 := range s.judgments {
		if j.Is(j2) {
			s.JudgmentCounts[i]++
			break
		}
	}
}

func (s Scorer) judgmentIndex(j mode.Judgment) int {
	for i, j2 := range s.judgments {
		if j.Is(j2) {
			return i
		}
	}
	return len(s.judgments) // blank
}
