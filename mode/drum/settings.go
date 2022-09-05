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
	// FieldDarkness float64 = 1

	FieldPosition    float64 = screenSizeY * 0.4115
	FieldHeight      float64 = screenSizeY * 0.26
	FieldInnerHeight float64 = screenSizeY * 0.23

	HitPosition float64 = screenSizeX * 0.1875
	minPosition float64 = -HitPosition
	maxPosition float64 = minPosition + screenSizeX
	posMargin   float64 = 100 // It should be larger than MaxSize/2 of all note sprites' width or height.

	// Height of notes are fixed.
	BigNoteHeight    float64 = FieldHeight * 0.725
	NormalNoteHeight float64 = BigNoteHeight * 0.65
	bodyLoss         float64 = 0 // No body loss in Drum mode.

	// Derived from other values.
	comboPosition         float64
	rollTickComboPosition float64 = HitPosition

	DancerPosX float64 = screenSizeX * 0.05
	DancerPosY float64 = screenSizeY * 0.1
)

// Skin-dependent settings.
// JudgmentScale might have scaled by FieldHeight.
// Yet, Judgment is not circle image, actually.
var (
	RollDotScale          float64 = 1
	ShakeScale            float64 = 1
	KeyScale              float64 = 1
	DancerScale           float64 = 1
	ComboScale            float64 = 0.75
	ComboGap              float64 = screenSizeX * -0.001
	RollTickComboScale    float64 = 0.3
	RollTickComboDigitGap float64 = ComboGap * 0.4
	JudgmentScale         float64 = 1
)

func SwitchDirection() {
	max, min := maxPosition, minPosition
	maxPosition = -min
	minPosition = -max
}
