package piano

import "github.com/hndada/gosu/game"

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
	blank = iota - 1
	kool
	cool
	good
	miss
)

// keysFocusNoteIndex []int // targets of judging
type Scorer struct {
	notes      Notes
	judgments  []game.Judgment // To use game.Judge, slice is preferred.
	missWindow int32
	units      [3]float64
	factors    [3]float64
	maxFactors [3]float64

	judgmentBuffer []int
	JudgmentCounts [4]int
	Combo          float64
	Score          float64
}

func NewScorer(ns Notes, js []game.Judgment) (s Scorer) {
	s.notes = ns
	s.judgments = js
	s.missWindow = js[miss].Window
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

func (s *Scorer) update(ka game.KeyboardAction) {
	s.judgmentBuffer = s.judgmentBuffer[:0]
	s.markKeysUntouchedNote(ka.Time)
	for k, ni := range s.notes.keysFocus {
		if ni == len(s.notes.notes) {
			continue
		}
		n := s.notes.notes[ni]

		if ka.KeysAction[k] == game.Hit {
			s.notes.sampleBuffer = append(s.notes.sampleBuffer, n.Sample)
		}
		if n.scored {
			continue
		}
		e := n.Time - ka.Time
		if ji := s.judge(n.Type, e, ka.KeysAction[k]); ji != blank {
			s.markNote(ni, ji)
		}
	}
}

// marks the untouched note as missed.
func (s Scorer) markKeysUntouchedNote(now int32) {
	for k, start := range s.notes.keysFocus {
		n := s.notes.notes[start]
		for ni := start; ni < len(s.notes.notes); ni = n.next {
			if e := n.Time - now; e >= -s.missWindow {
				break
			}
			if n.scored {
				// Tail note may be focused even after being marked.
				// Other types of note should not.
				if n.Type == Tail {
					s.notes.keysFocus[k] = n.next
				} else {
					panic("remained marked note is not Tail")
				}
			} else {
				s.markNote(ni, miss)
			}
		}
	}
}

func (s Scorer) judge(nt int, e int32, a game.KeyActionType) int {
	switch nt {
	case Normal, Head:
		return game.Judge(s.judgments, e, a)
	case Tail:
		return s.judgeTail(e, a)
	default:
		panic("invalid note type")
	}
}

// Either Hold or Released when Tail is not scored
func (s Scorer) judgeTail(e int32, at game.KeyActionType) int {
	switch {
	case e > s.missWindow:
		if at == game.Released {
			return miss
		}
	case e < -s.missWindow:
		return miss
	default: // In range
		if at == game.Released { // a != Hold
			ji := game.Evaluate(s.judgments, e)
			// Beware: Cool at Tail note goes Kool.
			if ji == cool {
				ji = kool
			}
			return ji
		}
	}
	return blank
}

// Todo: no getting Flow when hands off the long note
func (s *Scorer) markNote(ni int, ji int) {
	n := s.notes.notes[ni]
	j := s.judgments[ji]
	switch ji {
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
	s.notes.notes[ni].scored = true
	s.JudgmentCounts[ji]++

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && ji == miss {
		s.markNote(n.next, miss)
	}

	// Tail is flushed separately at updateKeysFocusNote().
	if n.Type != Tail {
		s.notes.keysFocus[n.Key] = n.next
	}
	s.judgmentBuffer = append(s.judgmentBuffer, ji)
}

func (s *Scorer) incrementFactor(fi int) {
	s.factors[fi] = min(s.factors[fi]+1, s.maxFactors[fi])
}

// func (s Scorer) keysFocusNote() []Note {
// 	kn := make([]Note, s.notes.keyCount)
// 	for k, ni := range s.notes.keysFocus {
// 		if ni == len(s.notes.notes) {
// 			continue
// 		}
// 		kn[k] = s.notes.notes[ni]
// 	}
// 	return kn
// }
