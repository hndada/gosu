package gosu

import (
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

// type BaseNote struct {
// 	Time int64
// 	// Time2        int64  // Time of opposite note. Time2 is Time when no opposite note.
// 	SampleName   string // aka SampleFilename.
// 	SampleVolume float64
// 	Duration     int64
// 	Marked       bool
// }

// func NewBaseNote(f any) (n BaseNote) {
// 	switch f := f.(type) {
// 	case osu.HitObject:
// 		n = BaseNote{
// 			Time: int64(f.Time),
// 			// Time2:        int64(f.Time),
// 			SampleName:   f.HitSample.Filename,
// 			SampleVolume: float64(f.HitSample.Volume) / 100,
// 		}
// 	}
// 	return n
// }

type Sample struct {
	Name   string // aka SampleFilename.
	Volume float64
}

func NewSample(f any) (s Sample) {
	switch f := f.(type) {
	case osu.HitObject:
		return Sample{
			Name:   f.HitSample.Filename,
			Volume: float64(f.HitSample.Volume) / 100,
		}
	}
	return Sample{}
}

func (n Sample) Path(cpath string) (string, bool) {
	if n.Name == "" {
		return "", false
	}
	return filepath.Join(filepath.Dir(cpath), n.Name), true
}
