package gosu

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(i)
}

type Sprite struct {
	I          *ebiten.Image // Todo: Change the name to Image or i?
	W, H, X, Y float64       // Scaled W, H
}

// Op does simple calculation only using struct's field.
func (s Sprite) Op() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.ScaleW(), s.ScaleH())
	op.GeoM.Translate(s.X, s.Y)
	return op
}
func (s Sprite) ScaleW() float64 { return s.W / float64(s.I.Bounds().Dx()) }
func (s Sprite) ScaleH() float64 { return s.H / float64(s.I.Bounds().Dy()) }
