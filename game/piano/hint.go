package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type HintResources struct {
	img draws.Image
}

func (res *HintResources) Load(fsys fs.FS) {
	fname := "piano/stage/hint.png"
	res.img = draws.NewImageFromFile(fsys, fname)
}

type HintOptions struct {
	w float64
	H float64
	x float64
	y float64 // center bottom
}

func NewHintOptions(stage StageOptions) HintOptions {
	return HintOptions{
		w: stage.w,
		H: 24,
		x: stage.X,
		y: stage.H,
	}
}

type HintComponent struct {
	sprite draws.Sprite
}

func NewHintComponent(res HintResources, opts HintOptions) (cmp HintComponent) {
	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, opts.H)
	s.Locate(opts.x, opts.y, draws.CenterBottom)
	cmp.sprite = s
	return
}

func (cmp HintComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
