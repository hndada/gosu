package gosu

import (
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type AudioPlayer struct {
	context *audio.Context
	player  *audio.Player
	// current   time.Duration
	total     time.Duration
	seBytes   []byte
	seCh      chan []byte
	volume128 int
	//
	// sampleRate beep.SampleRate
	// streamer   beep.StreamSeeker
	// ctrl       *beep.Ctrl
	// resampler  *beep.Resampler
	// volume     *effects.Volume
}

// type MusicType int
//
// const (
// 	MusicTypeMP3 = iota
// 	MusicTypeOGG
// )

// todo: timeRate
func NewAudioPlayer(context *audio.Context, audioPath string) *AudioPlayer {
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
		s, err = mp3.Decode(context, audio.BytesReadSeekCloser(b))
	}
	p, err := audio.NewPlayer(context, s)
	if err != nil {
		panic(err)
	}

	const bytesPerSample = 4
	const sampleRate = 44100
	player := &AudioPlayer{
		context:   context,
		player:    p,
		total:     time.Millisecond * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		seCh:      make(chan []byte),
		volume128: 128,
	}
	if player.total == 0 {
		player.total = 1
	}
	return player
}

func (ap *AudioPlayer) Play() {
	ap.player.Play()
}
func (ap *AudioPlayer) Close() error {
	return ap.player.Close()
}
func (ap *AudioPlayer) Time() time.Duration {
	// ap.current = ap.player.Current()
	return ap.player.Current()
}
