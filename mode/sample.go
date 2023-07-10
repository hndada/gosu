package mode

import (
	"path/filepath"

	"github.com/hndada/gosu/audios"
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

// Todo: refactor?
func (s Sample) Path(cpath string) (string, bool) {
	if s.Name == "" {
		return "", false
	}
	return filepath.Join(filepath.Dir(cpath), s.Name), true
}

func (s Sample) Play(vol2, scale float64) {
	if s.Name == "" {
		return
	}
	if s.Volume == 0 {
		s.Volume = vol2
	}
	p := audios.Context.NewPlayerFromBytes(s.Sound)
	s.SoundPlayer.Play(s.Name, s.Volume*scale)
}
