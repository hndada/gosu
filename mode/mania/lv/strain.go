package lv

import (
	"github.com/hndada/gosu/internal/tools"
	"github.com/hndada/gosu/mode/mania"
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

var fingerBonus = [5]float64{
	1.5 * outerBonus, // thumb
	0,                // index finger
	1 * outerBonus,
	2 * outerBonus,
	3 * outerBonus}

func CalcStrain(ns []mania.Note, keymode int) {
	markAffect(ns)
	mania.hand(ns, keymode)

	setStrainBase(ns, keymode)
	calcChordPanelty(ns)
	calcJackBonus(ns)
	calcTrillBonus(ns)
	setHoldImpacts(ns)
	calcHoldBonus(ns, keymode)

	setStrain(ns)
}

func setStrainBase(ns []mania.Note, keymode int) {
	var base, lnDuration float64
	for i, n := range ns {
		base = 1 + fingerBonus[mania.finger(keymode, n.Key)]
		if n.NoteType == NtHoldTail { // a tail of hold note will get partial strain
			lnDuration = float64(n.Time - n.OpponentTime)
			base *= tools.SolveY(curveTail, lnDuration)
		}
		ns[i].strainBase = base
	}
}

func calcChordPanelty(ns []mania.Note) {
	var chordNote mania.Note
	var timeDelta, v, div float64
	for i, n := range ns {
		// for _, idx := range tools.Neighbors(n.chord, n.Key) {
		for _, idx := range n.chord {
			if idx == noFound {
				continue
			}
			chordNote = ns[idx]
			switch {
			case chordNote.Key == n.Key: // note itself
				continue
			case chordNote.hand == -n.hand:
				div = 2
			case tools.AbsInt(chordNote.Key-n.Key) == 1:
				div = 1
			default:
				div = 1.5
			}
			// if chordNote.hand == -n.hand {
			// 	continue
			// }
			timeDelta = math.Abs(float64(n.Time - chordNote.Time))
			v = tools.SolveY(curveTrillChord, timeDelta)
			// keyDistance = math.Max(1, float64(tools.AbsInt(n.Key-chordNote.Key)))
			ns[i].chordPenalty += v / div
		}
		if ns[i].chordPenalty < -1 {
			ns[i].chordPenalty = -1
		}
	}
}

func calcJackBonus(ns []mania.Note) {
	var jackNote mania.Note
	var timeDelta float64
	for i, n := range ns {
		if n.NoteType == NtHoldTail {
			continue // no jack bonus to hold note tail
		}
		if n.trillJack[n.Key] != noFound {
			jackNote = ns[n.trillJack[n.Key]]
			timeDelta = float64(n.Time - jackNote.Time)
			ns[i].jackBonus = tools.SolveY(curveJack, timeDelta)
		}
	}
}

func calcTrillBonus(ns []mania.Note) {
	// trill bonus is independent of other notes in same chord
	// a note can get trill bonus at most once per each side
	var trillNote mania.Note
	var timeDelta, v, div float64
	for i, n := range ns {
		if n.NoteType == NtHoldTail {
			continue // no trill bonus to hold n tail
		}
		if n.jackBonus <= 0 {
			continue // only anchor gets trill bonus
		}
		// for _, idx := range tools.Neighbors(n.trillJack, n.Key) {
		for _, idx := range n.trillJack {
			if idx == noFound {
				continue
			}
			trillNote = ns[idx]
			switch {
			case trillNote.Key == n.Key: // note itself
				continue
			case trillNote.hand == -n.hand:
				div = 2
			case tools.AbsInt(trillNote.Key-n.Key) == 1:
				div = 1
			default:
				div = 1.5
			}
			timeDelta = float64(n.Time - trillNote.Time)
			v = tools.SolveY(curveTrillChord, timeDelta)
			// keyDistance = math.Max(1, float64(tools.AbsInt(n.Key-trillNote.Key)))
			ns[i].trillBonus += v / div
		}
	}
}

func setHoldImpacts(ns []mania.Note) {
	// sign in value stands for hit hand
	// holding starts: no impact
	// at end of holding: partial impact
	// other else: fully impact

	var affected mania.Note
	var affectedIdx int
	var elapsedTime, remainedTime float64
	var impact float64

	for i, ln := range ns {
		if ln.NoteType != NtHoldHead {
			continue
		}
		affectedIdx = i + 1 // notes in same chord might have lower index but they arent affected anyway
		for affectedIdx < len(ns) {
			affected = ns[affectedIdx]
			elapsedTime = float64(affected.Time - ln.Time)
			remainedTime = float64(ln.OpponentTime - affected.Time)

			if elapsedTime >= holdAffectDelta {
				impact = math.Max(0, 0.5+math.Min(remainedTime, holdAffectDelta)/(2*holdAffectDelta))
				ns[affectedIdx].holdImpacts[ln.Key] = impact * float64(ln.hand)
				if ln.hand == mania.alter {
					panic("still alter")
				}
				if impact == 0 { // hold note will not affect further notes
					break
				}
			}
			affectedIdx++
		}
	}
}

func calcHoldBonus(ns []mania.Note, keymode int) {
	// suppose hold notes on the other hand don't affect value
	// and no altering hand during pressing hold note
	// algorithm itself supposes playing with kb; outer fingers always have higher strain
	var bonus float64
	var existOuter, existInner bool
	for i, n := range ns {
		bonus = 0
		existOuter, existInner = false, false // for adding main bonus only once
		for holdKey, impact := range n.holdImpacts {
			if impact == 0 || !mania.sameHand(float64(n.hand), impact) {
				continue
			}
			switch {
			case mania.isHoldOuter(holdKey, n.Key, keymode):
				if !existOuter {
					bonus += holdOuterOnceBonus
				}
				existOuter = true
			case mania.isHoldInner(holdKey, n.Key, keymode):
				if mania.isHoldInnerAdj(holdKey, n.Key, keymode) {
					bonus += holdInnerAdjOnceBonus
				}
				if !existInner {
					bonus += holdInnerOnceBonus
				}
				existInner = true
			}
			bonus += holdRemainBonus * fingerBonus[mania.finger(keymode, holdKey)] * math.Abs(impact)
		}
		ns[i].holdBonus = bonus
	}
}

// changed from multiplying to adding
func setStrain(ns []mania.Note) {
	var strain float64
	for i, n := range ns {
		strain = n.strainBase
		strain += n.trillBonus
		strain += n.jackBonus
		strain += n.holdBonus
		strain += n.chordPenalty
		if strain < 0 {
			strain = 0
			// panic("negative strain")
		}
		ns[i].Strain = strain
	}
}
