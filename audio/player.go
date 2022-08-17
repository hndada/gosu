package audio

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

// type Player struct {
// 	audio.Player
// }

func NewStreamer(path string) (io.ReadSeeker, io.Closer) { // (io.ReadSeekCloser, *audio.Player) {
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
	return s, io.Closer(f)
}

// func (p Player) PlaySoundEffect() {
// 	p.Play()
// 	p.Rewind()
// }
