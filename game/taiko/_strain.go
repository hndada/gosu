package taiko

import (
	"github.com/hndada/gosu/game"
	"math"
)


const defaultHand = right

// const deltaLimitFactor = 1.1 // int(float64(lastTimeDelta)*deltaLimitFactor)
const colorBonus = 0.75
const rhythmBonus = 1

var fingerBonus = [4]float64{0.02, 0, 0, 0.02}

func (beatmap *TaikoBeatmap) CalcStrain() {
	beatmap.calcBaseStrain()

	beatmap.markColorBonus()
	beatmap.markRhythmBonus()
	beatmap.setStrain()
}

// supposed kddk with alter; all notes can be hit alternatively, lead in efficiency
func (beatmap *TaikoBeatmap) calcBaseStrain() {
	var lastHitTime [4]int
	var f1, f2 int
	var v1, v2 float64
	prevNote := TaikoNote{hand: -defaultHand}
	// suppose slider and spinner hit a first note only for now
	for i, note := range beatmap.Notes {
		f1, v1 = beatmap.baseStrain(i, -prevNote.hand, lastHitTime)
		f2, v2 = beatmap.baseStrain(i, prevNote.hand, lastHitTime)
		// choose the hand which makes less strain
		if v1 <= v2 {
			beatmap.Notes[i].hand = -prevNote.hand
			beatmap.Notes[i].baseStrain = 1 + v1
			lastHitTime[f1] = note.Time
		} else {
			beatmap.Notes[i].hand = prevNote.hand
			beatmap.Notes[i].baseStrain = 1 + v2
			lastHitTime[f2] = note.Time
		}
	}
}

func (beatmap TaikoBeatmap) baseStrain(i, hand int, lastHitTime [4]int) (int, float64) {
	note := beatmap.Notes[i]
	f := finger(hand, note.Color)
	if i == 0 {
		return f, 1
	}
	prevNote := beatmap.Notes[i-1]
	v := game.SolveY(beatmap.Curves["Jack"], float64(note.Time-lastHitTime[f]))
	v += game.SolveY(beatmap.Curves["Trill"], float64(note.Time-prevNote.Time))
	v *= fingerBonus[f]
	return f, v
}

func (beatmap *TaikoBeatmap) markColorBonus() {
	var lastColorParity int = none
	var sameColorCount int
	var timeDelta, lastTimeDelta int
	var prevNote TaikoNote
	for i, note := range beatmap.Notes {
		if note.NoteType != donKat { // slider, spinner
			lastColorParity = none
			sameColorCount = 0
			goto ready
		}
		if i == 0 {
			sameColorCount = 1
			goto ready
		}
		prevNote = beatmap.Notes[i-1]
		timeDelta = note.Time - prevNote.Time
		if timeDelta > lastTimeDelta+1 { // stream ends
			lastColorParity = none
			sameColorCount = 1
			goto ready
		}
		if prevNote.Color == note.Color {
			sameColorCount++
		} else {
			if lastColorParity != none && lastColorParity != sameColorCount%2 {
				beatmap.Notes[i].hasColorBonus = true
			}
			lastColorParity = sameColorCount % 2
			sameColorCount = 1
		}
	ready:
		lastTimeDelta = timeDelta
	}
}

func (beatmap *TaikoBeatmap) markRhythmBonus() {
	const threshold = 0.15
	var timeDelta, lastTimeDelta float64
	var ratio, remainder float64
	var prevNote TaikoNote
	for i, note := range beatmap.Notes {
		if note.NoteType != donKat {
			lastTimeDelta = 0
			continue
		}
		if i == 0 {
			continue
		}
		prevNote = beatmap.Notes[i-1]
		timeDelta = float64(note.Time - prevNote.Time)

		if lastTimeDelta == 0 || timeDelta == 0 {
			lastTimeDelta = timeDelta
			continue // to avoid dividing by zero when there is duplicated notes
		}
		ratio = math.Max(lastTimeDelta/timeDelta, timeDelta/lastTimeDelta)
		if ratio < 4 {
			remainder = math.Mod(math.Log2(ratio), 1)
			if remainder >= threshold && remainder < 1-threshold {
				beatmap.Notes[i].hasRhythmBonus = true
			}
		}
		lastTimeDelta = timeDelta
	}
}

func (beatmap *TaikoBeatmap) setStrain() {
	var bonus float64
	for i, note := range beatmap.Notes {
		bonus = 0
		if note.hasColorBonus {
			bonus += colorBonus
		}
		if note.hasRhythmBonus {
			bonus += rhythmBonus
		}
		beatmap.Notes[i].Strain = note.baseStrain + bonus
	}
}
