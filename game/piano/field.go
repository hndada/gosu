package piano

import "github.com/hndada/gosu/draws"

type Field struct {
	sprite draws.Sprite
}

func NewField(res *Resources, opts *Options, keyCount int) Field {
	field := Field{}
	s := draws.NewSprite(res.FieldImage)
	s.SetSize(opts.StageWidths[keyCount], opts.screenSizeY)
	s.Locate(opts.StagePositionX, 0, draws.CenterTop)
	s.ColorScale.Scale(1, 1, 1, opts.FieldOpacity)
	field.sprite = s
	return field
}

func (cmp *Field) Update() {
	// Do nothing.
}

func (cmp Field) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
