package ebitenui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// todo: it fails to render string when it is multiple-lines
// todo: numbers and normal words are not coherently rendered
func RenderText(s string, font font.Face, clr color.Color) *ebiten.Image {
	b := text.BoundString(font, s)
	img, _ := ebiten.NewImage(b.Dx()+15, b.Dy()+9, ebiten.FilterDefault)
	text.Draw(img, s, font, 0, b.Dy()-3, clr)
	return img
}

const (
	boxPadding = 8
)

func RenderTextBox(t *ebiten.Image, clr color.Color) *ebiten.Image {
	tx, ty := t.Size()
	tbox, _ := ebiten.NewImage(tx+2*boxPadding, ty+boxPadding, ebiten.FilterDefault)
	tbox.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(boxPadding), float64(boxPadding/2))
	tbox.DrawImage(t, op)
	return tbox
}
