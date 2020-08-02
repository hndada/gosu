package beatmap

import (
	"errors"
	"github.com/hndada/gosu/game"
)

var ErrDuration = errors.New("duration is not a positive value")

type NoteBase struct {
	NoteType                    int
	Time, OpponentTime          int64 // mods' time rate has applied in advance
	Strain, Legibility, Stamina float64
}

func NewNoteBase(h HitObject, mods game.Mods) (NoteBase, error) {
	base := NoteBase{
		NoteType:     h.NoteType,
		Time:         int64(float64(h.StartTime) / mods.TimeRate),
		OpponentTime: int64(float64(h.EndTime) / mods.TimeRate),
	}
	if base.NoteType != NtNote {
		duration := base.OpponentTime - base.Time
		if duration < 0 {
			return base, ErrDuration
		}
	}
	return base, nil
}

func (note NoteBase) Aggregate() float64 {
	// return note.Strain
	return note.Strain + note.Stamina
	// return note.Strain*note.Legibility + note.Stamina
}
