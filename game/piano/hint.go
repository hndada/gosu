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
	w float64
	H float64
	x float64
	y float64 // center bottom
}

func NewHintOpts(stage StageOpts) HintOpts {
	return HintOpts{
		w: stage.w,
		H: 24,
		x: stage.X,
		y: stage.H,
	}
}

type HintComp struct {
	sprite draws.Sprite
}

func NewHintComp(res HintRes, opts HintOpts) (comp HintComp) {
	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, opts.H)
	s.Locate(opts.x, opts.y, draws.CenterBottom)
	comp.sprite = s
	return
}

func (comp HintComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst)
}
