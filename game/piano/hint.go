package piano

import "github.com/hndada/gosu/draws"

type HintComponent struct {
	sprite draws.Sprite
}

func NewHintComponent(res *Resources, opts *Options, keyCount int) (cmp HintComponent) {
	s := draws.NewSprite(res.HintImage)
	s.SetSize(opts.StageWidths[keyCount], opts.HintHeight)
	s.Locate(opts.StagePositionX, opts.KeyPositionY, draws.CenterBottom)
	cmp.sprite = s
	return
}

func (cmp *HintComponent) Update() {
	// Do nothing.
}

func (cmp HintComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
