package draws

import "github.com/hajimehoshi/ebiten/v2"

type Box struct {
	Point
	Origin Origin
	Scale  Point
	Filter ebiten.Filter
}

func NewBox() Box {
	return Box{
		Scale:  Point{1, 1},
		Filter: ebiten.FilterLinear, // In ebiten, default filter is FilterNearest.
	}
}

// func (b Box) Size(src Point) Point     { return src.Mul(b.Scale) }
// func (b *Box) SetSize(src, size Point) { b.Scale = size.Div(src) }
func (b Box) In(size, p Point) bool {
	min := b.Min(size)
	max := b.Max(size)
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}
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
func (b *Box) Move(p Point)        { b.Point.Add(p) }
