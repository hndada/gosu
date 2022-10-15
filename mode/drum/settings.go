package drum

import (
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

// Logical size of in-game screen.
const (
	screenSizeX = gosu.ScreenSizeX
	screenSizeY = gosu.ScreenSizeY
)

var SpeedScale float64 = 1.0
var KeySettings = [4]input.Key{input.KeyD, input.KeyF, input.KeyJ, input.KeyK}

const PositionMargin = 100

// Default values are derived from osu!taiko.
// Todo: generalize Dancer for all modes?
var (
	FieldDarkness float64 = 0.7
	FieldPosition float64 = screenSizeY * 0.4115
	FieldHeight   float64 = screenSizeY * 0.26

	// Height of notes are dependent of FieldHeight.
	bigNoteHeight     float64 = FieldHeight * 0.725
	regularNoteHeight float64 = bigNoteHeight * 0.65

	HitPosition     float64 = screenSizeX * 0.1875
	minPosition     float64 = -HitPosition - PositionMargin
	maxPosition     float64 = -HitPosition + screenSizeX + PositionMargin
	DancerPositionX float64 = screenSizeX * 0.1
	DancerPositionY float64 = screenSizeY * 0.175

	keyCenter float64 // Used in key sprites and combo position.
)

// Skin-dependent settings.
var (
	FieldInnerHeight float64 = FieldHeight * 0.95
	JudgmentScale    float64 = 0.75 // 1.25
	DotScale         float64 = 0.5
	ShakeScale       float64 = 1
	DancerScale      float64 = 0.75 // 0.6
	ComboScale       float64 = 0.75 // 1.25
	ComboDigitGap    float64 = screenSizeX * -0.001
)

func SwitchDirection() {
	max, min := maxPosition, minPosition
	maxPosition = -min
	minPosition = -max
}
