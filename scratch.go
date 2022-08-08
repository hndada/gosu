package main

// 40 = 32 + 8; 8 keyCount with left-scratch
const (
	LeftScratch  = 1 << 5 // 32
	RightScratch = 1 << 6 // 64
)
const ScratchMask = ^(LeftScratch | RightScratch)

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

// func WithScratch(keyCount int) int { return keyCount | ScratchMode[keyCount] }
