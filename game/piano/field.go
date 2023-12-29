package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type FieldRes struct {
	img draws.Image
}

func (res *FieldRes) Load(fsys fs.FS) {
	// Uses generated image.
	res.img = draws.NewImage(game.ScreenW, game.ScreenH)
}

type FieldOpts struct {
	w       float64
	x       float64
	Opacity float32
}

func NewFieldOpts(stage StageOpts) FieldOpts {
	return FieldOpts{
		w:       stage.w,
		x:       stage.X,
		Opacity: 0.8,
	}
}

type FieldComp struct {
	sprite draws.Sprite
}

func NewFieldComp(res FieldRes, opts FieldOpts) (comp FieldComp) {
	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, game.ScreenH)
	s.Locate(opts.x, 0, draws.CenterTop)
	s.ColorScale.Scale(1, 1, 1, opts.Opacity)
	comp.sprite = s
	return
}

func (comp FieldComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst)
}
