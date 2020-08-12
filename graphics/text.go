package graphics

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// dpi와 hinting은 font에 포함돼있음
// todo: 질문하기, Text를 이미지로 관리하려는 시도
func DrawText(t string, font font.Face, clr color.Color) *ebiten.Image {
	b := text.BoundString(font, t)
	img, _ := ebiten.NewImage(b.Dx(), b.Dy(), ebiten.FilterDefault)
	text.Draw(img, t, font, 0, 0, clr)
	return img
}

const (
	boxPadding = 8
)

func DrawTextBox(text *ebiten.Image, clr color.Color) *ebiten.Image {
	tx, ty := text.Size()
	img, _ := ebiten.NewImage(tx+2*boxPadding, ty+boxPadding, ebiten.FilterDefault)
	img.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(boxPadding), float64(boxPadding/2))
	var img1, img2 ebiten.Image
	img1 = *img
	img.DrawImage(text, op)
	img2 = *img
	fmt.Println(img1 == img2) // true ?

	// temp, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	// temp.Fill(color.Black)
	// img.DrawImage(temp, op)
	// img3 = *img
	// fmt.Println(img2 == img3) // 얘도 true가 나옴
	return img
}
