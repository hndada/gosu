package drum

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
)

type StageDrawer struct {
	Field      draws.Sprite
	Hints      [2]draws.Sprite
	Hightlight bool
}

func (d *StageDrawer) Update(highlight bool) {
	d.Hightlight = highlight
}

// Todo: might add some effect on StageDrawer
func (d StageDrawer) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, FieldDarkness)
	d.Field.Draw(screen, op)
	if d.Hightlight {
		d.Hints[1].Draw(screen, nil)
	} else {
		d.Hints[0].Draw(screen, nil)
	}
}

// BarDrawer in Drum mode is floating-type lane drawer.
// Todo: set draw order for performance?
type BarDrawer struct {
	Sprite draws.Sprite
	Time   int64
	Bars   []*Bar
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

var (
	DotColorReady = color.NRGBA{255, 255, 255, 255} // White.
	// DotColorHit   = color.NRGBA{255, 255, 0, 255}   // Yellow.
	DotColorHit = color.NRGBA{255, 255, 0, 0} // Transparent.
	// DotColorMiss = color.NRGBA{0, 32, 96, 255} // Navy.
	DotColorMiss = color.NRGBA{255, 0, 0, 255} // Red.
)

type ShakeDrawer struct {
	BorderSprite draws.Sprite
	ShakeSprite  draws.Sprite
	Time         int64
	Staged       *Note
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

// Todo: should BodyDrawer's Notes also be reversed at Draw()?
type RollDrawer struct {
	// OverlaySprites [2][2]draws.Sprite
	HeadSprites [2]draws.Sprite
	BodySprites [2]draws.Sprite
	TailSprites [2]draws.Sprite
	DotSprite   draws.Sprite
	Time        int64
	Rolls       []*Note
	Dots        []*Dot
	// StagedDot   *Dot
}

//	func (d *RollDrawer) Update(time int64, stagedDot *Dot) {
//		d.Time = time
//		d.StagedDot = stagedDot
//	}
func (d *RollDrawer) Update(time int64) {
	d.Time = time
}
func (d RollDrawer) Draw(screen *ebiten.Image) {
	max := len(d.Rolls) - 1
	for i := range d.Rolls {
		head := d.Rolls[max-i]
		if head.Position(d.Time) > maxPosition+bigNoteHeight {
			continue
		}
		tail := *head
		tail.Time += head.Duration
		if tail.Position(d.Time) < minPosition-bigNoteHeight {
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
		// if dot.RevealTime > d.Time {
		// 	continue
		// }
		pos := dot.Position(d.Time)
		if pos > maxPosition+100 || pos < minPosition-100 {
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

// type DotDrawer struct {
// 	Sprite draws.Sprite
// 	Time   int64
// 	Dots   []*Dot
// 	Staged *Dot
// }

//	func (d *DotDrawer) Update(time int64, staged *Dot) {
//		d.Time = time
//		d.Staged = staged
//	}
// func (d DotDrawer) Draw(screen *ebiten.Image) {
// 	for _, dot := range d.Dots {
// 		if dot.RevealTime > d.Time {
// 			continue
// 		}
// 		pos := dot.Position(d.Time)
// 		if pos > maxPosition+100 ||
// 			pos < minPosition-100 {
// 			continue
// 		}
// 		sprite := d.Sprite
// 		op := &ebiten.DrawImageOptions{}
// 		if dot.Marked {
// 			op.ColorM.ScaleWithColor(DotColorHit)
// 			// sprite.SetColor(DotColorHit)
// 		} else if d.Staged.Time > dot.Time {
// 			op.ColorM.ScaleWithColor(DotColorMiss)
// 			// sprite.SetColor(DotColorMiss)
// 		} else {
// 			op.ColorM.ScaleWithColor(DotColorReady)
// 			// sprite.SetColor(DotColorReady)
// 		}
// 		sprite.Move(dot.Position(d.Time), 0)
// 		sprite.Draw(screen, op)
// 	}
// }

type NoteDarwer struct {
	NoteSprites    [2][4]draws.Sprite
	OverlaySprites [2][2]draws.Sprite // 2 Overlays.
	// ShakeNoteSprite draws.Sprite
	// Overlay indicates which overlay goes drawn.
	// Draw first overlay at even beat, second at odd beat.
	Time                  int64
	OverlayDuration       int64
	Overlay               int // It shows which overlay sprite goes drawn.
	LastOverlayChangeTime int64
	Notes                 []*Note
	Rolls                 []*Note
	Shakes                []*Note
}

func (d *NoteDarwer) Update(time int64, bpm float64) {
	d.Time = time
	d.OverlayDuration = int64(60000 / ScaledBPM(bpm))
	if d.Time-d.LastOverlayChangeTime >= d.OverlayDuration {
		d.Overlay = (d.Overlay + 1) % 2
		d.LastOverlayChangeTime = d.Time
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
			if pos > maxPosition+bigNoteHeight ||
				pos < minPosition-bigNoteHeight {
				continue
			}
			// var note draws.Sprite
			// if mode == modeShake {
			// 	note = d.ShakeNoteSprite
			// } else {
			// 	note = d.NoteSprites[n.Size][n.Color-1]
			// }
			note := d.NoteSprites[n.Size][n.Color]
			op := &ebiten.DrawImageOptions{}
			switch mode {
			case modeShake:
				if n.Time < d.Time {
					op.ColorM.Scale(1, 1, 1, 0)
					// op.ColorM.ChangeHSV(0, 1, 0)
				}
			case modeRoll:
				rate := (pos - 0) / 400
				if rate > 1 {
					rate = 1
				}
				if rate < 0 {
					rate = 0
				}
				// op.ColorM.ChangeHSV(0, 1, value)
				op.ColorM.Scale(1, 1, 1, rate)
			case modeNote:
				if n.Marked {
					op.ColorM.Scale(1, 1, 1, 0)
				}
			}
			note.Move(pos, 0)
			note.Draw(screen, op)
			if mode == modeShake {
				continue
			}
			overlay := d.OverlaySprites[n.Size][d.Overlay]
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
	d.Field.Draw(screen, nil)
	for k, countdown := range d.countdowns {
		if countdown > 0 {
			d.Keys[k].Draw(screen, nil)
		}
	}
}

type DancerDrawer struct {
	Time             int64
	Duration         float64 // Yes's duration is 4 times of other modes'.
	Mode             int
	Frame            int
	LastFrameTimes   [4]int64
	MissDanceEndTime int64
	Sprites          [4][]draws.Sprite
	// MissCycleCountdown int
	// MissTime       int64
	// Hit                bool
	// Yes                bool // Whether dancer dances Yes based on combo achieving.
	// Highlight          bool
}

func (d *DancerDrawer) Update(time int64, bpm float64, miss, hit bool, combo int, highlight bool) {
	d.Time = time
	d.Duration = 4 * 60000 / ScaledBPM(bpm)
	var modeChange bool
	// if hit && d.MissCycleCountdown > 0 {
	// 	d.MissCycleCountdown = 0
	// 	d.Mode = DancerIdle
	// 	modeChange = true
	// }
	// dancerNoCycle := float64(d.Time-d.LastFrameTimes[DancerNo]) / d.Duration
	switch {
	case miss:
		d.MissDanceEndTime = d.Time + int64(d.Duration*4)
		// d.MissCycleCountdown = 4
		if d.Mode != DancerNo {
			d.Mode = DancerNo
			modeChange = true
		}
	// case !d.Yes && combo >= 50 && combo%50 < 5:
	case combo >= 50 && combo%50 < 5:
		if d.Mode != DancerYes {
			d.Mode = DancerYes
			// d.Duration *= 4
			modeChange = true
		}
		// d.Yes = true
	case hit || d.Time >= d.MissDanceEndTime:
		// case hit || (d.Mode == DancerNo && d.Time >= d.MissDanceEndTime):
		if highlight {
			if d.Mode != DancerHigh {
				d.Mode = DancerHigh
				modeChange = true
			}
		} else {
			if d.Mode != DancerIdle {
				d.Mode = DancerIdle
				modeChange = true
			}
		}
		// case hit || d.MissCycleCountdown == 0:
		// 	d.MissCycleCountdown = 0
		// 	if highlight {
		// 		d.Mode = DancerHigh
		// 	} else {
		// 		d.Mode = DancerIdle
		// 	}
		// 	modeChange = true
	}
	if modeChange {
		d.LastFrameTimes[d.Mode] = time
	}
	td := float64(d.Time - d.LastFrameTimes[d.Mode])
	// q := math.Floor(td / d.Duration)
	// td -= q * d.Duration
	// d.Frame = int(td) * len(d.Sprites[d.Mode])
	rate := math.Remainder(td, d.Duration) / d.Duration
	if rate < 0 {
		rate += 1
	}
	frames := float64(len(d.Sprites[d.Mode]))
	d.Frame = int(rate * frames)
}
func (d DancerDrawer) Draw(screen *ebiten.Image) {
	d.Sprites[d.Mode][d.Frame].Draw(screen, nil)
}

type JudgmentDrawer struct {
	draws.BaseDrawer
	Sprites     [2][3]draws.Sprite
	judgment    gosu.Judgment
	big         bool
	startRadian float64
	radian      float64
	// reverse     bool
}

func (d *JudgmentDrawer) Update(j gosu.Judgment, big bool) {
	if d.Countdown <= 0 {
		d.judgment = gosu.Judgment{}
		d.big = false
	} else {
		d.Countdown--
		age := d.Age()
		// d.startRadian = 5 * rand.Float64() / 24
		rate := 1.0
		if age >= 0.25 {
			rate = (1 + 0.6*(age-0.25)/0.75)
		}
		d.radian = d.startRadian * rate
	}
	if j.Window != 0 {
		d.judgment = j
		d.big = big
		d.Countdown = d.MaxCountdown
		d.startRadian = (5*rand.Float64() - 2.5) / 24
		d.radian = d.startRadian
	}
}

// dx := hw*math.Cos(d.radian) - hh*math.Sin(d.radian)
// dy := hw*math.Sin(d.radian) + hh*math.Cos(d.radian)
// op.GeoM.Translate(dx, dy)
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
		if j.Window == d.judgment.Window {
			sprite = sprites[i]
			break
		}
	}
	op := &ebiten.DrawImageOptions{}
	sw, sh := sprite.SrcSize()
	if d.judgment.Window == Miss.Window {
		op.GeoM.Translate(-float64(sw)/2, -float64(sh)/2)
		op.GeoM.Rotate(d.radian)
		op.GeoM.Translate(float64(sw)/2, float64(sh)/2)
	}
	ratio := 1.0
	age := d.Age()
	switch {
	case age < 0.15:
		ratio = 1 + (0.15 - age)
		// case age >= 0.15 && age < 0.2:
		// 	ratio = 1.15 * (1.2 - age)
		// case age > 0.9:
		// 	ratio = 1 - 1.15*(age-0.9)
	}
	sprite.SetScale(ratio)
	sprite.Draw(screen, op)
}

// case Leftward, Rightward:
// 	if d.direction == Rightward {
// 		ratio *= -1
// 	}
// 	srcRect := image.Rect(int(srcStart), 0, int(srcEnd), int(body.H()))
// 	sprite := body.SubSprite(srcRect)
// 	op.GeoM.Scale(ratio, 1)
// 	op.GeoM.Translate(srcStart, 0)
// 	sprite.Draw(screen, op)
