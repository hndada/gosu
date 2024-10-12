package piano

import "github.com/hndada/gosu/draws"

type FieldComponent struct {
	sprite draws.Sprite
}

func NewFieldComponent(res *Resources, opts *Options, keyCount int) (cmp FieldComponent) {
	s := draws.NewSprite(res.FieldImage)
	s.SetSize(opts.StageWidths[keyCount], opts.screenSizeY)
	s.Locate(opts.StagePositionX, 0, draws.CenterTop)
	s.ColorScale.Scale(1, 1, 1, opts.FieldOpacity)
	cmp.sprite = s
	return
}

func (cmp *FieldComponent) Update() {
	// Do nothing.
}

func (cmp FieldComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
