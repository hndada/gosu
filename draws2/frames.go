package draws

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hndada/gosu/times"
)

type Frames struct {
	Images    []Image
	StartTime time.Time
	Period    time.Duration
	MaxLoop   int
}

func NewFrames(imgs []Image, period time.Duration) Frames {
	return Frames{
		Images:    imgs,
		StartTime: times.Now(),
		Period:    period,
	}
}

func NewFramesFromFile(fsys fs.FS, name string, period time.Duration) Frames {
	imgs := NewFramesImagesFromFile(fsys, name)
	return NewFrames(imgs, period)
}

// NewFramesFromFile read a sequence of images if there is a directory
// and the directory has entries. Otherwise, read a single image
// if there is no directory or the directory has no entries.
func NewFramesImagesFromFile(fsys fs.FS, name string) []Image {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	paths := framePaths(fsys, base)
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

type frameName struct {
	num int
	ext string
}

// Avoid using filepath at fs.FS.
// It yields backslash, which is invalid.
func framePaths(fsys fs.FS, dirName string) []string {
	es, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return []string{}
	}

	fns := make([]frameName, 0, len(es))
	for _, f := range es {
		if f.IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name())
		numStr := strings.TrimSuffix(f.Name(), ext)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		fns = append(fns, frameName{num, ext})
	}

	sort.Slice(fns, func(i, j int) bool {
		return fns[i].num < fns[j].num
	})

	paths := make([]string, len(fns))
	for i, fn := range fns {
		name := strconv.Itoa(fn.num) + fn.ext
		paths[i] = path.Join(dirName, name)
	}
	return paths
}

func (fs *Frames) SetLoop(maxLoop int) { fs.MaxLoop = maxLoop }
func (fs Frames) Loop() int {
	if fs.Period == 0 {
		return 0
	}

	loop := int(times.Since(fs.StartTime) / fs.Period)
	if loop > fs.MaxLoop {
		loop = fs.MaxLoop
	}
	return loop
}

func (fs Frames) IsEmpty() bool { return len(fs.Images) == 0 }

func (fs Frames) Index() int {
	if fs.Period == 0 {
		return 0
	}
	r := times.Since(fs.StartTime) % fs.Period
	progress := float64(r) / float64(fs.Period)
	count := float64(len(fs.Images))
	index := int(progress * count)
	return index
}

func (fs Frames) Size() Vector2 {
	if fs.IsEmpty() {
		return Vector2{}
	}
	return fs.Images[fs.Index()].Size()
}

func (fs Frames) IsFinished() bool {
	if fs.Period == 0 || fs.MaxLoop == 0 {
		return false
	}
	return fs.Loop() >= fs.MaxLoop
}

func (fs *Frames) Reset() { fs.StartTime = times.Now() }
