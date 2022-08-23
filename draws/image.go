package draws

import (
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// NewImage returns nil when fails to load image from the path.
func NewImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
		// panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil
		// panic(err)
	}
	return ebiten.NewImageFromImage(i)
}

func FlipX(i *ebiten.Image) *ebiten.Image { return flip(i, true) }
func FlipY(i *ebiten.Image) *ebiten.Image { return flip(i, false) }
func flip(i *ebiten.Image, isX bool) *ebiten.Image {
	w, h := i.Size()
	i2 := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	if isX {
		op.GeoM.Scale(-1, 1)
	} else {
		op.GeoM.Scale(1, -1)
	}
	i2.DrawImage(i, op)
	return i2
}

func ApplyColor(i *ebiten.Image, clr color.Color) *ebiten.Image {
	w, h := i.Size()
	i2 := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.ScaleWithColor(clr)
	i2.DrawImage(i, op)
	return i2
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
		// panic(err)
	}
	return src
}
