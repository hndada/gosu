package drum

import (
	"image/color"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

// Todo: Tick judgment should be bound to MaxScaledBPM (->280)
// Todo: let them put custom window
var (
	Cool               = gosu.Judgment{Flow: 0.01, Acc: 1, Window: 25}
	Good               = gosu.Judgment{Flow: 0.01, Acc: 0.25, Window: 60}
	Miss               = gosu.Judgment{Flow: -1, Acc: 0, Window: 100}
	DotHitWindow int64 = 25 // Todo: need to find best value
)

var Judgments = []gosu.Judgment{Cool, Good, Miss}

var JudgmentColors = []color.NRGBA{
	gosu.ColorCool,
	gosu.ColorGood,
	gosu.ColorBad,
}

const BigHitTimeDifferenceBound = 20

// const (
// 	CoolSum = iota
// 	GoodSum
// 	MissSum

// 	CoolBig
// 	CoolPartial
// 	GoodBig
// 	GoodPartial
// 	MissBig

//	DotHit
//	DotDrop
//	ShakeHit
//	ShakeDrop
//
// )
const (
	Cools = iota // Stands for Cool counts.
	Goods
	Misses
	CoolPartials
	GoodPartials
	TickHits
	TickDrops
)

func ExtraScoreRate(nws, ews float64) float64 {
	const factor = 10 * 1.5
	if nws == 0 {
		if ews == 0 {
			return 0
		}
		return 1
	}
	rate := factor * ews / nws
	if rate > 1 {
		rate = 1
	}
	return rate
}
func (s *ScenePlay) SetMaxScores() {
	nws := s.MaxWeights[gosu.Flow]
	ews := s.MaxWeights[gosu.Extra]
	// extraScore := s.MaxScores[gosu.Extra]
	extraMax := gosu.DefaultMaxScores[gosu.Extra]
	extraRate := ExtraScoreRate(nws, ews)
	extraScore := extraMax * extraRate
	extraRemained := extraMax * (1 - extraRate)
	s.MaxScores[gosu.Flow] += extraRemained * 0.7
	s.MaxScores[gosu.Acc] += extraRemained * 0.3
	s.MaxScores[gosu.Extra] = extraScore
	if nws == 0 {
		s.MaxScores[gosu.Extra] += s.MaxScores[gosu.Flow] + s.MaxScores[gosu.Acc]
		s.MaxScores[gosu.Flow] = 0
		s.MaxScores[gosu.Acc] = 0
	}
	s.Scorer.SetMaxScores(s.MaxScores)
}

// var JudgmentCountKinds = []string{
// 	"CoolSum", "GoodSum", "MissSum",

// 	"CoolBig", "CoolPartial",
// 	"GoodBig", "GoodPartial",
// 	"MissBig",

//		"DotHit", "DotDrop",
//		"ShakeHit", "ShakeDrop",
//	}
var JudgmentCountKinds = []string{
	"Cools", "Goods", "Misses",
	"CoolPartials", "GoodPartials",
	"TickHits", "TickDrops",
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
func IsColorHit(actions [2]int, color int) bool {
	return actions[color] != SizeNone
}
func IsOtherColorHit(actions [2]int, color int) bool {
	switch color {
	case Red:
		return IsColorHit(actions, Blue)
	case Blue:
		return IsColorHit(actions, Red)
	}
	return false
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
func VerdictNote(n *Note, actions [2]int, td int64) (j gosu.Judgment, big bool) {
	if td > Miss.Window {
		return
	}
	if td < -Miss.Window {
		return Miss, false
	}
	if IsOtherColorHit(actions, n.Color) {
		return Miss, false
	}
	if !IsColorHit(actions, n.Color) {
		return
	}
	j = gosu.Verdict(Judgments, input.Hit, td)
	if n.Size == Big && actions[n.Color] == Big {
		big = true
	}
	return
}

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
		s.JudgmentCounts[Cools]++
		if n.Size == Big && !big {
			s.JudgmentCounts[CoolPartials]++
		}
	case Good.Window:
		s.JudgmentCounts[Goods]++
		if n.Size == Big && !big {
			s.JudgmentCounts[GoodPartials]++
		}
	case Miss.Window:
		s.JudgmentCounts[Misses]++
	}
	n.Marked = true
	s.StagedNote = s.StagedNote.Next
}
func VerdictDot(dot *Dot, actions [2]int, td int64) (marked int) {
	switch {
	case td < -DotHitWindow:
		return DotMiss
	case td < DotHitWindow:
		if actions[Red] != SizeNone || actions[Blue] != SizeNone {
			return DotHit
		}
	}
	return DotReady
}

//	func VerdictDot(dot *Dot, as [2]int, td int64) (marked, hit bool) {
//		switch {
//		case td < -DotHitWindow:
//			return true, false
//		case td < DotHitWindow:
//			if as[0] != None || as[1] != None {
//				return true, true
//			}
//		}
//		return false, false
//	}
func (s *ScenePlay) MarkDot(dot *Dot, marked int) {
	switch marked {
	case DotHit:
		s.JudgmentCounts[TickHits]++
		s.CalcScore(gosu.Extra, 1, dot.Weight())
		dot.Marked = DotHit
	case DotMiss:
		s.JudgmentCounts[TickDrops]++
		s.CalcScore(gosu.Extra, 0, dot.Weight())
		dot.Marked = DotMiss
	}
	if marked != DotReady {
		s.StagedDot = s.StagedDot.Next
	}
}

//	func (s *ScenePlay) MarkDot(dot *Dot, hit bool) {
//		if hit {
//			s.JudgmentCounts[TickHits]++
//			s.CalcScore(gosu.Extra, 1, dot.Weight())
//			dot.Marked = DotHit
//		} else {
//			s.JudgmentCounts[TickDrops]++
//			s.CalcScore(gosu.Extra, 0, dot.Weight())
//			dot.Marked = DotMiss
//		}
//		s.StagedDot = s.StagedDot.Next
//	}
func VerdictShake(shake *Note, actions [2]int, waitingColor int) (nextColor int) {
	if waitingColor == Red || waitingColor == ColorNone {
		if actions[Red] != SizeNone {
			return Blue
		}
	}
	if waitingColor == Blue || waitingColor == ColorNone {
		if actions[Blue] != SizeNone {
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
func (s *ScenePlay) MarkShake(shake *Note, flush bool) {
	if flush {
		remained := shake.Tick - shake.HitTick
		s.JudgmentCounts[TickDrops] += remained
		s.CalcScore(gosu.Extra, 0, shake.Weight()*float64(remained)/float64(shake.Tick))
		shake.Marked = true
		s.StagedShake = s.StagedShake.Next
	} else {
		shake.HitTick++
		s.JudgmentCounts[TickHits]++
		s.CalcScore(gosu.Extra, 1, shake.Weight()/float64(shake.Tick))
		if shake.HitTick == shake.Tick {
			shake.Marked = true
			s.StagedShake = s.StagedShake.Next
		}
	}
}
