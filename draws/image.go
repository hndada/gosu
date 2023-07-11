package draws

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	// Following imports are required.
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct{ *ebiten.Image }

func (i Image) Size() Vector2 {
	if i.IsEmpty() {
		return Vector2{}
	}
	return IntVec2(i.Image.Size())
}

func (i Image) Draw(dst Image, op Op) {
	dst.Image.DrawImage(i.Image, &op)
}

func (i Image) IsEmpty() bool { return i.Image == nil }

func NewImage(w, h float64) Image {
	return Image{ebiten.NewImage(int(w), int(h))}
}

// NewImageFromFile returns nil when fails to load file from given path.
// ebiten.NewImageFromImage will panic when input is nil.
func NewImageFromFile(fsys fs.FS, name string) Image {
	if i := NewImageImageFromFile(fsys, name); i != nil {
		return Image{ebiten.NewImageFromImage(i)}
	}
	return Image{}
}

func NewImageFromURL(url string) (Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return Image{}, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return Image{}, err
	}
	eimg := ebiten.NewImageFromImage(img)
	return Image{eimg}, err
}

func NewImageImageFromFile(fsys fs.FS, name string) image.Image {
	f, err := fsys.Open(name)
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

// NewImagesFromFile is for Animation.
func NewImagesFromFile(fsys fs.FS, name string) (is []Image) {
	const ext = ".png"

	// name supposed to have no extension when passed in NewImagesFromFile.
	name = strings.TrimSuffix(name, filepath.Ext(name))

	one := []Image{NewImageFromFile(fsys, name+ext)}
	fs, err := fs.ReadDir(fsys, name)
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
		// Avoid use filepath here; it yields backslash, which is invalid path for FS.
		name2 := path.Join(name, fmt.Sprintf("%d.png", num))
		is = append(is, NewImageFromFile(fsys, name2))
	}
	return
}

func NewImageXFlipped(src Image) Image {
	size := src.Size()
	dst := Image{ebiten.NewImage(size.XYInt())}
	op := Op{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(size.X, 0)
	src.Draw(dst, op)
	return dst
}

func NewImageYFlipped(src Image) Image {
	size := src.Size()
	dst := Image{ebiten.NewImage(size.XYInt())}
	op := Op{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, size.Y)
	src.Draw(dst, op)
	return dst
}

func NewImageColored(src Image, color color.Color) Image {
	size := src.Size()
	dst := Image{ebiten.NewImage(size.XYInt())}
	op := Op{}
	op.ColorM.ScaleWithColor(color)
	src.Draw(dst, op)
	return dst
}
