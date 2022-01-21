package mania

// applied at keyCount
// 40 = 32 + 8; 8 keyCount with left-scratch
const (
	LeftScratch  = 1 << 5 // 32
	RightScratch = 1 << 6 // 64
)
const ScratchMask = ^(LeftScratch | RightScratch)

type keyKind uint8

const (
	one keyKind = iota
	two
	middle
	pinky
)

var keyKindsMap = make(map[int][]keyKind)
var fingers = make(map[int][]int)

func init() {
	keyKindsMap[0] = []keyKind{}
	keyKindsMap[1] = []keyKind{middle}
	keyKindsMap[2] = []keyKind{one, one}
	keyKindsMap[3] = []keyKind{one, middle, one}
	keyKindsMap[4] = []keyKind{one, two, two, one}
	keyKindsMap[5] = []keyKind{one, two, middle, two, one}
	keyKindsMap[6] = []keyKind{one, two, one, one, two, one}
	keyKindsMap[7] = []keyKind{one, two, one, middle, one, two, one}
	keyKindsMap[8] = []keyKind{pinky, one, two, one, one, two, one, pinky}
	keyKindsMap[9] = []keyKind{pinky, one, two, one, middle, one, two, one, pinky}
	keyKindsMap[10] = []keyKind{pinky, one, two, one, middle, middle, one, two, one, pinky}

	for k := 2; k <= 8; k++ { // I'm proud of writing these code by myself uwu
		keyKindsMap[k|LeftScratch] = append([]keyKind{pinky}, keyKindsMap[k-1]...)
		keyKindsMap[k|RightScratch] = append(keyKindsMap[k-1], pinky)
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
		fingers[k|LeftScratch] = append([]int{fingers[k-1][0] + 1}, fingers[k-1]...)
		fingers[k|RightScratch] = append(fingers[k-1], fingers[k-1][k-2]+1)
	}
}

func WithScratch(keyCount int) int { return keyCount | Settings.ScratchMode[keyCount] }
