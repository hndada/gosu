package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
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
	y float64
}

func NewHintOpts(key KeyOpts) HintOpts {
	return HintOpts{
		w: key.stageW,
		H: 0.05 * mode.ScreenH,
		x: key.StageX,
		y: key.BaselineY,
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
	comp.sprite.Draw(dst, draws.Op{})
}
