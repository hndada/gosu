package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Image, Frames, Text, Color implement Source.
type Source interface {
	Size() Vector2
	IsEmpty() bool
	Draw(dst Image, op *ebiten.DrawImageOptions)
}

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.

// Z-index is not implemented because it is rather complicated.
type Box[T Source] struct {
	Befores []Drawable
	Afters  []Drawable

	// source is not embedded for avoiding ambiguous method calls.
	Source T
	Rectangle
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	Visible    bool
}

func NewBox[T Source](src Source) Box[T] {
	return Box[T]{
		Rectangle: NewRectangle(src.Size().XY()),
		Source:    src.(T),
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

// Extend vs Expand
// Extend: Make something larger by adding to it.
// Expand: Make something larger by stretching it
type ExtendOptions struct {
	Spacing       Length
	Direction     int
	CollapseFirst bool
}

// X, Y, Aligns, Parent will be newly set.
func (b *Box[T]) Extend(children []Box[Source], opts ExtendOptions) {

}

func a() {
	b := Box[Image]{}
	b2 := Box[Text]{}
	b.Extend([]Box[Source]{b2}, ExtendOptions{})
}

// Passing by pointer is economical because
// Op is big and passed several times.
func (b Box[T]) imageOp() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{
		GeoM:       b.geoM(b.Source.Size()),
		ColorScale: b.ColorScale,
		Blend:      b.Blend,
		Filter:     b.Filter,
	}
}

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
func (b Box[T]) Draw(dst Image) {
	if !b.Visible || !b.Exposed(dst) {
		return
	}

	for _, child := range b.Befores {
		child.Draw(dst)
	}
	b.Source.Draw(dst, b.imageOp())
	for _, child := range b.Afters {
		child.Draw(dst)
	}
}

// colorm.ColorM is overkill for this package.
// Abandoned: Draw(dst Image, draw func(dst Image)):
// This requires type assertion on every child.Draw(dst, child.draw).

// Objective: Box.draw looks not pretty.
func (b Box[T]) draw(dst Image) {
	switch src := b.Source.(type) {
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
