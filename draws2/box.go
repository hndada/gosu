package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Image, Frames, Text implement Source.
type Source interface {
	Size() Vector2
	IsEmpty() bool
}

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.
type Box struct {
	Parent  *Box
	Befores []*Box
	Afters  []*Box

	Source     Source
	w, h       Length
	Theta      float64
	x, y       Length
	Aligns     Aligns
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	Visible    bool
}

func NewBox(src Source) Box {
	size := src.Size()
	return Box{
		Source: src,
		w:      Length{size.X, Pixel},
		h:      Length{size.Y, Pixel},
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

func (b Box) Root() *Box {
	if b.Parent == nil {
		return &b
	}
	return b.Parent.Root()
}

func (b Box) W() float64                 { return b.Size().X }
func (b Box) H() float64                 { return b.Size().Y }
func (b Box) X() float64                 { return b.Position().X }
func (b Box) Y() float64                 { return b.Position().X }
func (b *Box) SetW(w float64, unit Unit) { b.w = Length{w, unit} }
func (b *Box) SetH(h float64, unit Unit) { b.h = Length{h, unit} }
func (b *Box) SetX(x float64, unit Unit) { b.x = Length{x, unit} }
func (b *Box) SetY(y float64, unit Unit) { b.y = Length{y, unit} }

const (
	w = iota
	h
	x
	y
)

func (b Box) whxy(kind int) float64 {
	l := [4]Length{b.w, b.h, b.x, b.y}[kind]
	switch l.Unit {
	case Pixel:
		return l.Value
	case Percent:
		if b.Parent != nil {
			return l.Value * b.Parent.whxy(kind) / 100.0
		}
	case RootPercent:
		if root := b.Root(); root != nil {
			return l.Value * root.whxy(kind) / 100.0
		}
	}
	return l.Value
}

func (b Box) Size() Vector2                    { return Vec2(b.whxy(w), b.whxy(h)) }
func (b *Box) SetSize(w, h float64, unit Unit) { b.SetW(w, unit); b.SetH(h, unit) }
func (b *Box) Scale(scale float64)             { b.w.Value *= scale; b.h.Value *= scale }

func (b Box) Position() Vector2                    { return Vec2(b.whxy(x), b.whxy(y)) }
func (b *Box) SetPosition(x, y float64, unit Unit) { b.SetX(x, unit); b.SetY(y, unit) }
func (b *Box) Move(x, y float64)                   { b.x.Value += x; b.y.Value += y }

// Min is the left-top position of the box.
func (b Box) Min() Vector2 {
	min := b.Position()
	min.X -= []float64{0, b.W() / 2, b.W()}[b.Aligns.X]
	min.Y -= []float64{0, b.H() / 2, b.H()}[b.Aligns.Y]
	if b.Parent != nil {
		min = min.Add(b.Parent.Min())
	}
	return min
}
func (b Box) Max() Vector2 { return b.Min().Add(b.Size()) }

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
func (b Box) Draw(dst Image) {
	if !b.Visible {
		return
	}

	for _, child := range b.Befores {
		child.Draw(dst)
	}
	b.draw(dst)
	for _, child := range b.Afters {
		child.Draw(dst)
	}
}

func (b Box) draw(dst Image) {
	if b.Source.IsEmpty() {
		return
	}
	switch src := b.Source.(type) {
	case Image:
		dst.DrawImage(src.Image, b.Op())
	case Frames:
		frame := src.Images[src.Index()]
		dst.DrawImage(frame.Image, b.Op())
	case Filler:
		sub := dst.Sub(b.Min(), b.Max())
		sub.Fill(src.Color)
	case Text:
		op := text.DrawOptions{
			DrawImageOptions: *b.Op(),
			LayoutOptions: text.LayoutOptions{
				LineSpacingInPixels: src.LineSpacing,
			},
		}
		text.Draw(dst.Image, src.Text, src.face, &op)
	}
}

// colorm.ColorM is overkill for this package.
type Op = ebiten.DrawImageOptions

// Passing by pointer is economical because
// Op is big and passed several times.
func (b Box) Op() *Op {
	geom := ebiten.GeoM{}
	scale := b.Size().Div(b.Source.Size())
	geom.Scale(scale.XY())
	geom.Rotate(b.Theta)
	geom.Translate(b.Min().XY())
	return &Op{
		GeoM:       geom,
		ColorScale: b.ColorScale,
		Blend:      b.Blend,
		Filter:     b.Filter,
	}
}
