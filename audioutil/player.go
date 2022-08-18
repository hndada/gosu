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

const SampleRate = 44100

var Context *audio.Context = audio.NewContext(SampleRate)

// type Player struct {
// 	audio.Player
// }

// (io.ReadSeeker, io.Closer, error) { // (io.ReadSeekCloser, *audio.Player) {
// No need io.Closer I think.
// ([]byte, io.Closer, error) {
func NewBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var s io.ReadSeeker
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		s, err = mp3.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			return nil, err
		}
	case ".wav":
		s, err = wav.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			return nil, err
		}
	case ".ogg":
		s, err = vorbis.DecodeWithSampleRate(SampleRate, io.ReadSeekCloser(f))
		if err != nil {
			return nil, err
		}
	}
	b, err := io.ReadAll(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//	func (p Player) PlaySoundEffect() {
//		p.Play()
//		p.Rewind()
//	}
// func PlayAudio(b []byte) {
// 	p := Context.NewPlayerFromBytes(b)
// 	p.Play()
// }
