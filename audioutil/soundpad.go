package audioutil

import (
	"io"
	"path/filepath"
)

type SoundPad struct {
	Bytes   map[string][]byte
	Closers []io.Closer
}

func (s *SoundPad) Register(path, key string) error {
	b, closer, err := NewBytes(path)
	if err != nil {
		return err
	}
	if key == "" {
		key = filepath.Base(path)
	}
	s.Bytes[key] = b
	s.Closers = append(s.Closers, closer)
	return nil
}
func (s *SoundPad) Play(name string) {
	b := s.Bytes[name]
	p := Context.NewPlayerFromBytes(b)
	p.Play()
}
func (s *SoundPad) Close() {
	for _, closer := range s.Closers {
		closer.Close()
	}
	// Zeroize.
	s.Bytes = make(map[string][]byte)
	s.Closers = make([]io.Closer, 0)
}
