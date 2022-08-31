package draws

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is for storing image and translate value.
// Rotate / Scale -> Translate.
type Sprite struct {
	i          *ebiten.Image
	w, h, x, y float64
	origin     Origin
	filter     ebiten.Filter

	scaleW, scaleH float64
}
type Origin int

const (
	OriginLeftTop      Origin = iota // Default Origin.
	OriginLeftCenter                 // e.g., drawing piano notes.
	OriginLeftBottom                 // e.g., back button.
	OriginCenterTop                  // e.g., drawing field.
	OriginCenter                     // Most of sprite's Origin.
	OriginCenterBottom               // e.g., Meter.
	OriginRightTop                   // e.g., score.
	OriginRightCenter                // e.g., chart info boxes.
	OriginRightBottom                // e.g., Play button.
)

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
func (s *Sprite) SetScale(scaleW, scaleH float64, filter ebiten.Filter) {
	s.w *= scaleW
	s.h *= scaleH
	s.scaleW = scaleW
	s.scaleH = scaleH
	// i := ebiten.NewImageFromImage(s.i)
	// s.i.Clear()
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(scaleW, scaleH)
	// s.i.DrawImage(i, op)
	s.filter = filter
}
func (s *Sprite) SetPosition(x, y float64, origin Origin) {
	s.x = x
	s.y = y
	s.origin = origin
}
func (s Sprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if op == nil {
		op = &ebiten.DrawImageOptions{}
	}
	op.GeoM.Scale(s.scaleW, s.scaleH)
	x, y := s.LeftTopPosition()
	op.Filter = s.filter
	op.GeoM.Translate(x, y)
	screen.DrawImage(s.i, op)
}
func (s Sprite) LeftTopPosition() (float64, float64) {
	switch s.origin {
	case OriginLeftTop:
		// Does nothing.
	case OriginLeftCenter:
		s.y -= s.h / 2
	case OriginLeftBottom:
		s.y -= s.h
	case OriginCenterTop:
		s.x -= s.w / 2
	case OriginCenter:
		s.x -= s.w / 2
		s.y -= s.h / 2
	case OriginCenterBottom:
		s.x -= s.w / 2
		s.y -= s.h
	case OriginRightTop:
		s.x -= s.w
	case OriginRightCenter:
		s.x -= s.w
		s.y -= s.h / 2
	case OriginRightBottom:
		s.x -= s.w
		s.y -= s.h
	}
	return s.x, s.y
}
func (s Sprite) W() float64               { return s.w }
func (s Sprite) H() float64               { return s.h }
func (s Sprite) X() float64               { return s.x }
func (s Sprite) Y() float64               { return s.y }
func (s Sprite) Origin() Origin           { return s.origin }
func (s Sprite) Filter() ebiten.Filter    { return s.filter }
func (s Sprite) Size() (float64, float64) { return s.w, s.h }
func (s Sprite) SrcSize() (int, int)      { return s.i.Size() }

//	func (s Sprite) SubImage(rect image.Rectangle) *ebiten.Image {
//		return s.i.SubImage(rect).(*ebiten.Image)
//	}
//
// SubSprite supposes no x, y and origin are changed.
func (s Sprite) SubSprite(rect image.Rectangle) Sprite {
	s2 := s
	s2.i = s.i.SubImage(rect).(*ebiten.Image)
	s2.w = float64(rect.Dx())
	s2.h = float64(rect.Dy())
	return s2
}
func (s *Sprite) Move(tx, ty float64) {
	s.x += tx
	s.y += ty
}

// Should I make the image field unexported?
type Sprite0 struct {
	I          *ebiten.Image
	W, H, X, Y float64
	Filter     ebiten.Filter
}

// SetWidth sets sprite's width as well as set height scaled.
func (s *Sprite0) SetWidth(w float64) {
	ratio := w / float64(s.I.Bounds().Dx())
	s.W = w
	s.H = ratio * float64(s.I.Bounds().Dy())
}

// SetWidth sets sprite's width as well as set height scaled.
func (s *Sprite0) SetHeight(h float64) {
	ratio := h / float64(s.I.Bounds().Dy())
	s.W = ratio * float64(s.I.Bounds().Dx())
	s.H = ratio * h
}

func (s *Sprite0) ApplyScale(scale float64) {
	s.W = float64(s.I.Bounds().Dx()) * scale
	s.H = float64(s.I.Bounds().Dy()) * scale
}
func (s Sprite0) ScaleW() float64 { return s.W / float64(s.I.Bounds().Dx()) }
func (s Sprite0) ScaleH() float64 { return s.H / float64(s.I.Bounds().Dy()) }
