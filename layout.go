package main

type NoteKind int

const (
	One NoteKind = iota
	Two
	Mid
	Tip
)

var NoteKindsMap = map[int][]NoteKind{
	0:  {},
	1:  {Mid},
	2:  {One, One},
	3:  {One, Mid, One},
	4:  {One, Two, Two, One},
	5:  {One, Two, Mid, Two, One},
	6:  {One, Two, One, One, Two, One},
	7:  {One, Two, One, Mid, One, Two, One},
	8:  {Tip, One, Two, One, One, Two, One, Tip},
	9:  {Tip, One, Two, One, Mid, One, Two, One, Tip},
	10: {Tip, One, Two, One, Mid, Mid, One, Two, One, Tip},
}

var FingerMap = map[int][]int{
	0:  {},
	1:  {0},
	2:  {1, 1},
	3:  {1, 0, 1},
	4:  {2, 1, 1, 2},
	5:  {2, 1, 0, 1, 2},
	6:  {3, 2, 1, 1, 2, 3},
	7:  {3, 2, 1, 0, 1, 2, 3},
	8:  {4, 3, 2, 1, 1, 2, 3, 4},
	9:  {4, 3, 2, 1, 0, 1, 2, 3, 4},
	10: {4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
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
