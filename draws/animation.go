package draws

import (
	"io/fs"
	"time"
)

type Animation struct {
	frames    Frames
	StartTime time.Time
	Period    int32 // milliseconds
	Box
}

func NewAnimation(frames Frames, period int32) Animation {
	return Animation{
		frames:    frames,
		StartTime: time.Now(),
		Period:    period,
		Box:       NewBox(frames[0]),
	}
}

func NewAnimationFromFile(fsys fs.FS, name string, period int32) Animation {
	fs := NewFramesFromFile(fsys, name)
	return NewAnimation(fs, period)
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
