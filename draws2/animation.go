package draws

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Sprites   []Sprite
	StartTime int64
	Time      int64
	Duration  int64 // float64
	Point
}

func (d *Animation) Update(time, duration int64, reset bool) {
	if reset {
		d.StartTime = time
	}
	d.Time = time
	d.Duration = duration
}
func (d *Animation) Move(x, y float64) { d.Point = d.Point.Add(Pt(x, y)) }
func (d Animation) Frame() int {
	if d.Duration == 0 {
		return 0
	}
	td := float64(d.Time - d.StartTime)
	duration := float64(d.Duration)
	rate := math.Remainder(td, duration) / duration
	if rate < 0 {
		rate += 1
	}
	return int(rate * float64(len(d.Sprites)))
}
func (d Animation) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if len(d.Sprites) == 0 {
		return
	}
	sprite := d.Sprites[d.Frame()]
	sprite.Move(d.Point.XY())
	sprite.Draw(screen, op)
}
