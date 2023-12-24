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
	// Field component requires no external resources.
	res.img = draws.NewImage(mode.ScreenW, mode.ScreenH)
}

type FieldOpts struct {
	w       float64
	x       float64
	Opacity float64
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
	sprite := draws.NewSprite(res.img)
	sprite.SetSize(opts.w, mode.ScreenH)
	sprite.Locate(opts.x, 0, draws.CenterTop)
	sprite.Color.Scale(1, 1, 1, opts.Opacity)
	comp.sprite = sprite
	return
}

func (comp FieldComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst)
}
