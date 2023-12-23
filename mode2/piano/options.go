package piano

import (
	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

// Fixing screen size enables users to feel UI elements
// are consistent regardless of the screen size.
const (
	ScreenW = 640
	ScreenH = 480
)

// Objective: manage UI components with each own struct.
// Options is for passing from game to mode.
type Options struct {
	mode.MusicOpts
	mode.SoundOpts

	KeyOpts
	StageOpts
	Key struct {
		Ws func() []float64
		Xs func() []float64
	}
	Stage struct {
		W func() float64
		X func() float64
	}
	FieldOpts
	HintOpts
	BarOpts
	NoteOpts
	LightOpts
	JudgmentOpts
	mode.ComboOpts
	mode.ScoreOpts
}

func NewOptions(keyCount int) Options {
	opts := Options{}
	return opts
}

func (opts Options) newStage() draws.WHXY {
	rw := opts.StageRWs[opts.KeyCount]
	return draws.WHXY{
		W: rw * opts.Screen.W,
		H: opts.Screen.H,
		X: opts.StageRX * opts.Screen.W,
		Y: 0,
	}
}

func (opts Options) BaselineY() float64 {
	return opts.BaselineRY * opts.Screen.H
}

// I'm personally proud of this code.
func (opts Options) KeyOrder() []KeyKind {
	order := opts.KeyOrders[opts.KeyCount]
	order_1 := opts.KeyOrders[opts.KeyCount-1]

	switch opts.Scratches[opts.KeyCount] {
	case ScratchNone:
		return order
	case ScratchLeft:
		return append([]KeyKind{Tip}, order_1...)
	case ScratchRight:
		return append(order_1, Tip)
	}
	return nil
}

func (opts Options) KeyWs() (ws []float64) {
	ws = make([]float64, opts.KeyCount)
	for k, kind := range opts.KeyOrder() {
		rw := opts.KeyRWs[kind]
		ws[k] = rw * opts.Screen.W
	}

	// Adjust key width to fit the stage width.
	var rawSum float64
	for _, w := range ws {
		rawSum += w
	}
	scale := opts.StageW() / rawSum
	for k := range ws {
		ws[k] *= scale
	}
	return ws
}

// KeyXs returns centered x positions.
func (opts Options) KeyXs() (xs []float64) {
	xs = make([]float64, opts.KeyCount)
	ws := opts.KeyWs()
	x := opts.StageX() - opts.StageW()/2
	for k, w := range ws {
		x += w / 2
		xs[k] = x
		x += w / 2
	}
	return
}

// NoteExposureDuration returns time in milliseconds
// that cursor takes to move 1 logical pixel.
func (opts Options) NoteExposureDuration(speed float64) int32 {
	y := opts.BaselineY()
	return int32(y / speed)
}

// NewXxxComponent() requires multiple arguments.
// XxxArgs is for wrapping required arguments.
// Config.NewXxxArgs() returns XxxArgs based on configuration values.
// Separating Config and Args is also a good idea for post-processing.

// No XxxArgs. It just makes the code too verbose.

// Introducing interface to field would make the code too verbose.
// Stage  struct {
// 	W func() draws.Px
// 	X func() draws.Px
// }

// To take both brevity and clarity, structs are named
// using 3 to 4-lettered abbreviation containing vowels:
// XxxRes
// XxxOpts
// XxxArgs
// XxxComp
