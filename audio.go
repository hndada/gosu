package gosu

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const SampleRate = 44100

var Context *audio.Context = audio.NewContext(SampleRate)

type AudioPlayer struct {
	*audio.Player
}

func NewAudioPlayer(path string) (io.ReadSeekCloser, AudioPlayer) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	var s io.ReadSeeker
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		s, err = mp3.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			panic(err)
		}
	case ".wav":
		s, err = wav.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			panic(err)
		}
	case ".ogg":
		s, err = vorbis.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			panic(err)
		}
	}
	ap, err := Context.NewPlayer(s)
	if err != nil {
		panic(err)
	}
	return f, AudioPlayer{ap}
}
func (ap AudioPlayer) PlaySoundEffect() {
	ap.Play()
	ap.Rewind()
}
