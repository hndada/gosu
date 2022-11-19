package choose

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/mode"
)

type Settings struct {
	// Group1               int
	// Group2               int
	// Sort                 int
	// Filter               int
	backgroundBrightness *float64
}

var (
	DefaultSettings = Settings{}
	UserSettings    = DefaultSettings
	S               = &UserSettings
)

func init() {
	DefaultSettings.process()
	UserSettings.process()
	DefaultSkin.Load(defaultskin.FS)
}
func (settings *Settings) Load(data string) {
	_, err := toml.Decode(data, settings)
	if err != nil {
		fmt.Println(err)
	}
	defer settings.process()
}
func (settings *Settings) process() {
	MS := &mode.UserSettings
	settings.backgroundBrightness = &MS.BackgroundBrightness
}
