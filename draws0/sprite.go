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
func (s Sprite) SrcSize() Point { return IntPt(s.i.Size()) }
func (s Sprite) Size() Point {
	if s.i == nil {
		return Point{}
	}
	return s.SrcSize().Mul(s.Scale)
}
func (s *Sprite) SetSize(w, h float64) {
	if s.i == nil {
		return
	}
	s.Scale = Pt(w, h).Div(s.SrcSize())
}
func (s *Sprite) SetScaleToW(w float64) { s.SetScale(Scalar(w / s.W())) }
func (s *Sprite) SetScaleToH(h float64) { s.SetScale(Scalar(h / s.H())) }
func (s Sprite) W() float64             { return s.Size().X }
func (s Sprite) H() float64             { return s.Size().Y }
func (s Sprite) XY() (float64, float64) { return s.Box.XY(s.Size()) } // s.Min(s.Size()).XY() }
func (s Sprite) Min() Point             { return s.Box.Min(s.Size()) }
func (s Sprite) Max() Point             { return s.Box.Max(s.Size()) }
func (s Sprite) IsValid() bool          { return s.i != nil }
func (s Sprite) In(p Point) bool        { return s.Box.In(s.Size(), p) }
func (s Sprite) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if s.i == nil {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(s.XY())
	op.Filter = s.Filter
	screen.DrawImage(s.i, &op)
}
