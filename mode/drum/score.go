package drum

import (
	"image/color"
	"math"

	"github.com/hndada/gosu"
)

// Todo: Tick judgment should be bound to MaxScaledBPM (->280)
// Todo: let them put custom window
var (
	Cool = gosu.Judgment{Flow: 0.01, Acc: 1, Window: 25}
	Good = gosu.Judgment{Flow: 0.01, Acc: 0.25, Window: 60}
	Miss = gosu.Judgment{Flow: -1, Acc: 0, Window: 100}
)

var Judgments = []gosu.Judgment{Cool, Good, Miss}

var JudgmentColors = []color.NRGBA{
	gosu.ColorCool,
	gosu.ColorGood,
	gosu.ColorBad,
}

// When hit big notes only with one press, the note gives half the score only.
// For example, when hit a Big note by one press with Good, it will gives 0.25 * 0.5 = 0.125.
// No Flow decrease for hitting Big note by one press.
// When one side of judgment is Cool, Good on the other hand, overall judgment of Big note goes Good.
// In other word, to get Cool at Big note, you have to hit it Cool with both sides.

// Roll / Shake note does not affect on Flow / Acc scores.
// For example, a Roll / Shake only chart has extra score only: max score is 100k.

func IsColorHit(color int, hits []bool) bool {
	if color == Red && (hits[1] || hits[2]) {
		return true
	}
	if color == Blue && (hits[0] || hits[3]) {
		return true
	}
	return false
}
func IsOtherColorHit(color int, hits []bool) bool {
	if color == Red && (hits[0] || hits[3]) {
		return true
	}
	if color == Blue && (hits[1] || hits[2]) {
		return true
	}
	return false
}

func Verdict(n *Note, hits []bool, td int64) (j gosu.Judgment, big bool) {
	if IsOtherColorHit(n.Color, hits) {
		j = Miss
	}
	return
}

//	func (s ScenePlay) IsColorHit(color int, hits []bool) bool {
//		var keys []int
//		switch color {
//		case Red:
//			keys = []int{1, 2}
//		case Blue:
//			keys = []int{0, 3}
//		}
//		for k, hit := range hits {
//			if hit {
//				return true
//			}
//		}
//		return false
//	}
//
// Todo: no getting Flow when hands off the long note
func (s *ScenePlay) MarkNote(n *Note, j gosu.Judgment) {
	var a = FlowScoreFactor
	if j == Miss {
		s.Combo = 0
	} else {
		s.Combo++
	}
	s.Flow += j.Flow * n.Weight()
	if s.Flow < 0 {
		s.Flow = 0
	} else if s.Flow > 1 {
		s.Flow = 1
	}
	s.Flows += math.Pow(s.Flow, a) * n.Weight()
	s.Accs += j.Acc * n.Weight()
	for i, jk := range Judgments {
		if jk.Window == j.Window {
			s.JudgmentCounts[i]++
			break
		}
	}
	s.NoteWeights += n.Weight()
	n.Marked = true
	s.StagedNote = s.StagedNote.Next
}
func (s *ScenePlay) MarkTick() {
	if j.Window == Cool.Window {
		s.Extras += n.Weight()
	}
}
