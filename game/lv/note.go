package lv

import (
	"errors"

	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/tools"
)

var ErrDuration = errors.New("duration is not a positive value")

type NoteBase struct {
	NoteType                    int
	Time, OpponentTime          int // mods' time rate has applied in advance
	Strain, Legibility, Stamina float64
}
type OldNoteBase struct {
	NoteType           int
	StartTime, EndTime int
}

func GetNoteBase(raw element.RawNote, mods Mods) NoteBase {
	base := NoteBase{
		NoteType:     raw.NoteType,
		Time:         int(float64(raw.StartTime) / mods.TimeRate),
		OpponentTime: int(float64(raw.EndTime) / mods.TimeRate),
	}
	if base.NoteType != element.NtNote {
		duration := base.OpponentTime - base.Time
		if duration < 0 {
			panic(&tools.ValError{"End time", tools.Itoa(duration), ErrDuration})
		}
	}
	return base
}

func GetOldNoteBase(raw element.RawNote) OldNoteBase {
	base := OldNoteBase{
		NoteType:  raw.NoteType,
		StartTime: raw.StartTime,
		EndTime:   raw.EndTime,
	}
	return base
}
