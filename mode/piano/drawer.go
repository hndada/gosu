package piano

import (
	"image"
	"math"

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

// Bars are fixed: lane itself moves, all bars move same amount.
type BarDrawer struct {
	Sprite   draws.Sprite
	Cursor   float64
	Farthest *Bar
	Nearest  *Bar
}

func (d *BarDrawer) Update(cursor float64) {
	d.Cursor = cursor
	for d.Farthest.Position-d.Cursor <= maxPosition+margin {
		if d.Farthest.Next == nil {
			break
		}
		d.Farthest = d.Farthest.Next
	}
	for d.Nearest.Position-d.Cursor <= minPosition-margin {
		if d.Nearest.Next == nil {
			break
		}
		d.Nearest = d.Nearest.Next
	}
}

func (d BarDrawer) Draw(screen *ebiten.Image) {
	for b := d.Farthest; b != d.Nearest.Prev; b = b.Prev {
		sprite := d.Sprite
		pos := b.Position - d.Cursor
		sprite.Move(0, -pos)
		sprite.Draw(screen, nil)
	}
}

// Notes are fixed: lane itself moves, all notes move same amount.
type NoteLaneDrawer struct {
	Sprites  [4]draws.Sprite
	Cursor   float64
	Farthest *Note
	Nearest  *Note
}

func (d *NoteLaneDrawer) Update(cursor float64) {
	d.Cursor = cursor
	for d.Farthest != nil &&
		d.Farthest.Position-d.Cursor <= maxPosition+margin {
		d.Farthest = d.Farthest.Next
	}
	for d.Nearest != nil &&
		d.Nearest.Position-d.Cursor <= minPosition-margin {
		d.Nearest = d.Nearest.Next
	}
}

// Draw from farthest to nearest to make nearer notes priorly exposed.
func (d NoteLaneDrawer) Draw(screen *ebiten.Image) {
	n := d.Farthest
	for ; n != d.Nearest.Prev; n = n.Prev {
		sprite := d.Sprites[n.Type]
		pos := n.Position - d.Cursor
		sprite.Move(0, -pos)
		op := &ebiten.DrawImageOptions{}
		if n.Marked {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		sprite.Draw(screen, op)
		if n.Type == Head {
			d.DrawLongBody(screen, n)
		}
	}
	if n != nil && n.Type == Tail {
		d.DrawLongBody(screen, n.Prev)
	}
}

// DrawLongBody draws scaled, corresponding sub-image of Body sprite.
func (d NoteLaneDrawer) DrawLongBody(screen *ebiten.Image, head *Note) {
	tail := head.Next
	body := d.Sprites[Body]
	length := tail.Position - head.Position
	length -= -bodyLoss
	ratio := length / body.H()

	top := tail.Position
	if top > maxPosition {
		top = maxPosition
	}
	bottom := head.Position
	if bottom < minPosition {
		bottom = minPosition
	}
	subTop := math.Ceil((top - head.Position) / ratio)
	subBottom := math.Floor((bottom - head.Position) / ratio)
	subRect := image.Rect(0, int(subBottom), int(body.W()), int(subTop))
	subBody := body.SubSprite(subRect)

	op := &ebiten.DrawImageOptions{}
	if head.Marked {
		op.ColorM.ChangeHSV(0, 0.3, 0.3)
	}
	if ReverseBody {
		op.GeoM.Scale(1, -ratio)
	} else {
		op.GeoM.Scale(1, ratio)
	}
	subBody.Move(0, -subBottom)
	subBody.Draw(screen, op)
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

func NewKeyDrawer(ups, downs []draws.Sprite) KeyDrawer {
	return KeyDrawer{
		MinCountdown:   gosu.TimeToTick(30),
		Countdowns:     make([]int, len(ups)),
		KeyUpSprites:   ups,
		KeyDownSprites: downs,
	}
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
			d.KeyDownSprites[k].Draw(screen, nil)
		} else {
			d.KeyUpSprites[k].Draw(screen, nil)
		}
	}
}

type ComboDrawer struct {
	draws.BaseDrawer
	Sprites    [10]draws.Sprite
	DigitWidth float64
	DigitGap   float64
	Combo      int
}

// Each number has different width. Number 0's width is used as standard.
func NewComboDrawer(sprites [10]draws.Sprite) ComboDrawer {
	return ComboDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(2000),
		},
		Sprites:    sprites,
		DigitWidth: sprites[0].W(),
		DigitGap:   ComboDigitGap,
	}
}
func (d *ComboDrawer) Update(combo int) {
	if d.Countdown > 0 {
		d.Countdown--
	}
	if d.Combo != combo {
		d.Combo = combo
		d.Countdown = d.MaxCountdown
	}
}

// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (d ComboDrawer) Draw(screen *ebiten.Image) {
	if d.MaxCountdown != 0 && d.Countdown == 0 {
		return
	}
	if d.Combo == 0 {
		return
	}
	vs := make([]int, 0)
	for v := d.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at origin, no need to care of two 0.5w.
	w := d.DigitWidth + d.DigitGap
	// tx := (float64((len(vs)-1))*w + d.DigitWidth) / 2
	tx := float64(len(vs)-1) * w / 2
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(tx, 0)
		age := d.Age()
		h := sprite.H()
		switch {
		case age < 0.05:
			sprite.Move(0, 0.85*age*h)
		case age >= 0.05 && age < 0.1:
			sprite.Move(0, 0.85*(0.1-age)*h)
		}
		sprite.Draw(screen, nil)
		tx -= w
	}
}

type JudgmentDrawer struct {
	draws.BaseDrawer
	Sprites  []draws.Sprite
	Judgment gosu.Judgment
}

func NewJudgmentDrawer() (d JudgmentDrawer) {
	return JudgmentDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(600),
		},
		Sprites: GeneralSkin.JudgmentSprites,
	}
}
func (d *JudgmentDrawer) Update(worst gosu.Judgment) {
	if d.Countdown <= 0 {
		d.Judgment = gosu.Judgment{}
	} else {
		d.Countdown--
	}
	if worst.Window != 0 {
		d.Judgment = worst
		d.Countdown = d.MaxCountdown
	}
}

func (d JudgmentDrawer) Draw(screen *ebiten.Image) {
	if d.Countdown <= 0 {
		return
	}
	var sprite draws.Sprite
	for i, j := range Judgments {
		if j.Window == d.Judgment.Window {
			sprite = d.Sprites[i]
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
	sprite.SetScale(ratio, ratio, ebiten.FilterLinear)
	sprite.Draw(screen, nil)
}

// var MaxComboCountdown int = gosu.TimeToTick(2000)

// type ComboDrawer struct {
// 	Combo     int
// 	Countdown int
// 	Sprites   []draws.Sprite
// }

// func (d *ComboDrawer) Update(combo int) {
// 	if d.Combo != combo {
// 		d.Countdown = MaxComboCountdown // + 1
// 	}
// 	d.Combo = combo
// 	if d.Countdown > 0 {
// 		d.Countdown--
// 	}
// }

// func (d ComboDrawer) Draw(screen *ebiten.Image) {
// 	var wsum int
// 	if d.Combo == 0 || d.Countdown == 0 {
// 		return
// 	}
// 	vs := make([]int, 0)
// 	for v := d.Combo; v > 0; v /= 10 {
// 		vs = append(vs, v%10) // Little endian
// 		// wsum += int(d.Sprites[v%10].W + ComboDigitGap)
// 		wsum += int(d.Sprites[0].W) + int(ComboDigitGap)
// 	}
// 	wsum -= int(ComboDigitGap)

// 	t := MaxComboCountdown - d.Countdown
// 	age := float64(t) / float64(MaxJudgmentCountdown)
// 	x := screenSizeX/2 + float64(wsum)/2 - d.Sprites[0].W/2
// 	for _, v := range vs {
// 		// x -= d.Sprites[v].W + ComboDigitGap
// 		x -= d.Sprites[0].W + ComboDigitGap
// 		sprite := d.Sprites[v]
// 		// sprite.X = x
// 		sprite.X = x + (d.Sprites[0].W - sprite.W/2)
// 		sprite.SetCenterY(ComboPosition)
// 		switch {
// 		case age < 0.1:
// 			sprite.Y += 0.85 * age * sprite.H
// 		case age >= 0.1 && age < 0.2:
// 			sprite.Y += 0.85 * (0.2 - age) * sprite.H
// 		}
// 		sprite.Draw(screen)
// 	}
// }
