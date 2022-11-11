package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	axisX = iota
	axisY
)

type Position = Vector2

type Rectangle struct {
	srcSize Vector2
	Scale   Vector2
	Filter  ebiten.Filter
	Position
	Origin       Origin
	AxisReversed [2]bool
}

func NewRectangle(srcSize Vector2) Rectangle {
	return Rectangle{
		srcSize: srcSize,
		Scale:   Vector2{1, 1},
		// In ebiten, FilterNearest is the default.
		Filter: ebiten.FilterLinear,
	}
}
func (r Rectangle) SrcSize() Vector2          { return r.srcSize }
func (r Rectangle) Size() Vector2             { return r.srcSize.Mul(r.Scale) }
func (r Rectangle) W() float64                { return r.Size().X }
func (r Rectangle) H() float64                { return r.Size().Y }
func (r *Rectangle) SetSize(w, h float64)     { r.Scale = Vec2(w, h).Div(r.srcSize) }
func (r *Rectangle) ApplyScale(scale float64) { r.Scale = r.Scale.Mul(Scalar(scale)) }
func (r *Rectangle) SetScaleToW(w float64)    { r.Scale = Scalar(w / r.W()) }
func (r *Rectangle) SetScaleToH(h float64)    { r.Scale = Scalar(h / r.H()) }
func (r *Rectangle) Locate(x, y float64, origin Origin) {
	r.X = x
	r.Y = y
	r.Origin = origin
}
func (r *Rectangle) Move(x, y float64) { r.Position = r.Position.Add(Vec2(x, y)) }
func (r Rectangle) Min() (min Vector2) {
	size := r.Size()
	min = r.Position
	if r.AxisReversed[axisX] {
		r.Origin.X = []int{Right, Center, Left}[r.Origin.X]
	}
	if r.AxisReversed[axisY] {
		r.Origin.Y = []int{Bottom, Center, Top}[r.Origin.Y]
	}
	switch r.Origin.X {
	case Left:
		min.X -= 0
	case Center:
		min.X -= size.X / 2
	case Right:
		min.X -= size.X
	}
	switch r.Origin.Y {
	case Top:
		min.Y -= 0
	case Middle:
		min.Y -= size.Y / 2
	case Bottom:
		min.Y -= size.Y
	}
	return
}
func (r Rectangle) Max() Vector2 { return r.Min().Add(r.Size()) }
func (r Rectangle) In(p Vector2) bool {
	min := r.Min()
	max := r.Max()
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}

func (r Rectangle) LeftTop(screenSize Vector2) (v Vector2) {
	v = r.Min()
	if r.AxisReversed[axisX] {
		v.X = screenSize.X - r.X - r.W()
	}
	if r.AxisReversed[axisY] {
		v.Y = screenSize.Y - r.Y - r.H()
	}
	return
}
func (r Rectangle) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, subject Subject) {
	op.GeoM.Scale(r.Scale.XY())
	leftTop := r.LeftTop(ImageSize(screen))
	op.GeoM.Translate(leftTop.XY())
	op.Filter = r.Filter
	switch s := subject.(type) {
	case Sprite:
		screen.DrawImage(s.i, &op)
	case Label:
		text.DrawWithOptions(screen, s.text, s.face, &op)
	}
}

// func (r Rectangle) applyDrawImageOptions(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
// 	op.GeoM.Scale(r.Scale.XY())
// 	leftTop := r.LeftTop(ImageSize(screen))
// 	op.GeoM.Translate(leftTop.XY())
// 	op.Filter = r.Filter
// }
