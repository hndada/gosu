package draws

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	axisX = iota
	axisY
)

// Unit of Location is percent.
type Location = Vector2

// Unit of Position is pixel.
type Position = Vector2

// type Filter = ebiten.Filter
// const (
// 	Linear  = iota
// 	Nearest
// )

// Sprite is an image or a text drawn in a screen based on its position and scale.
// DrawImageOptions is not commutative. Do Translate at the final stage.
type Sprite struct {
	Source
	Scale  Vector2
	Filter ebiten.Filter
	Position
	Origin       Origin
	AxisReversed [2]bool
	// Outer        *Sprite
	// Inners       []Sprite
}

func NewSprite(fsys fs.FS, name string) Sprite {
	return NewSpriteFromSource(LoadImage(fsys, name))
}
func NewSpriteFromSource(src Source) Sprite {
	return Sprite{
		Source: src,
		Scale:  Vector2{1, 1},
		Filter: ebiten.FilterLinear, // FilterNearest is the default in ebiten.
	}
}

// func (s *Sprite) Append(src Source, loc Location) {
// 	outer := s
// 	inner := NewSpriteFromSource(src)
// 	inner.Outer = outer
// 	if ratio := loc.X; ratio <= 1 {
// 		inner.X += (outer.W() - inner.W()) * ratio
// 	}
// 	if ratio := loc.Y; ratio <= 1 {
// 		inner.Y += (outer.H() - inner.H()) * ratio
// 	}
// 	s.Inners = append(s.Inners, inner)
// }

// func (s *Sprite) SetRelativePosition(outer Sprite, location Vector2) {}
func (s Sprite) SrcSize() Vector2          { return s.Source.Size() }
func (s Sprite) Size() Vector2             { return s.SrcSize().Mul(s.Scale) }
func (s Sprite) W() float64                { return s.Size().X }
func (s Sprite) H() float64                { return s.Size().Y }
func (s *Sprite) SetSize(w, h float64)     { s.Scale = Vec2(w, h).Div(s.SrcSize()) }
func (s *Sprite) ApplyScale(scale float64) { s.Scale = s.Scale.Mul(Scalar(scale)) }
func (s *Sprite) SetScaleToW(w float64)    { s.Scale = Scalar(w / s.W()) }
func (s *Sprite) SetScaleToH(h float64)    { s.Scale = Scalar(h / s.H()) }
func (s *Sprite) Locate(x, y float64, origin Origin) {
	s.X = x
	s.Y = y
	s.Origin = origin
}
func (s *Sprite) Move(x, y float64) { s.Position = s.Position.Add(Vec2(x, y)) }
func (s Sprite) Min() (min Vector2) {
	size := s.Size()
	min = s.Position
	if s.AxisReversed[axisX] {
		s.Origin.X = []int{Right, Center, Left}[s.Origin.X]
	}
	if s.AxisReversed[axisY] {
		s.Origin.Y = []int{Bottom, Center, Top}[s.Origin.Y]
	}
	min.X -= []float64{0, size.X / 2, size.X}[s.Origin.X]
	min.Y -= []float64{0, size.Y / 2, size.Y}[s.Origin.Y]
	return
}
func (s Sprite) Max() Vector2 { return s.Min().Add(s.Size()) }
func (s Sprite) In(p Vector2) bool {
	min := s.Min()
	max := s.Max()
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}
func (s Sprite) Draw(dst Image, op Op) {
	if !s.IsValid() {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	// if s.Outer != nil {
	// 	s.Add(s.Outer.Position)
	// }
	leftTop := s.LeftTop(dst.Size())
	op.GeoM.Translate(leftTop.XY())
	op.Filter = s.Filter
	s.Source.Draw(dst, op)
	// for _, inner := range s.Inners {
	// 	inner.Draw(dst, op)
	// }
}
func (s Sprite) LeftTop(screenSize Vector2) (v Vector2) {
	v = s.Min()
	if s.AxisReversed[axisX] {
		v.X = screenSize.X - s.X - s.W()
	}
	if s.AxisReversed[axisY] {
		v.Y = screenSize.Y - s.Y - s.H()
	}
	return
}

//	func (s Sprite) Op(screen *ebiten.Image, op Op) Op {
//		op.GeoM.Scale(s.Scale.XY())
//		leftTop := s.LeftTop(ImageSize(screen))
//		op.GeoM.Translate(leftTop.XY())
//		op.Filter = s.Filter
//		return op
//	}
// func NewSprite(path string) Sprite { return NewSprite(NewImage(path)) }
//
//	func NewSpriteFromImage(i *ebiten.Image) Sprite {
//		return Sprite{Source: Image{i}}
//	}

// Todo: let Sprite skip calculating Inners' Position when Outer's Size is fixed
