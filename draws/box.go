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
	// IndexX, IndexY   int // For sub boxes.
	// Boxes            []Box
	Grid  [][]Box
	Texts []Text
}

//	func NewBox(i *ebiten.Image, PadW, PadH float64) Box {
//		return Box{
//			Sprite: NewSpriteFromImage(i),
//			PadW:   PadW,
//			PadH:   PadH,
//		}
//	}
//
//	func (root Box) AppendBoxToRow(b Box) {
//		root.Grid = append(root.Grid, b)
//	}
func (root *Box) AppendBoxInNewRow(b Box) {
	root.Grid = append(root.Grid, []Box{b})
}
func (root *Box) AppendBoxInRow(b Box) {
	if len(root.Grid) == 0 {
		root.Grid = append(root.Grid, []Box{})
	}
	i := len(root.Grid) - 1
	root.Grid[i] = append(root.Grid[i], b)
}

//	func (root *Box) AddBox(b Box, i, j int) { //indexX, indexY int) {
//		if i > len(root.Boxes) {
//			i = len(root.Boxes)
//			root.Boxes = append(root.Boxes, )
//		}
//		if i == len(root.Boxes) {
//			j = 0
//		} else if j > len(root.Boxes[i]) {
//			j = len(root.Boxes[i])
//		}
//	}
func (root Box) Draw(screen *ebiten.Image) {
	root.Sprite.Draw(screen, nil)
	// Suppose origins of all boxes in Grid are LeftTop.
	// Suppose initial positions of all boxes in Grid are consistent with the root box.
	for _, row := range root.Grid {
		x := root.LeftTopX() + root.PadW
		y := root.LeftTopY() + root.PadH
		var maxH float64
		for _, b := range row {
			x += b.MarginW
			b.Move(x, y)
			for _, txt := range b.Texts {
				txt.Draw(screen, b)
			}
			x += b.w + b.MarginW
			if h := b.h + 2*b.MarginH; maxH < h {
				maxH = h
			}
		}
		y += maxH
	}
}

type Text struct {
	Text   string
	Color  color.Color // color.NRGBA
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
