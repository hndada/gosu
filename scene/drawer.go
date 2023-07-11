package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

func NewBackgroundFromFile(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSpriteFromFile(fsys, name)
	s.MultiplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}

type BackgroundDrawer struct {
	Sprite draws.Sprite
}

func (d BackgroundDrawer) Draw(dst draws.Image) {
	op := draws.Op{}
	value := TheSettings.BackgroundBrightness
	op.ColorM.ChangeHSV(0, 1, value)
	d.Sprite.Draw(dst, op)
}
