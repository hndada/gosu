package piano

import (
	"github.com/hndada/gosu/game"
)

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

type KeysOpts struct {
	Count     int
	StageWs   map[int]float64
	stageW    float64
	StageX    float64
	BaselineY float64

	Mappings  map[int][]string
	Orders    map[int][]KeyKind
	Scratches map[int]Scratch
	KindWs    [4]float64
	ws        []float64
	xs        []float64
}

func NewKeysOpts() KeysOpts {
	opts := KeysOpts{
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
		KindWs: [4]float64{
			32, // One
			31, // Two
			33, // Mid
			33, // Tip
		},

		StageWs: map[int]float64{
			1:  240,
			2:  260,
			3:  280,
			4:  300,
			5:  320,
			6:  340,
			7:  360,
			8:  380,
			9:  400,
			10: 420,
		},
		StageX:    0.50 * game.ScreenW,
		BaselineY: 0.90 * game.ScreenH,
	}

	// Set derived fields.
	opts.stageW = opts.StageWs[opts.Count]
	opts.setKeyWs()
	opts.setXs()
	return opts
}

func (opts *KeysOpts) setKeyWs() {
	ws := make([]float64, opts.Count)
	for k, kind := range opts.Order() {
		ws[k] = opts.KindWs[kind]
	}

	// Adjust key width to fit the stage width.
	var rawSum float64
	for _, w := range ws {
		rawSum += w
	}
	scale := opts.stageW / rawSum

	for k := range ws {
		ws[k] *= scale
	}
	opts.ws = ws
}

// KeyXs returns centered x positions.
func (opts *KeysOpts) setXs() {
	xs := make([]float64, opts.Count)
	ws := opts.ws
	x := opts.StageX - opts.stageW/2
	for k, w := range ws {
		x += w / 2
		xs[k] = x
		x += w / 2
	}
	opts.xs = xs
}

// I'm personally proud of this code.
func (opts KeysOpts) Order() []KeyKind {
	order := opts.Orders[opts.Count]
	order_1 := opts.Orders[opts.Count-1]

	switch opts.Scratches[opts.Count] {
	case ScratchNone:
		return order
	case ScratchLeft:
		return append([]KeyKind{Tip}, order_1...)
	case ScratchRight:
		return append(order_1, Tip)
	}
	return nil
}
