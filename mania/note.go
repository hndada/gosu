package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/beatmap"
	"sort"

	"github.com/hndada/gosu/game/tools"
)

const noFound = tools.NoFound
const (
	NtHoldHead = beatmap.LastNoteType << iota
	NtHoldTail
)

type Note struct {
	beatmap.NoteBase
	Key int // i-th column
	// StrainBase   float64
	// TrillBonus   float64
	// JackBonus    float64
	// HoldBonus    float64
	// ChordPenalty float64

	hand        int
	chord       []int
	trillJack   []int
	holdImpacts []float64
}

// hand 나중에 입력
// init slice 해줘야 하나?
func NewNotes(hs []beatmap.HitObject, keymode int, mods game.Mods) ([]Note, error) {
	notes := make([]Note, 0, 2*len(hs))
	for _, h := range hs {
		n := make([]Note, 1, 2) // put one or two Note to []Note for every line
		base, err := beatmap.NewNoteBase(h, mods)
		if err != nil {
			return notes, err
		}
		n[0] = Note{
			NoteBase: base,
			Key:      key(keymode, h.X),
		}
		if n[0].NoteType == beatmap.NtHoldNote {
			n[0].NoteType = NtHoldHead
			tail := Note{
				NoteBase: beatmap.NoteBase{
					NoteType:     NtHoldTail,
					Time:         n[0].OpponentTime,
					OpponentTime: n[0].Time,
				},
				Key: n[0].Key,
			}
			n = append(n, tail)
		}
		notes = append(notes, n...)
	}
	return notes, nil
}
func key(keymode int, x int) int {
	return keymode * x / 512 // starts with 0
}

func SortNotes(notes []Note) {
	sort.Slice(notes, func(i, j int) bool {
		if notes[i].Time == notes[j].Time {
			return notes[i].Key < notes[j].Key
		}
		return notes[i].Time < notes[j].Time
	})
}

// func (note *ManiaNote) initSlices(keymode int) {
// 	note.trillJack = tools.GetIntSlice(keymode, noFound)
// 	note.chord = tools.GetIntSlice(keymode, noFound)
// 	note.holdImpacts = make([]float64, keymode)
// }
