package draws

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationDrawer struct {
	Time      int64
	Duration  int64 // float64
	StartTime int64
	// EndTime   int64
	// Cycle   int
	Sprites []Sprite
}

//	func (d *AnimationDrawer) SetTime(time, duration int64) {
//		d.Time = time
//		d.StartTime = time
//		d.EndTime = time + duration
//	}
func (d *AnimationDrawer) Update(time, duration int64, reset bool) {
	d.Time = time
	d.Duration = duration
	if reset {
		d.StartTime = d.Time
	}
}
func (d AnimationDrawer) Frame() int {
	td := float64(d.Time - d.StartTime)
	duration := float64(d.Duration)
	rate := math.Remainder(td, duration) / duration
	if rate < 0 {
		rate += 1
	}
	return int(rate * float64(len(d.Sprites)))
}
func (d AnimationDrawer) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions, tx, ty float64) {
	sprite := d.Sprites[d.Frame()]
	sprite.Move(tx, ty)
	sprite.Draw(screen, op)
}
