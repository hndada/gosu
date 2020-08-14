package ebitenui

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

func RenderBox(size image.Point, clr color.Color) *ebiten.Image {
	img, _ := ebiten.NewImage(size.X, size.Y, ebiten.FilterDefault)
	img.Fill(clr)
	return img
}
