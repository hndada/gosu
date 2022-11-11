package draws

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func ImageSize(i *ebiten.Image) Vector2 { return IntVec2(i.Size()) }

type Image struct{ *ebiten.Image }

func (i Image) Size() Vector2 {
	if !i.IsValid() {
		return Vector2{}
	}
	return ImageSize(i.Image)
}
func (i Image) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	screen.DrawImage(i.Image, &op)
}
func (i Image) IsValid() bool { return i.Image != nil }

// NewImage returns nil when fails to load image from the path.
func NewImage(path string) Image {
	// ebiten.NewImageFromImage will panic when input is nil.
	if i := NewImageImage(path); i != nil {
		return Image{ebiten.NewImageFromImage(i)}
	}
	return Image{}
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

func NewImages(path string) (is []Image) {
	const ext = ".png"
	one := []Image{NewImage(path + ext)}
	dir, err := os.Open(path)
	if err != nil {
		return one
	}
	defer dir.Close()
	fs, err := dir.ReadDir(-1)
	if err != nil {
		return one
	}

	nums := make([]int, 0, len(fs))
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		num := strings.TrimSuffix(f.Name(), ext)
		if num, err := strconv.Atoi(num); err == nil {
			nums = append(nums, num)
		}
	}
	sort.Ints(nums)
	for _, num := range nums {
		path := filepath.Join(path, fmt.Sprintf("%d.png", num))
		is = append(is, NewImage(path))
	}
	return
}

func NewImageXFlipped(i Image) Image {
	w, h := i.Image.Size()
	i2 := ebiten.NewImage(w, h)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(w), 0)
	i2.DrawImage(i.Image, op)
	return Image{i2}
}
func NewImageYFlipped(i Image) Image {
	w, h := i.Image.Size()
	i2 := ebiten.NewImage(w, h)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(h))
	i2.DrawImage(i.Image, op)
	return Image{i2}
}
func NewImageScaled(i Image, scale float64) Image {
	size := i.Size().Mul(Scalar(scale))
	i2 := ebiten.NewImage(size.XYInt())

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	i2.DrawImage(i.Image, op)
	return Image{i2}
}
