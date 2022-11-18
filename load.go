package gosu

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/scene"
)

func load(fsys fs.FS) {
	settings, err := fs.ReadFile(fsys, "settings.toml")
	if err != nil {
		fmt.Println(err)
	}
	scene.LoadSettings(string(settings))

	skinFS, err := fs.Sub(fsys, "skin")
	if err != nil {
		fmt.Println(err)
	}
	scene.LoadSkin(skinFS, scene.LoadSkinUser)
	soundFS, err := fs.Sub(fsys, "skin/sound")
	if err != nil {
		fmt.Println(err)
	}
	scene.LoadSounds(soundFS, scene.LoadSkinUser)
}
