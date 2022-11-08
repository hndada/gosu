package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	axisX = iota
	axisY
)

type Position = Vector2

type Box struct {
	srcSize Vector2
	Scale   Vector2
	Filter  ebiten.Filter
	Position
	Origin       Origin
	AxisReversed [2]bool
}

func NewBox(srcSize Vector2) Box {
	return Box{
		srcSize: srcSize,
		Scale:   Vector2{1, 1},
		// In ebiten, FilterNearest is the default.
		Filter: ebiten.FilterLinear,
	}
}
func (b Box) SrcSize() Vector2          { return b.srcSize }
func (b Box) Size() Vector2             { return b.srcSize.Mul(b.Scale) }
func (b Box) W() float64                { return b.Size().X }
func (b Box) H() float64                { return b.Size().Y }
func (b *Box) SetSize(w, h float64)     { b.Scale = Vec2(w, h).Div(b.srcSize) }
func (b *Box) ApplyScale(scale float64) { b.Scale = b.Scale.Mul(Scalar(scale)) }
func (b *Box) SetScaleToW(w float64)    { b.Scale = Scalar(w / b.W()) }
func (b *Box) SetScaleToH(h float64)    { b.Scale = Scalar(h / b.H()) }
func (b *Box) Locate(x, y float64, origin Origin) {
	b.X = x
	b.Y = y
	b.Origin = origin
}
func (b *Box) Move(x, y float64) { b.Position = b.Position.Add(Vec2(x, y)) }
func (b Box) Min() (min Vector2) {
	size := b.Size()
	min = b.Position
	if b.AxisReversed[axisX] {
		b.Origin.X = []int{Right, Center, Left}[b.Origin.X]
	}
	if b.AxisReversed[axisY] {
		b.Origin.Y = []int{Bottom, Center, Top}[b.Origin.Y]
	}
	switch b.Origin.X {
	case Left:
		min.X -= 0
	case Center:
		min.X -= size.X / 2
	case Right:
		min.X -= size.X
	}
	switch b.Origin.Y {
	case Top:
		min.Y -= 0
	case Middle:
		min.Y -= size.Y / 2
	case Bottom:
		min.Y -= size.Y
	}
	return
}
func (b Box) Max() Vector2 { return b.Min().Add(b.Size()) }
func (b Box) In(p Vector2) bool {
	min := b.Min()
	max := b.Max()
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}

func (b Box) LeftTop(screenSize Vector2) (v Vector2) {
	v = b.Min()
	if b.AxisReversed[axisX] {
		v.X = screenSize.X - b.X + b.W()
	}
	if b.AxisReversed[axisY] {
		v.Y = screenSize.Y - b.Y + b.H()
	}
	return
}

// func (b Box) applyDrawImageOptions(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
// 	op.GeoM.Scale(b.Scale.XY())
// 	leftTop := b.LeftTop(ImageSize(screen))
// 	op.GeoM.Translate(leftTop.XY())
// 	op.Filter = b.Filter
// }
