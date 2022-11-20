package piano

type KeyType int

const (
	One KeyType = iota
	Two
	Mid
	Tip // = Mid
)

// LeftScratch and RightScratch are bits for indicating scratch mode.
// For example, when key count is 40 = 32 + 8, it is 8-key with left scratch.
const (
	LeftScratch  = 32
	RightScratch = 64
	ScratchMask  = ^(LeftScratch | RightScratch)
)

var KeyTypes = map[int][]KeyType{
	general: {One, Two, Mid, Tip},
	1:       {Mid},
	2:       {One, One},
	3:       {One, Mid, One},
	4:       {One, Two, Two, One},
	5:       {One, Two, Mid, Two, One},
	6:       {One, Two, One, One, Two, One},
	7:       {One, Two, One, Mid, One, Two, One},
	8:       {Tip, One, Two, One, One, Two, One, Tip},
	9:       {Tip, One, Two, One, Mid, One, Two, One, Tip},
	10:      {Tip, One, Two, One, Mid, Mid, One, Two, One, Tip},
}

// I'm proud of this code.
func init() {
	for k := 2; k <= 8; k++ {
		KeyTypes[k|LeftScratch] = append([]KeyType{Tip}, KeyTypes[k-1]...)
		KeyTypes[k|RightScratch] = append(KeyTypes[k-1], Tip)
	}
}
