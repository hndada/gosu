package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/parser/osu"
)

const (
	TypeNote = 1 << (iota + 1)
	TypeReleaseNote
	TypeLongNote
)
const (
	LNHead = TypeLongNote | TypeNote
	LNTail = TypeLongNote | TypeReleaseNote
)

type Note struct {
	game.BaseNote
	Key int

	// 난이도 및 점수 관련은 나중에
}

func NewNotes(hs []osu.HitObject, keys int) ([]Note, error) {
	ns := make([]Note, 0, 2*len(hs))
	for _, h := range hs {
		// n := Note{
		// 	BaseNote: game.BaseNote{
		// 		Type: noteType(h.NoteType),
		// 		Time: int64(h.StartTime),
		// 	},
		// 	Key: key(keys, h.X),
		// }
		var n Note
		n.Type = noteType(h.NoteType)
		n.Key = key(keys, h.X)
		n.Time = int64(h.StartTime)
		ns = append(ns, n)

		if n.Type == TypeLongNote {
			n.Type = LNHead
			n.Time2 = int64(h.EndTime)

			var n2 Note
			n2.Type = LNTail
			n2.Key = n.Key
			n2.Time = n.Time2
			n2.Time2 = n.Time
			ns = append(ns, n2)
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
