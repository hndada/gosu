package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Box contains information to draw an 2D entity.
type Box struct {
	Size       Vector2 // pixel
	Theta      float64 // radian
	Position   Vector2 // pixel
	Anchor     Anchor
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
}

type Op = ebiten.DrawImageOptions

// Image, Frames, Text implement Source.
// Image -> Sprite
// Frames -> Animation
// Text -> TextBox
type Source interface {
	SourceSize() Vector2
	IsEmpty() bool
	Draw(dst Image, op *Op)
}

func NewBox(src Source) Box {
	return Box{
		Size: src.SourceSize(),
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

func (b Box) W() float64 { return b.Size.X }
func (b Box) H() float64 { return b.Size.Y }
func (b Box) X() float64 { return b.Position.X }
func (b Box) Y() float64 { return b.Position.Y }

func (b *Box) SetSize(w, h float64) {
	b.Size.X = w
	b.Size.Y = h
}
func (b *Box) Locate(x, y float64, anchor Anchor) {
	b.Position.X = x
	b.Position.Y = y
	b.Anchor = anchor
}
func (b *Box) Move(x, y float64) {
	b.Position = b.Position.Add(Vec2(x, y))
}

// Min is the left-top position of the box.
func (b Box) Min() Vector2 {
	min := b.Position
	min.X -= []float64{0, b.W() / 2, b.W()}[b.Anchor.X]
	min.Y -= []float64{0, b.H() / 2, b.H()}[b.Anchor.Y]
	return min
}
func (b Box) Max() Vector2 { return b.Min().Add(b.Size) }

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
func (b Box) Draw(dst Image, src Source) {
	if src == nil || src.IsEmpty() {
		return
	}
	op := ebiten.DrawImageOptions{}
	scale := b.Size.Div(src.SourceSize())
	op.GeoM.Scale(scale.XY())
	op.GeoM.Rotate(b.Theta)
	op.GeoM.Translate(b.Min().XY())
	op.ColorScale = b.ColorScale
	op.Blend = b.Blend
	op.Filter = b.Filter
	src.Draw(dst, &op)
}
