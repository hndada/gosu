package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Box struct {
	Sprite
	PadW, PadH float64
	Texts      []Text
}

type Text struct {
	Text   string
	Face   font.Face
	Color  color.NRGBA
	Origin Origin
	Align  int
}

const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

func NewBox(i *ebiten.Image, padW, padH float64) Box {
	return Box{
		Sprite: NewSpriteFromImage(i),
		PadW:   padW,
		PadH:   padH,
	}
}
func (b Box) Draw(screen *ebiten.Image) {
	b.Sprite.Draw(screen, nil)
	for _, txt := range b.Texts {
		var x, y int
		text.Draw(screen, txt.Text, txt.Face, x, y, txt.Color)
	}
}
