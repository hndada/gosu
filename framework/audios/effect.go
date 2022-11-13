package audios

import (
	"io"
	"io/fs"
)

func PlayEffect(src []byte, vol float64) {
	p := Context.NewPlayerFromBytes(src)
	p.SetVolume(vol)
	p.Play()
}

// NewBytes is for short sounds: long audio file will make the game stutter.
// No returns closer since NewPlayerFromBytes needs no Closer: no any files are open.
func NewBytes(fsys fs.FS, name string) ([]byte, error) {
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
