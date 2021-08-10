package game

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

const sampleRate = 44100

var AudioContext *audio.Context

func init() {
	var err error
	AudioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}
}

type AudioPlayer struct {
	context   *audio.Context
	player    *audio.Player
	total     time.Duration
	seBytes   []byte
	seCh      chan []byte
	volume128 int
}

// todo: timeRate
func NewAudioPlayer(audioPath string, volume int) *AudioPlayer {
	b, err := ioutil.ReadFile(audioPath)
	if err != nil {
		panic(err)
	}
	type audioStream interface {
		audio.ReadSeekCloser
		Length() int64
	}
	var s audioStream
	switch strings.ToLower(filepath.Ext(audioPath)) {
	case ".mp3":
		s, err = mp3.Decode(AudioContext, audio.BytesReadSeekCloser(b))
	case ".wav":
		s, err = wav.Decode(AudioContext, audio.BytesReadSeekCloser(b))
	}
	p, err := audio.NewPlayer(AudioContext, s)
	if err != nil {
		panic(err)
	}

	const bytesPerSample = 4
	const sampleRate = 44100
	player := &AudioPlayer{
		context:   AudioContext,
		player:    p,
		total:     time.Millisecond * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		seCh:      make(chan []byte),
		volume128: 128 / 4,
	}
	if player.total == 0 {
		player.total = 1
	}
	return player
}

func (ap *AudioPlayer) Play() {
	_ = ap.player.Play()
}
func (ap *AudioPlayer) Pause() {
	_ = ap.player.Pause()
}
func (ap *AudioPlayer) Rewind() {
	_ = ap.player.Rewind()
}
func (ap *AudioPlayer) Close() error {
	return ap.player.Close()
}
func (ap *AudioPlayer) Seek(offset time.Duration) {
	_ = ap.player.Seek(offset)
}
func (ap *AudioPlayer) Time() time.Duration {
	// ap.current = ap.player.Current()
	return ap.player.Current()
}
