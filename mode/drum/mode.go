package drum

import (
	"time"

	"github.com/hndada/gosu"
)

var ModeDrum = gosu.ModeProp{
	Name:           "Drum",
	Mode:           gosu.ModeDrum,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	LoadSkin:       LoadSkin,
	SpeedScale:     &SpeedScale,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
	KeySettings:    KeySettings,
}
