package choose

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(scene.ScreenSizeX / s.W())
	s.Locate(scene.ScreenSizeX/2, scene.ScreenSizeY/2, draws.CenterMiddle)
	return s
}
