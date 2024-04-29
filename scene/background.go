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
	Brightness   float32
	screenWidth  *float64
	screenHeight *float64
}

// Todo: *Options vs Options
// But I think, to use pointer, *Options is inevitable.
func NewBackgroundOptions(opts *Options) BackgroundOptions {
	return BackgroundOptions{
		Brightness:   0.8,
		screenWidth:  &opts.Resolution.X,
		screenHeight: &opts.Resolution.Y,
	}
}

type BackgroundComponent struct {
	sprite     draws.Sprite
	brightness *float32
}

func NewBackgroundComponent(res BackgroundResources, opts BackgroundOptions) (cmp BackgroundComponent) {
	// s := draws.NewSprite(res.defaultImage)
	s := draws.NewSpriteFromFile(fsys, name)
	if s.IsEmpty() {
		s = asset.DefaultBackgroundSprite
	}

	s.MultiplyScale(opts.screenWidth / s.Width())
	s.Locate(*opts.screenWidth/2, *opts.screenHeight/2, draws.CenterMiddle)
	cmp.sprite = s
	cmp.brightness = &opts.Brightness
	return
}

func (cmp BackgroundComponent) Draw(dst draws.Image) {
	// op.ColorM.ChangeHSV(0, 1, opts.Brightness)
	a := cmp.sprite.ColorScale.A()
	cmp.sprite.ColorScale.SetA(a * *cmp.brightness)
	cmp.sprite.Draw(dst)
}
