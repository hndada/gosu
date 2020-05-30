package taiko

import (
	"sort"

	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/lv"
	"github.com/hndada/gosu/game/tools"
)

const (
	don = iota + 1
	kat
)
const donKat = element.NtNote
const katBitMask = element.Whistle + element.Clap
const allowedNTMask = element.NtNote + element.NtSlider + element.NtSpinner

type TaikoNoteCommon struct {
	Color int
}

type TaikoNote struct {
	diff.NoteBase
	TaikoNoteCommon
	Big bool

	hand           int
	baseStrain     float64
	hasColorBonus  bool
	hasRhythmBonus bool
}

type TaikoOldNote struct {
	diff.OldNoteBase
	TaikoNoteCommon

	Strain float64
}

// need to handle converted case
func (beatmap *TaikoBeatmap) AddNotes() {
	beatmap.Notes = make([]TaikoNote, len(beatmap.RawNotes))
	for i, raw := range beatmap.RawNotes {
		note := TaikoNote{
			NoteBase: diff.GetNoteBase(raw, beatmap.Mods),
			TaikoNoteCommon: TaikoNoteCommon{
				Color: color(raw.HitSound),
			},
			Big: raw.HitSound&element.Finish != 0,
		}
		if note.NoteType&allowedNTMask == 0 {
			panic(&tools.ValError{"NoteType", tools.Itoa(note.NoteType), tools.ErrSyntax})
		}
		beatmap.Notes[i] = note
	}
}

func (beatmap *TaikoBeatmap) AddOldNotes() {
	beatmap.OldNotes = make([]TaikoOldNote, len(beatmap.RawNotes))
	for i, raw := range beatmap.RawNotes {
		note := TaikoOldNote{
			OldNoteBase: diff.GetOldNoteBase(raw),
			TaikoNoteCommon: TaikoNoteCommon{
				Color: color(raw.HitSound),
			},
		}
		beatmap.OldNotes[i] = note
	}
}

func (beatmap *TaikoBeatmap) SortNotes() {
	sort.Slice(beatmap.Notes, func(i, j int) bool {
		return beatmap.Notes[i].Time < beatmap.Notes[j].Time
	})
}

func color(hitSound int) int {
	if hitSound&katBitMask != 0 {
		return kat
	}
	return don
}
