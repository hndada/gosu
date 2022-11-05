package draws

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Animation2 []Sprite

type Animation3 struct {
	Sprites []Sprite
	Timer
}

func (a Animation3) Frame() Sprite {
	count := float64(len(a.Sprites))
	return a.Sprites[int(count*a.Age())]
}

func NewAnimation(path string) (a Animation2) {
	const ext = ".png"
	one := Animation2{NewSprite(path + ext)}
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

func NewAnimation3(path string) (a Animation3) {
	const ext = ".png"
	one := Animation2{NewSprite(path + ext)}
	dir, err := os.Open(path)
	if err != nil {
		a.Sprites = one
		return
	}
	defer dir.Close()
	fs, err := dir.ReadDir(-1)
	if err != nil {
		a.Sprites = one
		return
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
		a.Sprites = append(a.Sprites, NewSprite(path))
	}
	return
}
