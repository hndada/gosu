package draws

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Typeface: a style.
// Type Family: a set of relating styles.
// Font: a file of Typeface.
// Font Family: a set of relating fonts.
// Font Face: a file used for specific style.
type FontKey struct {
	Family string
	Face   string
}

type FaceKey struct {
	FontKey
	Size float64
}

var DefaultFontKey = FontKey{"go", "regular"}
var DefaultFont *truetype.Font

var Fonts = make(map[FontKey]*truetype.Font)
var Faces = make(map[FaceKey]font.Face)

func init() {
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	Fonts[DefaultFontKey] = f
	DefaultFont = f
}
func NewFace(fontKey FontKey, size float64) font.Face {
	var (
		_font *truetype.Font
		face  font.Face
		ok    bool
	)
	faceKey := FaceKey{fontKey, size}
	face, ok = Faces[faceKey]
	if !ok {
		_font, ok = Fonts[fontKey]
		if !ok {
			_font = DefaultFont
		}
		face = truetype.NewFace(_font, &truetype.Options{
			Size: size,
		})
		Faces[faceKey] = face
	}
	return face
}
func DefaultFace(size float64) font.Face {
	return NewFace(DefaultFontKey, size)
}
