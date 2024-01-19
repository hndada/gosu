package draws

import "github.com/hajimehoshi/ebiten/v2"

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
// colorm.ColorM is overkill for this package.
type Options struct {
	Rectangle
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	// Collapsed  bool
}

// Passing by value is convenient for Options, I guess.
func NewOptions(src sizer) Options {
	return Options{
		Rectangle: NewRectangle(src),
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

// Passing by pointer is economical for DrawImageOptions.
func (op Options) imageOp() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{
		GeoM:       op.geoM(),
		ColorScale: op.ColorScale,
		Blend:      op.Blend,
		Filter:     op.Filter,
	}
}
