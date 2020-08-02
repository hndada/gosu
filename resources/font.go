package resources

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"log"
)

var (
	SampleText      = `The quick brown fox jumps over the lazy dog.`
	MplusNormalFont font.Face
	MplusBigFont    font.Face
)

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	MplusNormalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	MplusBigFont = truetype.NewFace(tt, &truetype.Options{
		Size: 48,
		DPI:  dpi,
	})
}
