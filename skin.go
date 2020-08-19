package gosu

import (
	"github.com/hndada/gosu/mode/mania"
	"image/color"
)

type Skin struct {
}

const (
	even int = iota
	odd
	middle
	pinky
)
const (
	scratchLeft  = 1 << 5 // 32
	scratchRight = 1 << 6 // 64
)

func NewNoteImageInfo2() map[int][]int {
	info := make(map[int][]int)
	info[0] = []int{}
	info[1] = []int{middle}
	info[2] = []int{even, even}
	info[3] = []int{even, middle, even}
	info[4] = []int{even, odd, odd, even}
	info[5] = []int{even, odd, middle, odd, even}
	info[6] = []int{even, odd, even, even, odd, even}
	info[7] = []int{even, odd, even, middle, even, odd, even}
	info[8] = []int{pinky, even, odd, even, even, odd, even, pinky}
	info[9] = []int{pinky, even, odd, even, middle, even, odd, even, pinky}
	info[10] = []int{pinky, even, odd, even, middle, middle, even, odd, even, pinky}

	for i := 1; i <= 8; i++ { // 정말 잘 짠듯
		info[i|scratchLeft] = append([]int{pinky}, info[i-1]...)
		info[i|scratchRight] = append(info[i-1], pinky)
	}
	return info
}

func noteColor(n mania.Note, keys int) color.RGBA {
	switch n.Key {
	case 0, 2, 4, 6:
		return color.RGBA{239, 243, 247, 0xff} // white
	case 1, 5:
		return color.RGBA{66, 211, 247, 0xff} // blue
	case 3:
		return color.RGBA{255, 203, 82, 0xff} // yellow
	}
	panic("not reach")
}
