package drum

import (
	"sort"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

// Drum note has 3 components: Color, Size, Type(Note, Roll, Shake).
// cf. Piano note has 2 components: Key, Type(Note, Head, Tail)
const (
	Normal = iota
	Roll   // Roll with head and tail.
	Shake
)
const (
	Red    = iota // aka Don.
	Blue          // aka Kat.
	Yellow        // Roll
)
const (
	Regular = iota
	Big
)

// Todo: make Dots tell when to hit shake tick at auto mods?
type Note struct {
	Floater
	Duration   int64
	RevealTime int64 // The time when Roll body or Shake reveals.
	// DotDuration float64 // For calculating Roll tick density.
	Type  int
	Color int
	Size  int
	// Dots  []*Dot
	Tick int // The number of ticks left to hit in Roll or Shake.

	gosu.Sample
	Marked bool
	Next   *Note
	Prev   *Note
	length float64 // For compatibility with osu!.
}

// Length is for calculating Durations of Roll and its tick.
// NewNote temporarily set length to DotDuration.
func NewNote(f any) (n *Note) {
	switch f := f.(type) {
	case osu.HitObject:
		n = &Note{
			Sample: gosu.NewSample(f),
		}
		n.Time = int64(f.Time)
		switch {
		case f.NoteType&osu.HitTypeSlider != 0:
			n.Type = Roll
			n.length = f.SliderLength()
			// Roll's duration should not rely on f.EndTime.
			n.Color = Yellow
		case f.NoteType&osu.HitTypeSpinner != 0:
			n.Type = Shake
			n.Duration = int64(f.EndTime) - n.Time
		default:
			n.Type = Normal
			if osu.IsDon(f) {
				n.Color = Red
			} else {
				n.Color = Blue
			}
		}
		if osu.IsBig(f) {
			n.Size = Big
		} else {
			n.Size = Regular
		}
	}
	return
}

func NewNotes(f any) (notes, shakes []*Note) {
	switch f := f.(type) {
	case *osu.Format:
		notes = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			n := NewNote(ho)
			if n.Type == Shake {
				shakes = append(shakes, n)
			} else {
				notes = append(notes, n)
			}
		}
	}
	// Sort notes only with their time.
	// Order of notes at the same time might be intentional for gimmicks.
	sort.SliceStable(notes, func(i, j int) bool {
		return notes[i].Time < notes[j].Time
	})
	sort.SliceStable(shakes, func(i, j int) bool {
		return shakes[i].Time < shakes[j].Time
	})
	prevs := make([]*Note, 3)
	for _, ns := range [][]*Note{notes, shakes} {
		for _, n := range ns {
			prev := prevs[n.Type]
			n.Prev = prev
			if prev != nil {
				prev.Next = n
			}
			prevs[n.Type] = n
		}
	}
	return
}
