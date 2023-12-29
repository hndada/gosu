package piano

import "github.com/hndada/gosu/game"

// Objective: manage UI cmponents with each own struct.
// Options is for passing from game to mode.
// A field is plural form if is drawn per key.

// No XxxArgs. It just makes the code too verbose.
// Introducing interface as a field would make the code too verbose.
type Options struct {
	KeyCount int
	Stage    StageOptions
	Key      KeysOptions

	Field         FieldOptions
	Bars          BarsOptions
	Hint          HintOptions
	KeysNotes     KeysNotesOptions
	KeysButton    KeysButtonOptions
	KeysBacklight KeysBacklightOptions
	KeysHitLight  KeysHitLightOptions
	KeysHoldLight KeysHoldLightOptions
	Judgment      JudgmentOptions
	Combo         game.ComboOptions
	Score         game.ScoreOptions
}
