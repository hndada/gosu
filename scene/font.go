package scene

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var goregularFont *truetype.Font

const dpi = 72

var Faces = make(map[int]font.Face)

func init() {
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	goregularFont = f
}

func Face(size int) font.Face {
	if face, ok := Faces[size]; ok {
		return face
	}
	face := truetype.NewFace(goregularFont, &truetype.Options{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	Faces[size] = face
	return face
}

// FaceDefault = basicfont.Face7x13
