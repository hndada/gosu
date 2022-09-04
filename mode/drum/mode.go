package drum

import (
	"time"

	"github.com/hndada/gosu"
)

var ModeDrum = gosu.ModeProp{
	Mode:           gosu.ModeDrum,
	ChartInfos:     make([]gosu.ChartInfo, 0),      // Zero value.
	Results:        make(map[[16]byte]gosu.Result), // Zero value.
	LastUpdateTime: time.Time{},                    // Zero value.
	SpeedHandler:   gosu.NewSpeedHandler(&SpeedScale),
	LoadSkin:       LoadSkin,
	NewChartInfo:   NewChartInfo,
	NewScenePlay:   NewScenePlay,
	ExposureTime:   ExposureTime,
}

func LoadSkin()
