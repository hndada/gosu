package drum

import (
	"image/color"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

var Hit = input.Hit

// Todo: Tick judgment should be bound to MaxScaledBPM (->280)
// Todo: let them put custom window
var (
	Cool               = gosu.Judgment{Flow: 0.01, Acc: 1, Window: 25}
	Good               = gosu.Judgment{Flow: 0.01, Acc: 0.25, Window: 60}
	Miss               = gosu.Judgment{Flow: -1, Acc: 0, Window: 100}
	DotHitWindow int64 = 50 // Todo: need to find best value
)

var Judgments = []gosu.Judgment{Cool, Good, Miss}

var JudgmentColors = []color.NRGBA{
	gosu.ColorCool,
	gosu.ColorGood,
	gosu.ColorBad,
}

const (
	CoolSum = iota
	GoodSum
	MissSum

	CoolBig
	CoolPartial
	GoodBig
	GoodPartial
	MissBig

	DotHit
	DotDrop
	ShakeHit
	ShakeDrop
)

var JudgmentCountKinds = []string{
	"CoolSum", "GoodSum", "MissSum",

	"CoolBig", "CoolPartial",
	"GoodBig", "GoodPartial",
	"MissBig",

	"DotHit", "DotDrop",
	"ShakeHit", "ShakeDrop",
}

// When hit big notes only with one press, the note gives half the score only.
// For example, when hit a Big note by one press with Good, it will gives 0.25 * 0.5 = 0.125.
// No Flow decrease for hitting Big note by one press.
// When one side of judgment is Cool, Good on the other hand, overall judgment of Big note goes Good.
// In other word, to get Cool at Big note, you have to hit it Cool with both sides.

// Roll / Shake note does not affect on Flow / Acc scores.
// For example, a Roll / Shake only chart has extra score only: max score is 100k.

//	func IsColorHit(color int, hits [4]bool) bool {
//		if color == Red && (hits[1] || hits[2]) {
//			return true
//		}
//		if color == Blue && (hits[0] || hits[3]) {
//			return true
//		}
//		return false
//	}
//
//	func IsOtherColorHit(color int, hits [4]bool) bool {
//		if color == Red && (hits[0] || hits[3]) {
//			return true
//		}
//		if color == Blue && (hits[1] || hits[2]) {
//			return true
//		}
//		return false
//	}
func IsColorHit(as [2]int, color int) bool {
	return as[color-1] != None
}
func IsOtherColorHit(as [2]int, color int) bool {
	switch color {
	case Red:
		return IsColorHit(as, Blue)
	case Blue:
		return IsColorHit(as, Red)
	}
	return false
}

func VerdictNote(n *Note, as [2]int, td int64) (j gosu.Judgment, big bool) {
	if IsOtherColorHit(as, n.Color) {
		return Miss, false
	}
	if !IsColorHit(as, n.Color) {
		return
	}
	j = gosu.Verdict(Judgments, input.Hit, td)
	if n.Size == Big && as[n.Color-1] == Big {
		big = true
	}
	return
}
func VerdictDot(dot *Dot, as [2]int, td int64) (marked, hit bool) {
	switch {
	case td < -DotHitWindow:
		return true, false
	case td < DotHitWindow:
		if as[0] != None || as[1] != None {
			return true, true
		}
	}
	return false, false
}
func VerdictShake(shake *Note, as [2]int, waitingColor int) (nextColor int) {
	const (
		red = iota
		blue
	)
	if waitingColor == Red || waitingColor == None {
		if as[red] != None {
			return Blue
		}
	}
	if waitingColor == Blue || waitingColor == None {
		if as[blue] != None {
			return Red
		}
	}
	return waitingColor
	// switch waitingColor {
	// case None:
	// 	if as[red] != None {
	// 		return Blue
	// 	}
	// 	if as[blue] != None {
	// 		return Red
	// 	}
	// case Red:
	// 	if as[red] != None {
	// 		return Blue
	// 	}
	// case Blue:
	// 	if as[blue] != None {
	// 		return Red
	// 	}
	// }
}

// func IsRedHit(as [2][2]bool) bool {
// 	return as[Red-1][Regular] || as[Red-1][Big]
// }
// func IsBlueHit(as [2][2]bool) bool {
// 	return as[Blue-1][Regular] || as[Blue-1][Big]
// }

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
func (s *ScenePlay) MarkNote(n *Note, j gosu.Judgment, big bool) {
	if j == Miss {
		s.BreakCombo()
	} else {
		s.AddCombo()
	}
	s.CalcScore(gosu.Flow, j.Flow, n.Weight())
	if n.Size == Big && !big {
		j.Acc /= 2
	}
	s.CalcScore(gosu.Acc, j.Acc, n.Weight())
	switch j.Window {
	case Cool.Window:
		s.JudgmentCounts[CoolSum]++
		if n.Size == Big {
			if big {
				s.JudgmentCounts[CoolBig]++
			} else {
				s.JudgmentCounts[CoolPartial]++
			}
		}
	case Good.Window:
		s.JudgmentCounts[GoodSum]++
		if n.Size == Big {
			if big {
				s.JudgmentCounts[GoodBig]++
			} else {
				s.JudgmentCounts[GoodPartial]++
			}
		}
	case Miss.Window:
		s.JudgmentCounts[MissSum]++
		if n.Size == Big {
			s.JudgmentCounts[MissBig]++
		}
	}
	n.Marked = true
	s.StagedNote = s.StagedNote.Next
}

func (s *ScenePlay) MarkDot(dot *Dot, hit bool) {
	if hit {
		s.JudgmentCounts[DotHit]++
		s.CalcScore(gosu.Extra, 1, dot.Weight())
	} else {
		s.JudgmentCounts[DotDrop]++
		s.CalcScore(gosu.Extra, 0, dot.Weight())
	}
	dot.Marked = true
	s.StagedDot = s.StagedDot.Next
}
func (s *ScenePlay) MarkShake(shake *Note, flush bool) {
	if flush {
		remained := shake.Tick - shake.HitTick
		s.JudgmentCounts[ShakeDrop] += remained
		s.CalcScore(gosu.Extra, 0, shake.Weight()*float64(remained)/float64(shake.Tick))
		shake.Marked = true
		s.StagedShake = s.StagedShake.Next
	} else {
		shake.HitTick++
		s.JudgmentCounts[ShakeHit]++
		s.CalcScore(gosu.Extra, 1, shake.Weight()/float64(shake.Tick))
		if shake.HitTick == shake.Tick {
			shake.Marked = true
			s.StagedShake = s.StagedShake.Next
		}
	}
}
