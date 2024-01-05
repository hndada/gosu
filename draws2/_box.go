package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Rectangle struct {
	Parent   *Rectangle
	Viewport Rectangle
	W, H     Length
	Theta    float64
	X, Y     Length
	Aligns   Aligns
}

func (r Rectangle) Size() Vector2 {

}

// type Shape interface {
// 	Size() Vector2
// 	Position() Vector2
// 	Min() Vector2
// 	Max() Vector2
// }

type Drawable interface {
	Size() Vector2
	Position() Vector2
	Min() Vector2
	Max() Vector2
	Draw(dst Image)
}

type Box2 struct {
	Befores []Drawable
	Afters  []Drawable

	Source     Source
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	Visible    bool
}

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.

type Boxer interface {
	root() Boxer
	whxy(int) float64

	Size() Vector2
	Position() Vector2
	Min() Vector2
	Draw(dst Image)
}

// Image, Frames, Text, Color implement Source.
type Source interface {
	Size() Vector2
	IsEmpty() bool
}

func NewBox(src Source) Box {
	return Box{
		Source: src,
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
	}
}

const (
	w = iota
	h
	x
	y
)

func (b Box) root() Boxer {
	if b.Parent == nil {
		return b
	}
	return b.Parent.root()
}

func (b Box) whxy(kind int) float64 {
	l := [4]Length{b.W, b.H, b.X, b.Y}[kind]
	switch l.Unit {
	case Pixel:
		return l.Value
	case Percent:
		if b.Parent != nil {
			return l.Value * b.Parent.whxy(kind) / 100.0
		}
	case RootPercent:
		if root := b.root(); root != nil {
			return l.Value * root.whxy(kind) / 100.0
		}
	}
	return l.Value
}

func (b Box) w() float64 { return b.whxy(w) }
func (b Box) h() float64 { return b.whxy(h) }
func (b Box) x() float64 { return b.whxy(x) }
func (b Box) y() float64 { return b.whxy(y) }

func (b Box) Size() Vector2                    { return Vec2(b.w(), b.h()) }
func (b *Box) SetSize(w, h float64, unit Unit) { b.W = Length{w, unit}; b.H = Length{h, unit} }
func (b *Box) Scale(scale float64)             { b.W.Value *= scale; b.H.Value *= scale }

func (b Box) Position() Vector2                    { return Vec2(b.x(), b.y()) }
func (b *Box) SetPosition(x, y float64, unit Unit) { b.X = Length{x, unit}; b.Y = Length{y, unit} }
func (b *Box) Move(x, y float64)                   { b.X.Value += x; b.Y.Value += y }

// Min is the left-top position of the box.
func (b Box) Min() Vector2 {
	min := b.Position()
	min.X -= []float64{0, b.w() / 2, b.w()}[b.Aligns.X]
	min.Y -= []float64{0, b.h() / 2, b.h()}[b.Aligns.Y]
	if b.Parent != nil {
		min = min.Add(b.Parent.Min())
	}
	if b.Viewport != nil {
		min = min.Sub(b.Viewport.Min())
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
		// Draw only visible on screen.
		if child.Min().X > dst.Size().X || child.Min().Y > dst.Size().Y ||
			child.Max().X < 0 || child.Max().Y < 0 {
			continue
		}
		child.Draw(dst)
	}
	b.draw(dst)
	for _, child := range b.Afters {
		child.Draw(dst)
	}
}

// Abandoned: Draw(dst Image, draw func(dst Image)):
// This requires type assertion on every child.Draw(dst, child.draw).
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
	case Color:
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
