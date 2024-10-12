package game

import (
	"io/fs"
	"os"

	"github.com/hndada/gosu/plays"
	"github.com/hndada/gosu/plays/piano"
)

type Args interface{}

type PlayArgs struct {
	ChartFS        fs.FS // Music file exists in the same directory.
	ChartFilename  string
	Mods           plays.Mods
	ReplayFS       fs.FS
	ReplayFilename string
}

var testPlayArgs = PlayArgs{
	// ChartFS:       os.DirFS("C:/Users/hndada/Documents/GitHub/gosu/cmd/gosu/music/nekodex - circles!"),
	// ChartFilename: "nekodex - circles! (MuangMuangE) [Hard].osu",
	ChartFS:       os.DirFS("C:/Users/hndada/Documents/GitHub/gosu/cmd/gosu/music/cYsmix - triangles"),
	ChartFilename: "cYsmix - triangles (MuangMuangE) [Easy].osu",
	Mods:          piano.Mods{},
	// ReplayFS       fs.FS
	// ReplayFilename string
}
