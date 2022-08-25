package piano

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
)

var MaxComboCountdown int = gosu.TimeToTick(2000)

type ComboDrawer struct {
	Combo     int
	Countdown int
	Sprites   []draws.Sprite
}

func (d *ComboDrawer) Update(combo int) {
	if d.Combo != combo {
		d.Countdown = MaxComboCountdown // + 1
	}
	d.Combo = combo
	if d.Countdown > 0 {
		d.Countdown--
	}
}

// ComboDrawer's Draw draws each number at constant x regardless of their widths.
// Each number image has different size; The standard width is number 0's.
func (d ComboDrawer) Draw(screen *ebiten.Image) {
	var wsum int
	if d.Combo == 0 || d.Countdown == 0 {
		return
	}
	vs := make([]int, 0)
	for v := d.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		// wsum += int(d.Sprites[v%10].W + ComboGap)
		wsum += int(d.Sprites[0].W) + int(ComboGap)
	}
	wsum -= int(ComboGap)

	t := MaxComboCountdown - d.Countdown
	age := float64(t) / float64(MaxJudgmentCountdown)
	x := screenSizeX/2 + float64(wsum)/2 - d.Sprites[0].W/2
	for _, v := range vs {
		// x -= d.Sprites[v].W + ComboGap
		x -= d.Sprites[0].W + ComboGap
		sprite := d.Sprites[v]
		// sprite.X = x
		sprite.X = x + (d.Sprites[0].W - sprite.W/2)
		sprite.SetCenterY(ComboPosition)
		switch {
		case age < 0.1:
			sprite.Y += 0.85 * age * sprite.H
		case age >= 0.1 && age < 0.2:
			sprite.Y += 0.85 * (0.2 - age) * sprite.H
		}
		sprite.Draw(screen)
	}
}

var MaxJudgmentCountdown int = gosu.TimeToTick(600)

type JudgmentDrawer struct {
	LastJudgment gosu.Judgment
	Countdown    int
	Sprites      []draws.Sprite
}

func (d *JudgmentDrawer) Update(worst gosu.Judgment) {
	if worst.Window != 0 {
		d.LastJudgment = worst
		d.Countdown = MaxJudgmentCountdown
	}
	if d.Countdown == 0 {
		d.LastJudgment = gosu.Judgment{}
	} else {
		d.Countdown--
	}
}

// JudgmentDrawer's Draw draws the same judgment for a while.
func (d JudgmentDrawer) Draw(screen *ebiten.Image) {
	if d.Countdown <= 0 {
		return
	}
	var sprite draws.Sprite
	for i, j := range Judgments {
		if j.Window == d.LastJudgment.Window {
			sprite = d.Sprites[i]
			break
		}
	}
	t := MaxJudgmentCountdown - d.Countdown
	age := float64(t) / float64(MaxJudgmentCountdown)
	switch {
	case age < 0.1:
		sprite.ApplyScale(sprite.ScaleW() * 1.15 * (1 + age))
	case age >= 0.1 && age < 0.2:
		sprite.ApplyScale(sprite.ScaleW() * 1.15 * (1.2 - age))
	case age > 0.9:
		sprite.ApplyScale(sprite.ScaleW() * (1 - 1.15*(age-0.9)))
	}
	sprite.SetCenterX(screenSizeX / 2)
	sprite.SetCenterY(JudgmentPosition)
	sprite.Draw(screen)
}

var MinKeyDownTicks int = gosu.TimeToTick(30)

// Todo: Draw KeyUp Permanently, occassionally draw KeyDown
type KeyDrawer struct {
	KeyDownCountdowns []int
	Sprites           [2][]draws.Sprite
}

func (d *KeyDrawer) Update(lastPressed, pressed []bool) {
	for k, cd := range d.KeyDownCountdowns {
		if cd > 0 {
			d.KeyDownCountdowns[k]--
		}
	}

	for k, p := range pressed {
		lp := lastPressed[k]
		if !lp && p {
			d.KeyDownCountdowns[k] = MinKeyDownTicks
		}
	}
}
func (d KeyDrawer) Draw(screen *ebiten.Image, pressed []bool) {
	for k, p := range pressed {
		if p || d.KeyDownCountdowns[k] > 0 {
			d.Sprites[1][k].Draw(screen) // KeyDown
		} else {
			d.Sprites[0][k].Draw(screen) // KeyUp
		}
	}
}
