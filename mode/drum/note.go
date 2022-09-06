package drum

import (
	"sort"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

// Drum note has 3 components: Color, Size, Type(Note/Head/Tail/Shake).
// cf. Piano note has 2 components: Key, Type(Note/Head/Tail)
const (
	Normal = iota
	Head   // Roll head.
	Tail   // Roll tail.
	Shake
)
const (
	Red  = iota // Don
	Blue        // Kat
	// Yellow   // Roll
)
const (
	Regular = iota
	Big
)

// Todo: make Dots tell when to hit shake tick at auto mods?
type Note struct {
	Floater
	Duration int64
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
func NewNote(f any) (ns []*Note) {
	switch f := f.(type) {
	case osu.HitObject:
		n := Note{
			Sample: gosu.NewSample(f),
		}
		n.Time = int64(f.Time)
		switch {
		case f.NoteType&osu.HitTypeSlider != 0:
			n.Type = Head
			// Roll's duration should not rely on f.EndTime.
		case f.NoteType&osu.HitTypeSpinner != 0:
			n.Type = Shake
			n.Duration = int64(f.EndTime) - n.Time
		default:
			n.Type = Normal
		}
		if osu.IsDon(f) {
			n.Color = Red
		} else {
			n.Color = Blue
		}
		if osu.IsBig(f) {
			n.Size = Big
		} else {
			n.Size = Regular
		}
		if n.Type == Head {
			n.length = f.SliderLength()
			n2 := Note{
				// Time has not yet determined.
				Type:  Tail,
				Color: n.Color,
				Size:  n.Size,
			}
			ns = append(ns, &n, &n2)
		} else {
			ns = append(ns, &n)
		}
	}
	return ns
}

func NewNotes(f any) (ns []*Note) {
	switch f := f.(type) {
	case *osu.Format:
		ns = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			ns = append(ns, NewNote(ho)...)
		}
	}
	sort.Slice(ns, func(i, j int) bool {
		// Todo: additional sort for notes with same time?
		// if ns[i].Time == ns[j].Time {
		// 	return ns[i].Type < ns[j].Type
		// }
		return ns[i].Time < ns[j].Time
	})
	prevs := make([]*Note, 3)
	for _, n := range ns {
		var index int
		switch n.Type {
		case Normal:
			index = 0
		case Head, Tail:
			index = 1
		case Shake:
			index = 2
		}
		prev := prevs[index]
		n.Prev = prev
		if prev != nil {
			prev.Next = n
		}
		prevs[index] = n
	}
	return
}
