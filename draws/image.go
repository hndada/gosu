package draws

import (
	"image"
	"io/fs"
	"net/http"

	// Following imports are required.
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct{ *ebiten.Image }

func NewImage(w, h float64) Image {
	return Image{ebiten.NewImage(int(w), int(h))}
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

// NewImageFromFile returns nil when fails to load file from given path.
// ebiten.NewImageFromImage will panic when input is nil.
func NewImageFromFile(fsys fs.FS, name string) Image {
	if src := NewImageImageFromFile(fsys, name); src != nil {
		return Image{ebiten.NewImageFromImage(src)}
	}
	return Image{}
}

func NewImageFromURL(url string) (i Image, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	src, _, err := image.Decode(res.Body)
	if err != nil {
		return
	}
	i = Image{ebiten.NewImageFromImage(src)}
	return
}

func (i Image) IsEmpty() bool {
	return i.Image == nil
}

func (i Image) SourceSize() Vector2 {
	if i.IsEmpty() {
		return Vector2{}
	}
	size := i.Image.Bounds().Size()
	return NewVector2FromInts(size.X, size.Y)
}

// Passing by pointer is economical because
// Op is big and passed several times.
func (i Image) Draw(dst Image, op *Op) {
	dst.DrawImage(i.Image, op)
}

// func NewImageXFlipped(src Image) Image {
// 	size := src.Size()
// 	dst := Image{ebiten.NewImage(size.XYInts())}
// 	op := Op{}
// 	op.GeoM.Scale(-1, 1)
// 	op.GeoM.Translate(size.X, 0)
// 	src.Draw(dst, op)
// 	return dst
// }

// func NewImageYFlipped(src Image) Image {
// 	size := src.Size()
// 	dst := Image{ebiten.NewImage(size.XYInts())}
// 	op := Op{}
// 	op.GeoM.Scale(1, -1)
// 	op.GeoM.Translate(0, size.Y)
// 	src.Draw(dst, op)
// 	return dst
// }

// func NewImageColored(src Image, color color.Color) Image {
// 	size := src.Size()
// 	dst := Image{ebiten.NewImage(size.XYInts())}
// 	op := Op{}
// 	op.ColorM.ScaleWithColor(color)
// 	src.Draw(dst, op)
// 	return dst
// }
