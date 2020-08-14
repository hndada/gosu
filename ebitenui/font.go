package ebitenui

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"io/ioutil"
	"log"
	"os"
)

var (
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
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}


func LoadFontFromFile(path string, op *truetype.Options) font.Face {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(f, op)

}