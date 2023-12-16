var fingersList = [][]int{
	{},
	{0},
	{1, 1},
	{1, 0, 1},
	{2, 1, 1, 2},
	{2, 1, 0, 1, 2},
	{3, 2, 1, 1, 2, 3},
	{3, 2, 1, 0, 1, 2, 3},
	{4, 3, 2, 1, 1, 2, 3, 4},
	{4, 3, 2, 1, 0, 1, 2, 3, 4},
	{4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
}

func Fingers(keyCount int, scratchMode ScratchMode) []int {
	maxFinger := fingersList[keyCount-1][0] + 1
	switch scratchMode {
	case NoScratch:
		return fingersList[keyCount]
	case LeftScratch:
		return append([]int{maxFinger}, fingersList[keyCount-1]...)
	case RightScratch:
		return append(fingersList[keyCount-1], maxFinger)
	}
	return nil
}
