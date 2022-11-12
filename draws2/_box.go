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

// Inners' relative position is pre-calculated and fixed unless Box's Size has fixed.
func (b *Box) SetSize(w, h float64) {
	outer := b
	b.Sprite.SetSize(w, h)
	for i := range b.Inners {
		var offset Position
		// if i>=1{
		// 	offset.Y=b.Inners[i-1][0].Max
		// }
		for j, inner := range b.Inners[i] {
			p := b.Padding
			p.X += []float64{0, outer.W()/2 - inner.W()/2, outer.W() - inner.W()}[b.Align.X]
			p.Y += []float64{0, outer.H()/2 - inner.H()/2, outer.H() - inner.H()}[b.Align.Y]
			p.X += []float64{0, inner.W() / 2, inner.W()}[inner.Origin.X]
			p.Y += []float64{0, inner.H() / 2, inner.H()}[inner.Origin.Y]
			b.Inners[i][j].Position = p
		}
	}
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
