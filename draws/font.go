package draws

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	Face12 font.Face
	Face16 font.Face
	Face20 font.Face
	Face24 font.Face
)

// FaceDefault = basicfont.Face7x13

func init() {
	const dpi = 72
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	Face12 = truetype.NewFace(f, &truetype.Options{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	Face16 = truetype.NewFace(f, &truetype.Options{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	Face20 = truetype.NewFace(f, &truetype.Options{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	Face24 = truetype.NewFace(f, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
