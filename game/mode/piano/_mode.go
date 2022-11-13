package piano

import (
	"time"

	"github.com/hndada/gosu/game"
)

var ModePiano4 = game.ModeProp{
	Name:           "Piano4",
	Mode:           game.ModePiano4,
	ChartInfos:     make([]game.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]game.Result), // Zero value.
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

var ModePiano7 = game.ModeProp{
	Name:           "Piano7",
	Mode:           game.ModePiano7,
	ChartInfos:     make([]game.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]game.Result), // Zero value.
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
