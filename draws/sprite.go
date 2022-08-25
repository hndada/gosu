package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// type Point struct {
// 	X, Y float64
// }

// func Add(p1, p2 Point) Point     { return Point{p1.X + p2.X, p2.Y + p2.Y} }
// func Product(p1, p2 Point) Point { return Point{p1.X * p2.X, p2.Y * p2.Y} }

// Sprite is mainly for storing image and translate value.
// Filter is *not* Drawer's responsibility, since Sprite may scale.
// Todo: should Origin and OriginMode be unexported value?
type Sprite2 struct {
	i          *ebiten.Image
	w, h, x, y float64
	originMode OriginMode
	filter     ebiten.Filter
	// Size       Point
	// X, Y float64
	// Origin     struct{ X, Y float64 }
	// Offset     struct{ X, Y float64 }
	// Translater
	// OriginMode int
	// Effect func(float64)
}
type OriginMode int

const (
	OriginModeCenter OriginMode = iota
	OriginModeLeftTop
)

func NewSprite(src *ebiten.Image) Sprite2 {
	s := Sprite2{i: src}
	w, h := src.Size()
	s.w = float64(w)
	s.h = float64(h)
	// s.Size.X = float64(w) * scale.X
	// s.Size.Y = float64(h) * scale.Y
	return s
}
func (s *Sprite2) SetScale(scaleW, scaleH float64, filter ebiten.Filter) {
	s.w *= scaleW
	s.h *= scaleH
	s.filter = filter
}
func (s *Sprite2) SetPosition(x, y float64, originMode OriginMode) {
	s.x = x
	s.y = y
	s.originMode = originMode
}
func (s Sprite2) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if op == nil {
		op = &ebiten.DrawImageOptions{}
	}
	switch s.originMode {
	case OriginModeCenter:
		s.x -= s.w / 2
		s.y -= s.h / 2
	case OriginModeLeftTop:
		// Does nothing.
	}
	op.Filter = s.filter
	op.GeoM.Translate(s.x, s.y)
	screen.DrawImage(s.i, op)
}

// func NewSpriteWithoutScale(src *ebiten.Image) Sprite2 {
// 	return NewSprite(src, 1, 1)
// }

// func (s Sprite2) W() float64 { return s.Size.X }
// func (s Sprite2) H() float64 { return s.Size.Y }
func (s Sprite2) W() float64               { return s.w }
func (s Sprite2) H() float64               { return s.h }
func (s Sprite2) X() float64               { return s.x }
func (s Sprite2) Y() float64               { return s.y }
func (s Sprite2) OriginMode() OriginMode   { return s.originMode }
func (s Sprite2) Filter() ebiten.Filter    { return s.filter }
func (s Sprite2) Size() (float64, float64) { return s.w, s.h }
func (s Sprite2) SrcSize() (int, int)      { return s.i.Size() }

// Pos returns X and Y value of Sprite's Left top point.
// func (s Sprite2) Pos() Point { return Point{s.X(), s.Y()} }

// Should I make the image field unexported?
type Sprite struct {
	I          *ebiten.Image
	W, H, X, Y float64
	Filter     ebiten.Filter
}

// SetWidth sets sprite's width as well as set height scaled.
func (s *Sprite) SetWidth(w float64) {
	ratio := w / float64(s.I.Bounds().Dx())
	s.W = w
	s.H = ratio * float64(s.I.Bounds().Dy())
}

// SetWidth sets sprite's width as well as set height scaled.
func (s *Sprite) SetHeight(h float64) {
	ratio := h / float64(s.I.Bounds().Dy())
	s.W = ratio * float64(s.I.Bounds().Dx())
	s.H = ratio * h
}

// SetCenterX and SetCenterY suppose Sprite's width and height are set.
func (s *Sprite) SetCenterX(x float64) { s.X = x - s.W/2 }
func (s *Sprite) SetCenterY(y float64) { s.Y = y - s.H/2 }
func (s *Sprite) ApplyScale(scale float64) {
	s.W = float64(s.I.Bounds().Dx()) * scale
	s.H = float64(s.I.Bounds().Dy()) * scale
}
func (s Sprite) ScaleW() float64 { return s.W / float64(s.I.Bounds().Dx()) }
func (s Sprite) ScaleH() float64 { return s.H / float64(s.I.Bounds().Dy()) }

// DrawImageOptions is not commutative.
// Rotate / Scale -> Translate.
func (s Sprite) Op() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.ScaleW(), s.ScaleH())
	op.GeoM.Translate(s.X, s.Y)
	op.Filter = s.Filter
	return op
}
func (s *Sprite) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.I, s.Op())
}
