package draws

type TextBox struct {
	txt Text
	Box
}

func NewTextBox(txt Text) TextBox {
	return TextBox{
		txt: txt,
		Box: NewBox(txt),
	}
}

func (t *TextBox) SetText(txt string) {
	t.txt.Text = txt
	t.Box.Size = t.txt.SourceSize()
}

func (t TextBox) Draw(dst Image) {
	t.Box.Draw(dst, t.txt)
}
