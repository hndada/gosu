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

type Note struct {
	Floater
	Duration float64
	// DotDuration float64 // For calculating Roll tick density.
	Type  int
	Color int
	Size  int
	Dots  []*Dot // For shake note, it tells when to hit at auto mods.
	gosu.Sample
	Marked bool
	Next   *Note
	Prev   *Note

	length float64 // For compatibility with osu!.
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
			Sample: gosu.NewSample(f),
		}
		n.Time = int64(f.Time)
		switch {
		case f.NoteType&osu.HitTypeSlider != 0:
			n.Type = Head
			// Roll's duration should not rely on f.EndTime.
		case f.NoteType&osu.HitTypeSpinner != 0:
			n.Type = Shake
			n.Duration = int64(f.EndTime)
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
			speed := (tempMainBPM / 60000) * tp.Speed * (f.SliderMultiplier * 100)
			n.Duration = n.length / speed
			n.DotDuration = n.Duration / ScaledBPM(tp.BPM)
			switch n.Type {
			case Head:
				n.DotDuration /= 4
			case Shake:
				n.DotDuration /= 3
			}
		}
	}
	return
}

type Dot struct {
	Time   int64
	Speed  float64
	Marked bool
	Next   *Dot
	Prev   *Dot
}

func (d Dot) Position(time int64) float64 {
	return d.Speed * float64(d.Time-time)
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
func (n *Note) SetDots(tp *gosu.TransPoint) {
	if n.Type != Head {
		return
	}

}
