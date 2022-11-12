package audios

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// var bufferSize = 100 * time.Millisecond

func NewPlayer(path string) (*audio.Player, func() error, error) {
	s, closer, err := decode(path)
	if err != nil {
		return nil, nil, err
	}
	p, err := Context.NewPlayer(s)
	if err != nil {
		return nil, closer, err
	}
	// p.SetBufferSize(bufferSize)
	return p, closer, err
}
func PlayEffect(src []byte, vol float64) {
	p := Context.NewPlayerFromBytes(src)
	p.SetVolume(vol)
	p.Play()
}

// func NewStreamer(path string) (io.ReadSeeker, func() error, error) {
// 	s, closer, err := decode(path)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return s, f.Close, nil
// }
