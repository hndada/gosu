package piano

import "github.com/hndada/gosu/format/osu"

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

type Note struct {
	Type           NoteType
	Time           int64
	Time2          int64 // For Head note, it is tail's time; vice versa.
	Key            int
	SampleFilename string
	SampleVolume   int
}

// A sample sound file should be lazy loaded.
func NewNote(f any, keyCount int) []Note {
	ns := make([]Note, 0, 2)
	switch f := f.(type) {
	case osu.HitObject:
		n := Note{
			Type:           Normal,
			Time:           int64(f.Time),
			Time2:          int64(f.Time),
			Key:            f.Column(keyCount),
			SampleFilename: f.HitSample.Filename,
			SampleVolume:   f.HitSample.Volume,
		}
		if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
			n.Type = Head
			n.Time2 = int64(f.EndTime)
			n2 := Note{ // Tail has no sample sound.
				Type:  Tail,
				Time:  n.Time2,
				Time2: n.Time,
				Key:   n.Key,
			}
			ns = append(ns, n, n2)
		} else {
			ns = append(ns, n)
		}
	}
	return ns
}
