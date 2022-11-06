package drum

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	draws "github.com/hndada/gosu/draws2"
)

type StageDrawer struct {
	draws.Timer2
	Highlight    bool
	FieldSprites [2]draws.Sprite
	HintSprites  [2]draws.Sprite
}

func (d *StageDrawer) Update(highlight bool) {
	d.Ticker()
	if d.Highlight != highlight {
		d.Tick = 0
		d.Highlight = highlight
	}
}

func (d StageDrawer) Draw(screen *ebiten.Image) {
	const (
		idle = iota
		high
	)
	op := ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, FieldDarkness)
	d.FieldSprites[idle].Draw(screen, op)
	d.HintSprites[idle].Draw(screen, ebiten.DrawImageOptions{})
	if d.Highlight || d.Tick < d.MaxTick {
		var op1, op2 ebiten.DrawImageOptions
		if d.Highlight {
			op1.ColorM.Scale(1, 1, 1, FieldDarkness*d.Age())
			op2.ColorM.Scale(1, 1, 1, FieldDarkness*d.Age())
		} else {
			op1.ColorM.Scale(1, 1, 1, FieldDarkness*(1-d.Age()))
			op2.ColorM.Scale(1, 1, 1, FieldDarkness*(1-d.Age()))
		}
		d.FieldSprites[high].Draw(screen, op1)
		d.HintSprites[high].Draw(screen, op2)
	}
	// if d.Highlight || d.Tick < d.MaxTick {
	// 	{
	// 		op := ebiten.DrawImageOptions{}
	// 		op.ColorM.Scale(1, 1, 1, FieldDarkness*d.Age())
	// 		d.FieldSprites[high].Draw(screen, op)
	// 	}
	// 	{
	// 		op := ebiten.DrawImageOptions{}
	// 		op.ColorM.Scale(1, 1, 1, d.Age())
	// 		d.HintSprites[high].Draw(screen, op)
	// 	}
	// }
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
func (d BarDrawer) Draw(screen *ebiten.Image) {
	for _, b := range d.Bars {
		pos := b.Speed * float64(b.Time-d.Time)
		if pos <= maxPosition && pos >= minPosition {
			sprite := d.Sprite
			sprite.Move(pos, 0)
			sprite.Draw(screen, ebiten.DrawImageOptions{})
		}
	}
}

type ShakeDrawer struct {
	Time         int64
	Staged       *Note
	BorderSprite draws.Sprite
	ShakeSprite  draws.Sprite
}

func (d *ShakeDrawer) Update(time int64, staged *Note) {
	d.Time = time
	d.Staged = staged
}
func (d ShakeDrawer) Draw(screen *ebiten.Image) {
	if d.Staged == nil {
		return
	}
	if d.Staged.Time > d.Time {
		return
	}
	borderScale := 0.25 + 0.75*float64(d.Time-d.Staged.Time)/80
	if borderScale > 1 {
		borderScale = 1
	}
	d.BorderSprite.ApplyScale(borderScale)
	// d.BorderSprite.SetScale(draws.Scalar(borderScale))
	d.BorderSprite.Draw(screen, ebiten.DrawImageOptions{})

	shakeScale := float64(d.Staged.HitTick) / float64(d.Staged.Tick)
	d.ShakeSprite.ApplyScale(shakeScale)
	// d.ShakeSprite.SetScale(draws.Scalar(shakeScale))
	d.ShakeSprite.Draw(screen, ebiten.DrawImageOptions{})
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
	BodySprites [2]draws.Sprite
	TailSprites [2]draws.Sprite
	DotSprite   draws.Sprite
}

func (d *RollDrawer) Update(time int64) {
	d.Time = time
}
func (d RollDrawer) Draw(screen *ebiten.Image) {
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
		op := ebiten.DrawImageOptions{}
		op.ColorM.ScaleWithColor(ColorYellow)

		bodySprite := d.BodySprites[head.Size]
		length := tail.Position(d.Time) - head.Position(d.Time)
		bodySprite.SetSize(length, bodySprite.H())
		bodySprite.Move(head.Position(d.Time), 0)
		bodySprite.Draw(screen, op)

		headSprite := d.HeadSprites[head.Size]
		headSprite.Move(head.Position(d.Time), 0)
		headSprite.Draw(screen, op)

		tailSprite := d.TailSprites[tail.Size]
		tailSprite.Move(tail.Position(d.Time), 0)
		tailSprite.Draw(screen, op)
	}
	max = len(d.Dots) - 1
	for i := range d.Dots {
		dot := d.Dots[max-i]
		pos := dot.Position(d.Time)
		if pos > maxPosition || pos < minPosition {
			continue
		}
		sprite := d.DotSprite
		op := ebiten.DrawImageOptions{}
		switch dot.Marked {
		case DotReady:
			op.ColorM.ScaleWithColor(DotColorReady)
		case DotHit:
			op.ColorM.ScaleWithColor(DotColorHit)
		case DotMiss:
			op.ColorM.ScaleWithColor(DotColorMiss)
			op.GeoM.Scale(1.5, 1.5)
		}
		sprite.Move(dot.Position(d.Time), 0)
		sprite.Draw(screen, op)
	}
}

type NoteDrawer struct {
	draws.Timer2
	Time           int64
	Notes          []*Note
	Rolls          []*Note
	Shakes         []*Note
	NoteSprites    [2][4]draws.Sprite
	OverlaySprites [2][]draws.Sprite
}

func (d *NoteDrawer) Update(time int64, bpm float64) {
	d.Ticker()
	d.Time = time
	d.Period = int(2 * 60000 / ScaledBPM(bpm))
	// duration := 2 * 60000 / ScaledBPM(bpm)
	// for i := range d.OverlaySprites {
	// 	d.OverlaySprites[i].Update(time, int64(duration), false)
	// }
}

func (d NoteDrawer) Draw(screen *ebiten.Image) {
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
			op := ebiten.DrawImageOptions{}
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
			note.Draw(screen, op)
			// if mode == modeShake {
			// 	continue
			// }
			overlay := d.Frame(d.OverlaySprites[n.Size])
			overlay.Move(pos, 0)
			overlay.Draw(screen, op)
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
func (d KeyDrawer) Draw(screen *ebiten.Image) {
	d.Field.Draw(screen, ebiten.DrawImageOptions{})
	for k, countdown := range d.countdowns {
		if countdown > 0 {
			d.Keys[k].Draw(screen, ebiten.DrawImageOptions{})
		}
	}
}

type DancerDrawer struct {
	draws.Timer2
	Time        int64
	Sprites     [4][]draws.Sprite
	Mode        int
	ModeEndTime int64 // It extends when notes are continuously missed.
	// ModeMaxTick int // For Yes and No.
	// Duration int64
}

func (d *DancerDrawer) Update(time int64, bpm float64, combo int, miss, hit, high bool) {
	d.Ticker()
	d.Time = time
	// maxTick := 0
	period := 4 * 60000 / ScaledBPM(bpm)
	d.Period = int(period) // It should be updated even in constant mode.

	mode := d.Mode
	// var reset bool
	switch {
	case miss:
		mode = DancerNo
		// maxTick = int(4 * period)
		d.ModeEndTime = time + int64(4*period)
	case combo >= 50 && combo%50 <= 1:
		mode = DancerYes
		// maxTick = int(period)
		d.ModeEndTime = time + int64(period)
	// case hit || d.IsModeFinished():
	case d.Time >= d.ModeEndTime, d.Mode == DancerNo && hit: //, d.Mode == DancerIdle, d.Mode == DancerHigh:
		if high {
			mode = DancerHigh
		} else {
			mode = DancerIdle
		}
	}
	if d.Mode != mode {
		d.Tick = 0
		d.Mode = mode
	}
	// if d.Mode != mode { // && d.AnimationFinish()
	// 	if d.Mode == DancerYes && mode != DancerNo && !d.IsModeFinished() {
	// 	} else {
	// 		d.Tick = 0
	// 		// d.Timer2 = draws.NewTimer2(maxTick, int(period))
	// 		d.Mode = mode
	// 		// reset = true
	// 	}
	// }
	// d.Sprites[d.Mode].Update(time, int64(duration), reset)
}

// AnimationFinish infers if Dancer is ready to change its mode.
//
//	func (d DancerDrawer) IsModeFinished() bool {
//		if d.Mode == DancerIdle || d.Mode == DancerHigh {
//			return true
//		}
//		return d.Time >= d.ModeEndTime
//	}
func (d DancerDrawer) Draw(screen *ebiten.Image) {
	d.Frame(d.Sprites[d.Mode]).Draw(screen, ebiten.DrawImageOptions{})
}

type JudgmentDrawer struct {
	draws.Timer
	Sprites     [2][3]draws.Sprite
	judgment    gosu.Judgment
	big         bool
	startRadian float64
	radian      float64
}

func (d *JudgmentDrawer) Update(j gosu.Judgment, big bool) {
	if d.Countdown <= 0 {
		d.judgment = gosu.Judgment{}
		d.big = false
	} else {
		d.Countdown--
	}
	if j.Valid() {
		d.Countdown = d.MaxCountdown
		d.judgment = j
		d.big = big
		if j.Is(Miss) {
			d.startRadian = (5*rand.Float64() - 2.5) / 24
			d.radian = d.startRadian
		}
	}
}

func (d JudgmentDrawer) Draw(screen *ebiten.Image) {
	if d.Countdown <= 0 || d.judgment.Window == 0 {
		return
	}
	sprites := d.Sprites[0]
	if d.big {
		sprites = d.Sprites[1]
	}
	var sprite draws.Sprite
	for i, j := range Judgments {
		if d.judgment.Is(j) {
			sprite = sprites[i]
			break
		}
	}
	op := ebiten.DrawImageOptions{}
	if d.judgment.Is(Miss) {
		age := d.Age()
		if age < 0.25 {
			sprite.ApplyScale(1 + 0.2*(0.25-age)/0.25)
			alpha := 1 - 0.5*(0.2-age)/0.2
			op.ColorM.Scale(1, 1, 1, alpha)
		}
		if age >= 0.5 {
			scale := 1 + 0.6*(age-0.5)/0.5
			d.radian = d.startRadian * scale
		}
		op.GeoM.Translate(sprite.SrcSize().Div(draws.Scalar(-2)).XY())
		// op.GeoM.Translate(-float64(sw)/2, -float64(sh)/2)
		op.GeoM.Rotate(d.radian)
		op.GeoM.Translate(sprite.SrcSize().Div(draws.Scalar(2)).XY())
		// op.GeoM.Translate(float64(sw)/2, float64(sh)/2)
	} else {
		age := d.Age()
		if age < 0.25 {
			sprite.ApplyScale(1 + 0.2*(0.25-age)/0.25)
		}
		if age > 0.75 {
			alpha := 1 - (age-0.75)/0.25
			op.ColorM.Scale(1, 1, 1, alpha)
		}
	}
	sprite.Draw(screen, op)
}
