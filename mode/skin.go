package mode

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
	defaultskin "github.com/hndada/gosu/skin"
)

const (
	ScreenSizeX = scene.ScreenSizeX
	ScreenSizeY = scene.ScreenSizeY
)

type skinType struct {
	Score [13]draws.Sprite // number + sign(. , %)
	Combo [10]draws.Sprite // number only
}

// Unexported struct with exported function yields read-only feature.
var (
	defaultSkin skinType
	Skin        skinType
)

func init() { LoadSkin(defaultskin.FS, scene.LoadSkinDefault) }
func LoadSkin(fsys fs.FS, mode scene.LoadSkinMode) {
	skin := &Skin
	if mode == scene.LoadSkinDefault {
		skin = &defaultSkin
	}
	skin.Score = scene.Skin.Number1
	skin.Combo = scene.Skin.Number2
}
