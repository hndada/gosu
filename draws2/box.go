package draws

import "github.com/hajimehoshi/ebiten/v2"

type Box struct {
	Sprite
	Padding Vector2
	Align   Align
	Inners  [][]Box
}

func (b Box) InnerSize() Vector2 { return b.Size().Sub(b.Padding).Sub(b.Padding) }
func NewBox(inner [][]Sprite) Box {
	return Box{}
}

// Todo: recalculate inners' Position when Box's Size goes changed
func NewSimpleBox(sprite Sprite, padding Vector2, align Align) (b Box) {
	b.Sprite = sprite
	b.Source = nil
	b.X += []float64{-padding.X, 0, padding.X}[b.Origin.X]
	b.Y += []float64{-padding.Y, 0, padding.Y}[b.Origin.Y]
	b.Padding = padding
	b.Align = align // Align(sprite.Origin)

	sprite.Position = Position{}
	// Todo: calculate Sprite's Position by align
	inner := Box{Sprite: sprite}
	b.Inners = [][]Box{{inner}}
	return
}
func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	b.Sprite.Draw(screen, op)
	for i := range b.Inners {
		for j := range b.Inners[i] {
			b.Inners[i][j].Draw(screen, op)
		}
	}
}
