package draws

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Frames = []Image

// Read frames if there is a directory and the directory has entries.
// Read a single image if there is no directory or the directory has no entries.
func NewFramesFromFilename(fsys fs.FS, name string) Frames {
	one := []Image{NewImageFromFile(fsys, name)}
	dirName := strings.TrimSuffix(name, filepath.Ext(name))

	es, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return one
	}

	names := frameNames(es)
	if len(names) == 0 {
		return one
	}

	frames := make(Frames, len(names))
	for i, name := range names {
		frames[i] = NewImageFromFile(fsys, name)
	}
	return frames
}

// Avoid using filepath at fs.FS.
// It yields backslash, which is invalid.
func frameNames(es []fs.DirEntry) []string {
	type frameName struct {
		num int
		ext string
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

	names := make([]string, len(fns))
	for i, fn := range fns {
		names[i] = strconv.Itoa(fn.num) + fn.ext
	}
	return names
}
