package lv

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/tools"
	"sort"
)

// const noFound = tools.NoFound
const (
	NtHoldHead = beatmap.LastNoteType << iota
	NtHoldTail
)

type Note struct {
	beatmap.NoteBase
	Key int

	hand        int
	chord       []int
	trillJack   []int
	holdImpacts []float64

	strainBase   float64
	chordPenalty float64
	trillBonus   float64
	jackBonus    float64
	holdBonus    float64
}

func NewNotes(hs []beatmap.HitObject, keymode int, mods game.Mods) ([]Note, error) {
	ns := make([]Note, 0, 2*len(hs))
	for _, h := range hs {
		n := make([]Note, 1, 2) // put one or two Note to []Note for every line
		base, err := beatmap.NewNoteBase(h, mods)
		if err != nil {
			return ns, err
		}
		n[0] = Note{
			NoteBase: base,
			Key:      key(keymode, h.X),
		}
		n[0].hand = hand(n[0].Key, keymode)
		n[0].initSlices(keymode)

		if n[0].NoteType == beatmap.NtHoldNote {
			n[0].NoteType = NtHoldHead
			tail := Note{
				NoteBase: beatmap.NoteBase{
					NoteType:     NtHoldTail,
					Time:         n[0].OpponentTime,
					OpponentTime: n[0].Time,
				},
				Key:  n[0].Key,
				hand: n[0].hand,
			}
			tail.initSlices(keymode)
			n = append(n, tail)
		}
		ns = append(ns, n...)
	}
	return ns, nil
}
func key(keymode int, x int) int {
	return keymode * x / 512 // starts with 0
}

func (n *Note) initSlices(keymode int) {
	n.trillJack = tools.GetIntSlice(keymode, noFound)
	n.chord = tools.GetIntSlice(keymode, noFound)
	n.holdImpacts = make([]float64, keymode)
}

func SortNotes(ns []Note) {
	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})
}
