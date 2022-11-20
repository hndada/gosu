package gosu

import (
	"io/fs"
	"os"

	"github.com/hndada/gosu/mode/piano"
)

func NewGamePiano(fsys fs.FS) *game {
	load(fsys)
	scene, err := piano.NewScenePlay(ZipFS("test.osz"), "nekodex - circles! (MuangMuangE) [Hard].osu", nil, nil)
	if err != nil {
		panic(err)
	}
	g := &game{
		FS:    fsys,
		Scene: scene,
	}
	return g
}
func NewGameDrum(fsys fs.FS) *game {
	load(fsys)
	scene, err := piano.NewScenePlay(os.DirFS("asdf - 1223"), "asdf - 1223 (MuangMuangE) [Oni].osu", nil, nil)
	if err != nil {
		panic(err)
	}
	g := &game{
		FS:    fsys,
		Scene: scene,
	}
	return g
}
