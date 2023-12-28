package piano

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/input"
)

const (
	flow = iota
	acc
	extra
)

const (
	maxFlowFactor = 50
	maxAccFactor  = 20
)

type Judge struct {
	flowFactor float64
	accFactor  float64
	unitScores [3]float64
	score      *float64
}

// It is separated from ScenePlay because it can be used for score simulation.
func (s ScenePlay) setJudge() {
	unit := 1e6 / float64(len(s.notes.notes))
	unitScores := [3]float64{unit * 0.7, unit * 0.3, unit * 0.1}
	// Accumulating floating-point numbers may result in imprecise values.
	// To ensure that the maximum score is attainable,
	// we initialize the score with a small value in advance.
	s.score.Score = 0.01
	s.judge = Judge{
		flowFactor: maxFlowFactor,
		accFactor:  maxAccFactor,
		unitScores: unitScores,
		score:      &s.score.Score,
	}
}

func (s *Scorer) flushStagedNotes(now int32) (missed bool) { // return: for draw
	for k, n := range s.stagedNotes {
		for ; n != nil; n = n.Next {
			if e := n.Time - now; e >= -s.miss().Window {
				break
			}

			// Tail note may remain in staged even if it is missed.
			if !n.scored {
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

func (s *Scorer) tryJudge(ka input.KeyboardAction) []game.Judgment {
	js := make([]game.Judgment, len(s.stagedNotes)) // draw
	for k, n := range s.stagedNotes {
		if n == nil || n.scored {
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

func (s Scorer) judge(noteType int, e int32, a input.KeyActionType) game.Judgment {
	switch noteType {
	case Normal, Head:
		return game.Judge(s.judgments, e, a)
	case Tail:
		return s.judgeTail(e, a)
	default:
		panic("invalid note type")
	}
}

// Either Hold or Release when Tail is not scored
func (s Scorer) judgeTail(e int32, a input.KeyActionType) game.Judgment {
	switch {
	case e > s.miss().Window:
		if a == input.Release {
			return s.miss()
		}
	case e < -s.miss().Window:
		return s.miss()
	default: // In range
		if a == input.Release { // a != Hold
			j := game.Evaluate(s.judgments, e)
			if j.Is(s.cool()) { // Cool at Tail goes Kool
				j = s.kool()
			}
			return j
		}
	}
	return game.Judgment{}
}

// Todo: no getting Flow when hands off the long note
func (s *Scorer) mark(n *Note, j game.Judgment) {
	// Score consists of three parts: Flow, Acc, and Extra.
	// Ratios of Flow and Acc to their max values are multiplied to unit scores.
	// Flow drops to zero when Miss, and recovers when Kool, Cool, and Good.
	// Acc drops to zero when Miss or Good, and recovers when Kool and Cool.
	// Extra will be simply added to the score when hit Kool.

	if j == s.miss() {
		s.Combo = 0
		s.flow = 0
	} else { // Kool, Cool, Good
		s.Combo++
		s.flow++
		if s.flow > maxFlowFactor {
			s.flow = maxFlowFactor
		}

		if j.Is(s.good()) {
			s.acc = 0
		} else {
			s.acc++
			if s.acc > maxAccFactor {
				s.acc = maxAccFactor
			}
		}
	}

	flowScore := s.unitScores[flow] * (s.flow / maxFlowFactor)
	accScore := s.unitScores[acc] * (s.acc / maxAccFactor)
	var extraScore float64
	if j.Is(s.kool()) {
		extraScore = s.unitScores[extra]
	}

	s.Score += j.Weight * (flowScore + accScore + extraScore)
	s.addJugdmentCount(j)
	n.scored = true

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && j.Is(s.miss()) {
		s.mark(n.Next, s.miss())
	}

	// Tail is flushed separately at flushStagedNotes().
	if n.Type != Tail {
		s.stagedNotes[n.Key] = n.Next
	}
}
