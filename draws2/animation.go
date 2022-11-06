package draws

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Animation []Sprite

func NewAnimation(path string) (a Animation) {
	const ext = ".png"
	one := Animation{NewSprite(path + ext)}
	dir, err := os.Open(path)
	if err != nil {
		return one
	}
	defer dir.Close()
	fs, err := dir.ReadDir(-1)
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
		path := filepath.Join(path, fmt.Sprintf("%d.png", num))
		a = append(a, NewSprite(path))
	}
	return
}
