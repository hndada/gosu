package piano

type KeyType int

const (
	One KeyType = iota
	Two
	Mid
	Tip
)

var keyTypesList = [][]KeyType{
	{},
	{Mid},
	{One, One},
	{One, Mid, One},
	{One, Two, Two, One},
	{One, Two, Mid, Two, One},
	{One, Two, One, One, Two, One},
	{One, Two, One, Mid, One, Two, One},
	{Tip, One, Two, One, One, Two, One, Tip},
	{Tip, One, Two, One, Mid, One, Two, One, Tip},
	{Tip, One, Two, One, Mid, Mid, One, Two, One, Tip},
}

// I'm personally proud of this code.
func KeyTypes(keyCount int, scratchMode ScratchMode) []KeyType {
	switch scratchMode {
	case NoScratch:
		return keyTypesList[keyCount]
	case LeftScratch:
		return append([]KeyType{Tip}, keyTypesList[keyCount-1]...)
	case RightScratch:
		return append(keyTypesList[keyCount-1], Tip)
	}
	return nil
}
