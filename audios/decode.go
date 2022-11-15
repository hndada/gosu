package audios

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const SampleRate = 44100

var Context = audio.NewContext(SampleRate)

func decode(fsys fs.FS, name string) (streamer io.ReadSeeker, close func() error, err error) {
	// var f fs.File
	// f, err = fsys.Open(name)
	// if err != nil {
	// 	return
	// }
	// close = f.Close
	dat, err := fs.ReadFile(fsys, name)
	if err != nil {
		return
	}
	r := bytes.NewReader(dat)
	close = func() error { return nil } // Todo: need a test
	switch filepath.Ext(name) {
	case ".mp3", ".MP3":
		streamer, err = mp3.DecodeWithSampleRate(SampleRate, r)
	case ".wav", ".WAV":
		streamer, err = wav.DecodeWithSampleRate(SampleRate, r)
	case ".ogg", ".OGG":
		streamer, err = vorbis.DecodeWithSampleRate(SampleRate, r)
	}
	return
}
