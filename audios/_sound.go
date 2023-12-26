package audios

import (
	"io/fs"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

// Sound is for playing simple sound effects.
type Sound struct {
	buffer      *beep.Buffer
	volumeScale *float64
}

// NewSound returns an empty struct when error occurs
// because most of error just comes from file not found.
func NewSound(fsys fs.FS, name string, volumeScale *float64) Sound {
	streamer, format, err := DecodeFromFile(fsys, name)
	if err != nil {
		return Sound{}
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	return Sound{buffer, volumeScale}
}

// Play plays a random sound from the sound pod.
func (s Sound) Play(vol float64) {
	streamer := s.buffer.Streamer(0, s.buffer.Len())
	beepVol := beepVolume(vol * (*s.volumeScale))
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: beepVol}
	speaker.Play(volume)
}
