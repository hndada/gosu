package piano

import (
	"time"

	"github.com/hndada/gosu"
)

var ModePiano4 = gosu.ModeProp{
	Mode:           gosu.ModePiano4,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   gosu.NewSpeedHandler(&SpeedScale),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}

var ModePiano7 = gosu.ModeProp{
	Mode:           gosu.ModePiano7,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   gosu.NewSpeedHandler(&SpeedScale),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}
