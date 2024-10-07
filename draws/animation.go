package draws

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/hndada/gosu/times"
	"github.com/hndada/gosu/util"
)

type Frames = []Image

// NewFramesFromFile read a sequence of images if there is a directory
// and the directory has entries. Otherwise, read a single image
// if there is no directory or the directory has no entries.
func NewFramesFromFile(fsys fs.FS, name string) []Image {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	paths := util.DirElements(fsys, base)
	if len(paths) == 0 {
		one := NewImageFromFile(fsys, name)
		return []Image{one}
	}

	fs := make([]Image, len(paths))
	for i, name := range paths {
		fs[i] = NewImageFromFile(fsys, name)
	}
	return fs
}

type Animation struct {
	Frames    []Image
	StartTime time.Time
	Period    time.Duration
	MaxLoop   int
	Box
}

func NewAnimation(frames []Image, periodMS time.Duration) Animation {
	a := Animation{
		Frames:    frames,
		StartTime: times.Now(),
		Period:    periodMS * time.Millisecond,
	}
	if len(frames) > 0 {
		a.Box = NewBox(frames[0])
	}
	return a
}

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

func (a Animation) IsEmpty() bool { return len(a.Frames) == 0 }

func (a Animation) index() int {
	if a.Period == 0 {
		return 0
	}
	r := times.Since(a.StartTime) % a.Period
	progress := r.Seconds() / a.Period.Seconds()
	count := float64(len(a.Frames))
	idx := int(progress * count)
	return idx
}

func (a Animation) Size() XY {
	if a.IsEmpty() {
		return XY{}
	}
	return a.Frames[a.index()].Size()
}

func (a Animation) Draw(dst Image) {
	if a.IsEmpty() {
		return
	}
	src := a.Frames[a.index()]
	dst.DrawImage(src.Image, a.op())
}

func (a Animation) IsFinished() bool {
	if a.Period == 0 || a.MaxLoop == 0 {
		return false
	}
	return a.Loop() >= a.MaxLoop
}

func (a *Animation) Reset() { a.StartTime = times.Now() }
