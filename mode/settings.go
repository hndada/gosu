package mode

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/input"
)

type Settings struct {
	VolumeMusic          float64
	VolumeSound          float64
	Offset               int64
	BackgroundBrightness float64

	MeterUnit     float64 // number of pixels per 1ms
	MeterHeight   float64
	ScoreScale    float64
	ScoreDigitGap float64
}

var (
	DefaultSettings = Settings{
		VolumeMusic:          0.25,
		VolumeSound:          0.25,
		Offset:               -65,
		BackgroundBrightness: 0.6,

		MeterUnit:     4,
		MeterHeight:   50,
		ScoreScale:    0.65,
		ScoreDigitGap: 0,
	}
	UserSettings = DefaultSettings
)

func init() {
	DefaultSettings.process()
	DefaultSkin.Load(defaultskin.FS)
}

func (settings *Settings) Load(data string) {
	_, err := toml.Decode(data, &UserSettings)
	if err != nil {
		fmt.Println(err)
	}
	defer settings.process()

	Normalize(&UserSettings.VolumeMusic, 0, 1)
	Normalize(&UserSettings.VolumeSound, 0, 1)
	Normalize(&UserSettings.Offset, -300, 300)
	Normalize(&UserSettings.BackgroundBrightness, 0, 1)

	Normalize(&UserSettings.MeterUnit, 0, 5)
	Normalize(&UserSettings.MeterHeight, 0, 100)
	Normalize(&UserSettings.ScoreScale, 0, 5)
	Normalize(&UserSettings.ScoreDigitGap, -0.05, 0.05)
}

// process processes from raw settings to screen size-adjusted settings.
func (settings *Settings) process() {
	settings.ScoreDigitGap *= ScreenSizeX
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
func NormalizeKeys(names []string) {
	keys := input.NamesToKeys(names)
	m := make(map[input.Key]bool)
	for i, k := range keys {
		if m[k] {
			fmt.Printf("some keys are duplicated: %v\n", names)
			keys[i] = input.KeyNone
		}
		m[k] = true
	}
}
