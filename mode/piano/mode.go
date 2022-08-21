package piano

import (
	"time"

	"github.com/hndada/gosu"
)

var ModePiano4 = gosu.Mode{
	ModeType:       gosu.ModeTypePiano4,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	Mods:           gosu.Mods{},                    // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   gosu.NewSpeedHandler(&SpeedBase),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}

var ModePiano7 = gosu.Mode{
	ModeType:       gosu.ModeTypePiano7,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	Mods:           gosu.Mods{},                    // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   gosu.NewSpeedHandler(&SpeedBase),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}
