package audioutil

import (
	"io"
	"path/filepath"
)

// type SoundMap struct {
// 	Bytes map[string][]byte // Todo: Bytes -> bytes?
// 	// Closers []io.Closer
// }

type SoundMap map[string][]byte

// NewBytes is for short sounds: long audio file will make the game stutter.
// No returns closer since NewPlayerFromBytes needs no Closer: no any files are open.
func NewBytes(path string) ([]byte, error) {
	s, _, err := decode(path)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s SoundMap) Register(path, key string) error {
	b, err := NewBytes(path)
	if err != nil {
		return err
	}
	if key == "" {
		key = filepath.Base(path)
	}
	s[key] = b
	// s.Closers = append(s.Closers, closer)
	return nil
}
func (s SoundMap) Play(name string, vol float64) {
	p := Context.NewPlayerFromBytes(s[name])
	p.SetVolume(vol)
	p.Play()
}

// func (s *SoundPad) Close() {
// 	for _, closer := range s.Closers {
// 		closer.Close()
// 	}
// 	// Zeroize.
// 	s.Bytes = make(map[string][]byte)
// 	s.Closers = make([]io.Closer, 0)
// }
