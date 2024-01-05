package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Image, Frames, Text, Color implement Source.
type Source interface {
	Size() Vector2
	IsEmpty() bool
}

// Box implements Drawable.
// l.Source.Text = "Hello, World!" // It works!
type Drawable interface {
	Draw(dst Image)
}

// Box contains information to draw an 2D entity.
// Boxes consist of a tree structure, which is a flexible way to manage entities.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.

// Z-index is not implemented because it is rather complicated:
// It might need to see all the children's z-indexes to determine drawing orders.
type Box[T Source] struct {
	Source T
	Rectangle
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
	Filter     ebiten.Filter
	Visible    bool

	Befores []Drawable
	Afters  []Drawable
}

func NewBox[T Source](src T) Box[T] {
	return Box[T]{
		Source:    src,
		Rectangle: NewRectangle(src.Size().XY()),
		// Default filter value is FilterNearest in ebitengine,
		// but FilterLinear is more natural in my opinion.
		Filter: ebiten.FilterLinear,
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
	b.draw(dst)
	for _, child := range b.Afters {
		child.Draw(dst)
	}
}

// Box.draw looks not pretty. However, it was not
// trivial to unify a Draw method among structs.
func (b Box[T]) draw(dst Image) {
	if b.Source.IsEmpty() {
		return
	}
	switch src := Source(b.Source).(type) {
	case Image:
		dst.DrawImage(src.Image, b.imageOp())
	case Frames:
		frame := src.Images[src.Index()]
		dst.DrawImage(frame.Image, b.imageOp())
	case Color:
		sub := dst.Sub(b.Min(), b.Max())
		sub.Fill(src.Color)
	case Text:
		op2 := &text.DrawOptions{
			DrawImageOptions: *b.imageOp(),
			LayoutOptions: text.LayoutOptions{
				LineSpacingInPixels: src.LineSpacing,
			},
		}
		text.Draw(dst.Image, src.Text, src.face, op2)
	}
}

// Passing by pointer is economical because
// Op is big and passed several times.

// colorm.ColorM is overkill for this package.
// Abandoned: Draw(dst Image, draw func(dst Image)):
// This requires type assertion on every child.Draw(dst, child.draw).
func (b Box[T]) imageOp() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{
		GeoM:       b.geoM(b.Source.Size()),
		ColorScale: b.ColorScale,
		Blend:      b.Blend,
		Filter:     b.Filter,
	}
}

// Separate types are required to use Source's methods.
type Sprite = Box[Image]

func NewSprite(img Image) Sprite { return NewBox[Image](img) }

type Label = Box[Text]

func NewLabel(txt Text) Label { return NewBox[Text](txt) }

type Animation = Box[Frames]

func NewAnimation(frms Frames) Animation { return NewBox[Frames](frms) }

// Filler can realize background shadow, and maybe border too.
// By introducing an image, API becomes much simpler than Web's.
// However, it is hard to adjust the size of fillers automatically
// when its parent's size changes. Nevertheless, it won't be a problem
// UI components would not change their size drastically.
type Filler = Box[Color]

func NewFiller(clr color.Color, extra float64) Filler {
	return Box[Color]{
		Source: NewColor(clr),
		Rectangle: Rectangle{
			W:      Length{extra, Extra},
			H:      Length{extra, Extra},
			Aligns: CenterMiddle,
		},
	}
}

// Extend vs Expand
// Extend: Make something larger by adding to it.
// Expand: Make something larger by stretching it
type ExtendOptions struct {
	Spacing          Length
	Direction        int
	CollapseFirstBox bool
}

// Extend works as FlexBox.
// X, Y, Aligns, Parent will be newly set.
// Y: Height + Spacing
func (b *Box[T]) Extend(children []Drawable, opts ExtendOptions) {
	// for i := range children {
	// 	child := children[i].(Box[Source])
	// }
}
