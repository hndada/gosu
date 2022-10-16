package drum

import (
	"math"
	"sort"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

// Drum note has 3 components: Color, Size, Type(Note, Roll, Shake).
// Piano note has 2 components: Key, Type(Note, Head, Tail).
const (
	Normal = iota
	Roll
	Shake
)
const (
	ColorNone = iota - 1
	Red       // aka Don.
	Blue      // aka Kat.
	Yellow    // For Roll.
	Purple    // For Shake.
)
const (
	SizeNone = iota - 1
	Regular
	Big
)

type Note struct {
	Floater
	Type     int
	Color    int
	Size     int
	Duration int64
	length   float64 // For compatibility with osu!.
	Tick     int     // The number of ticks in Roll or Shake.
	gosu.Sample

	Marked  bool
	HitTick int // The number of ticks being hit.
	Next    *Note
	Prev    *Note
}

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
			n.length = f.SliderLength() // For calculating duration of Roll.
			n.Color = Yellow
		case f.NoteType&osu.HitTypeSpinner != 0:
			n.Type = Shake
			n.Duration = int64(f.EndTime) - n.Time
			n.Color = Purple
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

func NewNotes(f any) (notes, rolls, shakes []*Note) {
	switch f := f.(type) {
	case *osu.Format:
		notes = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			n := NewNote(ho)
			switch n.Type {
			case Normal:
				notes = append(notes, n)
			case Roll:
				rolls = append(rolls, n)
			case Shake:
				shakes = append(shakes, n)
			}
		}
	}
	// Sort notes only with their time.
	// Order of notes at the same time might be intentional for gimmicks.
	sort.SliceStable(notes, func(i, j int) bool {
		return notes[i].Time < notes[j].Time
	})
	sort.SliceStable(rolls, func(i, j int) bool {
		return rolls[i].Time < rolls[j].Time
	})
	sort.SliceStable(shakes, func(i, j int) bool {
		return shakes[i].Time < shakes[j].Time
	})
	prevs := make([]*Note, 3)
	for kind, ns := range [][]*Note{notes, rolls, shakes} {
		for _, n := range ns {
			prev := prevs[kind]
			n.Prev = prev
			if prev != nil {
				prev.Next = n
			}
			prevs[kind] = n
		}
	}
	return
}

func (n Note) Weight() float64 {
	switch n.Type {
	case Normal:
		switch n.Size {
		case Regular:
			return 1.0
		case Big:
			return 1.0 // 1.1
		}
	case Shake:
		// Shake is apparently easier than Roll, since it is free from beat.
		// https://www.desmos.com/calculator/nsogcrebx9
		return math.Pow(float64(32*n.Tick), 0.75) / 32 // 32 comes from 8 * 4.
	}
	return 0
}
func (d Dot) Weight() float64 { return 0.125 } // Assumes 8 ticks worth one normal note.
