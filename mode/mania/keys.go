package mania

// These values are applied at keys
// Example: 40 = 32 + 8 = Left-scratching 8 Key
// todo: leftScratch, rightScratch로 변경
const (
	ScratchLeft  = 1 << 5 // 32
	ScratchRight = 1 << 6 // 64
)
const ScratchMask = ^(ScratchLeft | ScratchRight)

type keyKind uint8

const (
	one keyKind = iota
	two
	middle
	pinky
)

// todo: struct로?
// todo: int -> type finger int
var keyKinds = make(map[int][]keyKind)
var fingers = make(map[int][]int)

func init() {
	keyKinds[0] = []keyKind{}
	keyKinds[1] = []keyKind{middle}
	keyKinds[2] = []keyKind{one, one}
	keyKinds[3] = []keyKind{one, middle, one}
	keyKinds[4] = []keyKind{one, two, two, one}
	keyKinds[5] = []keyKind{one, two, middle, two, one}
	keyKinds[6] = []keyKind{one, two, one, one, two, one}
	keyKinds[7] = []keyKind{one, two, one, middle, one, two, one}
	keyKinds[8] = []keyKind{pinky, one, two, one, one, two, one, pinky}
	keyKinds[9] = []keyKind{pinky, one, two, one, middle, one, two, one, pinky}
	keyKinds[10] = []keyKind{pinky, one, two, one, middle, middle, one, two, one, pinky}

	for k := 2; k <= 8; k++ { // 정말 잘 짠듯
		keyKinds[k|ScratchLeft] = append([]keyKind{pinky}, keyKinds[k-1]...)
		keyKinds[k|ScratchRight] = append(keyKinds[k-1], pinky)
	}

	fingers[0] = []int{}
	fingers[1] = []int{0}
	fingers[2] = []int{1, 1}
	fingers[3] = []int{1, 0, 1}
	fingers[4] = []int{2, 1, 1, 2}
	fingers[5] = []int{2, 1, 0, 1, 2}
	fingers[6] = []int{3, 2, 1, 1, 2, 3}
	fingers[7] = []int{3, 2, 1, 0, 1, 2, 3}
	fingers[8] = []int{4, 3, 2, 1, 1, 2, 3, 4}
	fingers[9] = []int{4, 3, 2, 1, 0, 1, 2, 3, 4}
	fingers[10] = []int{4, 3, 2, 1, 0, 0, 1, 2, 3, 4}

	for k := 2; k <= 8; k++ {
		fingers[k|ScratchLeft] = append([]int{fingers[k-1][0] + 1}, fingers[k-1]...)
		fingers[k|ScratchRight] = append(fingers[k-1], fingers[k-1][k-2]+1)
	}
}
