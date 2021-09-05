package audio

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const bytesPerSample = 4
const sampleRate = 44100

var Context = audio.NewContext(sampleRate)

type Player struct {
	*audio.Player
}

func NewPlayer(path string) *Player {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var s audioStream
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		s, err = mp3.Decode(Context, bytes.NewReader(b))
	case ".wav":
		s, err = wav.Decode(Context, bytes.NewReader(b))
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

// func (p *Player) SetVolume(v float64) { p.SetVolume(v) }
