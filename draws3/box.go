package draws

import "github.com/hajimehoshi/ebiten/v2"

// type Shape interface {
// 	Size() XY
// 	Position() XY
// 	Min() XY
// 	Max() XY
// }

// Image are not supposed to be changed.
// Text might be.
// type source interface {
// 	// IsEmpty() bool
// 	Size() XY
// }

type sizer interface {
	Size() XY
}

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
// colorm.ColorM is overkill for this package.
type Box struct {
	// for defining size and position.
	Base     *Box
	Viewport *Box
	srcSize  func() XY
	Size     Length2
	Theta    float64
	Position Length2
	Aligns   Aligns

	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	// Collapsed  bool
}

// Passing by value is convenient for Box, I guess.
func NewBox(src source) Box {
	w, h := src.Size().Values()
	return Box{
		src:  src,
		Size: NewLength2(nil, w, h, Pixel),
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

// Passing by pointer is economical for DrawImageOptions.
func (b Box) op() *ebiten.DrawImageOptions {
	geom := ebiten.GeoM{}
	scale := b.Size.Pixel().Div(b.src.Size())
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
func (b *Box) SetSize(w, h float64) { b.Size = NewLength2(nil, w, h, Pixel) }
func (b *Box) Scale(x, y float64)   { b.Size.Mul(XY{x, y}) }
func (b *Box) Move(x, y float64)    { b.Position.Add(XY{x, y}) }

// Min is the left-top position of the box.
func (b Box) Min() XY {
	min := b.Position.Pixel()
	w := b.Size.X.Pixel()
	h := b.Size.Y.Pixel()
	min.X -= []float64{0, w / 2, w}[b.Aligns.X]
	min.Y -= []float64{0, h / 2, h}[b.Aligns.Y]
	if b.Base != nil {
		min = min.Add(b.Base.Min())
		if vp := b.Base.Viewport; vp != nil {
			min = min.Sub(vp.Min())
		}
	}
	return min
}
func (b Box) Max() XY { return b.Min().Add(b.Size.Pixel()) }

func (b Box) In(p XY) bool {
	min := b.Min()
	max := b.Max()
	return min.X <= p.X && p.X < max.X &&
		min.Y <= p.Y && p.Y < max.Y
}

func (b Box) Exposed(dst Image) bool {
	return dst.In(b.Min()) || dst.In(b.Max())
}
