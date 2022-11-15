package audios

import (
	"io"
	"path/filepath"
	"strings"
)

// type SoundMap struct {
// 	Bytes map[string][]byte // Todo: Bytes -> bytes?
// 	// Closers []io.Closer
// }

type SoundMap struct {
	bytes map[string][]byte
	vol   *float64
}

func NewSoundMap(vol *float64) SoundMap {
	return SoundMap{
		bytes: make(map[string][]byte),
		vol:   vol,
	}
}
func (s SoundMap) Register(path string) error {
	b, err := NewBytes(path)
	if err != nil {
		return err
	}
	name := filepath.Base(path)
	if pos := strings.LastIndexByte(name, '.'); pos != -1 {
		name = name[:pos]
	}
	s.bytes[name] = b
	return nil
}

//	func (s SoundMap) Register(path, key string) error {
//		b, err := NewBytes(path)
//		if err != nil {
//			return err
//		}
//		if key == "" {
//			key = filepath.Base(path)
//		}
//		s.bytes[key] = b
//		// s.Closers = append(s.Closers, closer)
//		return nil
//	}
func (s SoundMap) Play(name string) {
	p := Context.NewPlayerFromBytes(s.bytes[name])
	p.SetVolume(*s.vol)
	p.Play()
}
func (s SoundMap) PlayWithVolume(name string, vol2 float64) {
	p := Context.NewPlayerFromBytes(s.bytes[name])
	p.SetVolume(*s.vol * vol2)
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
