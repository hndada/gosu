package choose

import (
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
func (settings *Settings) Load(src Settings) {
	*settings = src
	defer settings.process()
}
func (settings *Settings) process() {
	settings.backgroundBrightness = &mode.S.BackgroundBrightness
}
