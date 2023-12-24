package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

type FieldRes struct {
	// Field component requires no external resources.
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
	img := draws.NewImage(opts.w, mode.ScreenH)
	img.Fill(color.NRGBA{0, 0, 0, uint8(255 * opts.Opacity)})

	sprite := draws.NewSprite(img)
	sprite.Locate(opts.x, 0, draws.CenterTop)
	comp.sprite = sprite
	return
}

func (comp FieldComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst, draws.Op{})
}
