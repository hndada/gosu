package piano

// Convention: to organize types and structs in a file by
// defining the dependencies first and
// the types that utilize those dependencies later.
type KeyKind int

const (
	One KeyKind = iota
	Two
	Mid
	Tip
)

type Scratch int

const (
	ScratchNone = iota
	ScratchLeft
	ScratchRight
)

type KeyOpts struct {
	Count     int
	Mappings  map[int][]string
	Orders    map[int][]KeyKind
	Scratches map[int]Scratch
	Ws        [4]float64
	RY        float64 // Baseline
}

func NewKeyOpts() KeyOpts {
	return KeyOpts{
		Count: 4,
		Mappings: map[int][]string{
			1:  {"Space"},
			2:  {"F", "J"},
			3:  {"F", "Space", "J"},
			4:  {"D", "F", "J", "K"},
			5:  {"D", "F", "Space", "J", "K"},
			6:  {"S", "D", "F", "J", "K", "L"},
			7:  {"S", "D", "F", "Space", "J", "K", "L"},
			8:  {"A", "S", "D", "F", "Space", "J", "K", "L"},
			9:  {"A", "S", "D", "F", "Space", "J", "K", "L", "Semicolon"},
			10: {"A", "S", "D", "F", "V", "N", "J", "K", "L", "Semicolon"},
		},
		Orders: map[int][]KeyKind{
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
		},
		Scratches: map[int]Scratch{
			8: ScratchLeft,
		},
		Ws: [4]float64{
			80, // One
			78, // Two
			82, // Mid
			82, // Tip
		},
		RY: 0.90,
	}
}

type KeyComp struct {
	ws []float64
	xs []float64
}
