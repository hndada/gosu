package piano

import "github.com/hndada/gosu/draws"

type Hint struct {
	sprite draws.Sprite
}

func NewHint(res *Resources, opts *Options, keyCount int) Hint {
	hint := Hint{}
	s := draws.NewSprite(res.HintImage)
	s.SetSize(opts.StageWidths[keyCount], opts.HintHeight)
	s.Locate(opts.StagePositionX, opts.KeyPositionY, draws.CenterBottom)
	hint.sprite = s
	return hint
}

func (hint *Hint) Update() {
	// Do nothing.
}

func (hint Hint) Draw(dst draws.Image) {
	hint.sprite.Draw(dst)
}
