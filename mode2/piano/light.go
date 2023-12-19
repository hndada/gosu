package piano

import "image/color"

// Back, Hit, Hold
type LightConfig struct {
	BackColors  [4]color.NRGBA
	HitScale    float64
	HitOpacity  float64
	HoldScale   float64
	HoldOpacity float64
}

func NewLightConfig() LightConfig {
	return LightConfig{
		BackColors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
		HitScale:    1.0,
		HitOpacity:  0.5,
		HoldScale:   1.0,
		HoldOpacity: 1.2,
	}
}
