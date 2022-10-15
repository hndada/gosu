package drum

import (
	"time"

	"github.com/hndada/gosu"
)

// var SpeedKeyHandler = gosu.NewSpeedKeyHandler(&SpeedScale)
var ModeDrum = gosu.ModeProp{
	Name:           "Drum",
	Mode:           gosu.ModeDrum,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	// Loads:          []func(){LoadSkin, LoadHandlers},
	// SpeedHandler:   gosu.NewSpeedHandler(&SpeedScale),
	LoadSkin:   LoadSkin,
	SpeedScale: &SpeedScale,
	// SpeedKeyHandler: SpeedKeyHandler,
	NewChartInfo: NewChartInfo,
	NewScenePlay: NewScenePlay,
	ExposureTime: ExposureTime,
}
