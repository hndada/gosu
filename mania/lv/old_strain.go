package mania

import (
	"github.com/hndada/gosu/game/lv"
	"github.com/hndada/gosu/game/tools"
)

const (
	overallStrainDecayBase = diff.OldStrainDecayBase
	indStrainDecayBase     = 0.125
)

func (beatmap *ManiaBeatmap) CalcOldStrain() {
	var prevNote ManiaOldNote
	var elapsedTime float64
	timeRate := beatmap.Mods.TimeRate

	var holdAddition int
	var holdFactor float64
	prevIndividualStrains := make([]float64, beatmap.Keymode)
	holdEndTimes := tools.GetIntSlice(beatmap.Keymode, 0) // I think this should have been NoFound
	for i, note := range beatmap.OldNotes {
		if i == 0 {
			beatmap.OldNotes[i].overallStrain = 1
			continue
		}
		prevNote = beatmap.OldNotes[i-1]
		elapsedTime = float64(note.StartTime-prevNote.StartTime) / timeRate
		for key := range prevIndividualStrains { // each individual strain decays every note passes by
			prevIndividualStrains[key] *= tools.DecayFactor(indStrainDecayBase, elapsedTime)
		}

		// I found flaw of exisitng algorithm: the algorithm doesn't sorts only by (start) time
		// while the following for-loop doesn't break even if found the hold note which ends along with
		// so it has possibility of return different value with an identical beatmap
		holdAddition, holdFactor = 0, 1.0
		for key := range holdEndTimes {
			switch {
			case note.EndTime > holdEndTimes[key] && note.StartTime < holdEndTimes[key]:
				holdAddition = 1 // note is a hold note which starts first but lasts more than other hold note
			case note.EndTime == holdEndTimes[key]:
				holdAddition = 0 // note ends when hold note ends
			case note.EndTime < holdEndTimes[key]:
				holdFactor = 1.25 // other hold note is holding
			}
		}
		beatmap.OldNotes[i].individualStrain = prevIndividualStrains[note.Key]
		beatmap.OldNotes[i].individualStrain += 2 * holdFactor
		beatmap.OldNotes[i].overallStrain = prevNote.overallStrain * tools.DecayFactor(overallStrainDecayBase, elapsedTime)
		beatmap.OldNotes[i].overallStrain += float64(1+holdAddition) * holdFactor

		prevIndividualStrains[note.Key] = beatmap.OldNotes[i].individualStrain
		holdEndTimes[note.Key] = note.EndTime // simple handling since notes at end of hold note get no bonus
	}
}

func (beatmap ManiaBeatmap) InitOldStrainPeak(i int, sectionEndTime int, timeRate float64) float64 {
	// I once noticed that there's no dividing with time rate at original code
	if i == 0 {
		return 0
	}
	prevNote := beatmap.OldNotes[i-1]
	deltaTime := float64(sectionEndTime-prevNote.StartTime) / timeRate // different with elapsedTime
	strainPeak := prevNote.individualStrain * tools.DecayFactor(indStrainDecayBase, deltaTime)
	strainPeak += prevNote.overallStrain * tools.DecayFactor(overallStrainDecayBase, deltaTime)
	return strainPeak
}

func (beatmap ManiaBeatmap) GetOldStrain(i int) float64 {
	note := beatmap.OldNotes[i]
	return note.individualStrain + note.overallStrain
}
