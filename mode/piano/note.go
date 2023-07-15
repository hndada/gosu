package piano

import (
	"sort"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

const (
	Normal = iota
	Head
	Tail
	Body
)

type Note struct {
	Time     int32
	Type     int
	Key      int
	Sample   mode.Sample
	Duration int32

	Position float64 // Scaled x or y value.
	Next     *Note
	Prev     *Note // For accessing to Head from Tail.
	Marked   bool

	Strain float64 // Strain is for calculating difficulty.
}

func NewNotes(f any, keyCount int) (ns []*Note) {
	switch f := f.(type) {
	case *osu.Format:
		ns = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			ns = append(ns, newNoteFromOsu(ho, keyCount)...)
		}
	}

	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})

	// linking
	prevs := make([]*Note, keyCount)
	for _, n := range ns {
		prev := prevs[n.Key]
		n.Prev = prev
		if prev != nil {
			prev.Next = n
		}
		prevs[n.Key] = n
	}
	return
}

// The length of the returned slice is 1 or 2.
func newNoteFromOsu(f osu.HitObject, keyCount int) (ns []*Note) {
	n := &Note{
		Time:   int32(f.Time),
		Type:   Normal,
		Key:    f.Column(keyCount),
		Sample: mode.NewSample(f),
	}
	if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
		n.Type = Head
		n.Duration = int32(f.EndTime) - n.Time
		n2 := &Note{
			Time: n.Time + n.Duration,
			Type: Tail,
			Key:  n.Key,
			// Sample: mode.Sample{}, // Tail has no sample sound.
			// Duration: n.Duration, // Todo: 0 or n.Duration?
		}
		ns = append(ns, n, n2)
	} else {
		ns = append(ns, n)
	}
	return ns
}
