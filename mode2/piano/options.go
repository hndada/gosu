package piano

import (
	mode "github.com/hndada/gosu/mode2"
)

// Objective: manage UI components with each own struct.
// Options is for passing from game to mode.
type Options struct {
	// Music    mode.MusicOpts
	// Sound    mode.SoundOpts
	Key      KeyOpts
	Field    FieldOpts
	Hint     HintOpts
	Bar      BarOpts
	Note     NoteOpts
	Light    LightOpts
	Judgment JudgmentOpts
	Combo    mode.ComboOpts
	Score    mode.ScoreOpts
}

func NewOptions(keyCount int) Options {
	opts := Options{}
	return opts
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
