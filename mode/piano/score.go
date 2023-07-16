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
	maxFlow = 50
	maxAcc  = 20
)

type Scorer struct {
	// no changes after initialization
	Mods       Mods
	Judgments  []mode.Judgment // May change by mods.
	UnitScores [3]float64

	// exported to result
	Combo          int
	Score          float64
	JudgmentCounts []int
	// Todo: FlowPoint

	// score calculation
	flow        float64
	acc         float64
	stagedNotes []*Note

	// audio
	hitSoundQueue []mode.Sample

	// draw
	worstJudgment  mode.Judgment
	isNoteHits     []bool // for drawing hit lighting
	lastKeyActions []input.KeyActionType
}

// It is separated from ScenePlay because it can be used for score simulation.
func NewScorer(c *Chart) Scorer {
	unit := 1e6 / float64(len(c.Notes))
	js := Judgments
	unitScores := [3]float64{unit * 0.7, unit * 0.3, unit * 0.1}

	return Scorer{
		// no changes after initializations
		Mods:       c.Mods,
		Judgments:  js,
		UnitScores: unitScores,

		// exported to result
		Combo: 0,
		// Accumulating floating-point numbers may result in imprecise values.
		// To ensure that the maximum score is attainable,
		// we initialize the score with a small value in advance.
		Score:          0.01,
		JudgmentCounts: make([]int, len(Judgments)),

		// score calculation
		flow:        maxFlow,
		acc:         maxAcc,
		stagedNotes: newStagedNotes(c),

		// draw
		// These are assigned in every Update.
	}
}

func newStagedNotes(c *Chart) []*Note {
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

func (s *Scorer) Update(now int32, kas []input.KeyboardAction) {
	s.hitSoundQueue = s.hitSoundQueue[:0] // flush hit sounds in the queue
	s.worstJudgment = blank
	s.isNoteHits = make([]bool, len(s.stagedNotes))
	s.flushStagedNotes(now)

	for _, ka := range kas {
		s.appendHitSounds(ka)
		s.tryJudge(ka)
	}

	// Read guarantees that it length is at least one.
	s.lastKeyActions = kas[len(kas)-1].KeyActions
}

func (s *Scorer) flushStagedNotes(now int32) {
	for k, n := range s.stagedNotes {
		for ; n != nil; n = n.Next {
			if e := n.Time - now; e >= -Miss.Window {
				break
			}

			// Tail note may remain in staged even if it is missed.
			if !n.Marked {
				s.mark(n, Miss)
			} else {
				if n.Type != Tail {
					panic("remained marked note is not Tail")
				}
			}
		}
		s.stagedNotes[k] = n
	}
}

// if n.Type != Tail
func (s *Scorer) appendHitSounds(ka input.KeyboardAction) {
	for k, n := range s.stagedNotes {
		a := ka.KeyActions[k]
		if a != input.Hit {
			continue
		}
		if n != nil {
			s.hitSoundQueue = append(s.hitSoundQueue, n.Sample)
		} else {
			defaultSample := mode.Sample{Filename: "", Volume: 0.5}
			s.hitSoundQueue = append(s.hitSoundQueue, defaultSample)
		}
	}
}

func (s *Scorer) tryJudge(ka input.KeyboardAction) {
	for k, n := range s.stagedNotes {
		if n == nil {
			continue
		}
		e := n.Time - ka.Time
		j := s.judge(n.Type, e, ka.KeyActions[k])
		if j != blank { // Comparison between two structs is possible.
			s.mark(n, j)
			if !j.Is(Miss) && n.Type != Tail {
				s.isNoteHits[k] = true
			}
		}
	}
}

func (s Scorer) judge(noteType int, e int32, a input.KeyActionType) mode.Judgment {
	switch noteType {
	case Normal, Head:
		return mode.Judge(s.Judgments, e, a)
	case Tail:
		return judgeTail(e, a)
	default:
		panic("invalid note type")
	}
}

// Either Hold or Release when Tail is not scored
func judgeTail(e int32, a input.KeyActionType) mode.Judgment {
	switch {
	case e > Miss.Window:
		if a == input.Release {
			return Miss
		}
	case e < -Miss.Window:
		return Miss
	default: // In range
		if a == input.Release { // a != Hold
			j := mode.Evaluate(Judgments, e)
			if j.Is(Cool) { // Cool at Tail goes Kool
				j = Kool
			}
			return j
		}
	}
	return blank
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

	if j == Miss {
		s.Combo = 0
		s.flow = 0
	} else { // Kool, Cool, Good
		s.Combo++
		s.flow++
		if s.flow > maxFlow {
			s.flow = maxFlow
		}

		if j.Is(Good) {
			s.acc = 0
		} else {
			s.acc++
			if s.acc > maxAcc {
				s.acc = maxAcc
			}
		}
	}

	flowScore := s.UnitScores[flow] * (s.flow / maxFlow)
	accScore := s.UnitScores[acc] * (s.acc / maxAcc)
	var extraScore float64
	if j.Is(Kool) {
		extraScore = s.UnitScores[extra]
	}

	s.Score += j.Weight * (flowScore + accScore + extraScore)
	s.addJugdmentCount(j)
	n.Marked = true

	// when Head is missed, its tail goes missed as well.
	if n.Type == Head && j.Is(Miss) {
		s.mark(n.Next, Miss)
	}

	// Tail is flushed separately at flushStagedNotes().
	if n.Type != Tail {
		s.stagedNotes[n.Key] = n.Next
	}

	if s.worstJudgment.Window < j.Window {
		s.worstJudgment = j
	}
	// Todo: Add time error meter mark
	// Todo: Use different color for error meter of Tail
}

func (s *Scorer) addJugdmentCount(j mode.Judgment) {
	for i, j2 := range s.Judgments {
		if j.Is(j2) {
			s.JudgmentCounts[i]++
			break
		}
	}
}

func (s Scorer) judgmentIndex(j mode.Judgment) int {
	for i, j2 := range s.Judgments {
		if j.Is(j2) {
			return i
		}
	}
	return len(Judgments) // blank
}
