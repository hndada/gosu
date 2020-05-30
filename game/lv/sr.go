package lv

import (
	"reflect"

	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/tools"
)

const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)
const (
	srScaleMania = 0.025
	srScaleTaiko = 0.051
)

func CalcSR(beatmap Beatmap) float64 {
	mode := reflect.ValueOf(beatmap).Elem().FieldByName("Mode").Int()
	notes := getNoteBaseSlice(beatmap)
	if len(notes) == 0 {
		return 0
	}
	sectionCounts := notes[len(notes)-1].Time / sectionLength
	sectionEndTime := sectionLength + notes[0].Time

	var aggregate float64
	aggregates := make([]float64, 0, sectionCounts)
	for _, note := range notes {
		if note.Time > sectionEndTime {
			aggregates = append(aggregates, aggregate)
			aggregate = 0.0
			sectionEndTime += sectionLength
		}
		aggregate += note.aggregate()
	}
	difficulty := tools.WeightedSum(aggregates, diffWeightDecay)
	var sr float64
	switch mode {
	case element.ModeMania:
		sr = difficulty * srScaleMania
	case element.ModeTaiko:
		sr = difficulty * srScaleTaiko
	}
	reflect.ValueOf(beatmap).Elem().FieldByName("StarRating").SetFloat(sr)
	return sr
}

func (note NoteBase) aggregate() float64 {
	// return note.Strain
	return note.Strain + note.Stamina
	// return note.Strain*note.Legibility + note.Stamina
}

func getNoteBaseSlice(beatmap Beatmap) []NoteBase {
	notesValue := reflect.ValueOf(beatmap).Elem().FieldByName("Notes")
	notes := make([]NoteBase, notesValue.Len())
	for i := range notes {
		notes[i] = notesValue.Index(i).FieldByName("NoteBase").Interface().(NoteBase)
	}
	return notes
}
