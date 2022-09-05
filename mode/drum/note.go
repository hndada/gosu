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

type Dot struct {
	Time   int64
	Speed  float64
	Marked bool
}
type Note struct {
	gosu.BaseNote
	Type        int
	Color       int
	Size        int
	Speed       float64 // Each note has own speed.
	DotDuration float64 // For calculating Roll tick density.
	Dots        []*Dot  // For shake note, it tells when to hit at auto mods.
	Next        *Note
	Prev        *Note // For accessing to Head from Tail.
}

const (
	Red  = iota // Don
	Blue        // Kat
	// Yellow   // Roll
)
const (
	Regular = iota
	Big
)

// Length is for calculating Durations of Roll and its tick.
// NewNote temporarily set length to DotDuration.
func NewNote(f any) (ns []*Note) {
	switch f := f.(type) {
	case osu.HitObject:
		n := Note{
			BaseNote: gosu.NewBaseNote(f),
		}
		switch {
		case f.NoteType&osu.HitTypeSlider != 0:
			n.Type = Head
			// Roll's Time2 should not rely on f.EndTime.
		case f.NoteType&osu.HitTypeSpinner != 0:
			n.Type = Shake
			n.Time2 = int64(f.EndTime)
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
		}
		// n.Big = osu.IsBig(f)
		if n.Type == Head {
			length := float64(f.SliderParams.Slides) * f.SliderParams.Length
			n.DotDuration = length // Temporarily set to DotDuration.

			n2 := n
			n.Type = Tail
			// n2.Time = n.Time2 // It will be set later.
			n2.Time2 = n.Time
			n2.SampleName = "" // Tail has no sample sound.
			ns = append(ns, &n, &n2)
		} else {
			ns = append(ns, &n)
		}
	}
	return ns
}

func NewNotes(f any, transPoints []*gosu.TransPoint) (ns []*Note) {
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

	// Unit of speed is osupixel / 100ms.
	switch f := f.(type) {
	case *osu.Format:
		tp := transPoints[0]
		tempMainBPM := tp.BPM
		for _, n := range ns {
			if n.Type == Normal || n.Type == Tail {
				continue
			}
			for tp.Next != nil && n.Time >= tp.Next.Time {
				tp = tp.Next
			}
			length := n.DotDuration
			speed := (tempMainBPM / 60000) * tp.Speed * (f.SliderMultiplier * 100)
			duration := length / speed
			n.Time2 = n.Time + int64(duration)
			n.DotDuration = duration / ScaledBPM(tp.BPM)
			switch n.Type {
			case Head:
				n.DotDuration /= 4
				n.Next.Time = n.Time2
				n.Next.DotDuration = n.DotDuration
			case Shake:
				n.DotDuration /= 3
			}
		}
	}
	return
}

// It is proved that all BPMs are set into [MinScaledBPM, MaxScaledBPM) by v*2 or v/2
// if MinScaledBPM *2 >= MaxScaleBPM.
func ScaledBPM(bpm float64) float64 {
	if bpm < 0 {
		bpm = -bpm
	}
	switch {
	case bpm > MaxScaledBPM:
		for bpm > MaxScaledBPM {
			bpm /= 2
		}
	case bpm >= MinScaledBPM:
		return bpm
	case bpm < MinScaledBPM:
		for bpm < MinScaledBPM {
			bpm *= 2
		}
	}
	return bpm
}
