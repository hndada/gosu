package piano

import (
	"sort"

	"github.com/hndada/gosu/game/chart"
	"github.com/hndada/gosu/game/format/osu"
)

const (
	Normal = iota
	Head
	Tail
	Body // Todo: separate Body and other notes at Skin, Drawer?
)

type Note struct {
	Time     int64
	Duration int64
	Type     int
	Key      int
	Position float64 // Scaled x or y value.
	chart.Sample
	Marked bool
	Next   *Note
	Prev   *Note // For accessing to Head from Tail.
}

func NewNote(f any, keyCount int) (ns []*Note) {
	switch f := f.(type) {
	case osu.HitObject:
		n := Note{
			Time:   int64(f.Time),
			Type:   Normal,
			Key:    f.Column(keyCount),
			Sample: chart.NewSample(f),
		}
		if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
			n.Type = Head
			n.Duration = int64(f.EndTime) - n.Time
			n2 := Note{
				Time: n.Time + n.Duration,
				Type: Tail,
				Key:  n.Key,
				// Tail has no sample sound.
			}
			ns = append(ns, &n, &n2)
		} else {
			ns = append(ns, &n)
		}
	}
	return ns
}

// Brilliant idea: Make SpeedScale scaled by MainBPM.
func NewNotes(f any, keyCount int) (ns []*Note) {
	switch f := f.(type) {
	case *osu.Format:
		ns = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			ns = append(ns, NewNote(ho, keyCount)...)
		}
	}
	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})
	prevs := make([]*Note, keyCount&ScratchMask)
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
