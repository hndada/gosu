package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is for storing image and translate value.
// DrawImageOptions is not commutative. Do translate at final stage.
// Now Sprite has ColorM.
type Sprite struct {
	i              *ebiten.Image
	w, h, x, y     float64
	origin         Origin
	filter         ebiten.Filter
	scaleW, scaleH float64
	// color          color.Color
	// colorM         ebiten.ColorM
}
type Origin int

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
const (
	OriginLeftTop      Origin = iota // Default Origin.
	OriginLeftMiddle                 // e.g., Notes in Piano mode.
	OriginLeftBottom                 // e.g., back button.
	OriginCenterTop                  // e.g., drawing field.
	OriginCenterMiddle               // Most of sprite's Origin.
	OriginCenterBottom               // e.g., Meter.
	OriginRightTop                   // e.g., score.
	OriginRightMiddle                // e.g., chart info boxes.
	OriginRightBottom                // e.g., Play button.
)

// Position returns x-axis and y-axis position each.
//
//	func (origin Origin) Position() (int, int) {
//		return origin.PositionX(), origin.PositionY()
//	}
func (origin Origin) PositionX() int { return int(origin) / 3 }
func (origin Origin) PositionY() int { return int(origin) % 3 }
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

// func (s *Sprite) SetColor(clr color.Color) {
// 	// s.colorM.Reset()
// 	// s.colorM.ScaleWithColor(clr)
// 	s.color = clr
// }

// Todo: ColorM affects Translate.
func (s Sprite) Draw(screen *ebiten.Image, opSrc *ebiten.DrawImageOptions) {
	var op ebiten.DrawImageOptions
	if opSrc != nil {
		op = *opSrc
	}
	// op.ColorM = s.colorM
	op.GeoM.Scale(s.scaleW, s.scaleH)
	op.Filter = s.filter
	// x, y := s.LeftTopPosition()
	op.GeoM.Translate(s.LeftTopX(), s.LeftTopY())
	screen.DrawImage(s.i, &op)
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

//	func (s Sprite) LeftTopPosition() (float64, float64) {
//		switch s.origin {
//		case OriginLeftTop:
//			// Does nothing.
//		case OriginLeftMiddle:
//			s.y -= s.h / 2
//		case OriginLeftBottom:
//			s.y -= s.h
//		case OriginCenterTop:
//			s.x -= s.w / 2
//		case OriginCenterMiddle:
//			s.x -= s.w / 2
//			s.y -= s.h / 2
//		case OriginCenterBottom:
//			s.x -= s.w / 2
//			s.y -= s.h
//		case OriginRightTop:
//			s.x -= s.w
//		case OriginRightMiddle:
//			s.x -= s.w
//			s.y -= s.h / 2
//		case OriginRightBottom:
//			s.x -= s.w
//			s.y -= s.h
//		}
//		return s.x, s.y
//	}
func (s Sprite) W() float64               { return s.w }
func (s Sprite) H() float64               { return s.h }
func (s Sprite) X() float64               { return s.x }
func (s Sprite) Y() float64               { return s.y }
func (s Sprite) Origin() Origin           { return s.origin }
func (s Sprite) Filter() ebiten.Filter    { return s.filter }
func (s Sprite) Size() (float64, float64) { return s.w, s.h }
func (s Sprite) SrcSize() (int, int)      { return s.i.Size() }
func (s Sprite) IsValid() bool            { return s.i != nil }
func (s *Sprite) Move(tx, ty float64) {
	s.x += tx
	s.y += ty
}

//	func (s Sprite) SubSprite(propMinX, propMinY, propMaxX, propMaxY float64) Sprite {
//		w, h := s.SrcSize()
//		minX := math.Floor(propMinX * float64(w))
//		minY := math.Floor(propMinY * float64(h))
//		maxX := math.Ceil(propMaxX * float64(w))
//		maxY := math.Ceil(propMaxY * float64(h))
//		rect := image.Rect(int(minX), int(minY), int(maxX), int(maxY))
//		s2 := s
//		s2.i = s.i.SubImage(rect).(*ebiten.Image)
//		s2.w = float64(rect.Dx()) * s2.scaleW
//		s2.h = float64(rect.Dy()) * s2.scaleH
//		return s2
//	}
