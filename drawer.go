package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

type BackgroundDrawer struct {
	Dimness float64
	Sprite  draws.Sprite
}

func (d BackgroundDrawer) Draw(screen *ebiten.Image) {
	op := d.Sprite.Op()
	op.ColorM.ChangeHSV(0, 1, BgDimness)
	screen.DrawImage(d.Sprite.I, op)
}

type BarLineDrawer struct {
	Times      []int64
	Cursor     int     // Index of closest bar line.
	Offset     float64 // Bar line is drawn at bottom, not at the center.
	Sprite     draws.Sprite
	Horizontal bool
}

func (d *BarLineDrawer) Update(position func(time int64) float64) {
	bound := screenSizeY
	if d.Horizontal {
		bound = screenSizeX
	}
	t := d.Times[d.Cursor]
	// Bar line and Hint are anchored at the bottom.
	for d.Cursor < len(d.Times)-1 &&
		int(position(t)+d.Offset) >= bound {
		d.Cursor++
		t = d.Times[d.Cursor]
	}
}
func (d BarLineDrawer) Draw(screen *ebiten.Image, position func(time int64) float64) {
	for _, t := range d.Times[d.Cursor:] {
		sprite := d.Sprite
		sprite.Y = position(t) + d.Offset
		if sprite.Y < 0 {
			break
		}
		sprite.Draw(screen)
	}
}

//	type ScoreDrawer struct {
//		DelayedScore ctrl.Delayed
//		Sprites      []draws.Sprite
//	}
//
// ScoreDrawer.Update(int(math.Ceil(delayedScore)))
func NewScoreDrawer() draws.NumberDrawer {
	return draws.NumberDrawer{
		Sprites:       ScoreSprites,
		SignSprites:   SignSprites,
		DigitWidth:    ScoreSprites[0].W(),
		DigitGap:      0,
		FractionDigit: 0,
		Integer:       0,
		Fraction:      0,
		Effecter:      draws.Effecter{},
	}
}

// func (d *ScoreDrawer) Update(score float64) {
// 	d.DelayedScore.Set(score)
// 	d.DelayedScore.Update()
// }

// func (d ScoreDrawer) Draw(screen *ebiten.Image) {
// 	var wsum int
// 	vs := make([]int, 0)
// 	for v := int(math.Ceil(d.DelayedScore.Delayed)); v > 0; v /= 10 {
// 		vs = append(vs, v%10) // Little endian
// 		// wsum += int(d.Sprites[v%10].W)
// 		wsum += int(d.Sprites[0].W)
// 	}
// 	if len(vs) == 0 {
// 		vs = append(vs, 0) // Little endian
// 		wsum += int(d.Sprites[0].W)
// 	}
// 	x := float64(screenSizeX) - d.Sprites[0].W/2
// 	for _, v := range vs {
// 		// x -= d.Sprites[v].W
// 		x -= d.Sprites[0].W
// 		sprite := d.Sprites[v]
// 		sprite.X = x + (d.Sprites[0].W - sprite.W/2)
// 		sprite.Draw(screen)
// 	}
// }
