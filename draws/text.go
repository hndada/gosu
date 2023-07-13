package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Text struct {
	text string
	face font.Face
}

func NewText(t string, face font.Face) Text {
	return Text{
		text: t,
		face: face,
	}
}

// Append new line when each function has more than one line
// and functions are not strictly related.
func (t Text) Size() Vector2 {
	b := text.BoundString(t.face, t.text)
	return IntVec2(b.Max.X, -b.Min.Y)
}

// issue: ebiten/v2/text.DrawWithOptions does not support ColorM.
func (t Text) Draw(dst Image, op Op) {
	op2 := &ebiten.DrawImageOptions{
		GeoM:   op.GeoM,
		Blend:  op.Blend,
		Filter: op.Filter,
	}
	text.DrawWithOptions(dst.Image, t.text, t.face, op2)
}

func (t Text) IsEmpty() bool { return len(t.text) == 0 }
