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

var SpeedBase = 0.004
var KeySettings = [4]input.Key{input.KeyD, input.KeyF, input.KeyJ, input.KeyK}

// Todo: Should NoteHeight be separated into NoteHeight, HeadHeight, TailHeight?
var (
	// Default values are derived from osu! taiko.
	FieldPosition    float64 = screenSizeY * 0.4115
	FieldHeight      float64 = screenSizeY * 0.26
	FieldInnerHeight float64 = screenSizeY * 0.23
	HitPosition      float64 = screenSizeX * 0.1875

	// Height of notes are fixed.
	BigNoteHeight    float64 = FieldHeight * 0.725
	NormalNoteHeight float64 = BigNoteHeight * 0.65
	// bigNoteScale float64
	// normalNoteScale

	DancerPositionX float64 = screenSizeX * 0.05
	DancerPositionY float64 = screenSizeY * 0.1

	// BarLineWidth float64 = screenSizeX * 0.005

	comboPosition         float64
	rollTickComboPosition float64 = HitPosition
)

var FieldDarkness float64 = 1

var (
	ScoreScale         float64 = 0.65
	ComboScale         float64 = 0.75
	ComboGap           float64 = screenSizeX * -0.001
	RollTickComboScale float64 = 0.3
	RollTickComboGap   float64 = ComboGap * 0.4
	KeyScale           float64 = 1
	// JudgmentScale might have scaled by FieldHeight.
	// Yet, Judgment is not circle image, actually.
	JudgmentScale float64 = 1
	RollDotScale  float64 = 1
	ShakeScale    float64 = 1

	DancerScale float64 = 1
)

// 1 pixel is 1 millisecond.
// Todo: Separate NoteHeight / 2 at piano mode
func ExposureTime(speedBase float64) float64 {
	return (screenSizeX - HitPosition) / speedBase
}
func ExposureDegree(speedBase float64) (float64, float64) {
	return ExposureTime(speedBase), BigNoteHeight
}
