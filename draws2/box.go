package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Image, Frames, Text, Color implement Source.
type source interface {
	Size() Vector2
	IsEmpty() bool
}

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.
type Box struct {
	befores []Drawable
	afters  []Drawable

	// source is not embedded for avoiding ambiguous method calls.
	source source
	Rectangle
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	zIndex     int
	Visible    bool
}

func NewBox(src source) Box {
	return Box{
		Rectangle: NewRectangle(src.Size().XY()),
		source:    src,
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

func (b Box) ZIndex() int { return b.zIndex }
func (b *Box) SetZIndex(z int) {
	// Interpolate box into appropriate position.
	// for i, child := range b.befores {
	// 	if child.ZIndex() > z {
	// 		b.befores = append(b.befores[:i], append([]Drawable{b}, b.befores[i:]...)...)
	// 		return
	// 	}
	// }
	// for i, child := range b.afters {
	// 	if child.ZIndex() > z {
	// 		b.afters = append(b.afters[:i], append([]Drawable{b}, b.afters[i:]...)...)
	// 		return
	// 	}
	// }
	b.zIndex = z
}

type ExpandOptions struct {
	Spacing   Length
	Direction int
	Collapse  bool
}

func (b *Box) Expand(children []Drawable, opts ExpandOptions) {

}

// Passing by pointer is economical because
// Op is big and passed several times.
func (b Box) imageOp() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{
		GeoM:       b.geoM(b.source.Size()),
		ColorScale: b.ColorScale,
		Blend:      b.Blend,
		Filter:     b.Filter,
	}
}

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
func (b Box) Draw(dst Image) {
	for _, child := range b.befores {
		child.Draw(dst)
	}
	b.draw(dst)
	for _, child := range b.afters {
		child.Draw(dst)
	}
}

// colorm.ColorM is overkill for this package.
// Abandoned: Draw(dst Image, draw func(dst Image)):
// This requires type assertion on every child.Draw(dst, child.draw).
func (b Box) draw(dst Image) {
	if b.source.IsEmpty() || !b.Visible || !b.Exposed(dst) {
		return
	}
	switch src := b.source.(type) {
	case Image:
		dst.DrawImage(src.Image, b.imageOp())
	case Frames:
		frame := src.Images[src.Index()]
		dst.DrawImage(frame.Image, b.imageOp())
	case Color:
		sub := dst.Sub(b.Min(), b.Max())
		sub.Fill(src.Color)
	case Text:
		op := src.op(*b.imageOp())
		text.Draw(dst.Image, src.Text, src.face, op)
	}
}
