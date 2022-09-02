package piano

import (
	"sort"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

const (
	Normal = iota
	Head
	Tail
	Body
)

type Note struct {
	gosu.BaseNote
	Next *Note
	Prev *Note // For accessing to Head from Tail.
	Key  int
}

func NewNote(f any, keyCount int) (ns []*Note) {
	switch f := f.(type) {
	case osu.HitObject:
		n := Note{
			BaseNote: gosu.NewBaseNote(f),
		}
		n.Key = f.Column(keyCount)
		n.Type = Normal
		if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
			n.Type = Head
			n.Time2 = int64(f.EndTime)
			n2 := Note{
				BaseNote: gosu.BaseNote{
					Type:       Tail,
					Time:       n.Time2,
					Time2:      n.Time,
					SampleName: "", // Tail has no sample sound.
				},
				Key: n.Key,
			}
			ns = append(ns, &n, &n2)
		} else {
			ns = append(ns, &n)
		}
	}
	return ns
}

// Brilliant idea: Make SpeedScale scaled by MainBPM.
func NewNotes(f any, keyCount int, tp *gosu.TransPoint) (ns []*Note) {
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
	for _, n := range ns {
		for tp.Next != nil && (tp.Time < n.Time || tp.Time >= tp.Next.Time) {
			tp = tp.Next
		}
		n.Position = tp.Position + float64(n.Time-tp.Time)*tp.Speed()
	}
	return
}
