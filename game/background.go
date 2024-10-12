package game

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type BackgroundComponent struct {
	defaultSprite draws.Sprite
	screenSize    draws.XY
	sprite        draws.Sprite
	brightness    *float32
}

func (cmp BackgroundComponent) newSprite(img draws.Image) draws.Sprite {
	s := draws.NewSprite(img)
	s.Scale(cmp.screenSize.X / s.W())
	s.Locate(cmp.screenSize.X/2, cmp.screenSize.Y/2, draws.CenterMiddle)
	return s
}

// In osu!, background brightness at Song selectis 60% (153 / 255).
// However, for the sake of simplicity, gosu will use the option value.
func NewBackgroundComponent(res *Resources, opts *Options) (cmp BackgroundComponent) {
	cmp.defaultSprite = cmp.newSprite(res.DefaultBackgroundImage)
	cmp.brightness = &opts.BackgroundBrightness
	return
}

func (cmp *BackgroundComponent) UpdateBackground(fsys fs.FS, name string) {
	img := draws.NewImageFromFile(fsys, name)
	if img.IsEmpty() {
		cmp.sprite = cmp.defaultSprite
	} else {
		cmp.sprite = cmp.newSprite(img)
	}
}

func (cmp BackgroundComponent) Draw(dst draws.Image) {
	// op.ColorM.ChangeHSV(0, 1, opts.Brightness)
	cmp.sprite.ColorScale.ScaleAlpha(*cmp.brightness)
	cmp.sprite.Draw(dst)
}
