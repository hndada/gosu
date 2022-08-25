package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is for storing image and translate value.
// Rotate / Scale -> Translate.
type Sprite struct {
	i          *ebiten.Image
	w, h, x, y float64
	originMode OriginMode
	filter     ebiten.Filter
}
type OriginMode int

const (
	OriginModeLeftTop      OriginMode = iota // Default OriginMode.
	OriginModeLeftCenter                     // e.g., drawing piano notes.
	OriginModeLeftBottom                     // e.g., back button.
	OriginModeCenterTop                      // e.g., drawing field.
	OriginModeCenter                         // Most of sprite's OriginMode.
	OriginModeCenterBottom                   // e.g., TimingMeter.
	OriginModeRightTop                       // e.g., score.
	OriginModeRightCenter                    // e.g., chart info boxes.
	OriginModeRightBottom                    // e.g., Play button.
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
	return s
}
func (s *Sprite) SetScale(scaleW, scaleH float64, filter ebiten.Filter) {
	s.w *= scaleW
	s.h *= scaleH
	s.filter = filter
}
func (s *Sprite) SetPosition(x, y float64, originMode OriginMode) {
	s.x = x
	s.y = y
	s.originMode = originMode
}
func (s Sprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if op == nil {
		op = &ebiten.DrawImageOptions{}
	}
	switch s.originMode {
	case OriginModeLeftTop:
		// Does nothing.
	case OriginModeLeftCenter:
		s.y -= s.h / 2
	case OriginModeLeftBottom:
		s.y -= s.h
	case OriginModeCenterTop:
		s.x -= s.w / 2
	case OriginModeCenter:
		s.x -= s.w / 2
		s.y -= s.h / 2
	case OriginModeCenterBottom:
		s.x -= s.w / 2
		s.y -= s.h
	case OriginModeRightTop:
		s.x -= s.w
	case OriginModeRightCenter:
		s.x -= s.w
		s.y -= s.h / 2
	case OriginModeRightBottom:
		s.x -= s.w
		s.y -= s.h
	}
	op.Filter = s.filter
	op.GeoM.Translate(s.x, s.y)
	screen.DrawImage(s.i, op)
}

func (s Sprite) W() float64               { return s.w }
func (s Sprite) H() float64               { return s.h }
func (s Sprite) X() float64               { return s.x }
func (s Sprite) Y() float64               { return s.y }
func (s Sprite) OriginMode() OriginMode   { return s.originMode }
func (s Sprite) Filter() ebiten.Filter    { return s.filter }
func (s Sprite) Size() (float64, float64) { return s.w, s.h }
func (s Sprite) SrcSize() (int, int)      { return s.i.Size() }

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
