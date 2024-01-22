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

// _image and _text implement sizer.
type _image struct {
	*ebiten.Image
}

type Image struct {
	_image
	Box
}

// NewImage returns non-nil value even when fails to load file from given path.
// ebiten.NewImageFromImage will panic when input is nil.
func NewImage(r io.Reader) Image {
	iimg, _, err := image.Decode(r) // image.Image
	if err != nil {
		return Image{}
	}
	if r, ok := r.(io.Closer); ok {
		r.Close()
	}
	eimg := ebiten.NewImageFromImage(iimg) // *ebiten.Image
	return newImageFromEbitenImage(eimg)
}

func newImageFromEbitenImage(eimg *ebiten.Image) Image {
	img := Image{Image: eimg}
	img.Box = NewBox(img)
	return img
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
	eimg := ebiten.NewImage(int(w), int(h))
	return newImageFromEbitenImage(eimg)
}

func (i Image) IsEmpty() bool { return i.Image == nil }

func (i Image) Size() XY {
	if i.IsEmpty() {
		return XY{}
	}
	size := i.Image.Bounds().Size()
	return NewXYFromInts(size.X, size.Y)
}

// func (img Image) Sub(min, max XY) Image {
// 	if img.IsEmpty() {
// 		return Image{}
// 	}
// 	rect := image.Rectangle{
// 		Min: image.Pt(min.IntValues()),
// 		Max: image.Pt(max.IntValues()),
// 	}
// 	return img.SubImage(rect).(Image)
// }

// sub.Fill might fill the destination image permanently.
func (img Image) Draw(dst Image) {
	if img.IsEmpty() {
		return
	}
	dst.DrawImage(img.Image, img.op())
}

// func (img Image) In(p XY) bool {
// 	max := img.Size()
// 	return 0 <= p.X && p.X < max.X &&
// 		0 <= p.Y && p.Y < max.Y
// }
