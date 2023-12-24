package draws

import "github.com/hajimehoshi/ebiten/v2/text/v2"

type Text struct {
	Text string
	Face text.Face
}

func NewText(txt string, Face text.Face) Text {
	return Text{
		Text: txt,
		Face: Face,
	}
}

func (t Text) IsEmpty() bool { return len(t.Text) == 0 }

// Append new line when each function has more than one line
// and functions are not strictly related.
func (t Text) SourceSize() Vector2 {
	return Vec2(text.Measure(t.Text, t.Face, 1))
	// return NewVector2FromInts(b.Max.X, -b.Min.Y)
}

// issue: ebiten/v2/text.DrawWithOptions does not support ColorM.
func (t Text) Draw(dst Image, op *Op) {
	op2 := &text.DrawOptions{
		DrawImageOptions: *op,
	}
	text.Draw(dst.Image, t.Text, t.Face, op2)
}
