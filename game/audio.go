package game

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

const bytesPerSample = 4
const sampleRate = 44100

var AudioContext *audio.Context

func init() {
	var err error
	AudioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}
}

func NewAudioPlayer(path string) *audio.Player {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	type audioStream interface {
		audio.ReadSeekCloser
		Length() int64
	}
	var s audioStream
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		s, err = mp3.Decode(AudioContext, audio.BytesReadSeekCloser(b))
	case ".wav":
		s, err = wav.Decode(AudioContext, audio.BytesReadSeekCloser(b))
	}
	p, err := audio.NewPlayer(AudioContext, s)
	p.SetVolume(0.25) // temp
	if err != nil {
		panic(err)
	}
	return p
}
