package mode

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/render"
)

// Todo: BarLine color settings
var (
	Dark  = color.NRGBA{0, 0, 0, 128}
	Gray  = color.NRGBA{109, 120, 134, 255}
	White = color.NRGBA{255, 255, 255, 192}

	Red    = color.NRGBA{255, 0, 0, 128}
	Yellow = color.NRGBA{244, 177, 0, 255}
	Lime   = color.NRGBA{51, 255, 40, 255}
	Sky    = color.NRGBA{85, 251, 255, 255}
	Blue   = color.NRGBA{0, 170, 242, 255}
	Purple = color.NRGBA{213, 0, 242, 128}
)

// Timing meter. the height of colored rectangle is 1/4 of meter's.
// Anchor is a unit sprite constantly drawn at the middle of meter.
func NewTimingMeter(js []Judgment, colors []color.NRGBA) (render.Sprite, render.Sprite, render.Sprite) {
	Miss := js[len(js)-1]
	meterW := 1 + 2*int(TimingMeterWidth)*int(Miss.Window)
	meterH := int(TimingMeterHeight)
	meterSrc := image.NewRGBA(image.Rect(0, 0, meterW, meterH))
	draw.Draw(meterSrc, meterSrc.Bounds(), &image.Uniform{Dark}, image.Point{}, draw.Src)
	y1, y2 := int(float64(meterH)*0.375), int(float64(meterH)*0.625)
	for i, color := range colors {
		j := js[len(js)-i]
		w := 1 + 2*int(TimingMeterWidth)*int(j.Window)
		x1 := int(TimingMeterWidth) * (int(Miss.Window - j.Window))
		x2 := x1 + w
		rect := image.Rect(x1, y1, x2, y2)
		draw.Draw(meterSrc, rect, &image.Uniform{color}, image.Point{}, draw.Src)
	}

	meter := render.Sprite{
		I: ebiten.NewImageFromImage(meterSrc),
		W: float64(meterW),
		H: float64(meterH),
		Y: ScreenSizeY - TimingMeterHeight,
	}
	meter.SetCenterX(ScreenSizeX / 2)

	anchorSrc := ebiten.NewImage(int(TimingMeterWidth), int(TimingMeterHeight))
	anchorSrc.Fill(Red)
	anchor := render.Sprite{
		I: anchorSrc,
		W: TimingMeterWidth,
		H: TimingMeterHeight,
		Y: ScreenSizeY - TimingMeterHeight,
	}
	anchor.SetCenterX(ScreenSizeX / 2)

	unitSrc := ebiten.NewImage(int(TimingMeterWidth), int(TimingMeterHeight))
	unitSrc.Fill(White)
	unit := render.Sprite{
		I: unitSrc,
		W: TimingMeterWidth,
		H: TimingMeterHeight,
		// unit's x value is not fixed.
		Y: ScreenSizeY - TimingMeterHeight,
	}
	return meter, anchor, unit
}

type TimingMark struct {
	TimeDiff int64
	BornTick int
	IsTail   bool
}
