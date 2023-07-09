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
		base := draws.NewSprite(draws.Image{Image: i})
		base.Locate(ScreenSizeX/2, ScreenSizeY, draws.CenterBottom)
		d.Meter = base
	}
	{
		src := draws.NewImage(W, H)
		src.Fill(colorRed)
		sprite := draws.NewSprite(src)
		sprite.Locate(ScreenSizeX/2, ScreenSizeY, draws.CenterBottom)
		d.Anchor = sprite
	}
	{
		src := draws.NewImage(W, H)
		src.Fill(colorWhite)
		sprite := draws.NewSprite(src)
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
