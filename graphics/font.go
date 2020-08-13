package graphics

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"io/ioutil"
	"log"
	"os"
)

var (
	FontVarelaNormal font.Face
	mplusNormalFont  font.Face
	mplusBigFont     font.Face
)

func init() {
	const assetDir = "C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\asset\\" // todo: 경로
	f, err := os.Open(assetDir + "Varela-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	tt, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	FontVarelaNormal = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

}
func init(){
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	mplusBigFont = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
