package audios

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

func NewPlayer(path string) (*audio.Player, func() error, error) {
	s, closer, err := decode(path)
	if err != nil {
		return nil, nil, err
	}
	p, err := Context.NewPlayer(s)
	if err != nil {
		return nil, closer, err
	}
	return p, closer, err
}

// func NewStreamer(path string) (io.ReadSeeker, func() error, error) {
// 	s, closer, err := decode(path)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return s, f.Close, nil
// }
