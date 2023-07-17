package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
)

// Todo: SoundVolume -> SoundVolumeScale?
type Config struct {
	// settings exported from parent configuration
	ScreenSize  *draws.Vector2
	MusicVolume *float64
	SoundVolume *float64
	MusicOffset *int32

	// logic-affecting settings
	KeySettings       map[int][]string
	SpeedScale        float64
	HitPosition       float64
	TailExtraDuration float64

	// others
	KeyKindWidths         [4]float64
	FieldWidthScales      map[int]float64
	NoteHeigth            float64 // Applies to all types of notes.
	FieldPosition         float64
	ComboPosition         float64
	JudgmentPosition      float64
	ScratchColor          color.NRGBA
	FieldOpacity          float64
	KeyKindLightingColors [4]color.NRGBA
	HitLightingOpacity    float64
	HoldLightingOpacity   float64
	// BodyStyle          int // Stretch or Attach.
	// ReverseBody        int // Might have been used for stepmania skin.

	ScoreSpriteScale    float64
	ComboSpriteScale    float64
	ComboDigitGap       float64
	JudgmentSpriteScale float64
	HintHeight          float64
	LightingSpriteScale float64
}

// Todo: ScreenSize
func NewConfig() *Config {
	ScreenSize := draws.Vector2{X: 1600, Y: 900}
	return &Config{
		KeySettings: map[int][]string{
			4:  {"D", "F", "J", "K"},
			5:  {"D", "F", "Space", "J", "K"},
			6:  {"S", "D", "F", "J", "K", "L"},
			7:  {"S", "D", "F", "Space", "J", "K", "L"},
			8:  {"A", "S", "D", "F", "Space", "J", "K", "L"},
			9:  {"A", "S", "D", "F", "Space", "J", "K", "L", "Semicolon"},
			10: {"A", "S", "D", "F", "V", "N", "J", "K", "L", "Semicolon"},
		},
		SpeedScale:        1.0,
		HitPosition:       0.90 * ScreenSize.Y,
		TailExtraDuration: 0,

		KeyKindWidths: [4]float64{
			0.055 * ScreenSize.X, // One
			0.055 * ScreenSize.X, // Two
			0.055 * ScreenSize.X, // Mid
			0.055 * ScreenSize.X, // Tip
		},
		FieldWidthScales: map[int]float64{
			4: 1.2,  // 4.8
			5: 1.1,  // 5.5
			6: 1.05, // 6.3
			7: 1.0,  // 7.0
			8: 1.0,  // 8.0
			9: 1.0,  // 9.0
		},
		NoteHeigth:       0.03 * ScreenSize.Y, // 0.03: 27px
		FieldPosition:    0.50 * ScreenSize.X,
		ComboPosition:    0.40 * ScreenSize.Y,
		JudgmentPosition: 0.66 * ScreenSize.Y,
		ScratchColor:     color.NRGBA{224, 0, 0, 255},
		FieldOpacity:     0.8,
		KeyKindLightingColors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
		HitLightingOpacity:  0.5,
		HoldLightingOpacity: 1.5,

		ScoreSpriteScale:    0.65,
		ComboSpriteScale:    0.75,
		ComboDigitGap:       -1, // unit: pixel
		JudgmentSpriteScale: 0.33,
		HintHeight:          0.05 * ScreenSize.Y, // 0.06: 45px
		LightingSpriteScale: 1.0,
	}
}

func (cfg Config) KeyWidths(keyCount int, sm ScratchMode) []float64 {
	ws := make([]float64, keyCount)
	scale := cfg.FieldWidthScales[keyCount]
	for k, kk := range KeyKinds(keyCount, sm) {
		w := cfg.KeyKindWidths[kk] // Todo: math.Ceil()?
		ws[k] = w * scale
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

// NoteExposureDuration returns time in milliseconds
// that cursor takes to move 1 logical pixel.
func (cfg Config) NoteExposureDuration(speed float64) int32 {
	return int32(cfg.HitPosition / speed)
}
