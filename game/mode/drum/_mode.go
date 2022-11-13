package drum

import (
	"time"

	"github.com/hndada/gosu/game"
)

var ModeDrum = game.ModeProp{
	Name:           "Drum",
	Mode:           game.ModeDrum,
	ChartInfos:     make([]game.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]game.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	LoadSkin:       LoadSkin,
	SpeedScale:     &SpeedScale,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
	KeySettings:    KeySettings,
}
