package gosu

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
)

// type Direction int

// const (
// 	Upward   Direction = iota // e.g., Rhythm games using feet.
// 	Downward                  // e.g., Piano mode.
// 	Leftward                  // e.g., Drum mode.
// 	Rightward
// )

// Todo: use Effecter in draws.BaseDrawer
type BackgroundDrawer struct {
	Sprite  draws.Sprite
	Dimness *float64
}

func (d BackgroundDrawer) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.ChangeHSV(0, 1, *d.Dimness)
	d.Sprite.Draw(screen, op)
}

type ScoreDrawer struct {
	draws.NumberDrawer
	Score ctrl.Delayed
}

func NewScoreDrawer() (d ScoreDrawer) {
	d = ScoreDrawer{
		NumberDrawer: draws.NumberDrawer{
			Sprites:       ScoreSprites,
			SignSprites:   SignSprites,
			Origin:        ScoreSprites[0].Origin(),
			DigitWidth:    ScoreSprites[0].W(),
			DigitGap:      0,
			FractionDigit: 0,
			ZeroFill:      1,
		},
		Score: ctrl.Delayed{Mode: ctrl.DelayedModeExp},
	}
	return
}
func (d *ScoreDrawer) Update(score float64) {
	d.Score.Update(score)
	i := int(math.Ceil(score))
	d.NumberDrawer.Update(i, 0)
}

// Todo: merge with drawer.go
var (
	ColorKool = color.NRGBA{0, 170, 242, 255}   // Blue
	ColorCool = color.NRGBA{85, 251, 255, 255}  // Skyblue
	ColorGood = color.NRGBA{51, 255, 40, 255}   // Lime
	ColorBad  = color.NRGBA{244, 177, 0, 255}   // Yellow
	ColorMiss = color.NRGBA{109, 120, 134, 255} // Gray
)

var MeterMarkColors = []color.NRGBA{
	{255, 255, 255, 192}, // White
	{213, 0, 242, 192},   // Purple
	{252, 83, 6, 255},    // Orange
}

// Meter is also known as TimingMeter.
type MeterDrawer struct {
	MaxCountdown int
	Meter        draws.Sprite
	Anchor       draws.Sprite
	Unit         draws.Sprite
	Marks        []MeterMark
}
type MeterMark struct {
	Countdown int
	Offset    int // Derived from time error.
	ColorType int
}

// Anchor is a unit sprite constantly drawn at the middle of meter.
func NewMeterDrawer(js []Judgment, colors []color.NRGBA) (d MeterDrawer) {
	var (
		colorMeter = color.NRGBA{0, 0, 0, 128}       // Dark
		colorWhite = color.NRGBA{255, 255, 255, 192} // White
		colorRed   = color.NRGBA{255, 0, 0, 192}     // Red
	)
	d.MaxCountdown = TimeToTick(4000)
	{
		miss := js[len(js)-1]
		w := 1 + 2*math.Ceil(MeterWidth*float64(miss.Window))
		h := MeterHeight
		src := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		draw.Draw(src, src.Bounds(), &image.Uniform{colorMeter}, image.Point{}, draw.Src)

		// Height of colored range is 1/4 of meter's.
		y1, y2 := math.Ceil(h*0.375), math.Ceil(h*0.625)
		for i := range js { // In reverse order.
			j := js[len(js)-1-i]
			clr := colors[len(js)-1-i]
			w := 1 + 2*math.Ceil(MeterWidth*float64(j.Window))
			x1 := math.Ceil(MeterWidth * float64(miss.Window-j.Window))
			x2 := x1 + w
			rect := image.Rect(int(x1), int(y1), int(x2), int(y2))
			draw.Draw(src, rect, &image.Uniform{clr}, image.Point{}, draw.Src)
		}
		i := ebiten.NewImageFromImage(src)
		base := draws.NewSpriteFromImage(i)
		base.SetPosition(screenSizeX/2, screenSizeY, draws.OriginCenterBottom)
		d.Meter = base
	}
	{
		src := ebiten.NewImage(int(MeterWidth), int(MeterHeight))
		src.Fill(colorRed)
		i := ebiten.NewImageFromImage(src)
		anchor := draws.NewSpriteFromImage(i)
		anchor.SetPosition(screenSizeX/2, screenSizeY, draws.OriginCenterBottom)
		d.Anchor = anchor
	}
	{
		src := ebiten.NewImage(int(MeterWidth), int(MeterHeight))
		src.Fill(colorWhite)
		i := ebiten.NewImageFromImage(src)
		unit := draws.NewSpriteFromImage(i)
		unit.SetPosition(screenSizeX/2, screenSizeY, draws.OriginCenterBottom)
		d.Unit = unit
	}
	return
}
func (d *MeterDrawer) Update(newMarks []MeterMark) {
	d.Marks = append(d.Marks, newMarks...)
	cursor := 0
	for i, m := range d.Marks {
		if m.Countdown == 0 {
			cursor++
		} else {
			d.Marks[i].Countdown--
		}
	}
	d.Marks = d.Marks[cursor:] // Drop old marks.
}
func (d MeterDrawer) Draw(screen *ebiten.Image) {
	d.Meter.Draw(screen, nil)
	for _, m := range d.Marks {
		sprite := d.Unit
		op := &ebiten.DrawImageOptions{}
		clr := MeterMarkColors[m.ColorType]
		op.ColorM.ScaleWithColor(clr)
		age := float64(m.Countdown) / float64(d.MaxCountdown)
		draws.Fader(op, age)
		op.GeoM.Translate(float64(m.Offset)*MeterWidth, 0)
		sprite.Draw(screen, op)
	}
	d.Anchor.Draw(screen, nil)
}
