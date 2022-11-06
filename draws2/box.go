package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	Scale Point
	Point
	Origin Origin
	Filter ebiten.Filter
}

func NewBox() Box {
	return Box{
		Scale:  Point{1, 1},
		Filter: ebiten.FilterLinear, // In ebiten, default filter is FilterNearest.
	}
}
func (b *Box) SetScale(scale Point)     { b.Scale = scale }
func (b *Box) ApplyScale(scale float64) { b.Scale = b.Scale.Mul(Scalar(scale)) }
func (b *Box) SetPoint(x, y float64, origin Origin) {
	b.X = x
	b.Y = y
	b.Origin = origin
}
func (b *Box) Move(x, y float64)               { b.Point = b.Point.Add(Pt(x, y)) }
func (b Box) XY(size Point) (float64, float64) { return b.Min(size).XY() }
func (b Box) Min(size Point) (min Point) {
	min = b.Point
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
func (b Box) Max(size Point) Point { return b.Min(size).Add(size) }
func (b Box) In(size, p Point) bool {
	min := b.Min(size)
	max := b.Max(size)
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}
