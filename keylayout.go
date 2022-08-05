package main

type NoteKind int

const (
	One NoteKind = iota
	Two
	Mid
	Tip
)

var NoteKindsMap = map[int][]NoteKind{
	0:  []NoteKind{},
	1:  []NoteKind{Mid},
	2:  []NoteKind{One, One},
	3:  []NoteKind{One, Mid, One},
	4:  []NoteKind{One, Two, Two, One},
	5:  []NoteKind{One, Two, Mid, Two, One},
	6:  []NoteKind{One, Two, One, One, Two, One},
	7:  []NoteKind{One, Two, One, Mid, One, Two, One},
	8:  []NoteKind{Tip, One, Two, One, One, Two, One, Tip},
	9:  []NoteKind{Tip, One, Two, One, Mid, One, Two, One, Tip},
	10: []NoteKind{Tip, One, Two, One, Mid, Mid, One, Two, One, Tip},
}

var FingerMap = map[int][]int{
	0:  []int{},
	1:  []int{0},
	2:  []int{1, 1},
	3:  []int{1, 0, 1},
	4:  []int{2, 1, 1, 2},
	5:  []int{2, 1, 0, 1, 2},
	6:  []int{3, 2, 1, 1, 2, 3},
	7:  []int{3, 2, 1, 0, 1, 2, 3},
	8:  []int{4, 3, 2, 1, 1, 2, 3, 4},
	9:  []int{4, 3, 2, 1, 0, 1, 2, 3, 4},
	10: []int{4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
}

// 40 = 32 + 8; 8 keyCount with left-scratch
const (
	LeftScratch  = 1 << 5 // 32
	RightScratch = 1 << 6 // 64
)
const ScratchMask = ^(LeftScratch | RightScratch)

// func WithScratch(keyCount int) int { return keyCount | ScratchMode[keyCount] }

// I'm proud of these code.
func init() {
	for k := 2; k <= 8; k++ {
		NoteKindsMap[k|LeftScratch] = append([]NoteKind{Tip}, NoteKindsMap[k-1]...)
		NoteKindsMap[k|RightScratch] = append(NoteKindsMap[k-1], Tip)
	}
}
func init() {
	for k := 2; k <= 8; k++ {
		FingerMap[k|LeftScratch] = append([]int{FingerMap[k-1][0] + 1}, FingerMap[k-1]...)
		FingerMap[k|RightScratch] = append(FingerMap[k-1], FingerMap[k-1][k-2]+1)
	}
}
