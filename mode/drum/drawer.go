package drum

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
)

type StageDrawer struct {
	Field draws.Sprite
	Hint  draws.Sprite
}

// Todo: might add some effect on StageDrawer
func (d StageDrawer) Draw(screen *ebiten.Image) {
	d.Field.Draw(screen, nil)
	d.Hint.Draw(screen, nil)
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
	DotColorHit   = color.NRGBA{255, 255, 0, 255}   // Yellow.
	DotColorMiss  = color.NRGBA{0, 32, 96, 255}     // Navy.
)

type ShakeDrawer struct {
	BorderSprite draws.Sprite
	Sprite       draws.Sprite
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
	d.BorderSprite.Draw(screen, nil)
	op := &ebiten.DrawImageOptions{}
	scale := float64(d.Staged.HitTick) / float64(d.Staged.Tick)
	op.GeoM.Scale(scale, scale)
	d.Sprite.Draw(screen, op)
}

// Todo: should BodyDrawer's Notes also be reversed at Draw()?
type RollDrawer struct {
	// HeadSprites    [2]draws.Sprite
	// OverlaySprites [2][2]draws.Sprite
	BodySprites [2]draws.Sprite
	TailSprites [2]draws.Sprite
	DotSprite   draws.Sprite
	Time        int64
	Rolls       []*Note
	Dots        []*Dot
	StagedDot   *Dot
}

func (d *RollDrawer) Update(time int64, stagedDot *Dot) {
	d.Time = time
	d.StagedDot = stagedDot
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
		bodySprite := d.BodySprites[head.Size]
		length := tail.Position(d.Time) - head.Position(d.Time)
		ratio := length / bodySprite.W()
		bodySprite.SetScaleXY(ratio, 1, ebiten.FilterLinear)
		bodySprite.Move(head.Position(d.Time), 0)
		bodySprite.Draw(screen, nil)

		tailSprite := d.TailSprites[tail.Size]
		tailSprite.Move(tail.Position(d.Time), 0)
		tailSprite.Draw(screen, nil)
	}
	max = len(d.Dots) - 1
	for i := range d.Dots {
		dot := d.Dots[max-i]
		if dot.RevealTime > d.Time {
			continue
		}
		pos := dot.Position(d.Time)
		if pos > maxPosition+100 || pos < minPosition-100 {
			continue
		}
		sprite := d.DotSprite
		op := &ebiten.DrawImageOptions{}
		if dot.Marked {
			op.ColorM.ScaleWithColor(DotColorHit)
		} else if d.StagedDot.Time > dot.Time {
			op.ColorM.ScaleWithColor(DotColorMiss)
		} else {
			op.ColorM.ScaleWithColor(DotColorReady)
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
			note := d.NoteSprites[n.Size][n.Color-1]
			op := &ebiten.DrawImageOptions{}
			if mode == modeNote && n.Marked {
				op.ColorM.ChangeHSV(0, 1, 0)
			} else if mode == modeShake && n.Time < d.Time {
				op.ColorM.ChangeHSV(0, 1, 0)
			}
			note.Move(pos, 0)
			note.Draw(screen, op)
			if mode == modeShake {
				continue
			}
			op = &ebiten.DrawImageOptions{}
			if n.Type == Roll {
				value := (pos - 200) / 200
				if value > 1 {
					value = 1
				}
				if value < 0 {
					value = 0
				}
				op.ColorM.ChangeHSV(0, 1, value)
			}
			overlay := d.OverlaySprites[n.Size][d.Overlay]
			overlay.Move(pos, 0)
			overlay.Draw(screen, op)
			// fmt.Printf("note x, y: %.f %.f\noverlay x, y: %.f %.f\n", note.X(), note.Y(),
			// 	overlay.X(), overlay.Y())
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

type JudgmentDrawer struct {
	draws.BaseDrawer
	Sprites  [2][3]draws.Sprite
	judgment gosu.Judgment
	big      bool
}

func (d *JudgmentDrawer) Update(j gosu.Judgment, big bool) {
	if d.Countdown <= 0 {
		d.judgment = gosu.Judgment{}
		d.big = false
	} else {
		d.Countdown--
	}
	if j.Window != 0 {
		d.judgment = j
		d.big = big
		d.Countdown = d.MaxCountdown
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
		if j.Window == d.judgment.Window {
			sprite = sprites[i]
			break
		}
	}
	age := d.Age()
	ratio := 1.0
	switch {
	case age < 0.1:
		ratio = 1.15 * (1 + age)
	case age >= 0.1 && age < 0.2:
		ratio = 1.15 * (1.2 - age)
	case age > 0.9:
		ratio = 1 - 1.15*(age-0.9)
	}
	sprite.SetScale(ratio)
	sprite.Draw(screen, nil)
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
