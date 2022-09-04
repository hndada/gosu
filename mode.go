package gosu

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/format/osu"
)

var ModeNames = []string{"Piano4", "Piano7", "Drum", "Karaoke"}

type ModeProp struct { // Stands for Mode properties.
	Mode           int
	ChartInfos     []ChartInfo
	Cursor         int                 // Todo: custom chart infos - custom cursor
	Results        map[[16]byte]Result // md5.Size = 16
	LastUpdateTime time.Time
	SpeedHandler   ctrl.F64Handler
	LoadSkin       func()
	NewChartInfo   func(string) (ChartInfo, error)
	NewScenePlay   func(cpath string, rf *osr.Format, sh ctrl.F64Handler) (Scene, error)
	ExposureTime   func(float64) float64
}

// Mode consists of main mode and sub mode.
// For example, sub mode of Piano is Key count (with scratch mode bit adjusted).
const ModeUnknown = -1
const (
	ModePiano4 = iota // ~ 4 Key
	ModePiano7        // 5 ~ Key
	ModeDrum
	ModeKaraoke
)

// Mode determines a mode of chart file by its path.
func ChartFileMode(fpath string) int {
	switch strings.ToLower(filepath.Ext(fpath)) {
	case ".osu":
		mode, keyCount := osu.Mode(fpath)
		switch mode {
		case osu.ModeMania:
			if keyCount <= 4 {
				return ModePiano4
			}
			return ModePiano7
		case osu.ModeTaiko:
			return ModeDrum
		default:
			return ModeUnknown
		}
	case ".ojn", ".bms":
		return ModePiano7
	}
	return ModeUnknown
}
