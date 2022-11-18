package mode

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Settings struct {
	VolumeMusic          float64
	VolumeSound          float64
	Offset               int64
	BackgroundBrightness float64

	ScoreScale    float64
	ScoreDigitGap float64
	MeterUnit     float64 // number of pixels per 1ms
	MeterHeight   float64
}

var (
	DefaultSettings = Settings{
		VolumeMusic:          0.25,
		VolumeSound:          0.25,
		Offset:               -65,
		BackgroundBrightness: 0.6,

		ScoreScale:    0.65,
		ScoreDigitGap: 0,
		MeterUnit:     4,
		MeterHeight:   50,
	}
	UserSettings = DefaultSettings
)

func LoadSettings(data string) {
	_, err := toml.Decode(data, &UserSettings)
	if err != nil {
		fmt.Println(err)
	}
	Normalize(&UserSettings.VolumeMusic, 0, 1)
	Normalize(&UserSettings.VolumeSound, 0, 1)
	Normalize(&UserSettings.Offset, -300, 300)
	Normalize(&UserSettings.BackgroundBrightness, 0, 1)

	Normalize(&UserSettings.ScoreScale, 0, 5)
	Normalize(&UserSettings.ScoreDigitGap, -10, 10)
	Normalize(&UserSettings.MeterUnit, 0, 5)
	Normalize(&UserSettings.MeterHeight, 0, 100)
}

type Number interface{ int | int64 | float64 }

func Normalize[V Number](v *V, min, max V) {
	if *v > max {
		*v = max
	}
	if *v < min {
		*v = min
	}
}
