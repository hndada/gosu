package draws

import (
	"image"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// NewImage returns nil when fails to load image from the path.
func NewImage(path string) *ebiten.Image {
	return ebiten.NewImageFromImage(NewImageImage(path))
}

// NewImageImage returns image.Image.
func NewImageImage(path string) image.Image {
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

// // This is for when image.Image is needed
// func NewImageSrc(path string) image.Image {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil
// 	}
// 	defer f.Close()
// 	src, _, err := image.Decode(f)
// 	if err != nil {
// 		return nil
// 	}
// 	return src
// }

func NewXFlippedImage(i *ebiten.Image) *ebiten.Image {
	w, h := i.Size()
	i2 := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(w), 0)
	i2.DrawImage(i, op)
	return i2
}
func NewYFlippedImage(i *ebiten.Image) *ebiten.Image {
	w, h := i.Size()
	i2 := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(h))
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
