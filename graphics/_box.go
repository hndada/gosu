package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

type Box struct {
	W     int
	H     int
	Color color.Color
}

func (b *Box) Draw() *ebiten.Image {
	img, _ := ebiten.NewImage(b.W, b.H, ebiten.FilterDefault)
	ebitenutil.DrawRect(img, 0, 0, float64(b.W), float64(b.H), b.Color)
	return img
}

// box 내 text 간격 등도 있을 텐데 그런 거까지 넣어줘야 쓸모 있을 듯
type TextBox struct {
	Text *ebiten.Image
	Box  Box
}

func (tb *TextBox) Draw() *ebiten.Image {
	boxImg := tb.Box.Draw()
	boxImg.DrawImage(tb.Text, &ebiten.DrawImageOptions{})
	return boxImg
}
