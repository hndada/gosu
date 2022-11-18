package new

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// derived settings goes pointer with unexported.
type Settings struct {
	MusicRoots  []string
	WindowSize  int
	VolumeMusic float64
}

const (
	WindowSizeStandard = iota
	WindowSizeFull
)

var (
	DefaultSettings = Settings{
		MusicRoots:  []string{"music"},
		WindowSize:  WindowSizeStandard,
		VolumeMusic: 0.25,
	}
	UserSettings = DefaultSettings
)

func (settings *Settings) Load(data string) {
	_, err := toml.Decode(data, settings)
	if err != nil {
		fmt.Println(err)
	}
	// post-processing
}
