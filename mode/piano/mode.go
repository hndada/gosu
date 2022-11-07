package piano

import (
	"time"

	"github.com/hndada/gosu"
)

var ModePiano4 = gosu.ModeProp{
	Name:           "Piano4",
	Mode:           gosu.ModePiano4,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	LoadSkin:       LoadSkin,
	SpeedScale:     &SpeedScale,
	Settings: map[string]*float64{
		"TailExtraTime": &TailExtraTime,
	},
	NewChartInfo: NewChartInfo,
	NewScenePlay: NewScenePlay,
	ExposureTime: ExposureTime,
	KeySettings:  KeySettings,
}

var ModePiano7 = gosu.ModeProp{
	Name:           "Piano7",
	Mode:           gosu.ModePiano7,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	LoadSkin:       LoadSkin,
	SpeedScale:     &SpeedScale,
	Settings: map[string]*float64{
		"TailExtraTime": &TailExtraTime,
	},
	NewChartInfo: NewChartInfo,
	NewScenePlay: NewScenePlay,
	ExposureTime: ExposureTime,
	KeySettings:  KeySettings,
}
