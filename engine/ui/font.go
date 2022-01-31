package ui

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
)

var (
	FontBold     *truetype.Font
	FontBoldFace font.Face
)

func init() {
	var err error
	FontBold, err = truetype.Parse(gobold.TTF)
	if err != nil {
		panic(err)
	}
	opts := &truetype.Options{}
	opts.Size = 15
	FontBoldFace = truetype.NewFace(FontBold, opts)
}
