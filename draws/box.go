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
	PadW, PadH       float64
	MarginW, MarginH float64
	IndexX, IndexY   int // For sub boxes.
	Boxes            []Box
	Texts            []Text
}

//	func NewBox(i *ebiten.Image, PadW, PadH float64) Box {
//		return Box{
//			Sprite: NewSpriteFromImage(i),
//			PadW:   PadW,
//			PadH:   PadH,
//		}
//	}
func (root Box) AddBox(b Box, indexX, indexY int) {
	// sort.Slice()
}
func (root Box) Draw(screen *ebiten.Image) {
	root.Sprite.Draw(screen, nil)
	// Suppose all sub boxes have same origin and initial position with the root box.
	for _, b := range root.Boxes {
		// Todo: move regarding its margins and indexes
		b.Move(0, 0)
		for _, txt := range b.Texts {
			txt.Draw(screen, b)
		}
	}
}

type Text struct {
	Text   string
	Color  color.NRGBA
	face   font.Face
	origin Origin // Text's align is determined by its Origin.
	w, h   float64
	// x, y   float64
	// Align  int
}

//	func NewText(text string) (t Text) {
//		t.text = text
//		t.color = color.NRGBA{0, 0, 0, 0} // Black.
//		return
//	}
//
// func (t *Text) SetText(text string)      { t.text = text }
// func (t *Text) SetColor(clr color.NRGBA) { t.color = clr }
func (t *Text) SetFace(f font.Face, origin Origin) {
	t.face = f
	t.origin = origin
	bound := text.BoundString(t.face, t.Text)
	t.w = float64(bound.Max.X)
	t.h = float64(-bound.Min.Y)
}

//	func (t Text) BoundString() image.Rectangle {
//		return
//	}
func (t Text) Draw(screen *ebiten.Image, b Box) {
	var x, y float64
	switch t.origin.PositionX() {
	case OriginLeft:
		x = b.PadW
	case OriginCenter:
		x = b.w/2 - t.w/2
	case OriginRight:
		x = b.w - t.w - b.PadW
	}
	switch t.origin.PositionY() {
	case OriginTop:
		y = b.PadH + t.h
	case OriginMiddle:
		y = b.h/2 + t.h/2
	case OriginBottom:
		y = b.h - b.PadH
	}
	x = math.Floor(b.LeftTopX() + x)
	y = math.Ceil(b.LeftTopY() + y)
	text.Draw(screen, t.Text, t.face, int(x), int(y), t.Color)
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

// func DesiredBoxHeight(f font.Face, PadH float64) float64 {
// 	const standard = "0"
// 	rect := text.BoundString(f, standard)
// 	return PadH*2 - float64(rect.Min.Y)
// }
// // type Align int
// const (
// 	AlignLeft = iota
// 	AlignCenter
// 	AlignRight
// )
