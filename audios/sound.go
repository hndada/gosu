package audios

import (
	"io"
	"io/fs"
	"math/rand"
	"path"
	"path/filepath"
	"strings"
)

type Sound []byte

// NewSound is for effect sounds, which is short.
// Long audio file will make the game stutter.
// No closers are related since no any files are open.
func NewSound(fsys fs.FS, name string) Sound {
	s, _, err := decode(fsys, name)
	if err != nil {
		return nil
	}
	b, err := io.ReadAll(s)
	if err != nil {
		return nil
	}
	return b
}
func (s Sound) IsValid() bool { return s != nil }
func (s Sound) Play(vol float64) {
	p := Context.NewPlayerFromBytes(s)
	p.SetVolume(vol)
	p.Play()
}

// SoundBag is for playing one of effects in the slice. Useful for
// playing slightly different effect when doing same actions.
type SoundBag []Sound

// Todo: remove redundancy with NewImages()?
func NewSoundBag(fsys fs.FS, name string) SoundBag {
	// name supposed to have no extension when passed in NewSoundBag.
	var sb SoundBag
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	one := SoundBag{NewSound(fsys, name+ext)}
	fs, err := fs.ReadDir(fsys, name)
	if err != nil {
		return one
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		name := path.Join(name, f.Name())
		sb = append(sb, NewSound(fsys, name))
	}
	return sb
}
func (sb SoundBag) Play(vol float64) {
	i := rand.Intn(len(sb))
	sb[i].Play(vol)
}
