package piano

import (
	draws "github.com/hndada/gosu/draws5"
)

type HintResources struct {
	img draws.Image
}

type HintOptions struct {
	w float64
	H float64
	x float64
	y float64
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

func (cmp *HintComponent) Update() {
	// Do nothing.
}

func (cmp HintComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
