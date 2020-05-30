package lv

import (
	"math"
	"reflect"

	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/tools"
)

const (
	OldSectionLength   = 400
	OldDiffDecayWeight = 0.9

	OldStrainDecayBase = 0.3
)

// I found that same chord patterns can be calculated with different difficulty based on the order of notes
func CalcOldSR(beatmap Beatmap) {
	notes := oldNoteBaseSlice(beatmap)
	if len(notes) == 0 {
		return
	}
	var strain, strainPeak float64
	timeRate := reflect.ValueOf(beatmap).Elem().FieldByName("Mods").FieldByName("TimeRate").Float()
	actualSectionLength := int(OldSectionLength * timeRate) // const can be automatically converted to float64
	sectionCounts := BeatmapRawLength(beatmap) / actualSectionLength
	strainPeaks := make([]float64, 0, sectionCounts)

	sectionEndTime := actualSectionLength
	for i, note := range notes {
		for note.StartTime > sectionEndTime { // runs many times if a next note is far
			strainPeaks = append(strainPeaks, strainPeak)
			strainPeak = beatmap.InitOldStrainPeak(i, sectionEndTime, timeRate)
			sectionEndTime += actualSectionLength
		}
		strain = beatmap.GetOldStrain(i)
		if strain > strainPeak {
			strainPeak = strain
		}
	}
	if len(strainPeaks) != sectionCounts { // usually last peak won't be added
		strainPeaks = append(strainPeaks, strainPeak)
	}

	difficulty := tools.WeightedSum(strainPeaks, OldDiffDecayWeight)
	var sr float64
	mode := reflect.ValueOf(beatmap).Elem().FieldByName("Mode").Int()
	switch mode {
	case element.ModeMania:
		sr = difficulty * 0.018
	case element.ModeTaiko:
		sr = difficulty * 0.04125
	}
	reflect.ValueOf(beatmap).Elem().FieldByName("OldStarRating").SetFloat(sr)
}

func oldNoteBaseSlice(beatmap Beatmap) []OldNoteBase {
	notesValue := reflect.ValueOf(beatmap).Elem().FieldByName("OldNotes")
	notes := make([]OldNoteBase, notesValue.Len())
	for i := range notes {
		notes[i] = notesValue.Index(i).FieldByName("OldNoteBase").Interface().(OldNoteBase)
	}
	return notes
}

func BeatmapRawLength(beatmap Beatmap) int {
	notesValue := reflect.ValueOf(beatmap).Elem().FieldByName("RawNotes")
	notes := make([]element.RawNote, notesValue.Len())
	var note element.RawNote
	if len(notes) == 0 {
		return 0
	}
	beatmapLength := notes[len(notes)-1].EndTime
	for i := len(notes) - 1; i >= 0; i-- {
		note = notes[i]
		if note.NoteType == element.NtHoldNote {
			if note.EndTime > beatmapLength {
				beatmapLength = note.EndTime
			}
			break
		}
	}
	return beatmapLength
}

func GetDecayFactor(decayBase, elapsedTime float64) float64 {
	return math.Pow(decayBase, elapsedTime/1000)
}
