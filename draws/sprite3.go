package draws

import "github.com/hajimehoshi/ebiten/v2"

// Sprite is a scaled image.
// Sprite3's i is nil if Sprite3 is a placeholder or failed to load an image.
type Sprite3 struct {
	i     *ebiten.Image
	Scale Point
}

func NewSprite3(path string) *Sprite3 {
	return &Sprite3{NewImage(path), Point{1, 1}}
}
func NewSprite3FromImage(i *ebiten.Image) *Sprite3 {
	return &Sprite3{i, Point{1, 1}}
}

// func NewBlankSprite3(size Point) Sprite3 {
// 	return Sprite3{
// 		ebiten.NewImage(size.XYInt()),
// 		Point{1, 1},
// 	}
// }
func (s Sprite3) Size() Point {
	if s.i == nil {
		return Point{}
	}
	return IntPt(s.i.Size()).Mul(s.Scale)
}
func (s *Sprite3) SetSize(size Point) {
	if s.i == nil {
		return
	}
	s.Scale = size.Div(IntPt(s.i.Size()))
}

// Currently option's Filter is fixed with FilterLinear.
func (s Sprite3) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	if s.i == nil {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(p.XY())
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(s.i, &op)
}
