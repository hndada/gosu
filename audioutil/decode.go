package audioutil

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

var SampleRate = 44100

var Context *audio.Context = audio.NewContext(SampleRate)

// decode returns streamer, closer, and error.
func decode(apath string) (io.ReadSeeker, func() error, error) { // apath stands for audio path.
	var s io.ReadSeeker
	f, err := os.Open(apath)
	if err != nil {
		return nil, nil, err
	}
	switch strings.ToLower(filepath.Ext(apath)) {
	case ".mp3":
		s, err = mp3.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			return nil, f.Close, err
		}
	case ".wav":
		s, err = wav.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			return nil, f.Close, err
		}
	case ".ogg":
		s, err = vorbis.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			return nil, f.Close, err
		}
	}
	return s, f.Close, nil
}
