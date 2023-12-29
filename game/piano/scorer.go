package piano

import "github.com/hndada/gosu/game"

// e.g., Judgment counts
type ScorerRes struct {
}

type ScorerOpts struct {
	keyCount int
}

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

type Scorer struct {
	keyCount           int
	notes              []Note
	keysFocusNoteIndex []int // targets of judging

	judgments  []game.Judgment // To use game.Judge, slice is preferred.
	missWindow int32

	units      [3]float64
	factors    [3]float64
	maxFactors [3]float64

	JudgmentCounts [4]int
	Combo          float64
	Score          float64
}

func NewScorer(res ScorerRes, opts ScorerOpts, notes []Note, js [4]game.Judgment) (s Scorer) {
	s.keyCount = opts.keyCount
	s.notes = notes
	s.keysFocusNoteIndex = s.newKeysFocusNoteIndex(opts.keyCount, notes)
	s.judgments = js[:]
	s.missWindow = js[miss].Window

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

// Todo: move to Mods
func (Scorer) DefaultJudgments() []game.Judgment {
	return []game.Judgment{
		{Window: 20, Weight: 1},
		{Window: 40, Weight: 1},
		{Window: 80, Weight: 0.5},
		{Window: 120, Weight: 0},
	}
}

func (Scorer) newKeysFocusNoteIndex(keyCount int, ns []Note) []int {
	// Initialize with the max none value: length of notes.
	kfni := make([]int, keyCount)
	for k := range kfni {
		kfni[k] = len(ns)
	}

	for k := range kfni {
		for i, n := range ns {
			if k == n.Key {
				kfni[n.Key] = i
				break
			}
		}
	}
	return kfni
}

// markKeysUntouchedNote marks the untouched note as missed.
func (s *Scorer) markKeysUntouchedNote(now int32, kji []int) {
	for k, start := range s.keysFocusNoteIndex {
		n := s.notes[start]
		for ni := start; ni < len(s.notes); ni = n.next {
			if e := n.Time - now; e >= -s.missWindow {
				break
			}

			// Tail note may remain in keysFocusNote even if it is missed.
			if !n.scored {
				s.mark(ni, miss)
				kji[k] = miss
			} else {
				if n.Type != Tail {
					panic("remained marked note is not Tail")
				}
				// Other notes go flushed at mark().
				s.keysFocusNoteIndex[k] = n.next
			}
		}
	}
}

func (s Scorer) keysFocusNote() []Note {
	kn := make([]Note, s.keyCount)
	for k, ni := range s.keysFocusNoteIndex {
		if ni == len(s.notes) {
			continue
		}
		kn[k] = s.notes[ni]
	}
	return kn
}

// tryJudge requires index of the notes for marking them.
func (s *Scorer) tryJudge(ka game.KeyboardAction, kji []int) {
	for k, ni := range s.keysFocusNoteIndex {
		n := s.notes[ni]
		if ni == len(s.notes) || n.scored {
			continue
		}
		e := n.Time - ka.Time
		ji := s.judge(n.Type, e, ka.KeysAction[k])
		if ji != blank {
			s.mark(ni, ji)
		}
		kji[k] = ji
	}
	return
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
func (s *Scorer) mark(ni int, ji int) {
	n := s.notes[ni]
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
	s.notes[ni].scored = true
	s.JudgmentCounts[ji]++

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && ji == miss {
		s.mark(n.next, miss)
	}

	// Tail is flushed separately at updateKeysFocusNote().
	if n.Type != Tail {
		s.keysFocusNoteIndex[n.Key] = n.next
	}
}

func (s *Scorer) incrementFactor(fi int) {
	s.factors[fi] = min(s.factors[fi]+1, s.maxFactors[fi])
}
