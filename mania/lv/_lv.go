package lv

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/parser"
	"github.com/hndada/gosu/mania"
	"reflect"

	"github.com/hndada/gosu/tools"
)

const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)


func CalcLv(b mania.Beatmap) float64 {
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
	return difficulty
}

func getNoteBaseSlice(b parser.Beatmap) []game.NoteBase {
	notesValue := reflect.ValueOf(b).Elem().FieldByName("Notes")
	notes := make([]game.NoteBase, notesValue.Len())
	for i := range notes {
		notes[i] = notesValue.Index(i).FieldByName("NoteBase").Interface().(game.NoteBase)
	}
	return notes
}
