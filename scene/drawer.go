package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.LoadSprite(fsys, name)
	s.MultiplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}
