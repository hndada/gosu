package drum

import (
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

// Logical size of in-game screen
const (
	screenSizeX = gosu.ScreenSizeX
	screenSizeY = gosu.ScreenSizeY
)

var SpeedScale float64 = 1.1
var KeySettings = [4]input.Key{input.KeyD, input.KeyF, input.KeyJ, input.KeyK}

// Default values are derived from osu! taiko.
// Todo: Should NoteHeight be separated into NoteHeight, HeadHeight, TailHeight?
var (
	FieldDarkness float64 = 1

	FieldPosition float64 = screenSizeY * 0.4115
	FieldHeight   float64 = screenSizeY * 0.26
	// Height of notes are dependent of FieldHeight.
	bigNoteHeight     float64 = FieldHeight * 0.725
	regularNoteHeight float64 = bigNoteHeight * 0.65
	FieldInnerHeight  float64 = FieldHeight * 0.875 // For drawing bars. // screenSizeY * 0.23

	HitPosition float64 = screenSizeX * 0.1875
	minPosition float64 = -HitPosition
	maxPosition float64 = minPosition + screenSizeX
	ShakePosX   float64 = screenSizeX * 0.375
	ShakePosY   float64 = screenSizeY * 0.55

	DancerPosX float64 = screenSizeX * 0.05
	DancerPosY float64 = screenSizeY * 0.1

	keyCenter float64 // Used in key sprites and combo position.
	// Range of ShakeCountPosition is [0, 1].
	// Min: Right bottom of the middle of Shake spin.
	// Max: Right bottom of the border of Shake spin.
	ShakeCountPosition float64 = 0.5
)

// Skin-dependent settings.
// JudgmentScale might have scaled by FieldHeight.
// Yet, Judgment is not circle image, actually.
var (
	JudgmentScale float64 = 1
	DotScale      float64 = 1
	ShakeScale    float64 = 1
	KeyScale      float64 = 1
	DancerScale   float64 = 1
	ComboScale    float64 = 0.75
	ComboDigitGap float64 = screenSizeX * -0.001
	// Used at roll tick combo.
	DotCountScale    float64 = ComboScale * 0.4
	DotCountDigitGap float64 = ComboDigitGap * 0.4
	// Used at shake countdown.
	ShakeCountScale    float64 = ComboScale * 1
	ShakeCountDigitGap float64 = ComboDigitGap * 1
)

func SwitchDirection() {
	max, min := maxPosition, minPosition
	maxPosition = -min
	minPosition = -max
}
