package draws

import (
	"io/fs"
	"time"

	"github.com/hndada/gosu/times"
)

type Animation struct {
	frames Frames
	Box
	StartTime time.Time
	Period    time.Duration
	MaxLoop   int
}

func NewAnimation(frames Frames, period time.Duration) Animation {
	return Animation{
		frames:    frames,
		Box:       NewBox(frames[0]),
		StartTime: times.Now(),
		Period:    period,
	}
}

func NewAnimationFromFile(fsys fs.FS, name string, period time.Duration) Animation {
	fs := NewFramesFromFile(fsys, name)
	return NewAnimation(fs, period)
}

func (a *Animation) SetLoop(maxLoop int) { a.MaxLoop = maxLoop }
func (a Animation) Loop() int {
	if a.Period == 0 {
		return 0
	}

	loop := int(times.Since(a.StartTime) / a.Period)
	if loop > a.MaxLoop {
		loop = a.MaxLoop
	}
	return loop
}

func (a Animation) IsFinished() bool {
	if a.Period == 0 || a.MaxLoop == 0 {
		return false
	}
	return a.Loop() >= a.MaxLoop
}

func (a *Animation) Reset() { a.StartTime = times.Now() }

func (a Animation) Draw(dst Image) {
	if len(a.frames) == 0 {
		return
	}

	if a.Period == 0 {
		a.Box.Draw(dst, a.frames[0])
		return
	}

	r := times.Since(a.StartTime) % a.Period
	progress := float64(r) / float64(a.Period)
	count := float64(len(a.frames))
	index := int(progress * count)
	a.Box.Draw(dst, a.frames[index])
}
