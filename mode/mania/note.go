package mania

import (
	"errors"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/rg-parser/osugame/osu"
)

const (
	TypeNote = 1 << iota
	TypeReleaseNote
	TypeLongNote
)
const (
	LNHead = TypeLongNote | TypeNote
	LNTail = TypeLongNote | TypeReleaseNote
)

// 난이도 및 점수 관련은 나중에
type Note struct {
	mode.BaseNote
	Key int

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

	// Strain
	// Read
	// Stamina
}

var ErrDuration = errors.New("invalid duration: not a positive value")

func newNote(ho osu.HitObject, keys int, mods Mods) ([]Note, error) {
	ns := make([]Note, 0, 2)
	var n Note
	switch ho.NoteType & osu.ComboMask {
	case osu.HitTypeHoldNote:
		n.Type = TypeLongNote
	case osu.HitTypeNote:
		n.Type = TypeNote
	default:
		return ns, errors.New("invalid hit object")
	}
	n.Key = ho.Column(keys)
	n.Time = int64(float64(ho.Time) / mods.TimeRate)
	n.SampleFilename = ho.HitSample.Filename
	n.SampleVolume = uint8(ho.HitSample.Volume)

	if n.Type == TypeLongNote {
		n.Type = LNHead
		n.Time2 = int64(float64(ho.EndTime) / mods.TimeRate)
		ns = append(ns, n)
		if n.Time2-n.Time <= 0 {
			return ns, ErrDuration
		}

		var n2 Note
		n2.Type = LNTail
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