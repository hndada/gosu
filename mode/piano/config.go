package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
)

type Config struct {
	ScreenSize draws.Vector2

	// Settings that will be from scene.TheSettings
	MusicVolume float64
	SoundVolume float64
	Offset      int32

	// Logic settings
	KeySettings   map[int][]string
	SpeedScale    float64
	HitPosition   float64
	TailExtraTime float64
	// maxPosition   float64 // derived from HitPosition
	// minPosition   float64 // derived from HitPosition

	KeyTypeWidths      [4]float64
	NoteHeigth         float64 // Applies to all types of notes.
	FieldPosition      float64
	ComboPosition      float64
	JudgmentPosition   float64
	ScratchColor       color.NRGBA
	FieldOpaque        float64
	KeyLightingColors  [4]color.NRGBA
	HitLightingOpaque  float64
	HoldLightingOpaque float64
	// BodyStyle          int // Stretch or Attach.
	// ReverseBody        int // Might have been used for stepmania skin.

	ScoreScale    float64
	ComboScale    float64
	ComboDigitGap float64
	JudgmentScale float64
	HintHeight    float64
	LightingScale float64
}

func DefaultConfig() *Config {
	ScreenSize := draws.Vector2{X: 1600, Y: 900}
	return &Config{
		ScreenSize: ScreenSize,
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
		HitPosition:   0.90 * ScreenSize.Y,
		TailExtraTime: 0,

		KeyTypeWidths: [4]float64{
			0.06 * ScreenSize.X, // One
			0.06 * ScreenSize.X, // Two
			0.06 * ScreenSize.X, // Mid
			0.06 * ScreenSize.X, // Tip
		},
		NoteHeigth:       0.03 * ScreenSize.Y, // 0.03: 27px
		FieldPosition:    0.50 * ScreenSize.X,
		ComboPosition:    0.40 * ScreenSize.Y,
		JudgmentPosition: 0.66 * ScreenSize.Y,
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
		HintHeight:    0.05 * ScreenSize.Y, // 0.06: 45px
		LightingScale: 1.0,
	}
}

func (cfg Config) KeyWidths(keyCount int, sm ScratchMode) []float64 {
	ws := make([]float64, keyCount)
	for k, kt := range KeyTypes(keyCount, sm) {
		ws[k] = cfg.KeyTypeWidths[kt] // Todo: math.Ceil()?
	}
	return ws
}
func (cfg Config) FieldWidth(keyCount int, sm ScratchMode) float64 {
	ws := cfg.KeyWidths(keyCount, sm)
	var fw float64
	for _, width := range ws {
		fw += width
	}
	return fw
}

// KeyXs returns centered x positions.
func (cfg Config) KeyXs(keyCount int, sm ScratchMode) []float64 {
	xs := make([]float64, keyCount)
	ws := cfg.KeyWidths(keyCount, sm)
	x := cfg.FieldPosition - cfg.FieldWidth(keyCount, sm)/2
	for k, w := range ws {
		x += w / 2
		xs[k] = x
		x += w / 2
	}
	return xs
}

// const positionMargin = 100

// func (s *Settings) process() {
// 	max := s.HitPosition
// 	s.maxPosition = max + positionMargin
// 	s.minPosition = max - ScreenSizeY - positionMargin
// }
