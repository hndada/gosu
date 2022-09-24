package piano

import (
	"time"

	"github.com/hndada/gosu"
)

var (
	SpeedKeyHandler = gosu.NewSpeedKeyHandler(&SpeedScale)
	// Piano4SpeedKeyHandler = gosu.NewSpeedKeyHandler(&SpeedScale)
	// Piano7SpeedKeyHandler = gosu.NewSpeedKeyHandler(&SpeedScale)
)
var ModePiano4 = gosu.ModeProp{
	Name:           "Piano4",
	Mode:           gosu.ModePiano4,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	// SpeedHandler:   gosu.NewSpeedHandler(&SpeedScale),
	LoadSkin: LoadSkin,
	// Loads:        []func(){LoadSkin, LoadHandlers},
	SpeedScale:      &SpeedScale,
	SpeedKeyHandler: SpeedKeyHandler,
	NewChartInfo:    NewChartInfo,
	NewScenePlay:    NewScenePlay,
	ExposureTime:    ExposureTime,
}

var ModePiano7 = gosu.ModeProp{
	Name:           "Piano7",
	Mode:           gosu.ModePiano7,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	// SpeedHandler:   gosu.NewSpeedHandler(&SpeedScale),
	LoadSkin: LoadSkin,
	// Loads:        []func(){},
	SpeedScale:      &SpeedScale,
	SpeedKeyHandler: SpeedKeyHandler,
	NewChartInfo:    NewChartInfo,
	NewScenePlay:    NewScenePlay,
	ExposureTime:    ExposureTime,
}
