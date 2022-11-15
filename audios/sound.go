package audios

import (
	"io"
	"io/fs"
	"math/rand"
)

type Sound []byte

// NewSound is for effect sounds, which is short.
// Long audio file will make the game stutter.
// No closers are related since no any files are open.
func NewSound(fsys fs.FS, name string) (Sound, error) {
	s, _, err := decode(fsys, name)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (e Sound) Play(vol float64) {
	p := Context.NewPlayerFromBytes(e)
	p.SetVolume(vol)
	p.Play()
}

// Sounds is for playing one of effects in the slice. Useful for
// playing slightly different effect when doing same actions.
type Sounds []Sound

func (es Sounds) Play(vol float64) {
	i := rand.Intn(len(es))
	es[i].Play(vol)
}
