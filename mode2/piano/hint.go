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

func NewHintOpts(keys KeysOpts) HintOpts {
	return HintOpts{
		w: keys.stageW,
		H: 24,
		x: keys.StageX,
		y: keys.BaselineY,
	}
}

type HintComp struct {
	sprite draws.Sprite
}

func NewHintComp(res HintRes, opts HintOpts) (comp HintComp) {
	sprite := draws.NewSprite(res.img)
	sprite.SetSize(opts.w, opts.H)
	sprite.Locate(opts.x, opts.y, draws.CenterBottom)
	comp.sprite = sprite
	return
}

func (comp HintComp) Draw(dst draws.Image) {
	comp.sprite.Draw(dst)
}
