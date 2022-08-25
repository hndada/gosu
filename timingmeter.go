package gosu

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

const fadeInDuration = 800 // Duration of fade-in of timing meter marks.
var TimingMeterMarkDuration int = TimeToTick(4000)
var (
	ColorKool = color.NRGBA{0, 170, 242, 255}   // Blue
	ColorCool = color.NRGBA{85, 251, 255, 255}  // Skyblue
	ColorGood = color.NRGBA{51, 255, 40, 255}   // Lime
	ColorBad  = color.NRGBA{244, 177, 0, 255}   // Yellow
	ColorMiss = color.NRGBA{109, 120, 134, 255} // Gray
)

type TimingMeter struct {
	Meter  draws.Sprite
	Anchor draws.Sprite
	Unit   draws.Sprite
	Marks  []TimingMeterMark
}
type TimingMeterMark struct {
	Countdown int
	TimeDiff  int64
	Color     color.NRGBA
}

// Height of colored rectangle is 1/4 of Timing meter's.
// Anchor is a unit sprite constantly drawn at the middle of meter.
func NewTimingMeter(js []Judgment, colors []color.NRGBA) TimingMeter {
	var (
		dark  = color.NRGBA{0, 0, 0, 128}
		white = color.NRGBA{255, 255, 255, 192}
		red   = color.NRGBA{255, 0, 0, 128}
	)
	var tm TimingMeter
	Miss := js[len(js)-1]

	meterW := 1 + 2*int(TimingMeterWidth)*int(Miss.Window)
	meterH := int(TimingMeterHeight)
	meterSrc := image.NewRGBA(image.Rect(0, 0, meterW, meterH))
	draw.Draw(meterSrc, meterSrc.Bounds(), &image.Uniform{dark}, image.Point{}, draw.Src)
	y1, y2 := int(float64(meterH)*0.375), int(float64(meterH)*0.625)
	for i := 0; i < len(js); i++ {
		j := js[len(js)-1-i]
		color := colors[len(js)-1-i]
		w := 1 + 2*int(TimingMeterWidth)*int(j.Window)
		x1 := int(TimingMeterWidth) * (int(Miss.Window - j.Window))
		x2 := x1 + w
		rect := image.Rect(x1, y1, x2, y2)
		draw.Draw(meterSrc, rect, &image.Uniform{color}, image.Point{}, draw.Src)
	}
	tm.Meter = draws.Sprite{
		I: ebiten.NewImageFromImage(meterSrc),
		W: float64(meterW),
		H: float64(meterH),
		Y: ScreenSizeY - TimingMeterHeight,
	}
	tm.Meter.SetCenterX(ScreenSizeX / 2)

	anchorSrc := ebiten.NewImage(int(TimingMeterWidth), int(TimingMeterHeight))
	anchorSrc.Fill(red)
	tm.Anchor = draws.Sprite{
		I: anchorSrc,
		W: TimingMeterWidth,
		H: TimingMeterHeight,
		Y: ScreenSizeY - TimingMeterHeight,
	}
	tm.Anchor.SetCenterX(ScreenSizeX / 2)

	unitSrc := ebiten.NewImage(int(TimingMeterWidth), int(TimingMeterHeight))
	unitSrc.Fill(white)
	tm.Unit = draws.Sprite{
		I: unitSrc,
		W: TimingMeterWidth,
		H: TimingMeterHeight,
		// X is not fixed.
		Y: ScreenSizeY - TimingMeterHeight,
	}
	return tm
}
func (tm *TimingMeter) Update(newMarks []TimingMeterMark) {
	tm.Marks = append(tm.Marks, newMarks...)
	cursor := 0
	for i, m := range tm.Marks {
		if m.Countdown == 0 {
			i++
		} else {
			tm.Marks[i].Countdown--
		}
	}
	tm.Marks = tm.Marks[cursor:] // Drop old marks.
}
func (tm TimingMeter) Draw(screen *ebiten.Image) {
	tm.Meter.Draw(screen)
	for _, m := range tm.Marks {
		sprite := tm.Unit
		sprite.X = screenSizeX/2 + float64(m.TimeDiff)*TimingMeterWidth
		op := sprite.Op()
		clr := m.Color
		if TickToTime(m.Countdown) < fadeInDuration {
			clr.A = uint8(float64(clr.A) * float64(m.Countdown) / fadeInDuration)
		}
		op.ColorM.ScaleWithColor(clr)
		screen.DrawImage(sprite.I, op)
	}
	tm.Anchor.Draw(screen)
}
