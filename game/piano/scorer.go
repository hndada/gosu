package piano

import (
	"time"

	"github.com/hndada/gosu/game"
)

const (
	flow = iota
	acc
	extra
)

type Scorer struct {
	notes     []Note
	judgments [4]game.Judgment

	units      [3]float64
	factors    [3]float64
	maxFactors [3]float64

	JudgmentCount [4]int
	Combo         float64
	Score         float64
}

func NewScorer(notes []Note, js [4]game.Judgment) (s Scorer) {
	s.notes = notes
	s.judgments = js

	unit := 1e6 / float64(len(notes))
	s.units = [3]float64{unit * 0.7, unit * 0.3, unit * 0.1}
	s.maxFactors = [3]float64{50, 20, 1}
	s.factors = s.maxFactors

	// Accumulating floating-point numbers may result in imprecise values.
	// To ensure that the maximum score is attainable,
	// we initialize the score with a small value in advance.
	s.Score = 0.01
	return
}

// return: for draw
func (s *Play) flushStagedNotes(t time.Duration) (missed bool) {
	for k, i := range s.notes.stagedList {
		n := s.notes.notes[i]
		for ; n != nil; n = n.Next {
			if e := n.Time - t; e >= -s.judgment.miss().Window {
				break
			}

			// Tail note may remain in stagedList even if it is missed.
			if !n.scored {
				s.mark(i, s.judgment.miss())
				missed = true
			} else {
				if n.Type != Tail {
					panic("remained marked note is not Tail")
				}
				// Other notes go flushed at mark().
				s.notes.stagedList[k] = n.Next
			}
		}
	}
	return missed
}

func (s *Play) tryJudge(ka game.KeyboardAction) []game.Judgment {
	js := make([]game.Judgment, len(s.notes.stagedList)) // draw
	for k, n := range s.notes.stagedList {
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

func (s Play) judge(noteType int, e int32, a game.KeyActionType) game.Judgment {
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
func (s Play) judgeTail(e int32, a game.KeyActionType) game.Judgment {
	switch {
	case e > s.judgment.miss().Window:
		if a == game.Release {
			return s.judgment.miss()
		}
	case e < -s.judgment.miss().Window:
		return s.judgment.miss()
	default: // In range
		if a == game.Release { // a != Hold
			j := game.Evaluate(s.judgments, e)
			if j.Is(s.judgment.cool()) { // Cool at Tail goes Kool
				j = s.judgment.kool()
			}
			return j
		}
	}
	return game.Judgment{}
}

// Todo: no getting Flow when hands off the long note
func (s *Play) mark(nidx int, j game.Judgment) {
	// Score consists of three parts: Flow, Acc, and Extra.
	// Ratios of Flow and Acc to their max values are multiplied to unit scores.
	// Flow drops to zero when Miss, and recovers when Kool, Cool, and Good.
	// Acc drops to zero when Miss or Good, and recovers when Kool and Cool.
	// Extra will be simply added to the score when hit Kool.

	if j == s.judgment.miss() {
		s.Combo = 0
		s.flow = 0
	} else { // Kool, Cool, Good
		s.Combo++
		s.flow++
		if s.flow > maxFlowFactor {
			s.flow = maxFlowFactor
		}

		if j.Is(s.judgment.good()) {
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
	if j.Is(s.judgment.kool()) {
		extraScore = s.unitScores[extra]
	}

	s.score.Score += j.Weight * (flowScore + accScore + extraScore)
	s.addJugdmentCount(j)
	s.notes.notes[nidx].scored = true

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && j.Is(s.judgment.miss()) {
		s.mark(n.Next, s.judgment.miss())
	}

	// Tail is flushed separately at flushStagedNotes().
	if n.Type != Tail {
		s.stagedListNotes[n.Key] = n.Next
	}
}
