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
	Box
}

func NewLabel(text string, face font.Face) *Label {
	return &Label{
		Text:  text,
		Face:  face,
		Color: color.Black,
		Box:   NewBox(),
	}
}
func (l Label) SrcSize() Point {
	b := text.BoundString(l.Face, l.Text)
	return IntPt(b.Max.X, -b.Min.Y)
}
func (l Label) Size() Point {
	return l.SrcSize().Mul(l.Scale)
}
func (l *Label) SetSize(size Point) {
	l.Scale = size.Div(l.SrcSize())
}
func (l Label) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	op.GeoM.Scale(l.Scale.XY())
	op.GeoM.Translate(l.Point.XY())
	op.ColorM.ScaleWithColor(l.Color)
	op.Filter = l.Filter
	text.DrawWithOptions(screen, l.Text, l.Face, &op)
}
