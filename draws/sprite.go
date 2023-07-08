package draws

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// Unit of Position is pixel.
// cf. Location: unit is percent.
type Position = Vector2

// Sprite is an image or a text drawn in a screen based on its position and scale.
// DrawImageOptions is not commutative. Do Translate at the final stage.
type Sprite struct {
	Source
	Scale  Vector2
	Filter ebiten.Filter
	Position
	Anchor Anchor
}

func NewSprite(src Source) Sprite {
	return Sprite{
		Source: src,
		Scale:  Vector2{1, 1},
		Filter: ebiten.FilterLinear, // FilterNearest is the default in ebiten.
	}
}
func LoadSprite(fsys fs.FS, name string) Sprite {
	return NewSprite(LoadImage(fsys, name))
}
func LoadSpriteFromURL(url string) (Sprite, error) {
	img, err := LoadImageFromURL(url)
	if err != nil {
		return Sprite{}, err
	}
	return NewSprite(img), nil
}

func (s Sprite) SrcSize() Vector2             { return s.Source.Size() }
func (s Sprite) Size() Vector2                { return s.SrcSize().Mul(s.Scale) }
func (s Sprite) W() float64                   { return s.Size().X }
func (s Sprite) H() float64                   { return s.Size().Y }
func (s *Sprite) SetSize(w, h float64)        { s.Scale = Vec2(w, h).Div(s.SrcSize()) }
func (s *Sprite) SetScaleToW(w float64)       { s.Scale = Scalar(w / s.W()) }
func (s *Sprite) SetScaleToH(h float64)       { s.Scale = Scalar(h / s.H()) }
func (s *Sprite) MultiplyScale(scale float64) { s.Scale = s.Scale.Mul(Scalar(scale)) }
func (s *Sprite) Locate(x, y float64, anchor Anchor) {
	s.X = x
	s.Y = y
	s.Anchor = anchor
}
func (s *Sprite) Move(x, y float64) { s.Position = s.Position.Add(Vec2(x, y)) }
func (s Sprite) Min() (min Vector2) {
	size := s.Size()
	min = s.Position
	min.X -= []float64{0, size.X / 2, size.X}[s.Anchor.X]
	min.Y -= []float64{0, size.Y / 2, size.Y}[s.Anchor.Y]
	return
}
func (s Sprite) Max() Vector2 { return s.Min().Add(s.Size()) }
func (s Sprite) In(p Vector2) bool {
	min := s.Min()
	max := s.Max()
	p = p.Sub(min)
	return p.X >= 0 && p.X <= max.X && p.Y >= 0 && p.Y <= max.Y
}
func (s Sprite) LeftTop(screenSize Vector2) (v Vector2) { return s.Min() }
func (s Sprite) Draw(dst Image, op Op) {
	if s.Source == nil || !s.IsValid() {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	leftTop := s.LeftTop(dst.Size())
	op.GeoM.Translate(leftTop.XY())
	op.Filter = s.Filter
	s.Source.Draw(dst, op)
}
