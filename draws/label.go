package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Label may have expecting w and h by selecting specific Face.
type Label struct {
	Text  string
	Face  font.Face
	Color color.Color
	Scale Point
}

func NewLabel(text string, face font.Face, color color.Color) *Label {
	return &Label{
		Text:  text,
		Face:  face,
		Color: color,
		Scale: Point{1, 1},
	}
}
func (l Label) Size() Point {
	b := text.BoundString(l.Face, l.Text)
	return IntPt(b.Max.X, -b.Min.Y).Mul(l.Scale)
}
func (l *Label) SetSize(size Point) {
	l.Scale = size.Div(l.Size())
}
func (l Label) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	op.GeoM.Scale(l.Scale.XY())
	op.GeoM.Translate(p.XY())
	text.DrawWithOptions(screen, l.Text, l.Face, &op)
}
