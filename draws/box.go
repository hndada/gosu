package draws

import "github.com/hajimehoshi/ebiten/v2"

type source interface {
	// IsEmpty() bool
	Size() XY
	// draw(dst Image, op *ebiten.DrawImageOptions)
}

// Image are not supposed to be changed. Text might be.
// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
// colorm.ColorM is overkill for this package.
type Box struct {
	src      source
	Size     XY
	Theta    float64
	Position XY
	Aligns   Aligns

	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
}

// Returning by value is convenient for Box, I guess.
// Passed arguments should be pointers.
func NewBox(src source) Box {
	return Box{
		src:  src,
		Size: src.Size(),
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}
func (b Box) W() float64 { return b.Size.X }
func (b Box) H() float64 { return b.Size.Y }
func (b Box) X() float64 { return b.Position.X }
func (b Box) Y() float64 { return b.Position.Y }

func (b *Box) SetSize(w, h float64) { b.Size = XY{w, h} }
func (b *Box) Scale(s float64)      { b.Size = b.Size.Mul(XY{s, s}) }

func (b *Box) Locate(x, y float64, aligns Aligns) {
	b.Position = XY{x, y}
	b.Aligns = aligns
}
func (b *Box) Move(x, y float64) { b.Position = b.Position.Add(XY{x, y}) }

// Min is the left-top position of the box.
func (b Box) Min() XY {
	min := b.Position
	size := b.Size
	w := size.X
	h := size.Y
	min.X -= []float64{0, w / 2, w}[b.Aligns.X]
	min.Y -= []float64{0, h / 2, h}[b.Aligns.Y]
	return min
}
func (b Box) Max() XY { return b.Min().Add(b.Size) }

func (b Box) In(p XY) bool {
	min := b.Min()
	max := b.Max()
	return min.X <= p.X && p.X < max.X &&
		min.Y <= p.Y && p.Y < max.Y
}

// Passing by pointer is economical for DrawImageOptions.
func (b Box) op() *ebiten.DrawImageOptions {
	geom := ebiten.GeoM{}
	scale := b.Size.Div(b.src.Size())
	geom.Scale(scale.Values())
	geom.Rotate(b.Theta)
	geom.Translate(b.Min().Values())

	return &ebiten.DrawImageOptions{
		GeoM:       geom,
		ColorScale: b.ColorScale,
		Blend:      b.Blend,
		Filter:     b.Filter,
	}
}
