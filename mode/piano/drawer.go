package piano

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	draws "github.com/hndada/gosu/draws2"
)

type StageDrawer struct {
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite
}

// Todo: might add some effect on StageDrawer
func (d StageDrawer) Draw(screen *ebiten.Image) {
	d.FieldSprite.Draw(screen, ebiten.DrawImageOptions{})
	d.HintSprite.Draw(screen, ebiten.DrawImageOptions{})
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
		d.Farthest.Prev.Position-d.Cursor > maxPosition {
		d.Farthest = d.Farthest.Prev
	}
	// When Farthest is in screen, next note goes fetched if possible.
	for d.Farthest.Next != nil &&
		d.Farthest.Position-d.Cursor <= maxPosition {
		d.Farthest = d.Farthest.Next
	}
	// When Nearest is still in screen due to speed change.
	for d.Nearest.Prev != nil &&
		d.Nearest.Position-d.Cursor > minPosition {
		d.Nearest = d.Nearest.Prev
	}
	// When Nearest's next is still out of screen, next note goes fetched.
	for d.Nearest.Next != nil &&
		d.Nearest.Next.Position-d.Cursor <= minPosition {
		d.Nearest = d.Nearest.Next
	}
}

func (d BarDrawer) Draw(screen *ebiten.Image) {
	if d.Farthest == nil || d.Nearest == nil {
		return
	}
	for b := d.Farthest; b != d.Nearest.Prev; b = b.Prev {
		sprite := d.Sprite
		pos := b.Position - d.Cursor
		sprite.Move(0, -pos)
		sprite.Draw(screen, ebiten.DrawImageOptions{})
	}
}

// Notes are fixed. Lane itself moves, all notes move same amount.
type NoteDrawer struct {
	Cursor   float64
	Farthest *Note
	Nearest  *Note
	Sprites  [4]draws.Sprite
}

// Farthest and Nearest are borders of displaying notes.
// All in-screen notes are confirmed to be drawn when drawing from Farthest to Nearest.
func (d *NoteDrawer) Update(cursor float64) {
	d.Cursor = cursor
	if d.Farthest == nil || d.Nearest == nil {
		return
	}
	// When Farthest's prevs are still out of screen due to speed change.
	for d.Farthest.Prev != nil &&
		d.Farthest.Prev.Position-d.Cursor > maxPosition {
		d.Farthest = d.Farthest.Prev
	}
	// When Farthest is in screen, next note goes fetched if possible.
	for d.Farthest.Next != nil &&
		d.Farthest.Position-d.Cursor <= maxPosition {
		d.Farthest = d.Farthest.Next
	}
	// When Nearest is still in screen due to speed change.
	for d.Nearest.Prev != nil &&
		d.Nearest.Position-d.Cursor > minPosition {
		d.Nearest = d.Nearest.Prev
	}
	// When Nearest's next is still out of screen, next note goes fetched.
	for d.Nearest.Next != nil &&
		d.Nearest.Next.Position-d.Cursor <= minPosition {
		d.Nearest = d.Nearest.Next
	}
}

// Draw from farthest to nearest to make nearer notes priorly exposed.
func (d NoteDrawer) Draw(screen *ebiten.Image) {
	if d.Farthest == nil || d.Nearest == nil {
		return
	}
	for n := d.Farthest; n != nil && n != d.Nearest.Prev; n = n.Prev {
		if n.Type == Tail {
			d.DrawBody(screen, n)
		}
		sprite := d.Sprites[n.Type]
		pos := n.Position - d.Cursor
		sprite.Move(0, -pos)
		op := ebiten.DrawImageOptions{}
		if n.Marked {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		sprite.Draw(screen, op)
	}
}

// DrawBody draws scaled, corresponding sub-image of Body sprite.
func (d NoteDrawer) DrawBody(screen *ebiten.Image, tail *Note) {
	head := tail.Prev
	body := d.Sprites[Body]

	length := tail.Position - head.Position
	length -= -bodyLoss
	body.SetSize(body.W(), length)

	ty := head.Position - d.Cursor
	body.Move(0, -ty)

	op := ebiten.DrawImageOptions{}
	if tail.Marked {
		op.ColorM.ChangeHSV(0, 0.3, 0.3)
	}
	body.Draw(screen, op)
}

// KeyDrawer draws KeyDownSprite at least for 30ms, KeyUpSprite otherwise.
// KeyDrawer uses MinCountdown instead of MaxCountdown.
type KeyDrawer struct {
	MinCountdown   int
	Countdowns     []int
	KeyUpSprites   []draws.Sprite
	KeyDownSprites []draws.Sprite
	lastPressed    []bool
	pressed        []bool
}

func (d *KeyDrawer) Update(lastPressed, pressed []bool) {
	d.lastPressed = lastPressed
	d.pressed = pressed
	for k, countdown := range d.Countdowns {
		if countdown > 0 {
			d.Countdowns[k]--
		}
	}
	for k, now := range d.pressed {
		last := d.lastPressed[k]
		if !last && now {
			d.Countdowns[k] = d.MinCountdown
		}
	}
}
func (d KeyDrawer) Draw(screen *ebiten.Image) {
	for k, p := range d.pressed {
		if p || d.Countdowns[k] > 0 {
			d.KeyDownSprites[k].Draw(screen, ebiten.DrawImageOptions{})
		} else {
			d.KeyUpSprites[k].Draw(screen, ebiten.DrawImageOptions{})
		}
	}
}

type JudgmentDrawer struct {
	draws.Timer
	Sprites  [5][]draws.Sprite
	Judgment gosu.Judgment
}

func NewJudgmentDrawer() (d JudgmentDrawer) {
	const frameDuration = 1000.0 / 60
	count := float64(len(GeneralSkin.JudgmentSprites))
	period := int64(frameDuration * count)
	return JudgmentDrawer{
		Timer:   draws.NewTimer(gosu.TimeToTick(250), gosu.TimeToTick(period)),
		Sprites: GeneralSkin.JudgmentSprites,
	}
}
func (d *JudgmentDrawer) Update(worst gosu.Judgment) {
	if d.Countdown <= 0 {
		d.Judgment = gosu.Judgment{}
	} else {
		d.Countdown--
	}
	if worst.Valid() {
		d.Judgment = worst
		d.Countdown = d.MaxCountdown
	}
}

func (d JudgmentDrawer) Draw(screen *ebiten.Image) {
	if d.Countdown <= 0 || d.Judgment.Window == 0 {
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
	sprite := d.Frame(d.Sprites[idx])
	scale := 1.0
	switch {
	case age < 0.1:
		scale = 1.15 * (1 + age)
	case age >= 0.1 && age < 0.2:
		scale = 1.15 * (1.2 - age)
	case age > 0.9:
		scale = 1 - 1.15*(age-0.9)
	}
	sprite.ApplyScale(scale)
	sprite.Draw(screen, ebiten.DrawImageOptions{})
}
