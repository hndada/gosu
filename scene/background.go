package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type Background struct {
	defaultSprite draws.Sprite
	sprite        draws.Sprite
	screenSize    draws.XY
	brightness    *float32
}

// In osu!, background brightness at Song selectis 60% (153 / 255).
// However, for the sake of simplicity, gosu will use the option value.
func NewBackground(res *Resources, opts *Options) (bg Background) {
	bg.defaultSprite = bg.newSprite(res.DefaultBackgroundImage)
	bg.sprite = bg.defaultSprite
	bg.screenSize = opts.screenSize
	// fmt.Println(opts.BackgroundBrightness)
	bg.brightness = &opts.BackgroundBrightness
	return
}

func (bg Background) newSprite(img draws.Image) draws.Sprite {
	s := draws.NewSprite(img)
	s.Scale(bg.screenSize.X / s.W())
	s.Locate(bg.screenSize.X/2, bg.screenSize.Y/2, draws.CenterMiddle)
	return s
}

func (bg *Background) UpdateBackground(fsys fs.FS, name string) {
	img := draws.NewImageFromFile(fsys, name)
	if img.IsEmpty() {
		bg.sprite = bg.defaultSprite
	} else {
		bg.sprite = bg.newSprite(img)
	}
}

func (bg Background) Draw(dst draws.Image) {
	// op.ColorM.ChangeHSV(0, 1, opts.Brightness)
	bg.sprite.ColorScale.ScaleAlpha(*bg.brightness)
	bg.sprite.Draw(dst)
}
