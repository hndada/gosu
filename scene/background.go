package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type BackgroundResources struct {
	defaultImage draws.Image
}

func (res *BackgroundResources) Load(fsys fs.FS) {
	fname := "interface/default-bg.png"
	res.defaultImage = draws.NewImageFromFile(fsys, fname)
}

type BackgroundOptions struct {
	Brightness float64
}

type BackgroundComponent struct {
	sprite draws.Sprite
}

func NewBackgroundComponent(res BackgroundResources, opts BackgroundOptions) (cmp BackgroundComponent) {
	// s := draws.NewSprite(res.defaultImage)
	s := draws.NewSpriteFromFile(fsys, name)
	if s.IsEmpty() {
		s = asset.DefaultBackgroundSprite
	}

	s.MultiplyScale(cfg.ScreenSize.X / s.Width())
	s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.CenterMiddle)
	value := cfg.BackgroundBrightness
	op.ColorM.ChangeHSV(0, 1, value)
	cmp.sprite = s
	return
}

func (cmp BackgroundComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
