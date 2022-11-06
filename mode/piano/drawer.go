package piano

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
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
	draws.Timer
	Cursor   float64
	Farthest *Note
	Nearest  *Note
	Sprites  [4]draws.Animation
}

// Farthest and Nearest are borders of displaying notes.
// All in-screen notes are confirmed to be drawn when drawing from Farthest to Nearest.
func (d *NoteDrawer) Update(cursor float64) {
	d.Ticker()
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
		sprite := d.Frame(d.Sprites[n.Type])
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
	body := d.Frame(d.Sprites[Body])
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
func (d KeyDrawer) Draw(screen *ebiten.Image) {
	const (
		up = iota
		down
	)
	sprite := d.Sprites[up]
	// It still draws for a while even when pressed off very shortly.
	if d.lastPressed || d.Tick < d.MaxTick {
		sprite = d.Sprites[down]
	}
	sprite.Draw(screen, ebiten.DrawImageOptions{})
}

type JudgmentDrawer struct {
	draws.Timer
	Sprites  [5]draws.Animation
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
	d.Ticker()
	if worst.Valid() {
		d.Judgment = worst
		d.Timer.Reset()
	}
}

func (d JudgmentDrawer) Draw(screen *ebiten.Image) {
	if d.Done() {
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
	sprite.ApplyScale(scale)
	sprite.Draw(screen, ebiten.DrawImageOptions{})
}
