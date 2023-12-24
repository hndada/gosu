package piano

import (
	mode "github.com/hndada/gosu/mode2"
)

// Objective: manage UI components with each own struct.
// Options is for passing from game to mode.
// A field is plural form if is drawn per key.
type Options struct {
	// Music    mode.MusicOpts
	// Sound    mode.SoundOpts
	Key        KeysOpts
	Field      FieldOpts
	Hint       HintOpts
	Bar        BarOpts
	Notes      NotesOpts
	KeyButtons KeyButtonsOpts
	Backlights BacklightsOpts
	HitLights  HitLightsOpts
	HoldLights HoldLightsOpts
	Judgment   JudgmentOpts
	Combo      mode.ComboOpts
	Score      mode.ScoreOpts
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
