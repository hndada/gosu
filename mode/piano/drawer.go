package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

var (
	ColorKool = color.NRGBA{0, 170, 242, 255}   // Blue
	ColorCool = color.NRGBA{85, 251, 255, 255}  // Skyblue
	ColorGood = color.NRGBA{51, 255, 40, 255}   // Lime
	ColorBad  = color.NRGBA{244, 177, 0, 255}   // Yellow
	ColorMiss = color.NRGBA{109, 120, 134, 255} // Gray
)

// HCI: BackgroundRedDrawer
type BackgroundRedDrawer struct {
	draws.Timer
	Sprite draws.Sprite
}

func NewBackgroundRedDrawer() (d BackgroundRedDrawer) {
	image := draws.NewImage(ScreenSizeX, ScreenSizeX)
	image.Fill(color.RGBA{255, 128, 128, 255})
	sprite := draws.NewSprite(image)
	// timer := draws.NewTimer(draws.ToTick(100, TPS), draws.ToTick(100, TPS))
	// timer.Tick = timer.MaxTick
	return BackgroundRedDrawer{
		Timer:  draws.NewTimer(draws.ToTick(300, TPS), draws.ToTick(300, TPS)),
		Sprite: sprite,
	}
}
func (d *BackgroundRedDrawer) Update(passed bool) {
	d.Ticker()
	if passed {
		d.Timer.Reset()
	}
}

func (d BackgroundRedDrawer) Draw(dst draws.Image) {
	if d.Tick >= d.MaxTick {
		return
	}
	op := draws.Op{}
	// value := 1 - d.Age()
	// op.ColorM.ChangeHSV(0, 1, value)
	d.Sprite.Draw(dst, op)
}

type FieldDrawer struct {
	Sprite draws.Sprite
}

func (d FieldDrawer) Draw(dst draws.Image) {
	d.Sprite.Draw(dst, draws.Op{})
}

type HintDrawer struct {
	Sprite draws.Sprite
}

func (d HintDrawer) Draw(dst draws.Image) {
	d.Sprite.Draw(dst, draws.Op{})
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
type BarDrawer struct {
	Cursor   float64
	Farthest *Bar
	Nearest  *Bar
	Sprite   draws.Sprite
}

func (d *BarDrawer) Update(cursor float64) {
	d.Cursor = cursor
	// When Farthest's prevs are still out of screen due to speed change.
	for d.Farthest.Prev != nil &&
		d.Farthest.Prev.Position-d.Cursor > S.maxPosition {
		d.Farthest = d.Farthest.Prev
	}
	// When Farthest is in screen, next note goes fetched if possible.
	for d.Farthest.Next != nil &&
		d.Farthest.Position-d.Cursor <= S.maxPosition {
		d.Farthest = d.Farthest.Next
	}
	// When Nearest is still in screen due to speed change.
	for d.Nearest.Prev != nil &&
		d.Nearest.Position-d.Cursor > S.minPosition {
		d.Nearest = d.Nearest.Prev
	}
	// When Nearest's next is still out of screen, next note goes fetched.
	for d.Nearest.Next != nil &&
		d.Nearest.Next.Position-d.Cursor <= S.minPosition {
		d.Nearest = d.Nearest.Next
	}
}

func (d BarDrawer) Draw(dst draws.Image) {
	if d.Farthest == nil || d.Nearest == nil {
		return
	}
	for b := d.Farthest; b != d.Nearest.Prev; b = b.Prev {
		pos := b.Position - d.Cursor
		s := d.Sprite
		s.Move(0, -pos)
		s.Draw(dst, draws.Op{})
	}
}

// Notes are fixed. Lane itself moves, all notes move same amount.
type NoteDrawer struct {
	draws.Timer
	Cursor   float64
	Farthest *Note
	Nearest  *Note
	Sprites  [4]draws.Animation
	holding  bool
}

// Farthest and Nearest are borders of displaying notes.
// All in-screen notes are confirmed to be drawn when drawing from Farthest to Nearest.
func (d *NoteDrawer) Update(cursor float64, holding bool) {
	d.Ticker()
	d.Cursor = cursor
	defer func() { d.holding = holding }()
	if d.Farthest == nil || d.Nearest == nil {
		return
	}
	// When Farthest's prevs are still out of screen due to speed change.
	for d.Farthest.Prev != nil &&
		d.Farthest.Prev.Position-d.Cursor > S.maxPosition {
		d.Farthest = d.Farthest.Prev
	}
	// When Farthest is in screen, next note goes fetched if possible.
	for d.Farthest.Next != nil &&
		d.Farthest.Position-d.Cursor <= S.maxPosition {
		d.Farthest = d.Farthest.Next
	}
	// When Nearest is still in screen due to speed change.
	for d.Nearest.Prev != nil &&
		d.Nearest.Position-d.Cursor > S.minPosition {
		d.Nearest = d.Nearest.Prev
	}
	// When Nearest's next is still out of screen, next note goes fetched.
	for d.Nearest.Next != nil &&
		d.Nearest.Next.Position-d.Cursor <= S.minPosition {
		d.Nearest = d.Nearest.Next
	}
}

// Draw from farthest to nearest to make nearer notes priorly exposed.
func (d NoteDrawer) Draw(dst draws.Image) {
	if d.Farthest == nil || d.Nearest == nil {
		return
	}
	for n := d.Farthest; n != nil && n != d.Nearest.Prev; n = n.Prev {
		if n.Type == Tail {
			d.DrawBody(dst, n)
		}
		s := d.Frame(d.Sprites[n.Type])
		pos := n.Position - d.Cursor
		s.Move(0, -pos)
		op := draws.Op{}
		if n.Marked {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		// HCI
		if n.passed {
			op.ColorM.ScaleWithColor(color.NRGBA{235, 69, 44, 255}) // Red
		}
		s.Draw(dst, op)
	}
}

// DrawBody draws scaled, corresponding sub-image of Body sprite.
func (d NoteDrawer) DrawBody(dst draws.Image, tail *Note) {
	head := tail.Prev
	body := d.Sprites[Body][0]
	if d.holding {
		body = d.Frame(d.Sprites[Body])
	}
	length := tail.Position - head.Position // + BodyGain
	length += S.NoteHeigth                  // - bodyLoss
	if length < 0 {
		length = 0
	}
	body.SetSize(body.W(), length)
	ty := head.Position - d.Cursor
	body.Move(0, -ty)

	op := draws.Op{}
	if tail.Marked {
		op.ColorM.ChangeHSV(0, 0.3, 0.3)
	}
	body.Draw(dst, op)
}

type KeyDrawer struct {
	draws.Timer
	Sprites     [2]draws.Sprite
	lastPressed bool
}

func (d *KeyDrawer) Update(pressed bool) {
	d.Ticker()
	if !d.lastPressed && pressed {
		d.Timer.Reset()
	}
	d.lastPressed = pressed
}

// KeyDrawer draws for a while even when pressed off very shortly.
func (d KeyDrawer) Draw(dst draws.Image) {
	const (
		up = iota
		down
	)
	s := d.Sprites[up]
	if d.lastPressed || d.Tick < d.MaxTick {
		s = d.Sprites[down]
	}
	s.Draw(dst, draws.Op{})
}

type KeyLightingDrawer struct {
	draws.Timer
	Sprite      draws.Sprite
	Color       color.NRGBA
	lastPressed bool
}

func (d *KeyLightingDrawer) Update(pressed bool) {
	d.Ticker()
	if !d.lastPressed && pressed {
		d.Timer.Reset()
	}
	d.lastPressed = pressed
}

// KeyLightingDrawer draws for a while even when pressed off very shortly.
func (d KeyLightingDrawer) Draw(dst draws.Image) {
	if d.lastPressed || d.Tick < d.MaxTick {
		op := draws.Op{}
		op.ColorM.ScaleWithColor(d.Color)
		d.Sprite.Draw(dst, op)
	}
}

type HitLightingDrawer struct {
	draws.Timer
	Sprites draws.Animation
}

// HitLightingDrawer draws when Normal is Hit or Tail is Release.
func (d *HitLightingDrawer) Update(hit bool) {
	d.Ticker()
	if hit {
		d.Timer.Reset()
	}
}
func (d HitLightingDrawer) Draw(dst draws.Image) {
	if d.IsDone() {
		return
	}
	op := draws.Op{}
	// opaque := UserSettings.HitLightingOpaque * (1 - d.Progress(0.75, 1))
	op.ColorM.Scale(1, 1, 1, S.HitLightingOpaque)
	d.Frame(d.Sprites).Draw(dst, op)
}

type HoldLightingDrawer struct {
	draws.Timer
	Sprites     draws.Animation
	lastPressed bool
}

func (d *HoldLightingDrawer) Update(pressed bool) {
	d.Ticker()
	if d.lastPressed != pressed {
		d.Timer.Reset()
		d.lastPressed = pressed
	}
}
func (d HoldLightingDrawer) Draw(dst draws.Image) {
	if !d.lastPressed {
		return
	}
	op := draws.Op{}
	op.ColorM.Scale(1, 1, 1, S.HoldLightingOpaque)
	d.Frame(d.Sprites).Draw(dst, op)
}

type JudgmentDrawer struct {
	draws.Timer
	Sprites  []draws.Animation
	Judgment mode.Judgment
}

func NewJudgmentDrawer(sprites []draws.Animation) (d JudgmentDrawer) {
	const frameDuration = 1000.0 / 60
	count := float64(len(sprites))
	period := int64(frameDuration * count)
	return JudgmentDrawer{
		Timer:   draws.NewTimer(draws.ToTick(250, TPS), draws.ToTick(period, TPS)),
		Sprites: sprites,
	}
}
func (d *JudgmentDrawer) Update(worst mode.Judgment) {
	d.Ticker()
	if worst.IsValid() {
		d.Judgment = worst
		d.Timer.Reset()
	}
}

func (d JudgmentDrawer) Draw(dst draws.Image) {
	if d.IsDone() {
		return
	}
	var idx int
	for i, j := range Judgments {
		if j.Window == d.Judgment.Window {
			idx = i
			break
		}
	}
	age := d.Age()
	s := d.Frame(d.Sprites[idx])
	const (
		bound0 = 0.1
		bound1 = 0.2
		bound2 = 0.9
	)
	scale := 1.0
	if age < bound0 {
		scale = 1 + 0.15*d.Progress(0, bound0)
	}
	if age >= bound0 && age < bound1 {
		scale = 1.15 - 0.15*d.Progress(bound0, bound1)
	}
	if age >= bound2 {
		scale = 1 - 0.25*d.Progress(bound2, 1)
	}
	s.MultiplyScale(scale)
	s.Draw(dst, draws.Op{})
}
