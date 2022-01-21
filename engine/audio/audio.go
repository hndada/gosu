package audio

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const bytesPerSample = 4
const sampleRate = 44100

var Context *audio.Context = audio.NewContext(sampleRate)

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
		if err != nil {
			panic(err)
		}
	case ".wav":
		s, err = wav.Decode(Context, bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
	case ".ogg":
		s, err = vorbis.Decode(Context, bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
	}
	p, err := audio.NewPlayer(Context, s)
	if err != nil {
		panic(err)
	}
	p2 := &Player{
		Player: p,
	}
	return p2
}
func NewSEPlayer(path string, vol int) func() {
	ap := NewPlayer(path)
	ap.SetVolume(float64(vol) / 100)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
