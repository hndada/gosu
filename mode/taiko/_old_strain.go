package taiko

import (
	"math"

	"github.com/hndada/gosu/internal/tools"
)

const strainDecayBase = diff.OldStrainDecayBase
const none = -1

func (beatmap *TaikoBeatmap) CalcOldStrain() {
	var prevNote TaikoOldNote
	var elapsedTime float64
	timeRate := beatmap.Mods.TimeRate

	var addition, additionFactor float64
	var prevElapsedTime float64
	var sameColorCount int = 1 // might work without initializing with 1
	var lastColorParity int = none
	// slider and spinner do make strain
	for i, note := range beatmap.OldNotes {
		if i == 0 {
			beatmap.OldNotes[i].Strain = 1
			continue
		}
		prevNote = beatmap.OldNotes[i-1]
		elapsedTime = float64(note.StartTime-prevNote.StartTime) / timeRate
		beatmap.OldNotes[i].Strain = prevNote.Strain * tools.DecayFactor(strainDecayBase, elapsedTime)

		addition = 1
		if prevNote.NoteType == donKat && note.NoteType == donKat && elapsedTime < 1000 {
			// no bonus on ddkkdd, dkdkdk, but on ddkddk
			if prevNote.Color == note.Color {
				sameColorCount++
			} else {
				if lastColorParity != none && lastColorParity != sameColorCount%2 {
					addition += 0.75
				}
				lastColorParity = sameColorCount % 2
				sameColorCount = 1
			}
			if hasRhythmChange(prevElapsedTime, elapsedTime) {
				addition += 1
			}
		} else {
			lastColorParity = none
			sameColorCount = 1
		}

		additionFactor = 1.0
		if elapsedTime < 50 {
			additionFactor = 0.4 + 0.6*elapsedTime/50
		}

		beatmap.OldNotes[i].Strain += addition * additionFactor
		prevElapsedTime = elapsedTime
	}
}

func hasRhythmChange(prev, current float64) bool {
	const threshold = 0.2
	if prev == 0 || current == 0 {
		return false // to avoid dividing by zero when there is duplicated notes
	}
	ratio := math.Max(prev/current, current/prev)
	if ratio >= 8 {
		return false
	}
	difference := math.Mod(math.Log2(ratio), 1) // would this need round function for precise decimal?
	return difference >= threshold && difference < 1-threshold
}

func (beatmap TaikoBeatmap) InitOldStrainPeak(i int, sectionEndTime int, timeRate float64) float64 {
	// I once noticed that there's no dividing with time rate at original code
	if i == 0 {
		return 0
	}
	prevNote := beatmap.OldNotes[i-1]
	deltaTime := float64(sectionEndTime-prevNote.StartTime) / timeRate // different with elapsedTime
	strainPeak := prevNote.Strain * tools.DecayFactor(strainDecayBase, deltaTime)
	return strainPeak
}

func (beatmap TaikoBeatmap) GetOldStrain(i int) float64 {
	note := beatmap.OldNotes[i]
	return note.Strain
}
