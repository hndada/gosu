package gosu

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/format/osu"
)

// Todo: implement Mode() in each mode.
type Mode struct {
	Type           int
	ChartInfos     []ChartInfo
	Cursor         int                 // Todo: custom chart infos - custom cursor
	Results        map[[16]byte]Result // md5.Size = 16
	Mods           Mods
	LastUpdateTime time.Time
	SpeedHandler   ctrl.F64Handler
	LoadSkin       func()
	NewChartInfo   func(string, Mods) (ChartInfo, error)
	NewScenePlay   func(string, Mods, *osr.Format) (Scene, error)
	ExposureTime   func(float64) float64
}

// Mode consists of main mode + sub mode.
// Piano mode's sub mode is Key count (with scratch mode bit adjusted), for example.
// type ModeType int

const (
	ModeTypePiano4 = iota // ~ 4 Key
	ModeTypePiano7        // 5 ~ Key
	ModeTypeDrum
	ModeTypeKaraoke // aka jjava
)

// const DefaultModeType ModeType = ModeTypePiano4
const ModeTypeUnknown = -1

// Mode determines a mode of chart file by its path.
// Todo: should I make a new type Mode?
func FileModeType(fpath string) int {
	switch strings.ToLower(filepath.Ext(fpath)) {
	case ".osu":
		mode, keyCount := osu.Mode(fpath)
		switch mode {
		case osu.ModeMania:
			if keyCount <= 4 {
				return ModeTypePiano4
			}
			return ModeTypePiano7
		case osu.ModeTaiko:
			return ModeTypeDrum
		default:
			return ModeTypeUnknown
		}
	case ".ojn", ".bms":
		return ModeTypePiano7
	}
	return ModeTypeUnknown
}
