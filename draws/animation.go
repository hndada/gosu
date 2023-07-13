package draws

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Frames = []Image

// Read frames if there is a directory and the directory has entries.
// Read a single image if there is no directory or the directory has no entries.
func NewFramesFromFile(fsys fs.FS, name string) Frames {
	var frames Frames

	oneExt := filepath.Ext(name)
	one := []Image{NewImageFromFile(fsys, name)}

	dirName := strings.TrimSuffix(name, oneExt)
	fs, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return one
	}

	nums := make([]int, 0, len(fs))
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		frameExt := filepath.Ext(f.Name())
		num := strings.TrimSuffix(f.Name(), frameExt)
		if num, err := strconv.Atoi(num); err == nil {
			nums = append(nums, num)
		}
	}
	sort.Ints(nums)
	if len(nums) == 0 {
		return one
	}

	for _, num := range nums {
		// Avoid use filepath here.
		// It yields backslash, which is invalid path for FS.
		frameName := path.Join(name, fmt.Sprintf("%d.png", num))
		frames = append(frames, NewImageFromFile(fsys, frameName))
	}
	return frames
}

type Animation []Sprite

func NewAnimation(srcs any) Animation {
	switch srcs := srcs.(type) {
	case Frames:
		return newAnimationFromFrames(srcs)
	}
	return nil
}
func newAnimationFromFrames(frames Frames) Animation {
	a := make(Animation, len(frames))
	for i, img := range frames {
		a[i] = NewSprite(img)
	}
	return a
}

func NewAnimationFromFile(fsys fs.FS, name string) Animation {
	return NewAnimation(NewFramesFromFile(fsys, name))
}

func (a Animation) IsEmpty() bool {
	return len(a) <= 1 && a[0].IsEmpty()
}
