package util

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type BatchElement struct {
	num int
	ext string
}

// Avoid using filepath at fs.FS.
// It yields backslash, which is invalid.
func BatchElements(fsys fs.FS, dirName string) []string {
	es, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return []string{}
	}

	fns := make([]BatchElement, 0, len(es))
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

		fns = append(fns, BatchElement{num, ext})
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
