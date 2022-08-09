package gosu

import "github.com/hndada/gosu/parse/osu"

type NoteType int

const (
	press = 1 << iota
	release
)
const (
	normal = 4 << iota
	longNote
)
const (
	Normal NoteType = normal + press
	Head            = longNote + press
	Tail            = longNote + release
)

// Prev와 Next는 Note 내에 있으면 안될 것 같다. Scene 등에서 즉석에서 생성 및 관리해야지.
// Prev int // Next note's index
type Note struct {
	Type           NoteType
	Time           int64
	Time2          int64 // For Head note, it is tail's time; vice versa.
	Key            int
	SampleFilename string
	SampleVolume   int
}

// Samples should be lazy loaded
func NewNoteFromOsu(ho osu.HitObject, keyCount int) []Note {
	ns := make([]Note, 0, 2)
	n := Note{
		Type:  Normal,
		Time:  int64(ho.Time),
		Time2: int64(ho.Time),
		Key:   ho.Column(keyCount),
		// SampleFilename: ,
		// SampleVolume: ,
	}
	if ho.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
		n.Type = Head
		n.Time2 = int64(ho.EndTime)
		n2 := Note{ // Tail has no sample sound
			Type:  Tail,
			Time:  n.Time2,
			Time2: n.Time,
			Key:   n.Key,
		}
		ns = append(ns, n, n2)
	} else {
		ns = append(ns, n)
	}
	return ns
}
