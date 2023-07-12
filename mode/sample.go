package mode

import "github.com/hndada/gosu/format/osu"

type Sample struct {
	Name   string // aka SampleFilename.
	Volume float64
}

func NewSample(f any) (s Sample) {
	switch f := f.(type) {
	case osu.HitObject:
		return newSampleFromOsu(f)
	}
	return Sample{}
}
func newSampleFromOsu(f osu.HitObject) (s Sample) {
	return Sample{
		Name:   f.HitSample.Filename,
		Volume: float64(f.HitSample.Volume) / 100,
	}
}
