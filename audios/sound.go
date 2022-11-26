package audios

import (
	"io"
	"io/fs"
	"math/rand"
	"path"
	"path/filepath"
	"strings"
)

type Sounder interface {
	Play(vol float64)
}
type Sound []byte

// NewSound is for effect sounds, which is short.
// Long audio file will make the game stutter.
// No closers are related since no any files are open.
func NewSound(fsys fs.FS, name string) Sound {
	streamer, _, err := decode(fsys, name)
	if err != nil {
		return nil
	}
	b, err := io.ReadAll(streamer)
	if err != nil {
		return nil
	}
	return b
}
func (s Sound) Play(vol float64) {
	p := Context.NewPlayerFromBytes(s)
	p.SetVolume(vol)
	p.Play()
}
func (s Sound) IsValid() bool { return s != nil }

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
	if len(sb) == 0 {
		return
	}
	i := rand.Intn(len(sb))
	sb[i].Play(vol)
}

// A player for sample sound is generated at a place.
type SoundPlayer struct {
	Sounds map[string]Sound
	Volume *float64
	// Player *audio.Player
}

// Todo: need a test
func NewSoundPlayer(fsys fs.FS, vol *float64) (s SoundPlayer) {
	const oneMB = 1024 * 1024
	s.Sounds = make(map[string]Sound)
	s.Volume = vol
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			return err
		}
		if info.Size() > oneMB {
			return nil
		}
		switch ext := filepath.Ext(path); ext {
		case ".wav", ".WAV", ".ogg", ".OGG":
			s.Sounds[path] = NewSound(fsys, path)
		}
		return nil
	})
	return
}
func (s SoundPlayer) Play(name string, vol2 float64) {
	vol := (*s.Volume) * vol2
	s.Sounds[name].Play(vol)
}
