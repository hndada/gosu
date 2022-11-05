package draws

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Subjects  []Subject
	StartTime int64
	Time      int64
	Duration  int64 // float64
}

func (d *Animation) Update(time, duration int64, reset bool) {
	if reset {
		d.StartTime = time
	}
	d.Time = time
	d.Duration = duration
}
func (d Animation) Frame() int {
	td := float64(d.Time - d.StartTime)
	duration := float64(d.Duration)
	rate := math.Remainder(td, duration) / duration
	if rate < 0 {
		rate += 1
	}
	return int(rate * float64(len(d.Subjects)))
}
func (d Animation) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	subject := d.Subjects[d.Frame()]
	subject.Draw(screen, *op)
}
