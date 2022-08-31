package gosu

import (
	"sort"

	"github.com/hndada/gosu/format/osu"
)

// type NoteType int

const (
	Normal = iota
	Head
	Tail
	Body
	BodyTick // e.g., Roll tick in Drum mode.
	Extra    // e.g., Shake in Drum mode.
	ExtraTick
)

// Strategy of Piano mode
// Calculate position of each note in advance
// Parameter: SpeedScale, BPM Ratio, BeatLengthScale
// Calculate current HitPosition only.
// For other notes, just calculate the difference between HitPosition.
type Note struct {
	Type         int
	Time         int64
	Time2        int64
	SampleName   string // aka SampleFilename.
	SampleVolume float64
	Key          int // Not used in Drum mode.

	Marked bool
	LaneObject

	// Following Next and Prev are for calculating scores.
	// Note's Next/Prev and LaneObject's Next/Prev are consistent.
	Next *Note
	Prev *Note // For accessing to Head from Tail.
	// Position float64 // Scaled x or y value.
	// NextTail *Note // For drawing long body faster.
}

func NewNotes(f any, transPoints []*TransPoint, mode, subMode int) (ns []*Note) {
	switch f := f.(type) {
	case *osu.Format:
		// keyCount := int(f.CircleSize)
		// if keyCount <= 4 {
		// 	c.Mode = ModePiano4
		// } else {
		// 	c.Mode = ModePiano7
		// }
		// c.SubMode = keyCount
		ns = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			ns = append(ns, NewNote(ho, mode, subMode)...)
		}
	}
	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})
	var prevs []*Note
	// indexs := []int{Normal, Head, Extra}
	switch mode {
	case ModePiano4, ModePiano7:
		prevs = make([]*Note, subMode) // Todo: ScratchMask?
	case ModeDrum:
		prevs = make([]*Note, 3)
	}
	for _, n := range ns {
		switch mode {
		case ModePiano4, ModePiano7:
			prev := prevs[n.Key]
			n.Prev = prev
			if prev != nil {
				prev.Next = n
			}
			prevs[n.Key] = n
		case ModeDrum:
			var i int
			switch n.Type {
			case Head, Tail: // Head/Tail of Roll/BigRoll
				i = 1
			case Extra: // Shake
				i = 2
			default: // Don, Kat, BigDon, BigKat
				i = 0
			}
			prev := prevs[i]
			n.Prev = prev
			if prev != nil {
				prev.Next = n
			}
			prevs[i] = n
		}
	}

	tp := transPoints[0]
	for i, n := range ns {
		for tp.Next != nil && (tp.Time < n.Time || tp.Time >= tp.Next.Time) {
			tp = tp.Next
		}
		ns[i].LaneObject = LaneObject{
			Type:     n.Type,
			Position: tp.Position + float64(n.Time-tp.Time)*tp.Speed(),
			Speed:    tp.Speed(),
			Next:     &ns[i].Next.LaneObject,
			Prev:     &ns[i].Prev.LaneObject,
			Marked:   &ns[i].Marked,
		}
	}
	return
}

// Todo: NewNote -> newNote
func NewNote(f any, mode, subMode int) []*Note {
	ns := make([]*Note, 0, 2)
	switch f := f.(type) {
	case osu.HitObject:
		n := &Note{
			Type:         Normal,
			Time:         int64(f.Time),
			Time2:        int64(f.Time),
			SampleName:   f.HitSample.Filename,
			SampleVolume: float64(f.HitSample.Volume) / 100,
			// Key:          f.Column(keyCount),
		}
		if mode == ModePiano4 || mode == ModePiano7 {
			n.Key = f.Column(subMode)
		}
		if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
			n.Type = Head
			n.Time2 = int64(f.EndTime)
			n2 := &Note{ // Tail has no sample sound.
				Type:  Tail,
				Time:  n.Time2,
				Time2: n.Time,
				// Key:   n.Key,
			}
			if mode == ModePiano4 || mode == ModePiano7 {
				n2.Key = f.Column(subMode)
			}
			ns = append(ns, n, n2)
		} else {
			ns = append(ns, n)
		}
	}
	return ns
}
