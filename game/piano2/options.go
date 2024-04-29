package piano

import (
	"image/color"

	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/game"
)

// SpeedScale is universal for all key counts.
// If a player wants to use different speed scales for different key counts,
// use 'Option Profile' feature.
type Options struct {
	// CurrentKeyCount int     // old: KeyCount
	SpeedScale float64 // Added

	StageWidths       map[int]float64
	StageBasePosition float64 // virtual value; nor X or Y
	StagePositionX    float64

	// If there are multiple iterable types in a single type,
	// The name should expose such information; e.g., XxxsMap
	// Unless the name itself contains the information.
	KeyMappings      map[int][]string
	KeyOrders        map[int][]KeyKind
	KeyScratchModes  map[int]ScratchMode
	KeyKindWidths    [4]float64
	keyWidthsMap     map[int][]float64 // derived
	keyPositionXsMap map[int][]float64 // derived

	FieldOpacity        float32
	BarHeight           float64
	HintHeight          float64
	NoteHeight          float64
	TailNoteOffset      int32
	NoteColors          [4]color.NRGBA
	BacklightColors     [4]color.NRGBA
	HitLightImageScale  float64
	HitLightOpacity     float32
	HoldLightImageScale float64
	HoldLightOpacity    float32
	JudgmentImageScale  float64
	JudgmentPositionY   float64
	Combo               ComboOptions
	Score               ScoreOptions
}

// gosu/game
type ComboOptions struct {
	ImageScale float64
	PositionX  float64
	DigitGap   float64
	PositionY  float64
	IsPersist  bool
	Bounce     float64
}

// gosu/game
type ScoreOptions struct {
	ImageScale float64
	DigitGap   float64
}

type KeyKind int

const (
	One KeyKind = iota
	Two
	Mid
	Tip
)

type ScratchMode int

const (
	ScratchModeNone = iota
	ScratchModeLeft
	ScratchModeRight
)

// piano.Options has all key count options so that
// it can handle scratch options smoothly.
func NewOptions() *Options {
	halfScreen := draws.XY{
		X: game.ScreenSizeX / 2,
		Y: game.ScreenSizeY / 2,
	}

	opts := &Options{
		SpeedScale: 1.0,

		StageWidths: map[int]float64{
			1:  halfScreen.X - 80,
			2:  halfScreen.X - 60,
			3:  halfScreen.X - 40,
			4:  halfScreen.X - 20,
			5:  halfScreen.X,
			6:  halfScreen.X + 20,
			7:  halfScreen.X + 40,
			8:  halfScreen.X + 60,
			9:  halfScreen.X + 80,
			10: halfScreen.X + 100,
		},
		StageBasePosition: 0.90 * game.ScreenSizeY,
		StagePositionX:    halfScreen.X,

		KeyMappings: map[int][]string{
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
		KeyOrders: map[int][]KeyKind{
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
		KeyScratchModes: map[int]ScratchMode{
			8: ScratchModeLeft,
		},
		KeyKindWidths: [4]float64{
			32, // One
			31, // Two
			33, // Mid
			33, // Tip
		},

		FieldOpacity:   0.8,
		BarHeight:      1,
		HintHeight:     24,
		NoteHeight:     20,
		TailNoteOffset: 0,
		NoteColors: [4]color.NRGBA{
			{255, 255, 255, 255}, // One: white
			{239, 191, 226, 255}, // Two: pink
			{218, 215, 103, 255}, // Mid: yellow
			{224, 0, 0, 255},     // Tip: red
		},
		BacklightColors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
		HitLightImageScale:  1.0,
		HitLightOpacity:     0.5,
		HoldLightImageScale: 1.0,
		HoldLightOpacity:    1.2,
		JudgmentImageScale:  0.33,
		JudgmentPositionY:   0.66 * game.ScreenSizeY,
		Combo: ComboOptions{
			ImageScale: 0.75,
			// PositionX should not be set by user.
			// It will be handled at Normalize().
			PositionX: halfScreen.X,
			DigitGap:  -1,
			PositionY: 0.40 * game.ScreenSizeY,
			IsPersist: false,
			Bounce:    0.85,
		},
		Score: ScoreOptions{
			ImageScale: 0.65,
			DigitGap:   0,
		},
	}

	opts.keyWidthsMap = make(map[int][]float64)
	opts.keyPositionXsMap = make(map[int][]float64)
	for keyCount := 1; keyCount <= 10; keyCount++ {
		ws := opts.keyWidths(keyCount)
		opts.keyWidthsMap[keyCount] = ws
		opts.keyPositionXsMap[keyCount] = opts.keyPositionXs(keyCount, ws)
	}

	return opts
}

// I'm personally proud of this code.
func (opts Options) KeyOrder(keyCount int) []KeyKind {
	order := opts.KeyOrders[keyCount]
	if keyCount == 1 {
		return order
	}
	order_1 := opts.KeyOrders[keyCount-1]

	m, ok := opts.KeyScratchModes[keyCount]
	if !ok {
		return order
	}

	switch m {
	case ScratchModeNone:
		return order
	case ScratchModeLeft:
		return append([]KeyKind{Tip}, order_1...)
	case ScratchModeRight:
		return append(order_1, Tip)
	}
	return nil
}

func (opts Options) keyWidths(keyCount int) []float64 {
	keysW := make([]float64, keyCount)
	for k, kind := range opts.KeyOrder(keyCount) {
		keysW[k] = opts.KeyKindWidths[kind]
	}

	// Adjust key width to fit the stage width.
	var rawSum float64
	for _, w := range keysW {
		rawSum += w
	}
	scale := opts.StageWidths[keyCount] / rawSum

	for k := range keysW {
		keysW[k] *= scale
	}
	return keysW
}

func (opts Options) keyPositionXs(keyCount int, ws []float64) []float64 {
	keysX := make([]float64, keyCount)
	x := opts.StagePositionX - opts.StageWidths[keyCount]/2
	for k, w := range ws {
		x += w / 2
		keysX[k] = x
		x += w / 2
	}
	return keysX
}
