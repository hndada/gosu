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

// Let's draw as same as possible of osu! does.
// Todo: should BodyDrawer's Notes also be reversed at Draw()?
type BodyDrawer struct {
	Sprites [2]draws.Sprite
	Time    int64
	Notes   []*Note
}

func (d *BodyDrawer) Update(time int64) {
	d.Time = time
}

func (d BodyDrawer) Draw(screen *ebiten.Image) {
	for _, tail := range d.Notes {
		if tail.Type != Tail {
			continue
		}
		if tail.Position(d.Time) < minPosition-bigNoteHeight {
			continue
		}
		head := tail.Prev
		if head.Position(d.Time) > maxPosition+bigNoteHeight {
			continue
		}
		body := d.Sprites[tail.Size]
		length := tail.Position(d.Time) - head.Position(d.Time)
		ratio := length / body.W()
		body.SetScaleXY(ratio, 1, ebiten.FilterLinear)
		body.Move(head.Position(d.Time), 0)

		// op := &ebiten.DrawImageOptions{}
		// if tail.Marked {
		// 	op.ColorM.ChangeHSV(0, 0.3, 0.3)
		// }
		// body.Draw(screen, op)
		body.Draw(screen, nil)
	}
}

type DotDrawer struct {
	Sprite draws.Sprite
	Time   int64
	Dots   []*Dot
	Staged *Dot
}

func (d *DotDrawer) Update(time int64, staged *Dot) {
	d.Time = time
	d.Staged = staged
}
func (d DotDrawer) Draw(screen *ebiten.Image) {
	for _, dot := range d.Dots {
		if dot.Time-Showtime > d.Time {
			continue
		}
		pos := dot.Position(d.Time)
		if pos > maxPosition+100 ||
			pos < minPosition-100 {
			continue
		}
		sprite := d.Sprite
		if dot.Marked {
			sprite.SetColor(DotColorHit)
		} else if d.Staged.Time > dot.Time {
			sprite.SetColor(DotColorMiss)
		} else {
			sprite.SetColor(DotColorReady)
		}
		sprite.Move(dot.Position(d.Time), 0)
		sprite.Draw(screen, nil)
	}
}

type NoteDarwer struct {
	// RedSprites     [2]draws.Sprite
	// BlueSprites    [2]draws.Sprite
	NoteSprites     [2][2]draws.Sprite
	HeadSprites     [2]draws.Sprite
	TailSprites     [2]draws.Sprite
	OverlaySprites  [2][2]draws.Sprite // 2 Overlays.
	ShakeNoteSprite draws.Sprite
	// Overlay indicates which overlay goes drawn.
	// Draw first overlay at even beat, second at odd beat.
	Overlay int
	Time    int64
	Notes   []*Note
}

func (d *NoteDarwer) Update(time int64) {
	d.Time = time
}

// Todo: should Shake note be fade-in in specific time?
func (d NoteDarwer) Draw(screen *ebiten.Image) {
	max := len(d.Notes) - 1
	for i := range d.Notes {
		n := d.Notes[max-i]
		pos := n.Position(d.Time)
		if pos > maxPosition+bigNoteHeight ||
			pos < minPosition-bigNoteHeight {
			continue
		}
		var note draws.Sprite
		switch n.Type {
		case Normal:
			note = d.NoteSprites[n.Size][n.Color]
		case Head:
			note = d.HeadSprites[n.Size]
		case Tail:
			note = d.TailSprites[n.Size]
		case Shake:
			note = d.ShakeNoteSprite
		}
		op := &ebiten.DrawImageOptions{}
		if n.Type == Normal && n.Marked {
			op.ColorM.ChangeHSV(0, 1, 0)
		}
		note.Move(pos, 0)
		note.Draw(screen, nil)
		if n.Type == Tail || n.Type == Shake {
			continue
		}
		overlay := d.OverlaySprites[n.Size][d.Overlay]
		overlay.Move(pos, 0)
		// fmt.Printf("note x, y: %.f %.f\noverlay x, y: %.f %.f\n", note.X(), note.Y(),
		// 	overlay.X(), overlay.Y())
		overlay.Draw(screen, nil)
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

type ShakeDrawer struct{}

// case Leftward, Rightward:
// 	if d.direction == Rightward {
// 		ratio *= -1
// 	}
// 	srcRect := image.Rect(int(srcStart), 0, int(srcEnd), int(body.H()))
// 	sprite := body.SubSprite(srcRect)
// 	op.GeoM.Scale(ratio, 1)
// 	op.GeoM.Translate(srcStart, 0)
// 	sprite.Draw(screen, op)
