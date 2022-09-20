package gosu

import (
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

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
