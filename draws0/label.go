package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Label struct {
	text string
	face font.Face
	Box
}

// if either size's X or Y is 0, Label use BoundString for setting size.
func NewLabel(txt string, face font.Face, size Vector2) (l Label) {
	if size.X == 0 || size.Y == 0 {
		b := text.BoundString(face, txt)
		size = IntVec2(b.Max.X, -b.Min.Y)
	}
	return Label{
		text: txt,
		face: face,
		Box:  NewBox(size),
	}
}

func (l Label) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	op.GeoM.Scale(l.Scale.XY())
	leftTop := l.LeftTop(ImageSize(screen))
	op.GeoM.Translate(leftTop.XY())
	op.Filter = l.Filter
	text.DrawWithOptions(screen, l.text, l.face, &op)
}
