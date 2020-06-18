package game

import (
	"github.com/hndada/gosu/game/beatmap"
	"reflect"

	"github.com/hndada/gosu/game/tools"
)

const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)
// const (
// 	srScaleMania = 0.025
// 	srScaleTaiko = 0.051
// )

type Beatmap interface {
	SetBase(path string, modsBits int)
	AddNotes()
	SortNotes()

	SetCurves()
	CalcStrain()
	CalcStamina()
	// CalcLegibility()
}

func CalcSR(beatmap beatmap.Beatmap) float64 {
	// mode := reflect.ValueOf(beatmap).Elem().FieldByName("Mode").Int()
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
		aggregate += note.Aggregate()
	}
	difficulty := tools.WeightedSum(aggregates, diffWeightDecay)
	// var lv float64
	// switch mode {
	// case element.ModeMania:
	// 	lv = difficulty * srScaleMania
	// case element.ModeTaiko:
	// 	lv = difficulty * srScaleTaiko
	// }
	// reflect.ValueOf(beatmap).Elem().FieldByName("Lv").SetFloat(sr)
	return difficulty
}

func getNoteBaseSlice(b beatmap.Beatmap) []beatmap.NoteBase {
	notesValue := reflect.ValueOf(b).Elem().FieldByName("Notes")
	notes := make([]beatmap.NoteBase, notesValue.Len())
	for i := range notes {
		notes[i] = notesValue.Index(i).FieldByName("NoteBase").Interface().(beatmap.NoteBase)
	}
	return notes
}
