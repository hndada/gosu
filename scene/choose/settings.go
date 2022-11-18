package choose

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu/scene"
)

type settings struct {
	// Group1               int
	// Group2               int
	// Sort                 int
	// Filter               int
	BackgroundBrightness float64
}

var defaultSettings = settings{
	BackgroundBrightness: scene.Settings.BackgroundBrightness,
}
var Settings = defaultSettings

func ResetSettings() { Settings = defaultSettings }
func LoadSettings(data string) {
	_, err := toml.Decode(data, &Settings)
	if err != nil {
		fmt.Println(err)
	}
	scene.Normalize(&Settings.BackgroundBrightness, 0, 1)
}

// Todo: Settings -> settings, CurrentSettings() -> Settings(),
// settings -> userSettings?
// `toml:"BackgroundBrightness_Choose"`
