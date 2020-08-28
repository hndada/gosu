package mania

import (
	"errors"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/rg-parser/osugame/osu"
)

const (
	typeNote mode.NoteType = 1 << iota
	typeReleaseNote
	typeLongNote
)
const (
	TypeNote   = typeNote
	TypeLNHead = typeLongNote | typeNote
	TypeLNTail = typeLongNote | typeReleaseNote
)

// 난이도 및 점수 관련은 나중에
// 아래와 같이 난이도 계산에만 쓰이는 값들은 unexported로 할듯
type Note struct {
	mode.BaseNote
	Key int
	// Strain
	// Read
	// Stamina

	// hand        int
	// chord       []int
	// trillJack   []int
	// holdImpacts []float64
	//
	// strainBase   float64
	// chordPenalty float64
	// trillBonus   float64
	// jackBonus    float64
	// HoldBonus    float64 // score 필요함
}

var ErrDuration = errors.New("invalid duration: not a positive value")

func newNote(ho osu.HitObject, keys int) ([]Note, error) {
	ns := make([]Note, 0, 2)
	var n Note
	switch ho.NoteType & osu.ComboMask {
	case osu.HitTypeHoldNote:
		n.Type = typeLongNote
	case osu.HitTypeNote:
		n.Type = TypeNote
	default:
		return ns, errors.New("invalid hit object")
	}
	n.Key = ho.Column(keys)
	// n.Time = int64(float64(ho.Time) / mods.TimeRate)
	n.Time = int64(ho.Time)
	n.SampleFilename = ho.HitSample.Filename
	n.SampleVolume = uint8(ho.HitSample.Volume)

	if n.Type == typeLongNote {
		n.Type = TypeLNHead
		// n.Time2 = int64(float64(ho.EndTime) / mods.TimeRate)
		n.Time2 = int64(ho.EndTime)
		ns = append(ns, n)
		if n.Time2-n.Time <= 0 {
			return ns, ErrDuration
		}

		var n2 Note
		n2.Type = TypeLNTail
		n2.Key = n.Key
		n2.Time = n.Time2
		n2.Time2 = n.Time
		ns = append(ns, n2)
	} else {
		ns = append(ns, n)
	}
	return ns, nil
}

// func SortNotes(ns []Note) {
// 	sort.Slice(ns, func(i, j int) bool {
// 		if ns[i].Time == ns[j].Time {
// 			return ns[i].Key < ns[j].Key
// 		}
// 		return ns[i].Time < ns[j].Time
// 	})
// }

// func (n *Note) initSlices(keymode int) {
//	n.trillJack = tools.GetIntSlice(keymode, noFound)
//	n.chord = tools.GetIntSlice(keymode, noFound)
//	n.holdImpacts = make([]float64, keymode)
// }
