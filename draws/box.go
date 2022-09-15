package draws

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Box struct {
	Sprite
	padW, padH float64
	Texts      []Text
}

func NewBox(i *ebiten.Image, padW, padH float64) Box {
	return Box{
		Sprite: NewSpriteFromImage(i),
		padW:   padW,
		padH:   padH,
	}
}
func (b Box) Draw(screen *ebiten.Image) {
	b.Sprite.Draw(screen, nil)
	for _, txt := range b.Texts {
		txt.Draw(screen, b)
	}
}

type Text struct {
	text   string
	color  color.NRGBA
	face   font.Face
	origin Origin // Text's align is determined by its Origin.

	w, h float64
	// x, y   float64
	// Align  int
}

func NewText(text string) (t Text) {
	t.text = text
	t.color = color.NRGBA{0, 0, 0, 0} // Black.
	return
}
func (t *Text) SetText(text string)      { t.text = text }
func (t *Text) SetColor(clr color.NRGBA) { t.color = clr }
func (t *Text) SetFace(f font.Face, origin Origin) {
	t.face = f
	t.origin = origin
	bound := text.BoundString(t.face, t.text)
	t.w = float64(bound.Max.X)
	t.h = float64(-bound.Min.Y)
}

//	func (t Text) BoundString() image.Rectangle {
//		return
//	}
func (t Text) Draw(screen *ebiten.Image, b Box) {
	var x, y float64
	originX, originY := t.origin.Position()
	switch originX {
	case OriginLeft:
		x = b.padW
	case OriginCenter:
		x = b.w/2 - t.w/2
	case OriginRight:
		x = b.w - t.w - b.padW
	}
	switch originY {
	case OriginTop:
		y = b.padH + t.h
	case OriginMiddle:
		y = b.h/2 + t.h/2
	case OriginBottom:
		y = b.h - b.padH
	}
	leftTopX, leftTopY := b.LeftTopPosition()
	x = math.Floor(leftTopX + x)
	y = math.Ceil(leftTopY + y)
	text.Draw(screen, t.text, t.face, int(x), int(y), t.color)
}

// func NewRectImage(w, h, border int, outerColor, innerColor color.NRGBA) *ebiten.Image {
// 	// w = math.Ceil(w)
// 	// h = math.Ceil(h)
// 	// outer := ebiten.NewImage(int(w), int(h))
// 	// inner := ebiten.NewImage(w - 2*border)
// 	outer := ebiten.NewImage(w, h)
// 	outer.Fill(outerColor)
// 	inner := ebiten.NewImage(w-2*border, h-2*border)
// 	inner.Fill(innerColor)
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(float64(border), float64(border))
// 	outer.DrawImage(inner, op)
// 	return outer
// }

// func DesiredBoxHeight(f font.Face, padH float64) float64 {
// 	const standard = "0"
// 	rect := text.BoundString(f, standard)
// 	return padH*2 - float64(rect.Min.Y)
// }
// // type Align int
// const (
// 	AlignLeft = iota
// 	AlignCenter
// 	AlignRight
// )
