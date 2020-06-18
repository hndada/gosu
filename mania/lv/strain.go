package mania

import (
	"github.com/hndada/gosu/mania"
	"math"

	"github.com/hndada/gosu/game/tools"
)

const holdAffectDelta = 16
const (
	holdOuterOnceBonus    = 0.08
	holdInnerOnceBonus    = 0.08
	holdInnerAdjOnceBonus = 0.04
	holdRemainBonus       = 0.03
)

const outerBonus = 0.025

var fingerBonus [5]float64 = [5]float64{1.5 * outerBonus, 0, 1 * outerBonus, 2 * outerBonus, 3 * outerBonus}

func (beatmap *ManiaBeatmap) CalcStrain() {
	beatmap.markAffect()
	beatmap.determineAlters()

	beatmap.setStrainBase()
	beatmap.calcChordPanelty()
	beatmap.calcJackBonus()
	beatmap.calcTrillBonus()
	beatmap.setHoldImpacts()
	beatmap.calcHoldBonus()
	beatmap.setStrain()
}

func (beatmap *ManiaBeatmap) setStrainBase() {
	var base, holdNoteDuration float64
	for i, note := range beatmap.Notes {
		base = 1 + fingerBonus[finger(note.Key, beatmap.Keymode)]
		if note.NoteType == mania.NtHoldTail { // a tail of hold note will get partial strain
			holdNoteDuration = float64(note.Time - note.OpponentTime)
			base *= tools.SolveY(beatmap.Curves["HoldTail"], holdNoteDuration)
		}
		beatmap.Notes[i].StrainBase = base
	}
}

func (beatmap *ManiaBeatmap) calcChordPanelty() {
	var chordNote mania.ManiaNote
	var timeDelta, v, keyDistance float64
	for i, note := range beatmap.Notes {
		for _, idx := range tools.Neighbors(note.chord, note.Key) {
			if idx == noFound {
				continue
			}
			chordNote = beatmap.Notes[idx]
			if chordNote.hand == -note.hand {
				continue
			}
			timeDelta = math.Abs(float64(note.Time - chordNote.Time))
			v = tools.SolveY(beatmap.Curves["TrillChord"], timeDelta)
			keyDistance = math.Max(1, float64(tools.AbsInt(note.Key-chordNote.Key)))
			beatmap.Notes[i].ChordPenalty += v / keyDistance
		}
		// beatmap.Notes[i].ChordPenalty = math.Max(-0.8, penalty)
	}
}

func (beatmap *ManiaBeatmap) calcJackBonus() {
	var jackNote mania.ManiaNote
	var timeDelta float64
	for i, note := range beatmap.Notes {
		if note.NoteType == mania.NtHoldTail {
			continue // no jack bonus to hold note tail
		}
		if note.trillJack[note.Key] != noFound {
			jackNote = beatmap.Notes[note.trillJack[note.Key]]
			timeDelta = float64(note.Time - jackNote.Time)
			beatmap.Notes[i].JackBonus = tools.SolveY(beatmap.Curves["Jack"], timeDelta)
		}
	}
}

func (beatmap *ManiaBeatmap) calcTrillBonus() {
	// trill bonus is independent of other notes in same chord
	// a note can get trill bonus at most once per each side
	var trillNote mania.ManiaNote
	var timeDelta, v, keyDistance float64
	for i, note := range beatmap.Notes {
		if note.NoteType == mania.NtHoldTail {
			continue // no trill bonus to hold note tail
		}
		if note.JackBonus <= 0 {
			continue // only anchor gets trill bonus
		}
		for _, idx := range tools.Neighbors(note.trillJack, note.Key) {
			if idx == noFound {
				continue
			}
			trillNote = beatmap.Notes[idx]
			if trillNote.hand == -note.hand {
				continue
			}
			timeDelta = float64(note.Time - trillNote.Time)
			v = tools.SolveY(beatmap.Curves["TrillChord"], timeDelta)
			keyDistance = math.Max(1, float64(tools.AbsInt(note.Key-trillNote.Key)))
			beatmap.Notes[i].TrillBonus += v / keyDistance
		}
	}
}

func (beatmap *ManiaBeatmap) setHoldImpacts() {
	// sign in value stands for hit hand
	// holding starts: no impact
	// at end of holding: partial impact
	// other else: fully impact
	var affected mania.ManiaNote
	var affectedIdx int
	var elapsedTime, remainedTime float64
	var impact float64
	for i, holdNote := range beatmap.Notes {
		if holdNote.NoteType != mania.NtHoldHead {
			continue
		}
		affectedIdx = i + 1
		for affectedIdx < len(beatmap.Notes) {
			affected = beatmap.Notes[affectedIdx]
			elapsedTime = float64(affected.Time - holdNote.Time)
			remainedTime = float64(holdNote.OpponentTime - affected.Time)
			if elapsedTime >= holdAffectDelta {
				impact = math.Max(0, 0.5+math.Min(remainedTime, holdAffectDelta)/(2*holdAffectDelta))
				beatmap.Notes[affectedIdx].holdImpacts[holdNote.Key] = impact * float64(holdNote.hand)
				if holdNote.hand == alter {
					panic(&tools.ValError{"Hold impact hand", tools.Itoa(holdNote.hand), tools.ErrFlow})
				}
				if impact == 0 { // hold note will not affect further notes
					break
				}
			}
			affectedIdx++
		}
	}
}

func (beatmap *ManiaBeatmap) calcHoldBonus() {
	// suppose hold notes on the other hand don't affect value
	// and no altering hand during pressing hold note
	// algorithm itself supposes playing with kb; outer fingers always have higher strain
	var bonus float64
	var existOuter, existInner bool
	for i, note := range beatmap.Notes {
		bonus = 0
		existOuter, existInner = false, false // for adding main bonus only once
		for holdKey, impact := range note.holdImpacts {
			if impact == 0 || !isSameHand(note.hand, int(impact)) {
				continue
			}
			switch {
			case isHoldOuter(holdKey, note.Key, beatmap.Keymode):
				if !existOuter {
					bonus += holdOuterOnceBonus
				}
				existOuter = true
			case isHoldInner(holdKey, note.Key, beatmap.Keymode):
				if isHoldInnerAdj(holdKey, note.Key, beatmap.Keymode) {
					bonus += holdInnerAdjOnceBonus
				}
				if !existInner {
					bonus += holdInnerOnceBonus
				}
				existInner = true
			}
			bonus += holdRemainBonus * fingerBonus[finger(holdKey, beatmap.Keymode)] * math.Abs(impact)
		}
		beatmap.Notes[i].HoldBonus = bonus
	}
}

func (beatmap *ManiaBeatmap) setStrain() {
	var strain float64
	for i, note := range beatmap.Notes {
		strain = note.StrainBase
		strain *= 1 + note.TrillBonus
		strain *= 1 + note.JackBonus
		strain *= 1 + note.HoldBonus
		strain *= 1 + note.ChordPenalty
		if strain < 0 {
			panic(&tools.ValError{"Strain", tools.Ftoa(strain), tools.ErrFlow})
		}
		beatmap.Notes[i].Strain = strain
	}
}
