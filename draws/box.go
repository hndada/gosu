package draws

import "github.com/hajimehoshi/ebiten/v2"

type Subject interface {
	Size() Point
	SetSize(size Point)
	Draw(*ebiten.Image, ebiten.DrawImageOptions, Point)
}

// Box is wrapped Subject(Sprite, Label) with Position data.
type Box struct {
	Inner Subject
	Pad   Point
	Point
	Origin2 ModeXY
	Align   ModeXY
	Outer   Subject
}

// func NewBox(inner Subject, pad Point) Box {
// 	size := inner.Size()
// 	size.Add(pad.Mul(Scalar(2)))
// 	return Box{
// 		Outer: NewSprite3FromImage(
// 			ebiten.NewImage(size.XYInt()),
// 		),
// 		Inner: inner,
// 		Pad:   pad,
// 	}
// }
func OuterSize(inner Subject, pad Point) Point {
	inSize := inner.Size()
	return inSize.Add(pad.Mul(Scalar(2)))
}

func (b Box) OuterSize() Point {
	return OuterSize(b.Inner, b.Pad)
}

func (b *Box) SetSize(size Point) {
	b.Outer.SetSize(size)
}
func (b Box) OuterMin() Point {
	min := b.Point
	w, h := b.Outer.Size().XY()
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

// func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) { //, p Point) {
// 	b.Outer.Draw(screen, op, b.OuterMin())
// 	b.Inner.Draw(screen, op, b.InnerMin())
// }

// Box may be a inner subject.
// Input point passes external translate values.
func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	if b.Outer != nil {
		b.Outer.Draw(screen, op, b.OuterMin().Add(p))
	}
	if b.Inner != nil {
		b.Inner.Draw(screen, op, b.InnerMin().Add(p))
	}
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

const (
	ModeMin = iota
	ModeMid
	ModeMax
)

type ModeXY struct{ X, Y int }

var (
	AtMin = ModeXY{ModeMin, ModeMin}
	AtMid = ModeXY{ModeMid, ModeMid}
	AtMax = ModeXY{ModeMax, ModeMax}
)

// var (
// 	AlignLeftTop      = Align{AtMin, AtMin}
// 	AlignCenterMiddle = Align{AtMid, AtMid}
// )

// var (
// 	Origin2LeftTop      = Origin2{AtMin, AtMin}
// 	Origin2CenterMiddle = Origin2{AtMid, AtMid}
// )
