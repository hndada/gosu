package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
)

type FieldRes struct {
	// Field component requires no external resources.
}

type FieldOpts struct {
	stage   draws.WHXY // parent element
	Opacity float64
}

func NewFieldOpts(stage draws.WHXY) FieldOpts {
	return FieldOpts{
		stage:   stage,
		Opacity: 0.8,
	}
}

type FieldComp struct {
	sprite draws.Sprite
}

func NewFieldComp(res FieldRes, opts FieldOpts) (comp FieldComp) {
	img := draws.NewImage(opts.stage.W, opts.stage.H)
	img.Fill(color.NRGBA{0, 0, 0, uint8(255 * opts.Opacity)})

	sprite := draws.NewSprite(img)
	sprite.Locate(opts.stage.X, opts.stage.Y, draws.CenterTop)
	comp.sprite = sprite
	return
}

func (comp FieldComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst, draws.Op{})
}
