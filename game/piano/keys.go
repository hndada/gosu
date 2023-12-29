package piano

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

type KeysOptions struct {
	keyCount  int
	Mappings  map[int][]string
	Orders    map[int][]KeyKind
	Scratches map[int]Scratch
	KindWs    [4]float64
	kw        []float64
	kx        []float64
	y         float64 // center bottom
}

func NewKeysOptions(stage StageOptions) KeysOptions {
	opts := KeysOptions{
		keyCount: stage.keyCount,
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
		y: stage.H,
	}
	opts.kw = opts.newKeysW(stage)
	opts.kx = opts.newKeysX(stage, opts.kw)
	return opts
}

// I'm personally proud of this code.
func (opts KeysOptions) Order() []KeyKind {
	order := opts.Orders[opts.keyCount]
	order_1 := opts.Orders[opts.keyCount-1]

	switch opts.Scratches[opts.keyCount] {
	case ScratchNone:
		return order
	case ScratchLeft:
		return append([]KeyKind{Tip}, order_1...)
	case ScratchRight:
		return append(order_1, Tip)
	}
	return nil
}

func (opts KeysOptions) newKeysW(stage StageOptions) []float64 {
	kw := make([]float64, opts.keyCount)
	for k, kind := range opts.Order() {
		kw[k] = opts.KindWs[kind]
	}

	// Adjust key width to fit the stage width.
	var rawSum float64
	for _, w := range kw {
		rawSum += w
	}
	scale := stage.w / rawSum

	for k := range kw {
		kw[k] *= scale
	}
	return kw
}

// centered x positions
func (opts *KeysOptions) newKeysX(stage StageOptions, kw []float64) []float64 {
	kx := make([]float64, opts.keyCount)
	x := stage.X - stage.w/2
	for k, w := range kw {
		x += w / 2
		kx[k] = x
		x += w / 2
	}
	return kx
}
