package play

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu/scene"
)

type settings struct {
	ScoreScale           float64
	ScoreDigitGap        float64
	MeterWidth           float64 // number of pixels per 1ms
	MeterHeight          float64
	Offset               int64
	BackgroundBrightness float64
}

var defaultSettings = settings{
	ScoreScale:           0.65,
	ScoreDigitGap:        0,
	MeterWidth:           4,
	MeterHeight:          50,
	Offset:               -65,
	BackgroundBrightness: scene.Settings.BackgroundBrightness,
}
var Settings = defaultSettings

func ResetSettings() { Settings = defaultSettings }
func LoadSettings(data string) {
	_, err := toml.Decode(data, &Settings)
	if err != nil {
		fmt.Println(err)
	}
	scene.Normalize(&Settings.ScoreScale, 0, 10)
	scene.Normalize(&Settings.ScoreDigitGap, -10, 10)
	scene.Normalize(&Settings.MeterWidth, 0, 5)
	scene.Normalize(&Settings.MeterHeight, 0, 100)
	scene.Normalize(&Settings.Offset, -300, 300)
	scene.Normalize(&Settings.BackgroundBrightness, 0, 1)
}
