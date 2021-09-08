package audio

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

// var
func NewPlayer2(path string) beep.StreamSeekCloser {
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
		err      error
	)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			panic(err)
		}
	case ".wav":
		streamer, format, err = wav.Decode(f)
		if err != nil {
			panic(err)
		}
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
		if err != nil {
			panic(err)
		}
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// done := make(chan bool)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// done <- true
	// })))
	// <-done
	// speaker.Play(streamer)
	return streamer
}
