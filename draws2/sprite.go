package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is an image drawn in a screen based on its position and scale.
// DrawImageOptions is not commutative. Do Translate at the final stage.
type Sprite struct {
	i *ebiten.Image
	Box
}

func NewSprite(path string) Sprite {
	return Sprite{
		i:   NewImage(path),
		Box: NewBox(),
	}
}
func NewSpriteFromImage(i *ebiten.Image) Sprite {
	return Sprite{
		i:   i,
		Box: NewBox(),
	}
}
func (s Sprite) IsValid() bool    { return s.i != nil }
func (s Sprite) SrcSize() Vector2 { return IntVec2(s.i.Size()) }
func (s Sprite) Size() Vector2 {
	if s.i == nil {
		return Vector2{}
	}
	return s.SrcSize().Mul(s.Scale)
}
func (s *Sprite) SetSize(w, h float64) {
	if s.i == nil {
		return
	}
	s.Scale = Vec2(w, h).Div(s.SrcSize())
}
func (s *Sprite) SetScaleToW(w float64) { s.Scale = Scalar(w / s.W()) }
func (s *Sprite) SetScaleToH(h float64) { s.Scale = Scalar(h / s.H()) }
func (s Sprite) W() float64             { return s.Size().X }
func (s Sprite) H() float64             { return s.Size().Y }

func (s Sprite) Min() Vector2      { return s.Box.Min(s.Size()) }
func (s Sprite) Max() Vector2      { return s.Box.Max(s.Size()) }
func (s Sprite) In(p Vector2) bool { return s.Box.In(s.Size(), p) }
func (s Sprite) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if s.i == nil {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(s.XY())
	op.Filter = s.Filter
	screen.DrawImage(s.i, &op)
}
