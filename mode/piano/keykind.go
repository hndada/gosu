package piano

type KeyKind int

const (
	One KeyKind = iota
	Two
	Mid
	Tip
)

var keyKindsList = [][]KeyKind{
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
func KeyKinds(keyCount int, scratchMode ScratchMode) []KeyKind {
	switch scratchMode {
	case NoScratch:
		return keyKindsList[keyCount]
	case LeftScratch:
		return append([]KeyKind{Tip}, keyKindsList[keyCount-1]...)
	case RightScratch:
		return append(keyKindsList[keyCount-1], Tip)
	}
	return nil
}
