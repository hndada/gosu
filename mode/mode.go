package mode

import (
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/format/osu"
)

// Mode consists of main mode + sub mode.
// Piano mode's sub mode is Key count (with scratch mode bit adjusted), for example.
const (
	ModePiano4 = iota // ~ 4 Key
	ModePiano7        // 5 ~ Key
	ModeDrum
	ModeKaraoke // aka jjava
)
const ModeDefault = ModePiano4
const ModeUnknown = -1

// Mode determines a mode of chart file by its path.
// Todo: should I make a new type Mode?
func Mode(fpath string) int {
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
