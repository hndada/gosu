package gosu

import "github.com/hndada/gosu/mode"

var MusicPath = "music"

const (
	screenSizeX = mode.ScreenSizeX
	screenSizeY = mode.ScreenSizeY
)

var (
	WindowSizeX = 1600
	WindowSizeY = 900
)

// Todo: reset all tick-dependent variables.
// They are mostly at drawer.go or play.go, settings.go
// Keyword: TimeToTick
func SetTPS() {

}
