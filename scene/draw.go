package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

func NewBackgroundSprite(fsys fs.FS, name string,
	screenSize draws.Vector2) draws.Sprite {

	s := draws.NewSpriteFromFile(fsys, name)
	s.MultiplyScale(screenSize.X / s.W())
	s.Locate(screenSize.X/2, screenSize.Y/2, draws.CenterMiddle)
	return s
}

func NewDrawBackgroundFunc(s draws.Sprite,
	screenSize draws.Vector2, bgBrightness *float64) func(draws.Image) {

	return func(dst draws.Image) {
		op := draws.Op{}
		value := *bgBrightness
		op.ColorM.ChangeHSV(0, 1, value)
		s.Draw(dst, op)
	}
}
