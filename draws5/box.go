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
	Viewport *Box // Default is screen.
	Size     Length2
	Theta    float64
	Origin   *Box // Default is screen.
	Position Length2
	Aligns   Aligns

	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	// Collapsed  bool
}

// Returning by value is convenient for Box, I guess.
// Passed arguments should be pointers.
func NewBox(src source) Box {
	size := NewLength2(src.Size().Values())
	size.SetBase(&Screen.Size, Pixel)
	// Set relative position so that the box is centered regardless of the screen size change.
	pos := NewLength2(0.5, 0.5)
	pos.SetBase(&Screen.Size, Percent)
	// pos := NewLength2(Screen.Size.X.Value/2, Screen.Size.Y.Value/2)
	return Box{
		src:      src,
		Viewport: &Screen,
		Size:     size,
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}
func (b *Box) SetOrigin(org *Box) {
	b.Origin = org
	b.Position.SetBase(&org.Position, Pixel)
}

func (b *Box) SetViewport(vp *Box) {
	b.Viewport = vp
	b.Size.SetBase(&vp.Size, Pixel)
}

func (b *Box) Scale(x float64)      { b.Size.Mul(XY{x, x}) }
func (b *Box) ScaleXY(x, y float64) { b.Size.Mul(XY{x, y}) }
func (b *Box) Move(x, y float64)    { b.Position.Add(XY{x, y}) }

// Min is the left-top position of the box.
func (b Box) Min() XY {
	min := b.Position.Pixels()
	size := b.Size.Pixels()
	w := size.X
	h := size.Y
	min.X -= []float64{0, w / 2, w}[b.Aligns.X]
	min.Y -= []float64{0, h / 2, h}[b.Aligns.Y]
	if b.Viewport != nil {
		min = min.Sub(b.Viewport.Min())
	}
	return min
}
func (b Box) Max() XY { return b.Min().Add(b.Size.Pixels()) }

func (b Box) In(p XY) bool {
	min := b.Min()
	max := b.Max()
	return min.X <= p.X && p.X < max.X &&
		min.Y <= p.Y && p.Y < max.Y
}

// func (b Box) Exposed(dst Image) bool {
// 	return dst.In(b.Min()) || dst.In(b.Max())
// }

// func (b Box) Draw(dst Image) {
// 	if b.src.IsEmpty() {
// 		return
// 	}
// 	b.src.draw(dst, b.op())
// }

// Passing by pointer is economical for DrawImageOptions.
func (b Box) op() *ebiten.DrawImageOptions {
	geom := ebiten.GeoM{}
	scale := b.Size.Pixels().Div(b.src.Size())
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
