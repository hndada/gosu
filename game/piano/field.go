package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type FieldResources struct {
	img draws.Image
}

func (res *FieldResources) Load(fsys fs.FS) {
	// Uses generated image.
	res.img = draws.NewImage(game.ScreenW, game.ScreenH)
}

type FieldOptions struct {
	w       float64
	x       float64
	Opacity float32
}

func NewFieldOptions(stage StageOptions) FieldOptions {
	return FieldOptions{
		w:       stage.w,
		x:       stage.X,
		Opacity: 0.8,
	}
}

type FieldComponent struct {
	sprite draws.Sprite
}

func NewFieldComponent(res FieldResources, opts FieldOptions) (cmp FieldComponent) {
	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, game.ScreenH)
	s.Locate(opts.x, 0, draws.CenterTop)
	s.ColorScale.Scale(1, 1, 1, opts.Opacity)
	cmp.sprite = s
	return
}

func (cmp FieldComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
