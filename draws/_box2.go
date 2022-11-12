package draws

import "github.com/hajimehoshi/ebiten/v2"

// < 1 value: Percent, else Pixel.

type Box struct {
	Outer *Box
	Sprite
	Location
	Inners []Box
}

// Outer, Location, Sprite
// func NewBox(outer *Box, s Sprite) (b Box) {
// 	b.Outer = outer
// 	b.Sprite = s
// 	return
// }

// Inner, Location, Sprite
func (b *Box) Append(s Sprite, loc Location) {
	// s.Position = s.Min()
	s.Origin.X = Left
	if s.AxisReversed[axisX] {
		s.Origin.X = Right
	}
	s.Origin.Y = Top
	if s.AxisReversed[axisY] {
		s.Origin.Y = Bottom
	}
	inner := Box{
		Outer:    b,
		Sprite:   s,
		Location: loc,
	}
	b.Inners = append(b.Inners, inner)
}

func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	p := b.Outer.Min()
	if ratio := b.Location.X; ratio <= 1 {
		p.X += (b.Outer.W() - b.W()) * ratio
	}
	if ratio := b.Location.Y; ratio <= 1 {
		p.Y += (b.Outer.H() - b.H()) * ratio
	}
	b.Position = p
	b.Sprite.Draw(screen, op)
	for _, inner := range b.Inners {
		inner.Draw(screen, op)
	}
}

// Padding Vector2 // Exploited by excluded function.

// Supporting value is Anchor, Percent, Pixel.
// Anchor: aka Align; e.g. CenterMiddle. Use as String(Align).
// Percent: relative position in a percent. Use as 20%.
// Pixel: relative position.
// type Location string

// Todo: currently suppose AxisReversed are all consistent among Sprites
