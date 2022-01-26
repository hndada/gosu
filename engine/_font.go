package engine

import (
	"io/ioutil"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	titleFontSize = fontSize * 1.5
	fontSize      = 40
	smallFontSize = fontSize / 2
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func loadFont(cwd string) {
	dir := filepath.Join(cwd, "skin")
	name := "Raleway-Bold.ttf"
	path := filepath.Join(dir, name)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
}
