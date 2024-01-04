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

// All of these are for defining size and position.
type Box interface {
	Root() Box
	Size() Length2
	Position() Length2
	Draw(dst Image)
}

type Length2 struct{ X, Y Length }

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.
type BaseBox struct {
	Parent  Box
	Befores []Box
	Afters  []Box

	Draw func(dst Image)
	// Source     Source
	Size       Length2
	Theta      float64
	Position   Length2
	Aligns     Aligns
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	Visible    bool
}

func NewBaseBox() BaseBox {
	return BaseBox{
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

func (b BaseBox) Root() Box {
	if b.Parent == nil {
		return &b
	}
	return b.Parent.Root()
}

func (b BaseBox) W() float64 { return b.Size().X }
func (b BaseBox) H() float64 { return b.Size().Y }
func (b BaseBox) X() float64 { return b.Position().X }
func (b BaseBox) Y() float64 { return b.Position().X }

const (
	w = iota
	h
	x
	y
)

func (b BaseBox) whxy(kind int) float64 {
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

func (b BaseBox) Size() Vector2                    { return Vec2(b.whxy(w), b.whxy(h)) }
func (b *BaseBox) SetSize(w, h float64, unit Unit) { b.SetW(w, unit); b.SetH(h, unit) }
func (b *BaseBox) Scale(scale float64)             { b.w.Value *= scale; b.h.Value *= scale }

func (b BaseBox) Position() Vector2                    { return Vec2(b.whxy(x), b.whxy(y)) }
func (b *BaseBox) SetPosition(x, y float64, unit Unit) { b.SetX(x, unit); b.SetY(y, unit) }
func (b *BaseBox) Move(x, y float64)                   { b.x.Value += x; b.y.Value += y }

// Min is the left-top position of the box.
func (b BaseBox) Min() Vector2 {
	min := b.Position()
	min.X -= []float64{0, b.W() / 2, b.W()}[b.Aligns.X]
	min.Y -= []float64{0, b.H() / 2, b.H()}[b.Aligns.Y]
	if b.Parent != nil {
		min = min.Add(b.Parent.Min())
	}
	return min
}
func (b BaseBox) Max() Vector2 { return b.Min().Add(b.Size()) }

// Only four methods are required: Scale, Rotate, Translate, and ScaleWithColor.
// DrawImageOptions is not commutative: Do Translate at the final stage.
func (b BaseBox) Draw(dst Image) {
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

func (b BaseBox) draw(dst Image) {
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

func (b BaseBox) SizeInPixel() Vector2 {
	return b.Size().Mul(b.Source.Size())
}

// Passing by pointer is economical because
// Op is big and passed several times.
func (b BaseBox) Op(src Source) *Op {
	geom := ebiten.GeoM{}
	scale := b.SizeInPixel().Div(src.Size())
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

type Sprite struct {
	Image
	BaseBox
}

func (s Sprite) Draw(dst Image) {
	if !s.Visible {
		return
	}
	src := s.Image
	dst.DrawImage(src.Image, s.Op(src))
}
