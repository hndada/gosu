package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is a image drawn in screen based on its position and scale.
// DrawImageOptions is not commutative. Do translate at final stage.
type Sprite struct {
	i              *ebiten.Image
	w, h, x, y     float64
	origin         Origin
	filter         ebiten.Filter
	scaleW, scaleH float64
}

func NewSprite(path string) Sprite {
	return NewSpriteFromImage(NewImage(path))
}
func NewSpriteFromImage(src *ebiten.Image) Sprite {
	s := Sprite{i: src}
	if src == nil {
		return s
	}
	w, h := src.Size()
	s.w = float64(w)
	s.h = float64(h)
	s.scaleW = 1
	s.scaleH = 1
	return s
}

// func (s Sprite) Image() *ebiten.Image     { return s.i }
func (s Sprite) W() float64               { return s.w }
func (s Sprite) H() float64               { return s.h }
func (s Sprite) X() float64               { return s.x }
func (s Sprite) Y() float64               { return s.y }
func (s Sprite) Origin() Origin           { return s.origin }
func (s Sprite) Filter() ebiten.Filter    { return s.filter }
func (s Sprite) Size() (float64, float64) { return s.w, s.h }
func (s Sprite) SrcSize() (int, int)      { return s.i.Size() }
func (s Sprite) IsValid() bool            { return s.i != nil }
func (s Sprite) In(x, y float64) bool {
	x -= s.LeftTopX()
	y -= s.LeftTopY()
	return x >= 0 && x <= s.W() && y >= 0 && y <= s.H()
}
func (s Sprite) LeftTopX() float64 {
	switch s.origin.PositionX() {
	case OriginLeft:
		s.x -= 0
	case OriginCenter:
		s.x -= s.w / 2
	case OriginRight:
		s.x -= s.w
	}
	return s.x
}
func (s Sprite) LeftTopY() float64 {
	switch s.origin.PositionY() {
	case OriginTop:
		s.y -= 0
	case OriginMiddle:
		s.y -= s.h / 2
	case OriginBottom:
		s.y -= s.h
	}
	return s.y
}
func (s *Sprite) SetScale(scale float64) {
	s.SetScaleXY(scale, scale, ebiten.FilterLinear)
}
func (s *Sprite) SetScaleXY(scaleW, scaleH float64, filter ebiten.Filter) {
	s.w *= scaleW // / s.scaleW
	s.h *= scaleH // / s.scaleH
	s.scaleW *= scaleW
	s.scaleH *= scaleH
	s.filter = filter
}
func (s *Sprite) SetPosition(x, y float64, origin Origin) {
	s.x = x
	s.y = y
	s.origin = origin
}
func (s *Sprite) Move(tx, ty float64) {
	s.x += tx
	s.y += ty
}
func (s Sprite) Draw(screen *ebiten.Image, opSrc *ebiten.DrawImageOptions) {
	var op ebiten.DrawImageOptions
	if opSrc != nil {
		op = *opSrc
	}
	op.GeoM.Scale(s.scaleW, s.scaleH)
	op.Filter = s.filter
	op.GeoM.Translate(s.LeftTopX(), s.LeftTopY())
	screen.DrawImage(s.i, &op)
}
