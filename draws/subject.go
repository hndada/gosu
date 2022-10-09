package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Subject interface {
	WH() Point
	Draw(screen *ebiten.Image, x, y float64, op ebiten.DrawImageOptions)
}

type Sprite3 struct {
	i     *ebiten.Image
	Scale Point
}

func (s Sprite3) Size() Point {
	return Pt(s.i.Size()).Mul(s.Scale)
}

func (s Sprite3) Draw(screen *ebiten.Image, x, y float64, op ebiten.DrawImageOptions) {
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(x, y)
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

// Label's Draw ignores given op.
func (l Label) Draw(screen *ebiten.Image, x, y float64, op ebiten.DrawImageOptions) {
	text.Draw(screen, l.Text, l.Face, int(x), int(y), l.Color)
	// i := ebiten.NewImage(l.Size().XYInt())
	// text.Draw(i, l.Text, l.Face, int(x), int(y), l.Color)
	// screen.DrawImage(i, &op)
}
