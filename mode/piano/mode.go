package piano

import (
	"time"

	"github.com/hndada/gosu/mode"
)

var ModePiano4 = mode.Mode{
	ModeType:       mode.ModeTypePiano4,
	ChartInfos:     make([]mode.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]mode.Result), // Zero value.
	Mods:           mode.Mods{},                    // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   mode.NewSpeedHandler(&SpeedBase),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}

var ModePiano7 = mode.Mode{
	ModeType:       mode.ModeTypePiano7,
	ChartInfos:     make([]mode.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]mode.Result), // Zero value.
	Mods:           mode.Mods{},                    // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   mode.NewSpeedHandler(&SpeedBase),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}
