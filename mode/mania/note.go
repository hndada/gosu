package mania

import (
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/parser/osu"
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

type Note struct {
	mode.BaseNote
	Key int

	// 난이도 및 점수 관련은 나중에
}

func NewNotes(hs []osu.HitObject, keys int) ([]Note, error) {
	ns := make([]Note, 0, 2*len(hs))
	for _, h := range hs {
		// n := Note{
		// 	BaseNote: mode.BaseNote{
		// 		Type: noteType(h.NoteType),
		// 		Time: int64(h.StartTime),
		// 	},
		// 	Key: key(keys, h.X),
		// }
		var n Note
		n.Type = noteType(h.NoteType)
		n.Key = key(keys, h.X)
		n.Time = int64(h.StartTime)

		if n.Type == TypeLongNote {
			n.Type = LNHead
			n.Time2 = int64(h.EndTime)
			ns = append(ns, n)

			var n2 Note
			n2.Type = LNTail
			n2.Key = n.Key
			n2.Time = n.Time2
			n2.Time2 = n.Time
			ns = append(ns, n2)
		} else {
			ns = append(ns, n)
		}
	}
	return ns, nil
}

// func NewNotesFromOSU
// func NewNotesFromBMS

func noteType(t int) int16 {
	switch {
	case t&osu.NtNote != 0:
		return TypeNote
	case t&osu.NtHoldNote != 0:
		return TypeLongNote
	}
	panic("not reach")
}

func key(keys int, x int) int {
	return keys * x / 512 // starts with 0
}
