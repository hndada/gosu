package main

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
const SampleRate = 44100

var Context *audio.Context = audio.NewContext(SampleRate)

type AudioStreamer interface {
	io.ReadSeeker
	Length() int64
}

func NewAudioStreamer(path string) (AudioStreamer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		return mp3.Decode(Context, bytes.NewReader(b))
	case ".wav":
		return wav.Decode(Context, bytes.NewReader(b))
	case ".ogg":
		return vorbis.Decode(Context, bytes.NewReader(b))
	}
	p, err := audio.NewPlayer(Context, s)
	if err != nil {
		panic(err)
	}
	return &Player{
		Player: p,
	}
}
func NewSEPlayer(path string, vol int) func() {
	// ap := NewPlayer(path)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
