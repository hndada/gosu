package draws

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// NewImage returns nil when fails to load image from the path.
func NewImage(path string) *ebiten.Image {
	if i := NewImageImage(path); i != nil {
		return ebiten.NewImageFromImage(i)
	}
	return nil
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
