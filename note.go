package gosu

import (
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

type BaseNote struct {
	Type         int
	Time         int64
	Time2        int64  // Time of opposite note. Time2 is Time when no opposite note.
	SampleName   string // aka SampleFilename.
	SampleVolume float64

	Position float64 // Scaled x or y value.
	Marked   bool
}

func NewBaseNote(f any) (n BaseNote) {
	switch f := f.(type) {
	case osu.HitObject:
		n = BaseNote{
			Time:         int64(f.Time),
			Time2:        int64(f.Time),
			SampleName:   f.HitSample.Filename,
			SampleVolume: float64(f.HitSample.Volume) / 100,
		}
	}
	return n
}

func (n BaseNote) SamplePath(cpath string) (string, bool) {
	if n.SampleName == "" {
		return "", false
	}
	return filepath.Join(filepath.Dir(cpath), n.SampleName), true
}
