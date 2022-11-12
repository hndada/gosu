package drum

import (
	"image/color"
	"math/rand"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
)

type StageDrawer struct {
	draws.Timer
	Highlight    bool
	FieldSprites [2]draws.Sprite
	HintSprites  [2]draws.Sprite
}

func (d *StageDrawer) Update(highlight bool) {
	d.Ticker()
	if d.Highlight != highlight {
		d.Timer.Reset()
		d.Highlight = highlight
	}
}

func (d StageDrawer) Draw(dst draws.Image) {
	const (
		idle = iota
		high
	)
	op := draws.Op{}
	op.ColorM.Scale(1, 1, 1, FieldDarkness)
	d.FieldSprites[idle].Draw(dst, op)
	d.HintSprites[idle].Draw(dst, draws.Op{})
	if d.Highlight || d.Tick < d.MaxTick {
		var opField, opHint draws.Op
		if d.Highlight {
			opField.ColorM.Scale(1, 1, 1, FieldDarkness*d.Age())
			opHint.ColorM.Scale(1, 1, 1, FieldDarkness*d.Age())
		} else {
			opField.ColorM.Scale(1, 1, 1, FieldDarkness*(1-d.Age()))
			opHint.ColorM.Scale(1, 1, 1, FieldDarkness*(1-d.Age()))
		}
		opHint.ColorM.ScaleWithColor(ColorYellow)
		d.FieldSprites[high].Draw(dst, opField)
		d.HintSprites[high].Draw(dst, opHint)
	}
}

// Floating-type lane drawer.
type BarDrawer struct {
	Time   int64
	Bars   []*Bar
	Sprite draws.Sprite
}

func (d *BarDrawer) Update(time int64) {
	d.Time = time
}
func (d BarDrawer) Draw(dst draws.Image) {
	for _, b := range d.Bars {
		pos := b.Speed * float64(b.Time-d.Time)
		if pos <= maxPosition && pos >= minPosition {
			d.Sprite.Move(pos, 0)
			d.Sprite.Draw(dst, draws.Op{})
		}
	}
}

type ShakeDrawer struct {
	draws.Timer
	Time    int64
	Staged  *Note
	Sprites [2]draws.Sprite
}

func (d *ShakeDrawer) Update(time int64, staged *Note) {
	d.Ticker()
	d.Time = time
	if d.Staged != staged {
		if d.Staged.HitTick == d.Staged.Tick {
			d.Timer.Reset()
		}
		d.Staged = staged
	}
}
func (d ShakeDrawer) Draw(dst draws.Image) {
	const (
		outer = iota
		inner
	)
	if d.Tick < d.MaxTick {
		scale := 1 + 0.25*d.Progress(0, 1)
		alpha := 1 - d.Progress(0, 1)
		op := draws.Op{}
		op.ColorM.Scale(1, 1, 1, alpha)
		{
			sprite := d.Sprites[outer]
			sprite.ApplyScale(scale)
			sprite.Draw(dst, op)
		}
		{
			sprite := d.Sprites[inner]
			sprite.ApplyScale(scale)
			sprite.Draw(dst, op)
		}
	}
	if d.Staged == nil {
		return
	}
	if d.Staged.Time > d.Time {
		return
	}
	{
		op := draws.Op{}
		scale := 0.25 + 0.75*float64(d.Time-d.Staged.Time)/200
		if scale > 1 {
			scale = 1
		}
		op.ColorM.Scale(1, 1, 1, scale)
		sprite := d.Sprites[outer]
		sprite.ApplyScale(scale)
		sprite.Draw(dst, op)
	}
	{
		scale := 1.0
		if d.Staged.Tick > 0 {
			scale = float64(d.Staged.HitTick) / float64(d.Staged.Tick)
		}
		sprite := d.Sprites[inner]
		sprite.ApplyScale(scale)
		sprite.Draw(dst, draws.Op{})
	}
}

var (
	DotColorReady = color.NRGBA{255, 255, 255, 255} // White.
	DotColorHit   = color.NRGBA{255, 255, 0, 0}     // Transparent.
	DotColorMiss  = color.NRGBA{255, 0, 0, 255}     // Red.
)

type RollDrawer struct {
	Time        int64
	Rolls       []*Note
	Dots        []*Dot
	HeadSprites [2]draws.Sprite
	TailSprites [2]draws.Sprite
	BodySprites [2]draws.Sprite
	DotSprite   draws.Sprite
}

func (d *RollDrawer) Update(time int64) {
	d.Time = time
}
func (d RollDrawer) Draw(dst draws.Image) {
	max := len(d.Rolls) - 1
	for i := range d.Rolls {
		head := d.Rolls[max-i]
		if head.Position(d.Time) > maxPosition {
			continue
		}
		tail := *head
		tail.Time += head.Duration
		if tail.Position(d.Time) < minPosition {
			continue
		}
		op := draws.Op{}
		op.ColorM.ScaleWithColor(ColorYellow)
		{
			sprite := d.BodySprites[head.Size]
			length := tail.Position(d.Time) - head.Position(d.Time)
			sprite.SetSize(length, sprite.H())
			sprite.Move(head.Position(d.Time), 0)
			sprite.Draw(dst, op)
		}
		{
			sprite := d.HeadSprites[head.Size]
			sprite.Move(head.Position(d.Time), 0)
			sprite.Draw(dst, op)
		}
		{
			sprite := d.TailSprites[tail.Size]
			sprite.Move(tail.Position(d.Time), 0)
			sprite.Draw(dst, op)
		}
	}
	max = len(d.Dots) - 1
	for i := range d.Dots {
		dot := d.Dots[max-i]
		pos := dot.Position(d.Time)
		if pos > maxPosition || pos < minPosition {
			continue
		}
		op := draws.Op{}
		switch dot.Marked {
		case DotReady:
			op.ColorM.ScaleWithColor(DotColorReady)
		case DotHit:
			op.ColorM.ScaleWithColor(DotColorHit)
		case DotMiss:
			op.ColorM.ScaleWithColor(DotColorMiss)
			op.GeoM.Scale(1.5, 1.5)
		}
		d.DotSprite.Move(dot.Position(d.Time), 0)
		d.DotSprite.Draw(dst, op)
	}
}

type NoteDrawer struct {
	draws.Timer
	Time           int64
	Notes          []*Note
	Rolls          []*Note
	Shakes         []*Note
	NoteSprites    [2][4]draws.Sprite
	OverlaySprites [2]draws.Animation
}

func (d *NoteDrawer) Update(time int64, bpm float64) {
	d.Ticker()
	d.Period = int(2 * 60000 / ScaledBPM(bpm))
	d.Time = time
}

func (d NoteDrawer) Draw(dst draws.Image) {
	const (
		modeShake = iota
		modeRoll
		modeNote
	)
	for mode, notes := range [][]*Note{d.Shakes, d.Rolls, d.Notes} {
		max := len(notes) - 1
		for i := range notes {
			n := notes[max-i]
			pos := n.Position(d.Time)
			if pos > maxPosition || pos < minPosition {
				continue
			}
			note := d.NoteSprites[n.Size][n.Color]
			op := draws.Op{}
			switch mode {
			case modeShake:
				if n.Time < d.Time {
					op.ColorM.Scale(1, 1, 1, 0)
				}
			case modeRoll:
				alpha := pos / 400
				if alpha > 1 {
					alpha = 1
				}
				if alpha < 0 {
					alpha = 0
				}
				op.ColorM.Scale(1, 1, 1, alpha)
			case modeNote:
				if n.Marked {
					op.ColorM.Scale(1, 1, 1, 0)
				}
			}
			note.Move(pos, 0)
			note.Draw(dst, op)
			overlay := d.Frame(d.OverlaySprites[n.Size])
			overlay.Move(pos, 0)
			overlay.Draw(dst, op)
		}
	}
}

type KeyDrawer struct {
	MaxCountdown int
	Field        draws.Sprite
	Keys         [4]draws.Sprite
	countdowns   [4]int
	lastPressed  []bool
	pressed      []bool
}

func (d *KeyDrawer) Update(lastPressed, pressed []bool) {
	d.lastPressed = lastPressed
	d.pressed = pressed
	for k, countdown := range d.countdowns {
		if countdown > 0 {
			d.countdowns[k]--
		}
	}
	for k, now := range d.pressed {
		last := d.lastPressed[k]
		if !last && now {
			d.countdowns[k] = d.MaxCountdown
		}
	}
}
func (d KeyDrawer) Draw(dst draws.Image) {
	d.Field.Draw(dst, draws.Op{})
	for k, countdown := range d.countdowns {
		if countdown > 0 {
			d.Keys[k].Draw(dst, draws.Op{})
		}
	}
}

type DancerDrawer struct {
	draws.Timer
	Time        int64
	Sprites     [4]draws.Animation
	Mode        int
	ModeEndTime int64 // It extends when notes are continuously missed.
}

func (d *DancerDrawer) Update(time int64, bpm float64, combo int, miss, hit, high bool) {
	d.Ticker()
	d.Time = time
	period := 4 * 60000 / ScaledBPM(bpm)
	d.Period = int(period) // It should be updated even in constant mode.

	mode := d.Mode
	switch {
	case miss:
		mode = DancerNo
		d.ModeEndTime = time + int64(4*period)
	case combo >= 50 && combo%50 == 0:
		mode = DancerYes
		d.ModeEndTime = time + int64(period)
	case d.Time >= d.ModeEndTime, d.Mode == DancerNo && hit:
		if high {
			mode = DancerHigh
		} else {
			mode = DancerIdle
		}
	}
	if d.Mode != mode {
		if time < 0 && combo == 0 && !miss { // No update before start.
			return
		}
		d.Timer.Reset()
		d.Mode = mode
	}
}
func (d DancerDrawer) Draw(dst draws.Image) {
	d.Frame(d.Sprites[d.Mode]).Draw(dst, draws.Op{})
}

type JudgmentDrawer struct {
	draws.Timer
	Sprites     [2][3]draws.Animation
	judgment    gosu.Judgment
	big         bool
	startRadian float64
	radian      float64
}

func (d *JudgmentDrawer) Update(j gosu.Judgment, big bool) {
	d.Ticker()
	if !j.Valid() {
		return
	}
	d.Timer.Reset()
	d.judgment = j
	d.big = big
	if j.Is(Miss) {
		d.startRadian = (5*rand.Float64() - 2.5) / 24
		d.radian = d.startRadian
	}
}

func (d JudgmentDrawer) Draw(dst draws.Image) {
	if d.Done() || d.judgment.Window == 0 {
		return
	}
	sprites := d.Sprites[Regular]
	if d.big {
		sprites = d.Sprites[Big]
	}
	var sprite draws.Sprite
	for i, j := range Judgments {
		if d.judgment.Is(j) {
			sprite = d.Frame(sprites[i])
			break
		}
	}
	op := draws.Op{}
	age := d.Age()
	if bound := 0.25; age < bound {
		sprite.ApplyScale(1.2 - 0.2*d.Progress(0, bound))
		alpha := 0.5 + 0.5*d.Progress(0, bound)
		op.ColorM.Scale(1, 1, 1, alpha)
	}
	if bound := 0.75; age > bound {
		alpha := 1 - d.Progress(bound, 1)
		op.ColorM.Scale(1, 1, 1, alpha)
	}
	if d.judgment.Is(Miss) {
		if bound := 0.5; age >= bound {
			scale := 1 + 0.6*d.Progress(bound, 1)
			d.radian = d.startRadian * scale
		}
		sw, sh := sprite.SrcSize().XY()
		op.GeoM.Translate(-sw/2, -sh/2)
		op.GeoM.Rotate(d.radian)
		op.GeoM.Translate(sw/2, sh/2)
	}
	sprite.Draw(dst, op)
}
