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

// Sprite, Label, Animation, Filler implement Drawable.
// type Drawable interface {}

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.
type Box struct {
	Befores []Box
	Afters  []Box

	// source is not embedded for avoiding ambiguous method calls.
	source source
	Rectangle
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
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

// colorm.ColorM is overkill for this package.

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
	if !b.Visible {
		return
	}
	dr := NewRectangle(dst.Size().XY())

	for _, child := range b.Befores {
		if child.Intersect(dr) {
			child.Draw(dst)
		}
	}

	b.draw(dst)

	for _, child := range b.Afters {
		if child.Intersect(dr) {
			child.Draw(dst)
		}
	}
}

// Abandoned: Draw(dst Image, draw func(dst Image)):
// This requires type assertion on every child.Draw(dst, child.draw).
func (b Box) draw(dst Image) {
	if b.source.IsEmpty() {
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
