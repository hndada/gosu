package draws

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	i *ebiten.Image
	Box
}

func NewSprite(path string) *Sprite {
	return &Sprite{
		i:   NewImage(path),
		Box: NewBox(),
	}
}
func NewSpriteFromImage(i *ebiten.Image) *Sprite {
	return &Sprite{
		i:   i,
		Box: NewBox(),
	}
}
func (s Sprite) srcSize() Point { return IntPt(s.i.Size()) }
func (s Sprite) Size() Point {
	if s.i == nil {
		return Point{}
	}
	return s.srcSize().Mul(s.Scale)
	// return s.Box.Size(IntPt(s.i.Size()))
}
func (s *Sprite) SetSize(size Point) {
	if s.i == nil {
		return
	}
	s.Scale = size.Div(s.srcSize())
	// s.Box.SetSize(IntPt(s.i.Size()), size)
}
func (s Sprite) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if s.i == nil {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(s.Point.XY())
	op.Filter = s.Filter
	screen.DrawImage(s.i, &op)
}
