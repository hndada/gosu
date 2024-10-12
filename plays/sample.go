package plays

import "github.com/hndada/gosu/format/osu"

type Sample struct {
	Filename string
	Volume   float64
}

var DefaultSample = Sample{Filename: "", Volume: 1.0}

func NewSample(f any) (s Sample) {
	switch f := f.(type) {
	case osu.HitObject:
		return newSampleFromOsu(f)
	}
	return
}

func newSampleFromOsu(f osu.HitObject) (s Sample) {
	return Sample{
		Filename: f.SampleFilename(),
		Volume:   float64(f.HitSample.Volume) / 100,
	}
}
