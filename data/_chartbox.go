package db

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func NewChartInfoSprite(info ChartInfo) draws.Sprite { // h mode.ChartHeader, SubMode int, level float64
	const (
		dx = 20 // dot x
		dy = 30 // dot y
	)
	img := image.NewRGBA(image.Rect(0, 0, BoxWidth, BoxHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{Purple}, image.Point{}, draw.Src)
	inRect := image.Rect(border, border, BoxWidth-border, BoxHeight-border)
	draw.Draw(img, inRect, &image.Uniform{color.White}, image.Point{}, draw.Src)
	t := fmt.Sprintf("(%dK Lv %.1f) %s [%s]", info.SubMode, info.Level, info.Header.MusicName, info.Header.ChartName)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(dx * 64), Y: fixed.Int26_6(dy * 64)},
	}
	d.DrawString(t)
	return draws.Sprite{
		I: ebiten.NewImageFromImage(img),
		W: float64(BoxWidth),
		H: float64(BoxHeight),
		X: mode.ScreenSizeX - float64(BoxWidth),
		// Y is not fixed.
	}
}
