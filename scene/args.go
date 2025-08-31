package scene

import (
	"io/fs"

	"github.com/hndada/gosu/game"
)

type Args interface{}

type PlayArgs struct {
	ChartFS        fs.FS // Music file exists in the same directory.
	ChartFilename  string
	Mods           game.Mods
	ReplayFS       fs.FS
	ReplayFilename string
}
