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

// Order of fields of drawer: updating fields, others fields, sprites.
type BackgroundDrawer struct {
	Brightness *float64
	Sprite     draws.Sprite
}

func (d BackgroundDrawer) Draw(dst draws.Image) {
	op := draws.Op{}
	op.ColorM.ChangeHSV(0, 1, *d.Brightness)
	d.Sprite.Draw(dst, op)
}

const (
	SignDot = iota
	SignComma
	SignPercent
)

type NumberDrawer struct {
	draws.Timer
	DigitWidth float64
	DigitGap   float64
	Combo      int
	Bounce     float64
	Sprites    [10]draws.Sprite
}

// Each number has different width. Number 0's width is used as standard.
func (d *NumberDrawer) Update(combo int) {
	d.Ticker()
	if d.Combo != combo {
		d.Combo = combo
		d.Timer.Reset()
	}
}

// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (d NumberDrawer) Draw(dst draws.Image) {
	if d.Done() {
		return
	}
	if d.Combo == 0 {
		return
	}
	vs := make([]int, 0)
	for v := d.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at origin, no need to care of two 0.5w.
	w := d.DigitWidth + d.DigitGap
	tx := float64(len(vs)-1) * w / 2
	const (
		bound0 = 0.05
		bound1 = 0.1
	)
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(tx, 0)
		age := d.Age()
		if age < bound0 {
			scale := 0.1 * d.Progress(0, bound0)
			sprite.Move(0, d.Bounce*sprite.H()*scale)
		}
		if age >= bound0 && age < bound1 {
			scale := 0.1 - 0.1*d.Progress(bound0, bound1)
			sprite.Move(0, d.Bounce*sprite.H()*scale)
		}
		sprite.Draw(dst, draws.Op{})
		tx -= w
	}
}

// Todo: DigitWidth -> digitWidth with ScoreSprites[0].W()
type ScoreDrawer struct {
	draws.Timer
	DigitWidth float64 // Use number 0's width.
	DigitGap   float64
	ZeroFill   int
	Score      ctrl.Delayed
	Sprites    [10]draws.Sprite
}

func NewScoreDrawer() ScoreDrawer {
	return ScoreDrawer{
		DigitWidth: ScoreSprites[0].W(),
		DigitGap:   ScoreDigitGap,
		ZeroFill:   1,
		Score:      ctrl.Delayed{Mode: ctrl.DelayedModeExp},
		Sprites:    ScoreSprites,
	}
}
func (d *ScoreDrawer) Update(score float64) {
	d.Score.Update(score)
}

// NumberDrawer's Draw draws each number at the center of constant-width bound.
func (d ScoreDrawer) Draw(dst draws.Image) {
	if d.Done() {
		return
	}
	vs := make([]int, 0)
	score := int(math.Floor(d.Score.Value() + 0.1))
	for v := score; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}
	for i := len(vs); i < d.ZeroFill; i++ {
		vs = append(vs, 0)
	}
	w := d.DigitWidth + d.DigitGap
	var tx float64
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(tx, 0)
		sprite.Move(-w/2+sprite.W()/2, 0) // Need to set at center since origin is RightTop.
		sprite.Draw(dst, draws.Op{})
		tx -= w
	}
}

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

type MeterDrawer struct {
	MaxCountdown int
	Marks        []MeterMark
	Meter        draws.Sprite
	Anchor       draws.Sprite
	Unit         draws.Sprite
}
type MeterMark struct {
	Countdown int
	Offset    int // Derived from time error.
	ColorType int
}

// Todo: draw Meter with ebiten.Image
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
		base := draws.NewSpriteFromSource(draws.Image{Image: i})
		base.Locate(screenSizeX/2, screenSizeY, draws.CenterBottom)
		d.Meter = base
	}
	{
		src := ebiten.NewImage(int(MeterWidth), int(MeterHeight))
		src.Fill(colorRed)
		i := ebiten.NewImageFromImage(src)
		anchor := draws.NewSpriteFromSource(draws.Image{Image: i})
		anchor.Locate(screenSizeX/2, screenSizeY, draws.CenterBottom)
		d.Anchor = anchor
	}
	{
		src := ebiten.NewImage(int(MeterWidth), int(MeterHeight))
		src.Fill(colorWhite)
		i := ebiten.NewImageFromImage(src)
		unit := draws.NewSpriteFromSource(draws.Image{Image: i})
		unit.Locate(screenSizeX/2, screenSizeY, draws.CenterBottom)
		d.Unit = unit
	}
	return
}

func (d *MeterDrawer) AddMark(offset int, colorType int) {
	mark := MeterMark{
		Countdown: d.MaxCountdown,
		Offset:    offset,
		ColorType: colorType,
	}
	d.Marks = append(d.Marks, mark)
}
func (d *MeterDrawer) Update() {
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
func (d MeterDrawer) Draw(dst draws.Image) {
	d.Meter.Draw(dst, draws.Op{})
	d.Anchor.Draw(dst, draws.Op{})
	for _, m := range d.Marks {
		sprite := d.Unit
		op := draws.Op{}
		color := MeterMarkColors[m.ColorType]
		op.ColorM.ScaleWithColor(color)
		if age := d.MarkAge(m); age >= 0.8 {
			op.ColorM.Scale(1, 1, 1, 1-(age-0.8)/0.2)
		}
		op.GeoM.Translate(-float64(m.Offset)*MeterWidth, 0)
		sprite.Draw(dst, op)
	}
}

func (d MeterDrawer) MarkAge(m MeterMark) float64 {
	return 1 - float64(m.Countdown)/float64(d.MaxCountdown)
}
