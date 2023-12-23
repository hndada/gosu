package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type HintRes struct {
	img draws.Image
}

func (res *HintRes) Load(fsys fs.FS) {
	fname := "piano/stage/hint.png"
	res.img = draws.NewImageFromFile(fsys, fname)
	return
}

type HintOpts struct {
	stage    draws.WHXY // parent element
	baseline draws.Position
	RH       float64
}

func NewHintOpts(stage draws.WHXY, baseline draws.Position) HintOpts {
	return HintOpts{
		stage:    stage,
		baseline: baseline,
		RH:       0.05,
	}
}

type HintComp struct {
	sprite draws.Sprite
}

func NewHintComp(res HintRes, opts HintOpts) (comp HintComp) {
	sprite := draws.NewSprite(res.img)
	sprite.SetSize(opts.stage.W, opts.RH*opts.stage.H)
	sprite.Locate(opts.baseline.X, opts.baseline.Y, draws.CenterBottom)
	comp.sprite = sprite
	return
}

func (comp HintComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst, draws.Op{})
}
