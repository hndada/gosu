package draws

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	Text
	Options
}

func NewLabel(src Text) Label {
	return Label{src, NewOptions(src)}
}

func (l Label) Draw(dst Image) {
	if l.IsEmpty() {
		return
	}
	src := l.Text.Text
	op := &text.DrawOptions{
		DrawImageOptions: *l.imageOp(),
		LayoutOptions: text.LayoutOptions{
			LineSpacingInPixels: l.LineSpacing,
		},
	}
	text.Draw(dst.Image, src, l.face, op)
}
