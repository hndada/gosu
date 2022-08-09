package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Todo: should I make the image field unexported?
type Sprite struct {
	I          *ebiten.Image
	W, H, X, Y float64
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

// SetCenterXY assumes Sprite's width and height are set.
func (s *Sprite) SetCenterXY(x, y float64) {
	s.X = x - s.W/2
	s.Y = y - s.H/2
}
func (s *Sprite) ApplyScale(scale float64) {
	s.W = float64(s.I.Bounds().Dx()) * scale
	s.H = float64(s.I.Bounds().Dy()) * scale
}
func (s Sprite) ScaleW() float64 { return s.W / float64(s.I.Bounds().Dx()) }
func (s Sprite) ScaleH() float64 { return s.H / float64(s.I.Bounds().Dy()) }

// DrawImageOptions is not commutative.
// Rotate -> Scale -> Translate.
func (s Sprite) Op() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.ScaleW(), s.ScaleH())
	op.GeoM.Translate(s.X, s.Y)
	return op
}
func (s *Sprite) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.I, s.Op())
}
