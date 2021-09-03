package audio

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

// todo: oto v2 import
const bytesPerSample = 4
const sampleRate = 44100

var Context *audio.Context

type Player struct {
	*audio.Player
}

func init() {
	var err error
	Context, err = audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}
}

func NewPlayer(path string) *Player {
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
		s, err = mp3.Decode(Context, audio.BytesReadSeekCloser(b))
	case ".wav":
		s, err = wav.Decode(Context, audio.BytesReadSeekCloser(b))
	}
	p, err := audio.NewPlayer(Context, s)
	p.SetVolume(0.25) // temp
	if err != nil {
		panic(err)
	}
	p2 := &Player{
		Player: p,
	}
	return p2
}
