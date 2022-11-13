package drum

import (
	"image/color"

	"github.com/hndada/gosu/framework/input"
	"github.com/hndada/gosu/game/mode"
)

// Todo: let users use custom windows
// Todo: need to find best value for DotHitWindow
var (
	Cool               = mode.Judgment{Flow: 0.01, Acc: 1, Window: 25}
	Good               = mode.Judgment{Flow: 0.01, Acc: 0.25, Window: 60}
	Miss               = mode.Judgment{Flow: -1, Acc: 0, Window: 100}
	DotHitWindow int64 = 25
)
var Judgments = []mode.Judgment{Cool, Good, Miss}

var JudgmentColors = []color.NRGBA{
	mode.ColorCool,
	mode.ColorGood,
	mode.ColorBad,
}

const (
	Cools = iota // Stands for Cool counts.
	Goods
	Misses
	CoolPartials
	GoodPartials
	TickHits
	TickDrops
)

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
// Todo: judge for Big note when judgment goes different depending on selecting 2 key actions
const MaxBigHitDuration = 25

func (s *ScenePlay) UpdateKeyActions() {
	var hits [4]bool
	for k := range hits {
		if s.KeyLogger.KeyAction(k) == input.Hit {
			hits[k] = true
			s.LastHitTimes[k] = s.Now
		}
	}
	for color, keys := range [][]int{{1, 2}, {0, 3}} {
		if hits[keys[0]] || hits[keys[1]] {
			if hits[keys[0]] && s.Now-s.LastHitTimes[keys[1]] < MaxBigHitDuration ||
				hits[keys[1]] && s.Now-s.LastHitTimes[keys[0]] < MaxBigHitDuration {
				s.KeyActions[color] = Big
			} else {
				s.KeyActions[color] = Regular
			}
		} else {
			s.KeyActions[color] = SizeNone
		}
	}
	// for color, a := range s.KeyActions {
	// 	if a != SizeNone {
	// 		fmt.Printf("%d: %s at color %s\n", s.time, []string{"regular", "hit"}[a], []string{"red", "blue"}[color])
	// 	}
	// }
}

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

func VerdictNote(n *Note, actions [2]int, td int64) (j mode.Judgment, big bool) {
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
	j = mode.Verdict(Judgments, input.Hit, td)
	if n.Size == Big && actions[n.Color] == Big {
		big = true
	}
	// fmt.Println(n.Time, n.Size, n.Color, actions, j, big)
	return
}
func (s *ScenePlay) MarkNote(n *Note, j mode.Judgment, big bool) {
	if j == Miss {
		s.BreakCombo()
	} else {
		s.AddCombo()
	}
	s.CalcScore(mode.Flow, j.Flow, n.Weight())
	if n.Size == Big && !big {
		j.Acc /= 2
	}
	s.CalcScore(mode.Acc, j.Acc, n.Weight())
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

// Roll affects only at Extra score.
func (s *ScenePlay) MarkDot(dot *Dot, marked int) {
	switch marked {
	case DotHit:
		s.JudgmentCounts[TickHits]++
		s.CalcScore(mode.Extra, 1, dot.Weight())
		dot.Marked = DotHit
	case DotMiss:
		s.JudgmentCounts[TickDrops]++
		s.CalcScore(mode.Extra, 0, dot.Weight())
		dot.Marked = DotMiss
	}
	if marked != DotReady {
		s.StagedDot = s.StagedDot.Next
	}
}
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
}

// Shake affects only at Extra score.
func (s *ScenePlay) MarkShake(shake *Note, flush bool) {
	if flush {
		remained := shake.Tick - shake.HitTick
		s.JudgmentCounts[TickDrops] += remained
		s.CalcScore(mode.Extra, 0, shake.Weight()*float64(remained)/float64(shake.Tick))
	} else {
		shake.HitTick++
		s.JudgmentCounts[TickHits]++
		s.CalcScore(mode.Extra, 1, shake.Weight()/float64(shake.Tick))
	}
	if flush || shake.HitTick == shake.Tick {
		shake.Marked = true
		s.StagedShake = s.StagedShake.Next
	}
}

// If a chart has not enough Shake and Roll, max score of Extra will shrink.
// Margin will be distributed to Flow and Acc score by 7:3.
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
	nws := s.MaxWeights[mode.Flow]
	ews := s.MaxWeights[mode.Extra]

	extraMax := mode.DefaultMaxScores[mode.Extra]
	extraRate := ExtraScoreRate(nws, ews)
	extraScore := extraMax * extraRate
	extraRemained := extraMax * (1 - extraRate)

	s.MaxScores[mode.Flow] += extraRemained * 0.7
	s.MaxScores[mode.Acc] += extraRemained * 0.3
	s.MaxScores[mode.Extra] = extraScore
	if nws == 0 {
		s.MaxScores[mode.Extra] += s.MaxScores[mode.Flow] + s.MaxScores[mode.Acc]
		s.MaxScores[mode.Flow] = 0
		s.MaxScores[mode.Acc] = 0
	}
	s.Scorer.SetMaxScores(s.MaxScores)
}
