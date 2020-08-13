package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// dpi와 hinting은 font에 포함돼있음
func DrawText(s string, font font.Face, clr color.Color) *ebiten.Image {
	b := text.BoundString(font, s)
	img, _ := ebiten.NewImage(b.Dx()+3, b.Dy(), ebiten.FilterDefault) // needs a bit more pixels at x-axis
	text.Draw(img, s, font, 0, b.Dy(), clr)
	return img
}

const (
	boxPadding = 8
)

func DrawTextBox(t *ebiten.Image, clr color.Color) *ebiten.Image {
	tx, ty := t.Size()
	tbox, _ := ebiten.NewImage(tx+2*boxPadding, ty+boxPadding, ebiten.FilterDefault)
	tbox.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(boxPadding), float64(boxPadding/2))
	tbox.DrawImage(t, op)
	return tbox
}
