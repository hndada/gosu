package drum

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
)

type StageDrawer struct {
	Hightlight  bool
	FieldSprite draws.Sprite
	HintSprites [2]draws.Sprite
}

func (d *StageDrawer) Update(highlight bool) {
	d.Hightlight = highlight
}

func (d StageDrawer) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, FieldDarkness)
	d.FieldSprite.Draw(screen, op)
	if d.Hightlight {
		d.HintSprites[1].Draw(screen, nil)
	} else {
		d.HintSprites[0].Draw(screen, nil)
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
func (d BarDrawer) Draw(screen *ebiten.Image) {
	for _, b := range d.Bars {
		pos := b.Speed * float64(b.Time-d.Time)
		if pos <= maxPosition && pos >= minPosition {
			sprite := d.Sprite
			sprite.Move(pos, 0)
			sprite.Draw(screen, nil)
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
	d.BorderSprite.SetScale(borderScale)
	d.BorderSprite.Draw(screen, nil)

	shakeScale := float64(d.Staged.HitTick) / float64(d.Staged.Tick)
	d.ShakeSprite.SetScale(shakeScale)
	d.ShakeSprite.Draw(screen, nil)
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
		op := &ebiten.DrawImageOptions{}
		op.ColorM.ScaleWithColor(ColorYellow)
		length := tail.Position(d.Time) - head.Position(d.Time)

		bodySprite := d.BodySprites[head.Size]
		ratio := length / bodySprite.W()
		bodySprite.SetScaleXY(ratio, 1, ebiten.FilterNearest)
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
		op := &ebiten.DrawImageOptions{}
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

type NoteDarwer struct {
	Time           int64
	Notes          []*Note
	Rolls          []*Note
	Shakes         []*Note
	NoteSprites    [2][4]draws.Sprite
	OverlayDrawers [2]draws.AnimationDrawer
}

func (d *NoteDarwer) Update(time int64, bpm float64) {
	d.Time = time
	duration := 2 * 60000 / ScaledBPM(bpm)
	for i := range d.OverlayDrawers {
		d.OverlayDrawers[i].Update(time, int64(duration), false)
	}
}

func (d NoteDarwer) Draw(screen *ebiten.Image) {
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
			op := &ebiten.DrawImageOptions{}
			switch mode {
			case modeShake:
				if n.Time < d.Time {
					op.ColorM.Scale(1, 1, 1, 0)
				}
			case modeRoll:
				rate := pos / 400
				if rate > 1 {
					rate = 1
				}
				if rate < 0 {
					rate = 0
				}
				op.ColorM.Scale(1, 1, 1, rate)
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
			d.OverlayDrawers[n.Size].Draw(screen, op, pos, 0)
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
	d.Field.Draw(screen, nil)
	for k, countdown := range d.countdowns {
		if countdown > 0 {
			d.Keys[k].Draw(screen, nil)
		}
	}
}

type DancerDrawer struct {
	AnimationDrawers [4]draws.AnimationDrawer
	AnimationEndTime int64 // For Yes and No.
	Mode             int
}

func (d *DancerDrawer) Update(
	time int64, bpm float64, combo int, miss, hit, highlight bool) {
	duration := 4 * 60000 / ScaledBPM(bpm)
	mode := d.Mode
	var reset bool
	switch {
	case miss:
		mode = DancerNo
		d.AnimationEndTime = time + int64(4*duration)
	case combo >= 50 && combo%50 < 2:
		mode = DancerYes
		d.AnimationEndTime = time + int64(duration)
	case hit || d.AnimationFinish():
		if highlight {
			mode = DancerHigh
		} else {
			mode = DancerIdle
		}
	}
	if d.Mode != mode && d.AnimationFinish() {
		fmt.Println(d.Mode, mode)
		d.Mode = mode
		reset = true
	}
	d.AnimationDrawers[d.Mode].Update(time, int64(duration), reset)
}

// AnimationFinish infers if Dancer is ready to change its mode.
func (d DancerDrawer) AnimationFinish() bool {
	if d.Mode == DancerIdle || d.Mode == DancerHigh {
		return true
	}
	return d.AnimationDrawers[d.Mode].Time >= d.AnimationEndTime
}
func (d DancerDrawer) Draw(screen *ebiten.Image) {
	d.AnimationDrawers[d.Mode].Draw(screen, nil, 0, 0)
}

type JudgmentDrawer struct {
	draws.BaseDrawer
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
		if j.Is(Miss) {
			rate := 1.0
			if age := d.Age(); age >= 0.25 {
				rate = (1 + 0.6*(age-0.25)/0.75)
			}
			d.radian = d.startRadian * rate
		}
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
	op := &ebiten.DrawImageOptions{}
	sw, sh := sprite.SrcSize()
	if d.judgment.Is(Miss) {
		op.GeoM.Translate(-float64(sw)/2, -float64(sh)/2)
		op.GeoM.Rotate(d.radian)
		op.GeoM.Translate(float64(sw)/2, float64(sh)/2)
	}
	ratio := 1.0
	if age := d.Age(); age < 0.15 {
		ratio = 1 + (0.15 - age)
	}
	sprite.SetScale(ratio)
	sprite.Draw(screen, op)
}
