package draws

import (
	"image"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// NewImage returns nil when fails to load image from the path.
func NewImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(i)
}

// This is for when image.Image is needed
func NewImageSrc(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return src
}

func XFlippedImage(i *ebiten.Image) *ebiten.Image { return flipped(i, true) }
func YFlippedImage(i *ebiten.Image) *ebiten.Image { return flipped(i, false) }
func flipped(i *ebiten.Image, isX bool) *ebiten.Image {
	w, h := i.Size()
	i2 := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	if isX {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(w), 0)
	} else {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(h))
	}
	i2.DrawImage(i, op)
	return i2
}
func NewScaledImage(i *ebiten.Image, scale float64) *ebiten.Image {
	sw, sh := i.Size()
	w, h := math.Ceil(float64(sw)*scale), math.Ceil(float64(sh)*scale)
	i2 := ebiten.NewImage(int(w), int(h))
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	i2.DrawImage(i, op)
	return i2
}

// // Todo: should pass by input parameter instead of returning?
// func ApplyColor(i *ebiten.Image, clr color.Color) *ebiten.Image {
// 	w, h := i.Size()
// 	i2 := ebiten.NewImage(w, h)
// 	op := &ebiten.DrawImageOptions{}
// 	op.ColorM.ScaleWithColor(clr)
// 	i2.DrawImage(i, op)
// 	return i2
// }
