package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

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

type Box struct {
	Sprite
	PadW, PadH float64
	Texts      []Text
}

func NewBox(i *ebiten.Image, padW, padH float64) Box {
	return Box{
		Sprite: NewSpriteFromImage(i),
		PadW:   padW,
		PadH:   padH,
	}
}
func NewBoxImage(w, h, border int, outerColor, innerColor color.NRGBA) *ebiten.Image {
	// w = math.Ceil(w)
	// h = math.Ceil(h)
	// outer := ebiten.NewImage(int(w), int(h))
	// inner := ebiten.NewImage(w - 2*border)
	outer := ebiten.NewImage(w, h)
	outer.Fill(outerColor)
	inner := ebiten.NewImage(w-2*border, h-2*border)
	inner.Fill(innerColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(border), float64(border))
	outer.DrawImage(inner, op)
	return outer
}

func (b Box) Draw(screen *ebiten.Image) {
	b.Sprite.Draw(screen, nil)
	for _, txt := range b.Texts {
		var x, y int
		text.Draw(screen, txt.Text, txt.Face, x, y, txt.Color)
	}
}
func DesiredBoxHeight(f font.Face, padH float64) float64 {
	const standard = "0"
	rect := text.BoundString(f, standard)
	return padH*2 - float64(rect.Min.Y)
}
