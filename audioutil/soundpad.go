package audioutil

import (
	"path/filepath"
)

type SoundMap struct {
	Bytes map[string][]byte // Todo: Bytes -> bytes?
	// Closers []io.Closer
}

func (s *SoundMap) Register(path, key string) error {
	b, err := NewBytes(path)
	if err != nil {
		return err
	}
	if key == "" {
		key = filepath.Base(path)
	}
	s.Bytes[key] = b
	// s.Closers = append(s.Closers, closer)
	return nil
}
func (s *SoundMap) Play(name string) {
	b := s.Bytes[name]
	p := Context.NewPlayerFromBytes(b)
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
