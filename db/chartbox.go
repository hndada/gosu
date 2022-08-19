package db

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/render"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// https://github.com/vmihailenco/msgpack
// https://github.com/osuripple/cheesegull
type ChartInfo struct {
	Path string
	// Mods mode.Mods
	Header mode.ChartHeader
	Mode   int
	Mode2  int
	Level  float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
	// Box render.Sprite
}

func NewChartInfo(c *mode.Chart, fpath string, level float64) ChartInfo {
	mainBPM, minBPM, maxBPM := mode.BPMs(c.TransPoints, c.Duration)
	cb := ChartInfo{
		Path:   fpath,
		Header: c.ChartHeader,
		Mode:   c.Mode,
		Mode2:  c.Mode2,
		Level:  level,

		Duration:   c.Duration,
		NoteCounts: c.NoteCounts,
		MainBPM:    mainBPM,
		MinBPM:     minBPM,
		MaxBPM:     maxBPM,
	}
	// cb.Box = NewBoxSprite(c, level)
	return cb
}

const (
	BoxWidth  = 450 // Box width
	BoxHeight = 50  // Box height
)

var Purple = color.RGBA{172, 49, 174, 255}

func NewChartInfoSprite(info ChartInfo) render.Sprite { // h mode.ChartHeader, mode2 int, level float64
	const border = 3
	const (
		dx = 20 // dot x
		dy = 30 // dot y
	)
	img := image.NewRGBA(image.Rect(0, 0, BoxWidth, BoxHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{Purple}, image.Point{}, draw.Src)
	inRect := image.Rect(border, border, BoxWidth-border, BoxHeight-border)
	draw.Draw(img, inRect, &image.Uniform{color.White}, image.Point{}, draw.Src)
	t := fmt.Sprintf("(%dK Lv %.1f) %s [%s]", info.Mode2, info.Level, info.Header.MusicName, info.Header.ChartName)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(dx * 64), Y: fixed.Int26_6(dy * 64)},
	}
	d.DrawString(t)
	return render.Sprite{
		I: ebiten.NewImageFromImage(img),
		W: float64(BoxWidth),
		H: float64(BoxHeight),
		X: mode.ScreenSizeX - float64(BoxWidth),
		// Y is not fixed.
	}
}
