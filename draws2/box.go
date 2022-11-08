package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Position = Vector2
type Box struct {
	// srcSize Vector2
	Scale  Vector2
	Filter ebiten.Filter

	Position
	Origin            Origin
	DirectionReversed [2]bool
}

func NewBox() Box {
	return Box{
		Scale:  Vector2{1, 1},
		Filter: ebiten.FilterLinear, // In ebiten, default filter is FilterNearest.
	}
}

// func (b *Box) SetScale(scale Vector2)   { b.Scale = scale }
func (b *Box) ApplyScale(scale float64) { b.Scale = b.Scale.Mul(Scalar(scale)) }

// func (b Box) Pos(size Vector2) (float64, float64) { return b.XY() }
func (b *Box) SetPosition(x, y float64, origin Origin) {
	b.X = x
	b.Y = y
	b.Origin = origin
}
func (b *Box) Move(x, y float64) { b.Position = b.Position.Add(Vec2(x, y)) }
func (b Box) Min(size Vector2) (min Vector2) {
	min = b.Position
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
func (b Box) Max(size Vector2) Vector2 { return b.Min(size).Add(size) }
func (b Box) In(size, p Vector2) bool {
	min := b.Min(size)
	max := b.Max(size)
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}
