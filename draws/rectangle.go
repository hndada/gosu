package draws

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var blank = color.NRGBA{255, 255, 255, 255}

// Todo: nine-patch
type Rectangle struct {
	Size_ Point
	Color color.Color
	Outer *Rectangle
}

func NewRectangle(size Point) *Rectangle {
	return &Rectangle{
		Size_: size,
		Color: blank,
	}
}

//	func NewRectangle(size Point, color color.Color) *Rectangle {
//		return &Rectangle{
//			Size_: size,
//			Color: color,
//		}
//	}
func (r Rectangle) Size() Point         { return r.Size_ }
func (r *Rectangle) SetSize(size Point) { r.Size_ = size }
func (r Rectangle) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	r.draw(screen, p) // No DrawImageOptions are supported for Rectangle.
}
func (r Rectangle) draw(screen *ebiten.Image, p Point) {
	if r.Outer != nil {
		r.Outer.draw(screen, p)
		offset := r.Outer.Size_.Sub(r.Size_)
		offset = offset.Div(Scalar(2))
		p.Add(offset)
	}
	min := image.Pt(p.XYInt())
	max := image.Pt(p.Add(r.Size_).XYInt())
	rect := image.Rectangle{min, max}
	sub := screen.SubImage(rect).(*ebiten.Image)
	sub.Fill(r.Color)
}
