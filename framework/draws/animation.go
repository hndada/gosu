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

type Animation []Sprite

func NewAnimation(fsys fs.FS, name string) Animation {
	return NewAnimationFromImages(LoadImages(fsys, name))
}
func NewAnimationFromImages(images []Image) (a Animation) {
	a = make(Animation, len(images))
	for i, image := range images {
		a[i] = NewSpriteFromSource(image)
	}
	return
}
func LoadImages(fsys fs.FS, name string) (is []Image) {
	const ext = ".png"

	// name supposed to have no extension when passed in LoadImages.
	name = strings.TrimSuffix(name, filepath.Ext(name))

	one := []Image{LoadImage(fsys, name+ext)}
	// dir, err := fsys.Open(name)
	// if err != nil {
	// 	return one
	// }
	// defer dir.Close()
	// fs, err := dir.ReadDir(-1)
	fs, err := fs.ReadDir(fsys, name)
	if err != nil {
		return one
	}
	nums := make([]int, 0, len(fs))
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		num := strings.TrimSuffix(f.Name(), ext)
		if num, err := strconv.Atoi(num); err == nil {
			nums = append(nums, num)
		}
	}
	sort.Ints(nums)
	for _, num := range nums {
		// Avoid use filepath here; it yields backslash, which is invalid path for FS.
		name2 := path.Join(name, fmt.Sprintf("%d.png", num))
		is = append(is, LoadImage(fsys, name2))
	}
	return
}
