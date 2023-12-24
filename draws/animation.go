package draws

import (
	"io/fs"
	"time"
)

type Animation struct {
	frames Frames
	Box
	StartTime time.Time
	Period    int32 // milliseconds
	MaxLoop   int
}

func NewAnimation(frames Frames, period int32) Animation {
	return Animation{
		frames:    frames,
		Box:       NewBox(frames[0]),
		StartTime: time.Now(),
		Period:    period,
	}
}

func NewAnimationFromFile(fsys fs.FS, name string, period int32) Animation {
	fs := NewFramesFromFile(fsys, name)
	return NewAnimation(fs, period)
}

func (a *Animation) SetLoop(maxLoop int) { a.MaxLoop = maxLoop }
func (a Animation) Loop() int {
	if a.Period == 0 {
		return 0
	}

	d := time.Since(a.StartTime).Milliseconds()
	loop := int(d / int64(a.Period))
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

func (a *Animation) Reset() { a.StartTime = time.Now() }

func (a Animation) Draw(dst Image) {
	if len(a.frames) == 0 {
		return
	}

	if a.Period == 0 {
		a.Box.Draw(dst, a.frames[0])
		return
	}

	d := time.Since(a.StartTime).Milliseconds()
	progress := float64(d%int64(a.Period)) / float64(a.Period)
	count := float64(len(a.frames))
	index := int(progress * count)
	a.Box.Draw(dst, a.frames[index])
}
