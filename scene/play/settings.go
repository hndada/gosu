package play

import (
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
)

type Settings struct {
	MeterWidth  float64 // number of pixels per 1ms
	MeterHeight float64
	Offset      int64
}

var (
	defaultSettings Settings
	settings        Settings
)

func init() {
	initSettings()
	initSkin()
}

func initSettings() {
	defaultSettings = Settings{
		MeterWidth:  4,
		MeterHeight: 50,
		Offset:      -65,
	}
	settings = defaultSettings
}
func DefaultSettings() Settings { return defaultSettings }
func CurrentSettings() Settings { return settings }
func LoadSettings(data string) {
	_, err := toml.Decode(data, &settings)
	if err != nil {
		fmt.Println(err)
	}
}

func Load(fsys fs.FS) {
	data, _ := fs.ReadFile(fsys, "settings.toml")
	LoadSettings(string(data))
	LoadSkin(fsys, defaultSkin)
}
