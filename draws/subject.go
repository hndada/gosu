package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Subject interface {
	Size() Point
	SetSize(size Point)
	Draw(*ebiten.Image, ebiten.DrawImageOptions, Point)
}

// Sprite is a scaled image.
type Sprite3 struct {
	i     *ebiten.Image
	Scale Point
	// Scale_1 Point
}

func NewSprite3(path string) Sprite3 {
	return Sprite3{NewImage(path), Point{1, 1}}
}
func NewSprite3FromImage(i *ebiten.Image) Sprite3 {
	return Sprite3{i, Point{1, 1}}
}
func (s Sprite3) Size() Point {
	// scale := s.Scale_1.Add(Pt(1, 1))
	// return Pt(s.i.Size()).Mul(scale)
	return Pt(s.i.Size()).Mul(s.Scale)
}
func (s *Sprite3) SetSize(size Point) {
	// scale := size.Div(Pt(s.i.Size()))
	// s.Scale_1 = scale.Sub(Pt(1, 1))
	s.Scale = size.Div(Pt(s.i.Size()))
}

// Currently option's Filter is fixed with FilterLinear.
func (s Sprite3) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	// scale := s.Scale_1.Add(Pt(1, 1))
	// op.GeoM.Scale(scale.XY())
	op.GeoM.Scale(s.Scale.XY())
	op.GeoM.Translate(p.XY())
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(s.i, &op)
}

// Label may have expecting w and h by selecting specific Face.
type Label struct {
	Text  string
	Face  font.Face
	Color color.Color
	Scale Point
}

func NewLabel(text string, face font.Face, color color.Color) Label {
	return Label{
		Text:  text,
		Face:  face,
		Color: color,
		Scale: Point{1, 1},
	}
}
func (l Label) Size() Point {
	b := text.BoundString(l.Face, l.Text)
	return Pt(b.Max.X, -b.Min.Y).Mul(l.Scale)
}
func (l *Label) SetSize(size Point) {
	l.Scale = size.Div(l.Size())
}

func (l Label) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	op.GeoM.Scale(l.Scale.XY())
	op.GeoM.Translate(p.XY())
	text.DrawWithOptions(screen, l.Text, l.Face, &op)
}
