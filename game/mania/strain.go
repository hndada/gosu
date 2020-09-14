package mania

import (
	"github.com/hndada/gosu/game"
	"math"
)

const holdAffectDelta = 16
const (
	holdOuterOnceBonus    = 0.08 * 3
	holdInnerOnceBonus    = 0.08 * 3
	holdInnerAdjOnceBonus = 0.04 * 3
	holdRemainBonus       = 0.03 * 3
)

const outerBonus = 0.025 * 3.5

var fingerBonus = [5]float64{1.5, 0, 1, 2, 3} // from thumb to little finger

func init() {
	for i := range fingerBonus {
		fingerBonus[i] *= outerBonus
	}
}

// todo: finger, note 불러올 때 미리 계산
// todo: memory-less 인 애들은 루프 분리
func (c *Chart) CalcStrain() {
	c.markAffect()
	c.setHoldImpacts()
	for i, n := range c.Notes {
		c.Notes[i].hand = hand(c.Keys, n.Key)
		c.Notes[i].settleAlterHand()
		c.Notes[i].baseStrain = baseStrain(c.Keys, n)
		c.Notes[i].chordPenalty = c.chordPenalty(n)
		c.Notes[i].jackBonus = c.jackBonus(n)
		c.Notes[i].trillBonus = c.trillBonus(n)
		c.Notes[i].holdBonus = c.holdBonus(n)
		c.Notes[i].calcStrain()
	}
}

// todo: time2, prev/next를 이용하면 대체 가능
func baseStrain(keys int, n Note) float64 {
	base := 1 + fingerBonus[finger(keys, n.Key)]
	if n.Type == TypeLNTail { // a tail of hold note will get partial strain
		lnDuration := float64(n.Time - n.Time2)
		base *= game.SolveY(curveTail, lnDuration)
	}
	return base
}

func (c *Chart) chordPenalty(n Note) float64 {
	var penalty float64
	// for _, idx := range tools.Neighbors(n.chord, n.Key) {
	for _, idx := range n.chord {
		if idx == noFound {
			continue
		}
		chordNote := c.Notes[idx]
		var div float64
		switch {
		case chordNote.Key == n.Key: // note itself
			continue
		case chordNote.hand == -n.hand:
			div = 2
		case chordNote.Key-n.Key == 1, chordNote.Key-n.Key == -1:
			div = 1
		default:
			div = 1.5
		}
		// if chordNote.hand == -n.hand {
		// 	continue
		// }
		time := math.Abs(float64(n.Time - chordNote.Time))
		v := game.SolveY(curveTrillChord, time)
		// keyDistance = math.Max(1, float64(tools.AbsInt(n.Key-chordNote.Key)))
		penalty += v / div
	}
	if penalty < -1 {
		penalty = -1
	}
	return penalty
}

func (c *Chart) jackBonus(n Note) float64 {
	if n.Type == TypeLNTail {
		return 0 // no jack bonus to hold note tail
	}
	if n.trillJack[n.Key] != noFound {
		jackNote := c.Notes[n.trillJack[n.Key]]
		time := float64(n.Time - jackNote.Time)
		return game.SolveY(curveJack, time)
	}
	return 0
}

func (c *Chart) trillBonus(n Note) float64 {
	// trill bonus is independent of other notes in same chord
	// a note can get trill bonus at most once per each side
	var bonus float64
	if n.Type == TypeLNTail {
		return 0 // no trill bonus to hold n tail
	}
	if n.jackBonus <= 0 {
		return 0 // only anchor gets trill bonus
	}
	// for _, idx := range tools.Neighbors(n.trillJack, n.Key) {
	for _, idx := range n.trillJack {
		if idx == noFound {
			continue
		}
		trillNote := c.Notes[idx]
		var div float64
		switch {
		case trillNote.Key == n.Key: // note itself
			continue
		case trillNote.hand == -n.hand:
			div = 2
		case trillNote.Key-n.Key == 1, trillNote.Key-n.Key == -1:
			div = 1
		default:
			div = 1.5
		}
		time := float64(n.Time - trillNote.Time)
		v := game.SolveY(curveTrillChord, time)
		// keyDistance = math.Max(1, float64(tools.AbsInt(n.Key-trillNote.Key)))
		bonus += v / div
	}
	return bonus
}

func (c *Chart) setHoldImpacts() {
	// sign in value stands for hit hand
	// holding starts: no impact
	// at end of holding: partial impact
	// other else: fully impact
	for i, ln := range c.Notes {
		if ln.Type != TypeLNHead {
			continue
		}
		j := i + 1 // notes in same chord might have lower index but they arent affected anyway
		for j < len(c.Notes) {
			n := c.Notes[j]
			elapsedTime := float64(n.Time - ln.Time)
			remainedTime := float64(ln.Time2 - n.Time)
			if elapsedTime >= holdAffectDelta {
				impact := math.Max(0, 0.5+math.Min(remainedTime, holdAffectDelta)/(2*holdAffectDelta))
				c.Notes[j].holdImpacts[ln.Key] = impact * float64(ln.hand)
				if ln.hand == alter {
					panic("still alter")
				}
				if impact == 0 { // hold note will not affect further notes
					break
				}
			}
			j++
		}
	}
}

func (c *Chart) holdBonus(n Note) float64 {
	// suppose hold notes on the other hand don't affect value
	// and no altering hand during pressing hold note
	// algorithm itself supposes playing with kb; outer fingers always have higher strain
	var bonus float64
	existOuter, existInner := false, false // for adding main bonus only once
	for holdKey, impact := range n.holdImpacts {
		if impact == 0 || !sameHand(float64(n.hand), impact) {
			continue
		}
		switch {
		case isHoldOuter(holdKey, n.Key, c.Keys):
			if !existOuter {
				bonus += holdOuterOnceBonus
			}
			existOuter = true
		case isHoldInner(holdKey, n.Key, c.Keys):
			if isHoldInnerAdj(holdKey, n.Key, c.Keys) {
				bonus += holdInnerAdjOnceBonus
			}
			if !existInner {
				bonus += holdInnerOnceBonus
			}
			existInner = true
		}
		if impact < 0 {
			impact *= -1
		}
		bonus += holdRemainBonus * fingerBonus[finger(c.Keys, holdKey)] * impact
	}
	return bonus
}

// changed from multiplying to adding
func (n *Note) calcStrain() {
	v := n.baseStrain
	v += n.trillBonus
	v += n.jackBonus
	v += n.holdBonus
	v += n.chordPenalty
	if v < 0 {
		panic("negative strain")
	}
	n.strain = v
}

func neighbors(slice []int, i int) [2]int {
	ns := [2]int{-1, -1}
	uBound := len(slice)

	var cursor, v int
	for ni, direct := range [2]int{left, right} {
		for offset := 1; ; offset++ {
			cursor = i + offset*direct
			if cursor < 0 || cursor >= uBound {
				break
			}
			v = slice[cursor]
			if v == -1 {
				continue
			}
			ns[ni] = v
			break
		}
	}
	return ns
}
