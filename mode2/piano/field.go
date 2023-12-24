package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

type FieldRes struct {
	img draws.Image
}

func (res *FieldRes) Load(fsys fs.FS) {
	// Uses generated image.
	res.img = draws.NewImage(mode.ScreenW, mode.ScreenH)
}

type FieldOpts struct {
	w       float64
	x       float64
	Opacity float32
}

func NewFieldOpts(keys KeysOpts) FieldOpts {
	return FieldOpts{
		w:       keys.stageW,
		x:       keys.StageX,
		Opacity: 0.8,
	}
}

type FieldComp struct {
	sprite draws.Sprite
}

func NewFieldComp(res FieldRes, opts FieldOpts) (comp FieldComp) {
	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, mode.ScreenH)
	s.Locate(opts.x, 0, draws.CenterTop)
	s.ColorScale.Scale(1, 1, 1, opts.Opacity)
	comp.sprite = s
	return
}

func (comp FieldComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst)
}
