package mania

import (
	"sort"

	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/lv"
	"github.com/hndada/gosu/game/tools"
)

const (
	NtHoldHead = element.LastNoteType << iota
	NtHoldTail
)

type ManiaNoteCommon struct {
	Key int // i-th column
}

type ManiaNote struct {
	diff.NoteBase
	ManiaNoteCommon
	StrainBase   float64
	TrillBonus   float64
	JackBonus    float64
	HoldBonus    float64
	ChordPenalty float64

	hand        int
	chord       []int
	trillJack   []int
	holdImpacts []float64
}

type ManiaOldNote struct {
	diff.OldNoteBase
	ManiaNoteCommon

	individualStrain float64
	overallStrain    float64
}

// need to handle converted case
func (beatmap *ManiaBeatmap) AddNotes() {
	beatmap.Notes = make([]ManiaNote, 0, 2*len(beatmap.RawNotes))
	for _, raw := range beatmap.RawNotes {
		notes := make([]ManiaNote, 1, 2) // put one or two Note to []Note for every line
		notes[0] = ManiaNote{
			NoteBase: diff.GetNoteBase(raw, beatmap.Mods),
			ManiaNoteCommon: ManiaNoteCommon{
				Key: key(beatmap.Keymode, raw.X),
			},
			hand: hand(notes[0].Key, beatmap.Keymode),
		}
		notes[0].initSlices(beatmap.Keymode)

		note := notes[0]
		if note.NoteType == element.NtHoldNote {
			notes[0].NoteType = NtHoldHead
			holdNoteTail := ManiaNote{
				NoteBase: diff.NoteBase{
					NoteType:     NtHoldTail,
					Time:         note.OpponentTime,
					OpponentTime: note.Time,
				},
				ManiaNoteCommon: ManiaNoteCommon{
					Key: note.Key,
				},
				hand: note.hand,
			}
			holdNoteTail.initSlices(beatmap.Keymode)
			notes = append(notes, holdNoteTail)
		}
		beatmap.Notes = append(beatmap.Notes, notes...)
	}
}

func (beatmap *ManiaBeatmap) AddOldNotes() {
	beatmap.OldNotes = make([]ManiaOldNote, len(beatmap.RawNotes))
	for i, raw := range beatmap.RawNotes {
		note := ManiaOldNote{
			OldNoteBase: diff.GetOldNoteBase(raw),
			ManiaNoteCommon: ManiaNoteCommon{
				Key: key(beatmap.Keymode, raw.X),
			},
		}
		beatmap.OldNotes[i] = note
	}
}

func (beatmap *ManiaBeatmap) SortNotes() {
	sort.Slice(beatmap.Notes, func(i, j int) bool {
		if beatmap.Notes[i].Time == beatmap.Notes[j].Time {
			return beatmap.Notes[i].Key < beatmap.Notes[j].Key
		}
		return beatmap.Notes[i].Time < beatmap.Notes[j].Time
	})
}

func key(keymode int, x int) int {
	return keymode * x / 512 // starts with 0
}

func (note *ManiaNote) initSlices(keymode int) {
	note.trillJack = tools.GetIntSlice(keymode, noFound)
	note.chord = tools.GetIntSlice(keymode, noFound)
	note.holdImpacts = make([]float64, keymode)
}
