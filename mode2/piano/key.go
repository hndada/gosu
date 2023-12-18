package piano

import (
	"image/color"

	mode "github.com/hndada/gosu/mode2"
)

type KeyConfig struct {
	KeyboardMapping  map[int][]string
	SpotlightColors  [4]color.NRGBA
	HitLightScale    float64
	HitLightOpacity  float64
	HoldLightScale   float64
	HoldLightOpacity float64
}

func NewKeyConfig(screen mode.ScreenConfig, stage *StageConfig) KeyConfig {
	return KeyConfig{
		KeyboardMapping: map[int][]string{
			4:  {"D", "F", "J", "K"},
			5:  {"D", "F", "Space", "J", "K"},
			6:  {"S", "D", "F", "J", "K", "L"},
			7:  {"S", "D", "F", "Space", "J", "K", "L"},
			8:  {"A", "S", "D", "F", "Space", "J", "K", "L"},
			9:  {"A", "S", "D", "F", "Space", "J", "K", "L", "Semicolon"},
			10: {"A", "S", "D", "F", "V", "N", "J", "K", "L", "Semicolon"},
		},
		SpotlightColors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
		HitLightScale:    1.0,
		HitLightOpacity:  0.5,
		HoldLightScale:   1.0,
		HoldLightOpacity: 1.2,
	}
}
