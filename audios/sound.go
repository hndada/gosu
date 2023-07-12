package audios

import (
	"io/fs"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// Sound is for playing simple sound effects.
type Sound struct {
	buffer      *beep.Buffer
	volumeScale *float64
}

func NewSound(fsys fs.FS, name string, volumeScale *float64) (Sound, error) {
	streamer, format, err := DecodeFromFile(fsys, name)
	if err != nil {
		return Sound{}, err
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	return Sound{
		buffer:      buffer,
		volumeScale: volumeScale,
	}, nil
}

// Play plays a random sound from the sound pod.
func (s Sound) Play(vol float64) {
	streamer := s.buffer.Streamer(0, s.buffer.Len())
	beepVol := beepVolume(vol * (*s.volumeScale))
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: beepVol}
	speaker.Play(volume)
}
