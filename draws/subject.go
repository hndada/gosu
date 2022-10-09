package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Subject interface {
	Size() Point
	Draw(*ebiten.Image, ebiten.DrawImageOptions, Point)
}

// Sprite is scaled image.
type Sprite3 struct {
	i     *ebiten.Image
	Scale Point
}

func (s Sprite3) Size() Point {
	return Pt(s.i.Size()).Mul(s.Scale)
}

// Currently option's Filter is fixed with FilterLinear.
func (s Sprite3) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(p.XY())
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(s.i, &op)
}

// func (s *Sprite3) SetScale(scale float64) {
// 	s.SetScaleWH(scale, scale)
// }
// func (s *Sprite3) SetScaleWH(scaleW, scaleH float64) {
// 	s.Scale.W *= scaleW
// 	s.Scale.H *= scaleH
// }

// Label may have expecting w and h by selecting specific Face.
type Label struct {
	Text  string
	Face  font.Face
	Color color.Color
}

func (l Label) Size() Point {
	b := text.BoundString(l.Face, l.Text)
	return Pt(b.Max.X, -b.Min.Y)
}

func (l Label) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	op.GeoM.Translate(p.XY())
	text.DrawWithOptions(screen, l.Text, l.Face, &op)
}
