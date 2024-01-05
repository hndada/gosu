package draws

import (
	"image"
	"io"
	"io/fs"
	"net/http"

	// Following imports are required.
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct{ *ebiten.Image }

// NewImage returns non-nil value even when fails to load file from given path.
// ebiten.NewImageFromImage will panic when input is nil.
func NewImage(r io.Reader) Image {
	src, _, err := image.Decode(r)
	if err != nil {
		return Image{}
	}
	if r, ok := r.(io.Closer); ok {
		r.Close()
	}
	return Image{ebiten.NewImageFromImage(src)}
}

func NewImageFromFile(fsys fs.FS, name string) Image {
	f, err := fsys.Open(name)
	if err != nil {
		return Image{}
	}
	return NewImage(f)
}

func NewImageFromURL(url string) Image {
	res, err := http.Get(url)
	if err != nil {
		return Image{}
	}
	return NewImage(res.Body)
}

func CreateImage(w, h float64) Image {
	return Image{ebiten.NewImage(int(w), int(h))}
}

func (i Image) IsEmpty() bool {
	return i.Image == nil
}

func (i Image) Size() Vector2 {
	if i.IsEmpty() {
		return Vector2{}
	}
	size := i.Image.Bounds().Size()
	return NewVector2FromInts(size.X, size.Y)
}

func (i Image) Sub(min, max Vector2) Image {
	if i.IsEmpty() {
		return Image{}
	}
	rect := image.Rectangle{
		Min: image.Pt(min.XYInts()),
		Max: image.Pt(max.XYInts()),
	}
	return i.SubImage(rect).(Image)
}

// Separate types are required to use Source's methods.
type Sprite struct {
	Image
	Box
}

func NewSprite(img Image) Sprite {
	return Sprite{Image: img, Box: NewBox(img)}
}
