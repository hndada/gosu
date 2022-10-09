package draws

import "github.com/hajimehoshi/ebiten/v2"

// Box is wrapped Subject(Sprite, Label) with Position data.
type Box struct {
	Outer Sprite3 // Subject
	Inner Subject
	Point
	Origin2
	Pad Point
	Align
}

func (b Box) OuterMin() Point {
	min := b.Point
	w, h := b.Inner.Size().XY()
	switch b.Origin2.X {
	case OriginLeft:
		min.X -= 0
	case OriginCenter:
		min.X -= w / 2
	case OriginRight:
		min.X -= w
	}
	switch b.Origin2.Y {
	case OriginTop:
		min.Y -= 0
	case OriginMiddle:
		min.Y -= h / 2
	case OriginBottom:
		min.Y -= h
	}
	return min
}
func (b Box) InnerMin() Point {
	min := b.OuterMin()
	switch b.Align.X {
	case AlignLeft:
		min.X += b.Pad.X
	case AlignCenter:
	case AlignRight:
		min.X -= b.Pad.X
	}
	switch b.Align.Y {
	case AlignTop:
		min.Y += b.Pad.Y
	case AlignMiddle:
	case AlignBottom:
		min.Y -= b.Pad.Y
	}
	return min
}
func (b Box) OuterMax() Point {
	return b.OuterMin().Add(b.Outer.Size())
}
func (b Box) InnerMax() Point {
	return b.InnerMin().Add(b.Inner.Size())
}
func (b Box) In(p Point) bool { // Input is usually cursor's position.
	p = p.Sub(b.OuterMin())
	w, h := b.Outer.Size().XY()
	return p.X >= 0 && p.X <= w && p.Y >= 0 && p.Y <= h
}

// Input point passes external translate values.
func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) { //, p Point) {
	// b.Outer.Draw(screen, op, b.OuterMin().Add(p))
	// b.Inner.Draw(screen, op, b.InnerMin().Add(p))
	b.Outer.Draw(screen, op, b.OuterMin())
	b.Inner.Draw(screen, op, b.InnerMin())
}

type Origin2 struct{ X, Y int }

const (
	OriginLeft = iota
	OriginCenter
	OriginRight
)
const (
	OriginTop = iota
	OriginMiddle
	OriginBottom
)

type Align struct{ X, Y int }

const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)
const (
	AlignTop = iota
	AlignMiddle
	AlignBottom
)
