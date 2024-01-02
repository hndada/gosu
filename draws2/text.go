package draws

import "github.com/hajimehoshi/ebiten/v2/text/v2"

type Text struct {
	Text        string
	Face        text.Face
	LineSpacing float64
}

func NewText(txt string, Face text.Face) Text {
	return Text{
		Text:        txt,
		Face:        Face,
		LineSpacing: 1.6,
	}
}

func (t Text) IsEmpty() bool { return len(t.Text) == 0 }

func (t Text) Size() Vector2 {
	return Vec2(text.Measure(t.Text, t.Face, t.LineSpacing))
}

func (t Text) Draw(dst Image, op *Op) {
	text.Draw(dst.Image, t.Text, t.Face, &text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacingInPixels: t.LineSpacing,
		},
	})
}
