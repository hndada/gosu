package mode

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.MultiplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}

type BackgroundDrawer struct {
	Sprite draws.Sprite
}

func (d BackgroundDrawer) Draw(dst draws.Image) {
	op := draws.Op{}
	value := Settings.BackgroundBrightness
	op.ColorM.ChangeHSV(0, 1, value)
	d.Sprite.Draw(dst, op)
}

type ScoreDrawer struct {
	digitWidth float64 // Use number 0's width.
	DigitGap   float64
	ZeroFill   int
	Score      Delayed
	Sprites    []draws.Sprite
}

func NewScoreDrawer(sprites []draws.Sprite) ScoreDrawer {
	return ScoreDrawer{
		digitWidth: sprites[0].W(),
		DigitGap:   Settings.ScoreDigitGap,
		ZeroFill:   1,
		Score:      NewDelayed(),
		Sprites:    sprites[:10],
	}
}
func (d *ScoreDrawer) Update(score float64) {
	d.Score.Update(score)
}

func (d ScoreDrawer) Draw(dst draws.Image) {
	vs := make([]int, 0)
	score := int(math.Floor(d.Score.Delayed + 0.1))
	for v := score; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}
	for i := len(vs); i < d.ZeroFill; i++ {
		vs = append(vs, 0)
	}
	w := d.digitWidth + d.DigitGap
	var tx float64
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(tx, 0)
		sprite.Move(-w/2+sprite.W()/2, 0) // Need to set at center since anchor is RightTop.
		sprite.Draw(dst, draws.Op{})
		tx -= w
	}
}

// Todo: add combo *int and skip passing combo value?
type ComboDrawer struct {
	draws.Timer
	DigitWidth float64 // Use number 0's width.
	DigitGap   float64
	Combo      int
	Bounce     float64
	Sprites    [10]draws.Sprite
}

// Each number has different width. Number 0's width is used as standard.
func (d *ComboDrawer) Update(combo int) {
	d.Ticker()
	if d.Combo != combo {
		d.Combo = combo
		d.Timer.Reset()
	}
}

// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (d ComboDrawer) Draw(dst draws.Image) {
	if d.IsDone() {
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
	// Since sprites are already at anchor, no need to care of two 0.5w.
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

// Todo: refactor MeterDrawer code
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

// Anchor is a unit sprite constantly drawn at the middle of meter.
// Todo: should w and h be math.Ceil()?
func NewMeterDrawer(js []Judgment, colors []color.NRGBA) (d MeterDrawer) {
	var (
		colorMeter = color.NRGBA{0, 0, 0, 128}       // Dark
		colorWhite = color.NRGBA{255, 255, 255, 192} // White
		colorRed   = color.NRGBA{255, 0, 0, 192}     // Red
	)
	var (
		W = Settings.MeterUnit
		H = Settings.MeterHeight
	)
	d.MaxCountdown = ToTick(4 * time.Second)
	{
		miss := js[len(js)-1]
		w := 1 + 2*W*float64(miss.Window)
		h := H
		src := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		draw.Draw(src, src.Bounds(), &image.Uniform{colorMeter}, image.Point{}, draw.Src)
		// src := draws.NewImage(w, h)
		// src.Fill(colorMeter)
		y1, y2 := h*(0.5-0.125), h*(0.5+0.125) // Height of inner is 1/4 of meter's.
		for i := range js {
			j := js[len(js)-1-i] // In reverse order.
			clr := colors[len(js)-1-i]

			w := 1 + 2*W*float64(j.Window)
			x1 := W * float64(miss.Window-j.Window)
			x2 := x1 + w
			rect := image.Rect(int(x1), int(y1), int(x2), int(y2))
			draw.Draw(src, rect, &image.Uniform{clr}, image.Point{}, draw.Src)
			// dst := src.SubImage(rect).(*ebiten.Image)
			// ebitenutil.DrawRect(dst, 0, 0, w, h/4, clr)
		}
		i := ebiten.NewImageFromImage(src)
		base := draws.NewSpriteFromSource(draws.Image{Image: i})
		base.Locate(ScreenSizeX/2, ScreenSizeY, draws.CenterBottom)
		d.Meter = base
	}
	{
		src := draws.NewImage(W, H)
		src.Fill(colorRed)
		sprite := draws.NewSpriteFromSource(src)
		sprite.Locate(ScreenSizeX/2, ScreenSizeY, draws.CenterBottom)
		d.Anchor = sprite
	}
	{
		src := draws.NewImage(W, H)
		src.Fill(colorWhite)
		sprite := draws.NewSpriteFromSource(src)
		sprite.Locate(ScreenSizeX/2, ScreenSizeY, draws.CenterBottom)
		d.Unit = sprite
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
		op.GeoM.Translate(-float64(m.Offset)*Settings.MeterUnit, 0)
		sprite.Draw(dst, op)
	}
}

// Todo: remove
func (d MeterDrawer) MarkAge(m MeterMark) float64 {
	return 1 - float64(m.Countdown)/float64(d.MaxCountdown)
}
