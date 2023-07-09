package piano

import (
	"image/color"
)

// ScreenSize is a logical size of in-game screen.
// Todo: Make it be at one package only.
const (
	ScreenSizeX = 1600
	ScreenSizeY = 900
)

type Settings struct {
	// Settings that will be from scene.TheSettings
	MusicVolume float64
	SoundVolume float64
	Offset      int64

	// Logic settings
	KeySettings   map[int][]string
	SpeedScale    float64
	HitPosition   float64
	TailExtraTime float64
	// maxPosition   float64 // derived from HitPosition
	// minPosition   float64 // derived from HitPosition

	// Skin-independent settings
	NoteWidths         [4]float64
	NoteHeigth         float64 // Applies to all types of notes.
	BodyStyle          int
	FieldPosition      float64
	ComboPosition      float64
	JudgmentPosition   float64
	ScratchColor       color.NRGBA
	FieldOpaque        float64
	KeyLightingColors  [4]color.NRGBA
	HitLightingOpaque  float64
	HoldLightingOpaque float64

	// Skin-dependent settings
	ComboScale    float64
	ComboDigitGap float64
	JudgmentScale float64
	HintHeight    float64
	LightingScale float64
}

var TheSettings = Settings{
	// Settings that will be from scene.TheSettings
	// These values are placeholders.
	MusicVolume: 0.50,
	SoundVolume: 0.50,
	Offset:      -20,

	KeySettings: map[int][]string{
		4:  {"D", "F", "J", "K"},
		5:  {"D", "F", "Space", "J", "K"},
		6:  {"S", "D", "F", "J", "K", "L"},
		7:  {"S", "D", "F", "Space", "J", "K", "L"},
		8:  {"A", "S", "D", "F", "Space", "J", "K", "L"},
		9:  {"A", "S", "D", "F", "Space", "J", "K", "L", "Semicolon"},
		10: {"A", "S", "D", "F", "V", "N", "J", "K", "L", "Semicolon"},
	},
	SpeedScale:    1.0,
	HitPosition:   0.90 * ScreenSizeY,
	TailExtraTime: 0,

	NoteWidths: [4]float64{
		0.06 * ScreenSizeX, // One
		0.06 * ScreenSizeX, // Two
		0.06 * ScreenSizeX, // Mid
		0.06 * ScreenSizeX, // Tip
	},
	NoteHeigth:       0.03 * ScreenSizeY, // 0.03: 27px
	FieldPosition:    0.50 * ScreenSizeX,
	ComboPosition:    0.40 * ScreenSizeY,
	JudgmentPosition: 0.66 * ScreenSizeY,
	ScratchColor:     color.NRGBA{224, 0, 0, 255},
	FieldOpaque:      0.8,
	KeyLightingColors: [4]color.NRGBA{
		{224, 224, 224, 64}, // One: white
		{255, 170, 204, 64}, // Two: pink
		{224, 224, 0, 64},   // Mid: yellow
		{255, 0, 0, 64},     // Tip: red
	},
	HitLightingOpaque:  0.5,
	HoldLightingOpaque: 0.8,

	ComboScale:    0.75,
	ComboDigitGap: -1, // unit: pixel
	JudgmentScale: 0.33,
	HintHeight:    0.05 * ScreenSizeY, // 0.06: 45px
	LightingScale: 1.0,
}

// const positionMargin = 100

// func (s *Settings) process() {
// 	max := s.HitPosition
// 	s.maxPosition = max + positionMargin
// 	s.minPosition = max - ScreenSizeY - positionMargin
// }
