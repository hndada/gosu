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
	raw, _, err := image.Decode(r) // image.Image
	if err != nil {
		return Image{}
	}
	if r, ok := r.(io.Closer); ok {
		r.Close()
	}
	return Image{ebiten.NewImageFromImage(raw)}
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

func (i Image) IsEmpty() bool { return i.Image == nil }

func (i Image) Size() XY {
	if i.IsEmpty() {
		return XY{}
	}
	size := i.Image.Bounds().Size()
	return NewXYFromInts(size.X, size.Y)
}

func (i Image) SubImage(x1, y1, x2, y2 int) Image {
	if i.IsEmpty() {
		return Image{}
	}
	rect := image.Rect(x1, y1, x2, y2)
	return i.Image.SubImage(rect).(Image)
}

// sub.Fill might fill the destination image permanently.
func (i Image) draw(dst Image, op *ebiten.DrawImageOptions) {
	if i.IsEmpty() {
		return
	}
	dst.DrawImage(i.Image, op)
}

//	func (img Image) In(p XY) bool {
//		max := img.Size()
//		return 0 <= p.X && p.X < max.X &&
//			0 <= p.Y && p.Y < max.Y
//	}

type Sprite struct {
	Image
	Box
}

func NewSprite(img Image) Sprite {
	return Sprite{
		Image: img,
		Box:   NewBox(img),
	}
}
