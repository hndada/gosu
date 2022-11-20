package gosu

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
)

func NewGamePiano(fsys fs.FS) *game {
	load(fsys)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	scene, err := piano.NewScenePlay(ZipFS(filepath.Join(dir, "test.osz")), "nekodex - circles! (MuangMuangE) [Hard].osu", nil, nil)
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	scene, err := drum.NewScenePlay(os.DirFS(path.Join(dir, "asdf - 1223")), "asdf - 1223 (MuangMuangE) [Oni].osu", nil, nil)
	if err != nil {
		panic(err)
	}
	g := &game{
		FS:    fsys,
		Scene: scene,
	}
	return g
}
