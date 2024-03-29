package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/game"
)

// There are three kinds of factors: Flow, Acc, and Extra.
// Ratios of Flow and Acc to their max values are multiplied to unit scores.
// Flow drops to zero when Miss, and recovers when Kool, Cool, and Good.
// Acc drops to zero when Miss or Good, and recovers when Kool and Cool.
// Extra will be simply added to the score when hit Kool.
const (
	flow = iota
	acc
	extra
)

const (
	kool game.JudgmentKind = iota
	cool
	good
	miss
	blank
)

// Todo: FlowPoint (kind of HP)
type Scorer struct {
	Notes
	game.Judgments
	keysJudgmentKind []game.JudgmentKind
	Combo            int
	units            [3]float64
	factors          [3]float64
	maxFactors       [3]float64
	Score            float64
}

func NewScorer(ns Notes, mods Mods) (s Scorer) {
	s.Notes = ns
	js := mods.DefaultJudgments()
	s.Judgments = game.NewJudgments(js)

	unit := 1e6 / float64(len(ns.notes))
	s.units = [3]float64{unit * 0.7, unit * 0.3, unit * 0.1}
	s.maxFactors = [3]float64{50, 20, 1}
	s.factors = s.maxFactors

	// Accumulating floating-point numbers may result in imprecise values.
	// To ensure that the maximum score is attainable,
	// we initialize the score with a small value in advance.
	s.Score = 0.01
	return
}

// update returns the indices of the judgments.
func (s *Scorer) update(ka game.KeyboardAction) {
	s.keysJudgmentKind = make([]game.JudgmentKind, s.keyCount)
	for k := range s.keysJudgmentKind {
		s.keysJudgmentKind[k] = blank
	}

	s.markKeysUntouchedNote(ka.Time)

	for k, ni := range s.keysFocus {
		if ni == s.none {
			continue
		}
		n := s.notes[ni]
		if ka.KeysAction[k] == game.Hit {
			s.sampleBuffer = append(s.sampleBuffer, n.Sample)
		}
		if n.scored {
			continue
		}
		e := n.Time - ka.Time
		if jk := s.judge(n.Kind, e, ka.KeysAction[k]); jk != blank {
			s.markNote(ni, jk)
		}
	}
}

// marks the untouched note as missed.
func (s Scorer) markKeysUntouchedNote(now int32) {
	for k, lowest := range s.keysFocus {
		for ni := lowest; ni < len(s.notes); ni = s.notes[ni].next {
			n := s.notes[ni]
			if e := n.Time - now; s.IsTooLate(e) {
				break
			}
			if n.scored {
				// Tail note may be focused even after being marked.
				// Other types of note should not.
				if n.Kind == Tail {
					s.keysFocus[k] = n.next
				} else {
					panic("remained marked note is not Tail")
				}
			} else {
				s.markNote(ni, miss)
			}
		}
	}
}

func (s Scorer) judge(nk NoteKind, e int32, a game.KeyActionType) game.JudgmentKind {
	switch nk {
	case Normal, Head:
		return s.Judge(e, a)
	case Tail:
		return s.judgeTail(e, a)
	default:
		panic("invalid note type")
	}
}

// Either Hold or Released when Tail is not scored
func (s Scorer) judgeTail(e int32, at game.KeyActionType) game.JudgmentKind {
	switch {
	case s.IsTooEarly(e):
		if at == game.Released {
			return miss
		}
	case s.IsTooLate(e):
		return miss
	case s.IsInRange(e):
		if at == game.Released {
			// Cool goes Kool when judging Tail note.
			jk := s.Evaluate(e)
			if jk == cool {
				jk = kool
			}
			return jk
		}
	}
	return blank
}

// Todo: no getting Flow when hands off the long note
func (s *Scorer) markNote(ni int, jk game.JudgmentKind) {
	n := s.notes[ni]
	j := s.Judgments.Judgments[jk]
	switch jk {
	case kool:
		s.Combo++
		s.incrementFactor(flow)
		s.incrementFactor(acc)
		s.incrementFactor(extra)
	case cool:
		s.Combo++
		s.incrementFactor(flow)
		s.incrementFactor(acc)
		s.factors[extra] = 0
	case good:
		s.Combo++
		s.incrementFactor(flow)
		s.factors[acc] = 0
		s.factors[extra] = 0
	case miss:
		s.Combo = 0
		s.factors[flow] = 0
		s.factors[acc] = 0
		s.factors[extra] = 0
	}

	for i, unit := range s.units {
		ratio := s.factors[i] / s.maxFactors[i]
		score := j.Weight * (ratio * unit)
		s.Score += score
	}
	s.notes[ni].scored = true
	s.Judgments.Counts[jk]++

	// when Head is missed, its tail goes missed as well.
	if n.Kind == Head && jk == miss {
		s.markNote(n.next, miss)
	}

	// Tail is flushed separately at markKeysUntouchedNote.
	if n.Kind != Tail {
		s.keysFocus[n.Key] = n.next
	}

	s.keysJudgmentKind[n.Key] = jk
}

func (s *Scorer) incrementFactor(fi int) {
	s.factors[fi] = min(s.factors[fi]+1, s.maxFactors[fi])
}

func (s Scorer) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	f(&b, "Score: %.0f \n", s.Score)
	f(&b, "Combo: %d\n", s.Combo)
	f(&b, "Flow: %.0f/%.0f\n", s.factors[flow], s.maxFactors[flow])
	f(&b, " Acc: %.0f/%.0f\n", s.factors[acc], s.maxFactors[acc])
	f(&b, "Judgment counts: %v\n", s.Judgments.Counts)
	f(&b, "\n")
	return b.String()
}
